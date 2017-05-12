package remote

import (
	"bytes"
	"golang.org/x/crypto/ssh"
	"chatOperations/operations"
)

type Shell struct {
	config *ssh.ClientConfig
}

func NewShell(user string, keyBuffer []byte) (Shell, error) {
	key, err := ssh.ParsePrivateKey(keyBuffer)

	return Shell{
		config: &ssh.ClientConfig{
			User: user,
			Auth: []ssh.AuthMethod{
				ssh.PublicKeys(key),
			},
		},
	}, err
}

func (shell *Shell) Execute(o operations.Request) (string, error) {
	client, err := ssh.Dial("tcp", o.Server.Address + ":22", shell.config)
	if err != nil {
		return "failed to dial", err
	}

	defer client.Close()

	session, err := client.NewSession()
	if err != nil {
		return "failed to create session", err
	}

	defer session.Close()

	var b bytes.Buffer
	session.Stdout = &b

	err = session.Run(commandString(o))

	return b.String(), err
}

func commandString(o operations.Request) (string) {
	return "cd /home/sites/wise_config && " + o.Action.Command
}
