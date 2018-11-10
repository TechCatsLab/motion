package deploy

import (
	"encoding/json"
	"io/ioutil"
)

var EnvConfig struct {
	Format  string
	OS      string
	Servers []*server

	haveDocker bool
	haveNginx  bool
}

func init() {
	data, err := ioutil.ReadFile("./config/serverconfig.json")
	if err != nil {
		panic(err)
	}
	err = json.Unmarshal(data, &EnvConfig)
	if err != nil {
		panic(err)
	}

	ok, err := haveDocker()
	if err != nil {
		panic(err)
	}
	EnvConfig.haveDocker = ok

	ok, err = haveNginx()
	if err != nil {
		panic(err)
	}
	EnvConfig.haveNginx = ok
}

func haveDocker() (bool, error) {
	return haveSomething("docker version", "Version")
}

func haveNginx() (bool, error) {
	return haveSomething("netstat -anput | grep nginx", "nginx")
}
