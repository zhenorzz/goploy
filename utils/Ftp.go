package utils

import (
	"crypto/tls"
	"github.com/jlaffaye/ftp"
	"net/url"
	"strings"
	"time"
)

type FTPConfig struct {
	TLS      bool
	Addr     string
	Username string
	Password string
	Path     string
}

func (ftpConfig FTPConfig) Dial() (*ftp.ServerConn, error) {

	opts := []ftp.DialOption{
		ftp.DialWithTimeout(10 * time.Second),
	}

	if ftpConfig.TLS == true {
		opts = append(opts, ftp.DialWithTLS(&tls.Config{
			InsecureSkipVerify: true,
		}))
	}

	c, err := ftp.Dial(ftpConfig.Addr, opts...)
	if err != nil {
		return nil, err
	}

	if ftpConfig.Username != "" {
		if err = c.Login(ftpConfig.Username, ftpConfig.Password); err != nil {
			return nil, err
		}
	}

	if ftpConfig.Path != "" {
		if err = c.ChangeDir(ftpConfig.Path); err != nil {
			return nil, err
		}
	}
	return c, nil
}

func (ftpConfig FTPConfig) DialFromURL(_url string) (*ftp.ServerConn, error) {
	u, err := url.Parse(_url)
	if err != nil {
		return nil, err
	}
	if u.Scheme == "ftps" {
		ftpConfig.TLS = true
	}
	h := strings.Split(u.Host, ":")
	if len(h) == 1 {
		u.Host += ":21"
	}
	ftpConfig.Addr = u.Host
	ftpConfig.Username = u.User.Username()
	if ftpConfig.Username == "" {
		ftpConfig.Username = "anonymous"
		ftpConfig.Password = "anonymous@domain.com"
	} else {
		pwd, _ := u.User.Password()
		ftpConfig.Password = pwd
	}
	ftpConfig.Path = u.Path
	return ftpConfig.Dial()
}
