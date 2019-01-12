package ssh

import (
	"github.com/Zeng1999/motion/config"
	"github.com/Zeng1999/motion/config/sshconf"
	"github.com/Zeng1999/motion/constant"
	"golang.org/x/crypto/ssh"
)

type SSHConn struct {
	cli  *ssh.Client
	conf *sshconf.SSHConfig
}

func NewConn(conf config.ConnConfig) *SSHConn {
	con := conf.(sshconf.SSHConfig)
	return &SSHConn{
		conf: &con,
	}
}

func (con *SSHConn) Dial() error {
	if con.cli != nil {
		return constant.ConnHasBeDialed
	}
	c, err := ssh.Dial("tcp", con.conf.CombinedHost(), con.conf.ToClientConfig())
	con.cli = c
	return err
}

func (con *SSHConn) OutPut(cmd string) ([]byte, error) {
	if con.cli == nil {
		return nil, constant.ConnNonDial
	}

	ses, err := con.cli.NewSession()
	if err != nil {
		return nil, err
	}
	defer ses.Close()

	return ses.Output(cmd)
}

func (con *SSHConn) Close() error {
	return con.cli.Close()
}

func (con *SSHConn) NewSession() (*ssh.Session, error) {
	return con.cli.NewSession()
}

func (con SSHConn) GetConf() *sshconf.SSHConfig {
	return con.conf
}

func (con SSHConn) GetClient() *ssh.Client {
	return con.cli
}
