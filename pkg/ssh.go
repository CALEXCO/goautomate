package pkg

import (
	"golang.org/x/crypto/ssh"
)

func ConnectToServer(user, password, host string) (*ssh.Client, error) {
	config := &ssh.ClientConfig{
		User: user,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}

	client, err := ssh.Dial("tcp", host+":22", config)
	if err != nil {
		return nil, err
	}
	return client, nil
}

func ExecuteCommand(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	return string(output), err
}

func ConnectAndExecute(user, pass, host, command string) (string, error) {
	sc, err := ConnectToServer(user, pass, host)
	if err != nil {
		panic(err)
	}
	output, err := ExecuteCommand(sc, command)
	if err != nil {
		return "", err
	}
	return output, nil
}
