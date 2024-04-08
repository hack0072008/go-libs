package ssh2

import (
	"github.com/pkg/sftp"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var sf *sftp.Client

func init() {
	host := "123.249.27.227"
	port := 22
	user := "root"
	paasword := "Pass@Ecs"
	timeout := 60 * time.Second
	cli, _ := NewSftpClient(user, paasword, host, port, timeout)
	sf = cli
}

func TestUploadLocalFileBySftp(t *testing.T) {
	localFilePath := "/cpaas/installer-v3.10.1-x86.tar"
	remotePath := "/cpaas/"
	err := UploadLocalFileBySftp(sf, localFilePath, remotePath)
	assert.Nil(t, err)
}
