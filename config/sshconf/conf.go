package sshconf

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"golang.org/x/crypto/ssh"
)

type SSHConfig struct {
	Host string `json:"host"`
	Port int    `json:"port"`
	User string `json:"user"`
	Pass string `json:"pass"`
}

func (sc SSHConfig) ToClientConfig() *ssh.ClientConfig {
	ops := ssh.ClientConfig{
		User: sc.User,
		Auth: []ssh.AuthMethod{
			ssh.Password(sc.Pass),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	return &ops
}

func (sc SSHConfig) CombinedHost() string {
	return fmt.Sprintf("%s:%d", sc.Host, sc.Port)
}

func (sc SSHConfig) CombinedHostAndUser() string {
	return fmt.Sprintf("%s@%s:%d", sc.User, sc.Host, sc.Port)
}

func ReadConfig() ([]byte, error) {
	return ioutil.ReadFile("./config/sshconf/config.json")
}

func NewConfig() (*SSHConfig, error) {
	data, err := ReadConfig()
	if err != nil {
		return nil, err
	}

	var x SSHConfig
	err = json.Unmarshal(data, &x)
	return &x, err
}
