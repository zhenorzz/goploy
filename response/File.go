package response

import (
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"net/http"
	"strconv"
)

type SftpFile struct {
	Client   *ssh.Client
	Filename string
}

//JSON response
func (sf SftpFile) Write(w http.ResponseWriter) error {
	defer sf.Client.Close()

	sftpClient, err := sftp.NewClient(sf.Client)
	if err != nil {
		return err
	}
	defer sftpClient.Close()

	srcFile, err := sftpClient.Open(sf.Filename) //远程
	if err != nil {
		return err
	}
	defer srcFile.Close()

	fileStat, err := srcFile.Stat()
	if err != nil {
		return err
	}
	w.Header().Set("Content-Disposition", "attachment; filename="+fileStat.Name())
	w.Header().Set("Content-Type", "application/octet-stream")
	w.Header().Set("Content-Length", strconv.FormatInt(fileStat.Size(), 10))
	_, err = io.Copy(w, srcFile)

	if err != nil {
		return err
	}
	return nil
}
