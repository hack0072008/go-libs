package ssh

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"regexp"
	"time"

	"github.com/hack0072008/go-libs/log"
	"golang.org/x/crypto/ssh"
)

type SSHClient struct {
	Host           string
	Port           int
	User           string
	Password       string
	PrivateKeyFile string
	Timeout        int
	Config         *ssh.ClientConfig
	Env            []string
	Cmd            string
	Stdout         string
	Stderr         string

	Client  *ssh.Client
	Session *ssh.Session
}

func (client *SSHClient) newClient(ctx context.Context) error {
	err := errors.New("")
	if client.Client, err = ssh.Dial("tcp", fmt.Sprintf("%s:%d", client.Host, client.Port), client.Config); err != nil {
		log.Errorf("Failed to dial: ip[%s] port[%d] error:%s", client.Host, client.Port, err)
		return err
	}
	return nil
}

func (client *SSHClient) newSession(ctx context.Context) error {

	if client.Client == nil {
		if err := client.newClient(ctx); err != nil {
			return err
		}
	}

	err := errors.New("")
	client.Session, err = client.Client.NewSession()
	if err != nil {
		log.Errorf("Failed to create session: %s", err)
		return err
	}

	modes := ssh.TerminalModes{
		// ssh.ECHO:          0,     // disable echoing
		ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := client.Session.RequestPty("xterm", 80, 40, modes); err != nil {
		client.Session.Close()
		log.Errorf("set ssh terminal RequestPty Failed, error:%s", err)
		return err
	}
	log.Infof("create new session success!")

	return nil
}

func (client *SSHClient) closeSession(ctx context.Context) error {
	if client.Session != nil {
		// https://stackoverflow.com/questions/60879023/getting-eof-as-error-in-golang-ssh-session-close
		client.Session.Close()
		log.Infof("close session success!")
	}
	return nil
}

func (client *SSHClient) PublicKeyFile(file string) ssh.Signer {
	buffer, err := ioutil.ReadFile(file)
	if err != nil {
		return nil
	}

	key, err := ssh.ParsePrivateKey(buffer)
	if err != nil {
		return nil
	}
	return key
}

func (client *SSHClient) setClientConfig(ctx context.Context, config *ssh.ClientConfig) error {
	auth := []ssh.AuthMethod{}

	if client.Password != "" {
		auth = append(auth, ssh.Password(client.Password))
	}
	if client.PrivateKeyFile != "" {
		auth = append(auth, ssh.PublicKeys(client.PublicKeyFile(client.PrivateKeyFile)))
	}

	if config == nil {
		client.Config = &ssh.ClientConfig{
			Config: ssh.Config{},
			User:   client.User,
			Auth:   auth,
			HostKeyCallback: (func(hostname string, remote net.Addr, key ssh.PublicKey) error {
				return nil
			}),
			Timeout: time.Duration(client.Timeout) * time.Second,
		}
	} else {
		client.Config = config
	}

	log.Infof("set ssh client config success!")

	return nil
}

func (client *SSHClient) Run(ctx context.Context, cmd string) error {

	defer client.closeSession(ctx)

	// need remote sshd support set env, /etc/ssh/sshd_config[AcceptEnv]: https://groups.google.com/forum/#!topic/golang-nuts/OlEJtOjxdDw
	// set env
	/*	for _, env := range client.Env {
		variable := strings.Split(env, "=")
		if len(variable) != 2 {
			continue
		}

		if err := client.Session.Setenv(variable[0], variable[1]); err != nil {
			return err
		}
	}*/

	client.setClientConfig(ctx, nil)

	if err := client.newSession(ctx); err != nil {
		log.Errorf("create session failed! error:%s", err)
		return err
	}

	var (
		stdoutBuffer bytes.Buffer
		stderrBuffer bytes.Buffer
	)

	client.Session.Stdout = &stdoutBuffer
	client.Session.Stderr = &stderrBuffer

	// run command
	if cmd != "" {
		if err := client.Session.Run(cmd); err != nil {
			//fmt.Printf("host[%s] port[%d] user[%s] execute cmd[%s] failed! out:\n%v error:%s", client.Host, client.Port, client.User, cmd, client.Stderr, err)
			log.Errorf("host[%s] port[%d] user[%s] execute cmd[%s] failed! out:%v error:%s", client.Host, client.Port, client.User, cmd, client.Stderr, err)
			return err
		}
		client.Stdout = stdoutBuffer.String()
		client.Stderr = stderrBuffer.String()
		//fmt.Printf("host[%s] port[%d] user[%s] execute cmd[%s] success! \nout:\n%v err:\n%v", client.Host, client.Port, client.User, cmd, client.Stdout, client.Stderr)
		log.Debugf("host[%s] port[%d] user[%s] execute cmd[%s] success! out:%v err:%v", client.Host, client.Port, client.User, cmd, client.Stdout, client.Stderr)
	} else {
		if err := client.Session.Run(client.Cmd); err != nil {
			log.Errorf("host[%s] port[%d] user[%s] execute cmd[%s] failed! out:[%s] error:[%s]", client.Host, client.Port, client.User, client.Cmd, client.Stderr, err)
			return err
		}
		client.Stdout = stdoutBuffer.String()
		client.Stderr = stderrBuffer.String()
		//fmt.Printf("host[%s] port[%d] user[%s] execute cmd[%s] success! \nout:\n[%v]err:\n[%v]", client.Host, client.Port, client.User, client.Cmd, client.Stdout, client.Stderr)
		log.Debugf("host[%s] port[%d] user[%s] execute cmd[%s] success! out:[%v] err:[%v]", client.Host, client.Port, client.User, client.Cmd, client.Stdout, client.Stderr)
	}
	return nil
}

func TrimSpaceNewlineInString(s string, expr string) string {
	re, _ := regexp.Compile(expr)
	return re.ReplaceAllString(s, "")
}
