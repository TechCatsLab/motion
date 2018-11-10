package deploy

import (
	"fmt"

	"github.com/TechCatsLab/motion/nginxconfig"
)

type server struct {
	Path                  string
	ExportHost            string
	Host                  string
	Port                  int
	Image                 string
	Run_args              map[string]string
	Nginx_location_config map[string][]string

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
		if used, err := portIsUsed(port); err != nil {
			return err
		} else {
			if used {
				port += 1
			} else {
				break
			}
		}
	}
	out, err := runContainer(s.Image, s.Port, port, s.Run_args)
	if err != nil {
		return err
	}
	s.containerID = string(out)
	s.exportPort = port
	return nil
}

func (s server) Config() {
	sv := nginxconfig.AddServer(s.ExportHost, s.exportPort)
	for k, v := range s.Nginx_location_config {
		l := sv.AddLocation(k)
		l.AddSetSlice(v)
	}
}

func Generate() string {
	for _, e := range EnvConfig.Servers {
		e.Config()
	}
	return nginxconfig.GenerateConfig()
}
