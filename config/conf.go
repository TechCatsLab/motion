package config

import (
	"encoding/json"
	"io/ioutil"
)

type sshConf struct {
	User     string
	Address  string
	Port     string
	Password string
}

// SSHConf provide an approch to access to config information
var SSHConf *sshConf

func init() {
	data, err := ioutil.ReadFile("./config/ssh.json")
	if err != nil {
		panic(err)
	}

	err = json.Unmarshal(data, &SSHConf)
	if err != nil {
		panic(err)
	}
}
