package config

import (
	"golang.org/x/crypto/ssh"
)

type ConnConfig interface {
	CombinedHost() string
	CombinedHostAndUser() string
	ToClientConfig() *ssh.ClientConfig
}
