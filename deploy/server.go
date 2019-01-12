package deploy

import (
	"fmt"

	"github.com/TechCatsLab/motion/nginx/config"
)

type server struct {
	Path                string
	ExportHost          string
	Host                string
	Port                int
	Image               string
	RunArgs             map[string]string
	NginxLocationConfig map[string][]string

	containerID string
	exportPort  int
}

func (s server) String() string {
	return fmt.Sprintf("%s:%d", s.Host, s.Port)
}

func (s server) ExPort() int {
	return s.exportPort
}

func (s *server) RunServer() error {
	port := s.Port
	for {
		used, err := portIsUsed(port)
		if err != nil {
			return err
		}

		if !used {
			break
		}

		port += 1
	}

	out, err := runContainer(s.Image, s.Port, port, s.RunArgs)
	if err != nil {
		return err
	}
	s.containerID = string(out)
	s.exportPort = port
	return nil
}

func (s server) Config() {
	sv := config.AddServer(s.ExportHost, s.exportPort)
	for k, v := range s.NginxLocationConfig {
		l := sv.AddLocation(k)
		l.AddSetSlice(v)
	}
}

func Generate() string {
	for _, e := range EnvConfig.Servers {
		e.Config()
	}
	return config.GenerateConfig()
}
