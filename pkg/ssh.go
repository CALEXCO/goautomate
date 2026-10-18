package pkg

import (
	"context"
	"time"

	"golang.org/x/crypto/ssh"
)

type SshConfig struct {
	Config     ssh.ClientConfig
	Ctx        context.Context
	TargetHost string
}

//TODO: Wrap the ssh cinfiguration and everything into a struct to have better organization

func NewSshConfig(user, password, host string) SshConfig {
	return SshConfig{
		Config: ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         200 * time.Millisecond,
		},
		Ctx:        nil,
		TargetHost: host + ":22",
	}
}

func (sc SshConfig) ExecuteCommand(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}

	defer session.Close()

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-sc.Ctx.Done():
			client.Close()
		case <-done:
		}
	}()

	output, err := session.CombinedOutput(command)
	if sc.Ctx.Err() != nil {
		return "", sc.Ctx.Err()
	}

	return string(output), err
}
