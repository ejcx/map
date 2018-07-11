package ssh

import (
	"golang.org/x/crypto/ssh"
)

const (
	identifier = "ssh"
)

type SshDoer struct {
	Username string
	Password string
}

func (s *SshDoer) Identifier() string {
	return identifier
}

func (s *SshDoer) Do(addr string) (bool, []string) {
	username := "root"
	if s.Username != "" {
		username = s.Username
	}
	config := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(s.Password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	}
	// Dial your ssh server.
	conn, err := ssh.Dial("tcp", addr, config)
	if err != nil {
		return false, nil
	}
	defer conn.Close()
	return true, nil
}
