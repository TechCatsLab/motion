package config

import (
	"fmt"
	"strings"
)

type NginxConfig struct {
	name         string
	middleConfig string
	sets         []*set
	subs         []*NginxConfig
	main         bool
}

var (
	defaultConfig       *NginxConfig
	defaultHttpConfig   *NginxConfig
	defaultEventsConfig *NginxConfig
)

func init() {
	defaultConfig = NewNginxConfig("default")
	defaultConfig.SetMain()
	defaultConfig.AddSetString("worker_processes 1")

	defaultEventsConfig = NewNginxConfig("events")
	defaultEventsConfig.AddSetString("worker_connections 1024")

	defaultHttpConfig = NewNginxConfig("http")
	defaultHttpConfig.AddSetString("include mime.types")
	defaultHttpConfig.AddSetString("default_type application/octet-stream")
	defaultHttpConfig.AddSetString("sendfile on")
	defaultHttpConfig.AddSetString("keepalive_timeout 65")
	defaultHttpConfig.AddSetString("include servers/*")

	defaultConfig.AddSubs(defaultEventsConfig, defaultHttpConfig)
}

func NewNginxConfig(name string) *NginxConfig {
	return &NginxConfig{
		name:         name,
		middleConfig: "",
		sets:         nil,
		subs:         nil,
		main:         false,
	}
}

func AddEvents(sets ...*set) {
	defaultEventsConfig.AddSets(sets...)
}

func AddServer(host string, port int) *NginxConfig {
	s := NewNginxConfig("server")
	s.AddSetString(fmt.Sprintf("listen %d", port))
	s.AddSetString("server_name " + host)
	defaultHttpConfig.AddSubs(s)
	return s
}

func (n *NginxConfig) AddLocation(path string, sets ...*set) *NginxConfig {
	if n.name != "server" {
		return nil
	}
	l := NewNginxConfig("location")
	l.SetMiddleConfig(path)
	l.AddSets(sets...)
	n.AddSubs(l)
	return l
}

func (n *NginxConfig) SetName(name string) {
	n.name = name
}

func (n *NginxConfig) SetMiddleConfig(mc string) {
	n.middleConfig = mc
}

func (n *NginxConfig) SetMain() {
	n.main = true
}

func (n *NginxConfig) SetSets(sets ...*set) {
	n.sets = sets
}

func (n *NginxConfig) AddSets(sets ...*set) {
	if n.sets == nil {
		n.sets = make([]*set, 0, len(sets))
	}
	for _, e := range sets {
		if !includes(n.sets, e) {
			n.sets = append(n.sets, e)
		}
	}
}

func (n *NginxConfig) AddSetString(s string) {
	set := NewSet(s)
	n.AddSets(set)
}

func (n *NginxConfig) AddSetSlice(s []string) {
	n.AddSetsString(s...)
}

func (n *NginxConfig) AddSetsString(s ...string) {
	for _, e := range s {
		n.AddSetString(e)
	}
}

func (n *NginxConfig) SetSubs(subs ...*NginxConfig) {
	n.subs = subs
}

func (n *NginxConfig) AddSubs(subs ...*NginxConfig) {
	if n.subs == nil {
		n.subs = make([]*NginxConfig, 0, len(subs))
	}
	for _, e := range subs {
		if e != n {
			n.subs = append(n.subs, e)
		}
	}
}

func (n NginxConfig) String() string {
	var str string
	if n.name == "location" {
		str = fmt.Sprintf("%s %s {\n", n.name, n.middleConfig)
	} else {
		str = fmt.Sprintf("%s {\n", n.name)
	}
	if n.main {
		str = "\n"
	}
	for _, e := range n.sets {
		str += "    " + e.String() + ";\n"
	}
	if n.subs != nil {
		for _, e := range n.subs {
			str += strings.Replace("    "+e.String(), "\n", "\n    ", -1) + "\n"
		}
	}
	if n.main {
		return strings.Replace(str, "\n    ", "\n", -1)
	} else {
		return str + "}"
	}
}

func includes(slice []*set, set *set) bool {
	for _, e := range slice {
		if e == set {
			return true
		}
	}
	return false
}

func GenerateConfig() string {
	return defaultConfig.String()
}
