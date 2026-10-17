package pkg

import "golang.org/x/crypto/ssh"

func connectToServer(user, password, host string) (*ssh.Client, error) {
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

func executeCommand(client *ssh.Client, command string) (string, error) {
	session, err := client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()

	output, err := session.CombinedOutput(command)
	return string(output), err
}

// TODO: Create finction for Connect and Execute to a list of servers (sequencial)
func ConnectAndExecuteCommand(client *ssh.Client, command string) {
	panic("TODO: Not implemented")
}
