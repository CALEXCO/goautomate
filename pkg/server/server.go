package server

import (
	"fmt"
	"os"

	"github.com/CALEXCO/goautomate/pkg"
	"github.com/CALEXCO/goautomate/pkg/errors"
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

func BuildServersFromYaml(hostFile string) (TagetServers, error) {
	var ts TagetServers
	s, err := os.ReadFile(hostFile)
	if err != nil {
		return TagetServers{}, &errors.GoAutomateErrors{Code: errors.ReadingFileFailure, Msg: "ERROR: Could not read file"}
	}
	err = yaml.Unmarshal(s, &ts)
	if err != nil {
		fmt.Fprintln(os.Stderr, "ERROR: Could not unmarshal from yaml", err)
		return TagetServers{}, &errors.GoAutomateErrors{Code: errors.UnmarshalingFailure, Msg: "ERROR: Could not unmarshal"}
	}
	for i := range ts.Servers {
		s := &ts.Servers[i]
		s.SshConfig = pkg.NewSshConfig(s.User, s.Pass, s.Ip)
	}
	return ts, nil
}

func (s *Server) ConnectToServer() (*ssh.Client, error) {
	client, err := ssh.Dial("tcp", s.SshConfig.TargetHost, &s.SshConfig.Config)
	if err != nil {
		return nil, &errors.GoAutomateErrors{Code: errors.ConnectionFailure, Msg: "Error: Connection failure"}
	}
	return client, nil
}

func (s *Server) ConnectAndExecute(command string) (string, error) {
	sc, err := s.ConnectToServer()
	if err != nil {
		return "", &errors.GoAutomateErrors{Code: errors.ConnectionFailure, Msg: "Error: Connection failure"}
	}
	output, err := s.SshConfig.ExecuteCommand(sc, command)
	if err != nil {
		return "", &errors.GoAutomateErrors{Code: errors.CommandTimeout, Msg: "Error: Command did not finished correctly"}
	}
	return output, nil
}
