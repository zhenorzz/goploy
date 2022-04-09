// Copyright 2022 The Goploy Authors. All rights reserved.
// Use of this source code is governed by a GPLv3-style
// license that can be found in the LICENSE file.

package utils

import (
	"bytes"
	"fmt"
	"golang.org/x/crypto/ssh"
	"io/ioutil"
	"net"
	"strings"
	"time"
)

// SSHConfig -
type SSHConfig struct {
	User         string
	Password     string
	Path         string
	Host         string
	Port         int
	JumpUser     string
	JumpPassword string
	JumpPath     string
	JumpHost     string
	JumpPort     int
}

func (sshConfig SSHConfig) Dial() (*ssh.Client, error) {
	clientConfig, err := sshConfig.getConfig(sshConfig.User, sshConfig.Password, sshConfig.Path)
	if err != nil {
		return nil, err
	}
	// connect to ssh
	sshClient, err := ssh.Dial("tcp", sshConfig.addr(), clientConfig)
	if err != nil {
		return nil, err
	}

	if sshConfig.JumpHost != "" {
		conn, err := sshClient.Dial("tcp", sshConfig.jumpAddr())
		if err != nil {
			return nil, err
		}

		targetConfig, err := sshConfig.getConfig(sshConfig.JumpUser, sshConfig.JumpPassword, sshConfig.JumpPath)
		if err != nil {
			return nil, err
		}
		ncc, chans, reqs, err := ssh.NewClientConn(conn, sshConfig.jumpAddr(), targetConfig)
		if err != nil {
			return nil, err
		}

		sshClient = ssh.NewClient(ncc, chans, reqs)
	}

	return sshClient, err
}

func (sshConfig SSHConfig) ToRsyncOption() string {
	if sshConfig.JumpHost == "" {
		return fmt.Sprintf("ssh -o StrictHostKeyChecking=no -p %d -i %s", sshConfig.Port, sshConfig.Path)
	} else {
		return fmt.Sprintf("ssh -o StrictHostKeyChecking=no -o ProxyCommand='ssh -o StrictHostKeyChecking=no -W %%h:%%p -i %s -p %d %s@%s' -p %d -i %s", sshConfig.JumpPath, sshConfig.JumpPort, sshConfig.JumpUser, sshConfig.JumpHost, sshConfig.Port, sshConfig.Path)
	}
}

// version|cpu cores|mem

func (sshConfig SSHConfig) GetOSInfo() string {
	osInfoScript := `cat /etc/os-release | grep "PRETTY_NAME" | awk -F\" '{print $2}' && cat /proc/cpuinfo  | grep "processor" | wc -l && cat /proc/meminfo | grep MemTotal | awk '{print $2}'`
	client, err := sshConfig.Dial()
	if err != nil {
		return ""
	}
	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return ""
	}
	defer session.Close()

	var sshOutbuf, sshErrbuf bytes.Buffer
	session.Stdout = &sshOutbuf
	session.Stderr = &sshErrbuf
	if err := session.Run(osInfoScript); err != nil {
		return ""
	}

	// version|cpu cores|mem
	return strings.Replace(strings.Trim(sshOutbuf.String(), "\n"), "\n", "|", -1)
}

func (sshConfig SSHConfig) getConfig(user, password, path string) (*ssh.ClientConfig, error) {
	pemBytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var signer ssh.Signer
	if password == "" {
		signer, err = ssh.ParsePrivateKey(pemBytes)
	} else {
		signer, err = ssh.ParsePrivateKeyWithPassphrase(pemBytes, []byte(password))
	}
	if err != nil {
		return nil, err
	}

	auth := make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.PublicKeys(signer))

	config := ssh.Config{
		Ciphers: []string{"aes128-ctr", "aes192-ctr", "aes256-ctr", "aes128-gcm@openssh.com", "arcfour256", "arcfour128", "aes128-cbc", "3des-cbc", "aes192-cbc", "aes256-cbc"},
	}

	return &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: 30 * time.Second,
		Config:  config,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}, nil
}

func (sshConfig SSHConfig) jumpAddr() string {
	return fmt.Sprintf("%s:%d", sshConfig.JumpHost, sshConfig.JumpPort)
}

func (sshConfig SSHConfig) addr() string {
	return fmt.Sprintf("%s:%d", sshConfig.Host, sshConfig.Port)
}
