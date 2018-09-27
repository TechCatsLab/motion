package conn

import (
	"golang.org/x/crypto/ssh"

	"fmt"

	"github.com/zh1014/ssh/config"
)

// Client is initialized for create session
var Client *ssh.Client

func init() {
	cc := &ssh.ClientConfig{
		User:            config.SSHConf.User,
		Auth:            []ssh.AuthMethod{ssh.Password(config.SSHConf.Password)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	host := config.SSHConf.Address + ":" + config.SSHConf.Port
	var err error
	Client, err = ssh.Dial("tcp", host, cc)
	if err != nil {
		panic(fmt.Sprintf("failed to create client: %s", err.Error()))
	}
}
