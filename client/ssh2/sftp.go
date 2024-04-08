package ssh2

import (
	"fmt"
	"github.com/hack0072008/go-libs/log"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"io"
	"net"
	"os"
	"path"
	"time"
)

/*
 create new sftp client
*/
func NewSftpClient(user, password, host string, port int, timeoutSec time.Duration) (*sftp.Client, error) {
	var (
		auth         []ssh.AuthMethod
		addr         string
		clientConfig *ssh.ClientConfig
		sshClient    *ssh.Client
		sftpClient   *sftp.Client
		err          error
	)
	// get auth method
	auth = make([]ssh.AuthMethod, 0)
	auth = append(auth, ssh.Password(password))

	clientConfig = &ssh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: timeoutSec * time.Second,
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// connet to ssh
	addr = fmt.Sprintf("%s:%d", host, port)

	if sshClient, err = ssh.Dial("tcp", addr, clientConfig); err != nil {
		return nil, err
	}

	// create sftp client
	if sftpClient, err = sftp.NewClient(sshClient); err != nil {
		return nil, err
	}

	return sftpClient, nil
}

/*
 sftp client to upload file
*/
func UploadLocalFileBySftp(objClient *sftp.Client, srcFilePath, remotePath string) error {
	var fileName = path.Base(srcFilePath)
	var remoteFilePath = path.Join(remotePath, fileName)

	// open src file
	srcFile, err := os.Open(srcFilePath)
	if err != nil {
		log.Errorf("upload source file: %+v open error: %+v", srcFile, err)
		return err
	}
	defer srcFile.Close()

	// open remote file
	dstFile, err := objClient.Create(remoteFilePath)
	if err != nil {
		log.Errorf("upload remote file: %+v create error: %+v", remoteFilePath, err)
		return err
	}
	defer dstFile.Close()

	// copy
	size, err := io.Copy(dstFile, srcFile)
	if err != nil {
		log.Errorf("upload source file: %+v remote file: %+v copy error: %+v", srcFilePath, remoteFilePath, err)
		return err
	}
	log.Infof("upload source file: %+v remote file: %+v size: %s finish", srcFilePath, remoteFilePath, FormatFileSize(size))

	return nil
}

// 字节的单位转换 保留两位小数
func FormatFileSize(s int64) (size string) {
	if s < 1024 {
		return fmt.Sprintf("%.2fB", float64(s)/float64(1))
	} else if s < (1024 * 1024) {
		return fmt.Sprintf("%.2fKB", float64(s)/float64(1024))
	} else if s < (1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fMB", float64(s)/float64(1024*1024))
	} else if s < (1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fGB", float64(s)/float64(1024*1024*1024))
	} else if s < (1024 * 1024 * 1024 * 1024 * 1024) {
		return fmt.Sprintf("%.2fTB", float64(s)/float64(1024*1024*1024*1024))
	} else { // if s < (1024 * 1024 * 1024 * 1024 * 1024 * 1024)
		return fmt.Sprintf("%.2fEB", float64(s)/float64(1024*1024*1024*1024*1024))
	}
}
