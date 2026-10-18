package server

import (
	"os"

	"github.com/CALEXCO/goautomate/pkg"
	"github.com/goccy/go-yaml"
	"golang.org/x/crypto/ssh"
)

type Server struct {
	Ip        string `yaml:"ip"`
	User      string `yaml:"user"`
	Pass      string `yaml:"password"`
	SshConfig pkg.SshConfig
}

type TagetServers struct {
	Servers []Server `yaml:"target_servers"`
}

func NewServer(ip, user, password string) Server {
	return Server{
		Ip:        ip,
		User:      user,
		Pass:      password,
		SshConfig: pkg.NewSshConfig(user, password, ip),
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

func (s Server) ConnectToServer() (*ssh.Client, error) {
	client, err := ssh.Dial("tcp", s.SshConfig.TargetHost, &s.SshConfig.Config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func (s Server) ConnectAndExecute(command string) (string, error) {
	sc, err := s.ConnectToServer()
	if err != nil {
		panic(err)
	}
	output, err := s.SshConfig.ExecuteCommand(sc, command)
	if err != nil {
		return "", err
	}
	return output, nil
}
