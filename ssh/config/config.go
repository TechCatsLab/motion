package config

import (
	"log"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
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
	viper.SetConfigName("ssh")
	viper.AddConfigPath("./config/")
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}

	SSHConf = &sshConf{
		User:     viper.GetString("user"),
		Address:  viper.GetString("host.address"),
		Port:     viper.GetString("host.port"),
		Password: viper.GetString("password"),
	}

	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		log.Println("Config file changed:", e.Name)
	})
}
