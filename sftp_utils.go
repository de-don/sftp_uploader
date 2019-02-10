package main

import (
	"fmt"
	"github.com/pkg/sftp"
	"golang.org/x/crypto/ssh"
	"net"
)

func createSshSession(host, port, username, password string) (*ssh.Client, error) {
	// user data
	sshConfig := &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: func(hostname string, remote net.Addr, key ssh.PublicKey) error {
			return nil
		},
	}

	// create ssh connection
	return ssh.Dial("tcp", fmt.Sprintf("%s:%s", host, port), sshConfig)
}

func saveImageOnServer(sftpConn *sftp.Client, uploadDirectory string, screenName string, data []byte) error {
	// create new file
	f, err := sftpConn.Create(uploadDirectory + "/" + screenName)
	if err != nil {
		return err
	}
	// write in new fle
	if _, err := f.Write(data); err != nil {
		return err
	}
	return nil
}
