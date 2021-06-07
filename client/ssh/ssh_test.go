package ssh


import (
	"context"
	"testing"
)

func TestSSHClient(t *testing.T) {
	cli := SSHClient{
		Host:     "192.168.203.187",
		Port:     22,
		User:     "root",
		Password: "root",
		//PrivateKeyFile: "d:/.../.../id_rsa",
		Timeout: 60,
		Config:  nil,
		//Env:      []string{"a=b", "version=v1.1.1"},
		Cmd:     "",
		Stdout:  "",
		Stderr:  "",
		Session: nil,
	}
	ctx := context.Background()
	cli.Run(ctx, "cat /etc/profile")
	cli.Run(ctx, "df -h")
	cli.Run(ctx, "free -m")

}


