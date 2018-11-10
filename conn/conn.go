package conn

import (
	"github.com/TechCatsLab/motion/config"
	"golang.org/x/crypto/ssh"
)

var Client *ssh.Client

func init() {
	cc := &ssh.ClientConfig{
		User: config.SSHConf.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(config.SSHConf.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	host := config.SSHConf.Address + ":" + config.SSHConf.Port

	var err error

	Client, err = ssh.Dial("tcp", host, cc)
	if err != nil {
		panic(err)
	}
}

func Run(c string) ([]byte, error) {
	s, err := Client.NewSession()
	if err != nil {
		return nil, err
	}
	defer s.Close()

	// return s.Output("yes | " + c)
	return s.Output(c)
}
