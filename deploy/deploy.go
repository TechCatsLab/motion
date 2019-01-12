package deploy

import (
	"fmt"
	"os"
	"os/exec"
)

func Start() error {
	var err error

	err = before()
	if err != nil {
		return err
	}

	err = startTrans()
	if err != nil {
		return err
	}

	return after()
}

func before() error {
	var err error
	for i, e := range EnvConfig.Servers {
		err = deal(e.Path, i)
		if err != nil {
			break
		}
	}
	return err
}

func startTrans() error {
	var err error
	for i := 0; i < len(EnvConfig.Servers); i++ {
		str := fmt.Sprintf("~/ready%d.tar.gz", i)
		err = trans(str, str)
		if err != nil {
			break
		}
	}
	return err
}

func after() error {
	var err error
	for i, e := range EnvConfig.Servers {
		err = unzip(i)
		if err != nil {
			break
		}

		err = pullImage(e.Image)
		if err != nil {
			break
		}

		err = e.RunServer()
		if err != nil {
			break
		}
	}
	err = writeConfig(Generate())
	if err != nil {
		return err
	}
	return useConfig()
}

func trans(from, to string) error {
	err := writeScript("./script.sh", from, to)
	if err != nil {
		return err
	}
	err = exec.Command("./script.sh").Run()
	if err != nil {
		return err
	}
	return os.Remove("./script.sh")
}

func deal(from string, index int) error {
	f, err := os.Stat(from)
	if err != nil {
		return err
	}

	if f.IsDir() {
		err = zip(from, index)
	} else {
		err = rename(from, index)
	}
	return err
}
