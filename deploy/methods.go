package deploy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/TechCatsLab/motion/conn/ssh"
	"github.com/TechCatsLab/motion/conn/ssh/config"
)

var EnvConfig struct {
	Format  string
	OS      string
	Servers []*server

	haveDocker bool
	haveNginx  bool
}

var con ssh.Client

func init() {
	con = ssh.New()
	data, err := ioutil.ReadFile("./serverconfig.json")
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

func zip(dir string, index int) error {
	return exec.Command("tar", "-zcvf", fmt.Sprintf("~/ready%d.tar.gz", index), dir).Run()
}

func rename(file string, index int) error {
	return exec.Command("mv", file, fmt.Sprintf("~/ready%d.tar.gz", index)).Run()
}

func writeScript(name, from, to string) error {
	temp := fmt.Sprintf("%s@%s:"+to, config.SSHConf.User, config.SSHConf.Address)
	return ioutil.WriteFile(name, []byte(fmt.Sprintf(EnvConfig.Format, from, temp, config.SSHConf.Password)), os.ModePerm)
}

func writeConfig(data string) error {
	return ioutil.WriteFile("./nginx_config.conf", []byte(data), os.ModePerm)
}

func portIsUsed(port int) (bool, error) {
	if port < 0 || port > 100000 {
		return true, nil
	}
	hp := fmt.Sprintf(":%d ", port)
	out, err := output("netstat -anput")
	return strings.Index(string(out), hp) > -1, err
}

func useConfig() error {
	return run("nginx -c /root/nginx_config.conf -s reload")
}

func runContainer(image string, listenPort, exportPrt int, args map[string]string) ([]byte, error) {
	str := "docker run -d -p %d:%d %s %s"
	arg := ""
	if args != nil {
		arg = "-e "
		for key, value := range args {
			arg += key + "=" + value + " "
		}
	}
	return output(fmt.Sprintf(str, exportPrt, listenPort, arg, image))
}

func pullImage(what string) error {
	return run("docker pull " + what)
}

func unzip(index int) error {
	return run(fmt.Sprintf("tar -zxvf ~/ready%d.tar.gz", index))
}

func haveSomething(cd, flag string) (bool, error) {
	out, err := output(cd)
	if err != nil {
		if strings.Index(err.Error(), "exited with status 127") == -1 {
			return false, err
		}
		return false, nil
	}
	if strings.Index(string(out), flag) == -1 {
		return false, nil
	}
	return true, nil
}

func haveDocker() (bool, error) {
	return haveSomething("docker version", "Version")
}

func haveNginx() (bool, error) {
	return haveSomething("netstat -anput | grep nginx", "nginx")
}

func run(c string) error {
	_, err := con.Run(c)
	return err
}

func output(c string) ([]byte, error) {
	return con.Run(c)
}
