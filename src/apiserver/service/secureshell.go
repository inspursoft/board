package service

import (
	"bytes"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/tmc/scp"

	"github.com/astaxie/beego/logs"
	"golang.org/x/crypto/ssh"
)

type secureShell struct {
	client    *ssh.Client
	stdOutput bytes.Buffer
}

func NewSecureShell(host string, port int, username string, password string) (*secureShell, error) {
	client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, port), &ssh.ClientConfig{
		User: username,
		Auth: []ssh.AuthMethod{
			ssh.Password(password),
		},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(),
	})
	if err != nil {
		logs.Error("Failed to dial host: %+v", err)
		return nil, err
	}
	return &secureShell{client: client}, nil
}

func (s *secureShell) execute(callback func(stdOutput *bytes.Buffer, args ...string) error, commands ...string) (err error) {
	results := make(chan string, 10)
	timeout := time.After(time.Second * 10)
	go func() {
		var stdOutput bytes.Buffer
		err = callback(&stdOutput, commands...)
		if err != nil {
			logs.Error("Failed to execute via SSH: %+v", err)
		}
		results <- stdOutput.String()
	}()
	select {
	case res := <-results:
		logs.Info("Finished to run command via SSH: %+s", res)
	case <-timeout:
		logs.Debug("Timeout while executing command.")
	}
	return
}

func (s *secureShell) ExecuteCommand(cmd string) error {
	session, err := s.client.NewSession()
	if err != nil {
		logs.Error("Failed to create session: %+v", err)
		return err
	}
	defer session.Close()
	return s.execute(func(stdOutput *bytes.Buffer, args ...string) error {
		session.Stdout = stdOutput
		logs.Debug("Execute command: %s", cmd)
		return session.Run(args[0])
	}, cmd)
}

func (s *secureShell) SecureCopy(filePath string, destinationPath string) error {
	return s.execute(func(stdOuput *bytes.Buffer, args ...string) error {
		return filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			session, err := s.client.NewSession()
			if err != nil {
				logs.Error("Failed to create session: %+v", err)
				return err
			}
			defer session.Close()
			session.Stdout = stdOuput
			if info.IsDir() {
				logs.Debug("From path: %s to path: %s", path, args[1])
				// scp.CopyPath(path, args[1], session)
				return nil
			}
			logs.Debug("From path: %s to path: %s", path, filepath.Join(args[1], info.Name()))
			return scp.CopyPath(path, filepath.Join(args[1], info.Name()), session)
		})
	}, filePath, destinationPath)
}
