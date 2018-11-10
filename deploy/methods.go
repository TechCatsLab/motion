package deploy

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

	"github.com/TechCatsLab/motion/conn"
)

func zip(dir string, index int) error {
	return exec.Command("tar", "-zcvf", fmt.Sprintf("~/ready%d.tar.gz", index), dir).Run()
}

func rename(file string, index int) error {
	return exec.Command("mv", file, fmt.Sprintf("~/ready%d.tar.gz", index)).Run()
}

func writeScript(name, from, to, pass string) error {
	return ioutil.WriteFile(name, []byte(fmt.Sprintf(EnvConfig.Format, from, to, pass)), os.ModePerm)
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

func run(c string) error {
	_, err := conn.Run(c)
	return err
}

func output(c string) ([]byte, error) {
	return conn.Run(c)
}
