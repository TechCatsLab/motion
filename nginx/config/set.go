package config

import (
	"strings"
)

type set struct {
	cmd  string
	args []string
}

func NewSetWithSlice(cmd string, args ...string) *set {
	return &set{
		cmd:  cmd,
		args: args,
	}
}

func NewSet(cmd string) *set {
	list := strings.Split(cmd, " ")
	if len(list) < 1 {
		return nil
	}
	if len(list) == 1 {
		return &set{
			cmd:  list[0],
			args: nil,
		}
	}
	return &set{
		cmd:  list[0],
		args: list[1:],
	}
}

func (s set) String() string {
	res := s.cmd + " "
	for i := 0; i < len(s.args); i++ {
		res += s.args[i] + " "
	}
	return res[:len(res)-1]
}
