package server

import (
	"os"

	"github.com/CALEXCO/goautomate/pkg"
	"github.com/goccy/go-yaml"
)

type Server struct {
	Ip        string `yaml:"ip"`
	User      string `yaml:"user"`
	Pass      string `yaml:"password"`
	sshConfig pkg.SshConfig
}

type TagetServers struct {
	Servers []Server `yaml:"target_servers"`
}

func NewServer(ip, user, password string) Server {
	return Server{
		Ip:        ip,
		User:      user,
		Pass:      password,
		sshConfig: pkg.NewSshConfig(user, password, ip),
	}
}

func NewTargetServers() TagetServers {
	return TagetServers{
		Servers: make([]Server, 0),
	}
}

func BuildServersFromYaml(hostFile string) TagetServers {
	var ts TagetServers
	s, err := os.ReadFile(hostFile)
	if err != nil {
		panic(err)
	}
	err = yaml.Unmarshal(s, &ts)
	if err != nil {
		panic(err)
	}
	return ts
}
