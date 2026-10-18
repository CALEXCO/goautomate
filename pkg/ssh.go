package pkg

import (
	"context"
	"time"

	"golang.org/x/crypto/ssh"
)

type SshConfig struct {
	config ssh.ClientConfig
	ctx    context.Context
}

//TODO: Wrap the ssh cinfiguration and everything into a struct to have better organization

func NewSshConfig(user, password, host string) SshConfig {
	return SshConfig{
		config: ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
			Timeout:         200 * time.Millisecond,
		},
		ctx: context.Background(),
	}
}

func ConnectToServer(user, password, host string) (*ssh.Client, error) {
	config := NewSshConfig(user, password, host)

	client, err := ssh.Dial("tcp", host+":22", &config.config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func ExecuteCommand(ctx context.Context, client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}

	defer session.Close()

	done := make(chan struct{})
	defer close(done)

	go func() {
		select {
		case <-ctx.Done():
			client.Close()
		case <-done:
		}
	}()

	output, err := session.CombinedOutput(command)
	if ctx.Err() != nil {
		return "", ctx.Err()
	}

	return string(output), err
}

func ConnectAndExecute(ctx context.Context, user, pass, host, command string) (string, error) {
	sc, err := ConnectToServer(user, pass, host)
	if err != nil {
		panic(err)
	}
	output, err := ExecuteCommand(ctx, sc, command)
	if err != nil {
		return "", err
	}
	return output, nil
}
