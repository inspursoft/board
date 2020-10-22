package secureShell

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/tmc/scp"
	"golang.org/x/crypto/ssh"
)

const maxSSHRetries = 10
const maxSSHDelay = 50 * time.Microsecond

type SecureShell struct {
	client *ssh.Client
	output io.Writer
}

func NewSecureShell(output io.Writer, host, username, password string, port ...int) (*SecureShell, error) {
	sshPort := 22
	if port != nil {
		sshPort = port[0]
	}
	// Retry few times if ssh connection fails
	for i := 0; i < maxSSHRetries; i++ {
		client, err := ssh.Dial("tcp", fmt.Sprintf("%s:%d", host, sshPort), &ssh.ClientConfig{
			User: username,
			Auth: []ssh.AuthMethod{
				ssh.Password(password),
			},
			HostKeyCallback: ssh.InsecureIgnoreHostKey(),
		})
		if err != nil {
			time.Sleep(maxSSHDelay)
			errMessage := fmt.Sprintf("Failed to dial host: %+v\n", err)
			output.Write([]byte(errMessage))
			logs.Info(errMessage)
			continue
		}
		s, err := client.NewSession()
		if err != nil {
			client.Close()
			time.Sleep(maxSSHDelay)
			continue
		}
		modes := ssh.TerminalModes{
			ssh.ECHO:          0,
			ssh.TTY_OP_ISPEED: 14400,
			ssh.TTY_OP_OSPEED: 14400,
		}
		// Request pseudo terminal
		if err := s.RequestPty("xterm", 40, 80, modes); err != nil {
			return nil, fmt.Errorf("failed to get pseudo-terminal: %v", err)
		}
		output.Write([]byte(fmt.Sprintf("---Connected to %s successfully.---\n", host)))
		return &SecureShell{client: client, output: output}, nil
	}
	output.Write([]byte(fmt.Sprintf("---Connected to %s. unsuccessfully---\n", host)))
	return nil, fmt.Errorf("retry times was exceeded 10")
}

func (s *SecureShell) ExecuteCommand(cmd string) error {
	session, err := s.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	w := io.MultiWriter(s.output)
	session.Stdout = w
	session.Stderr = w
	if err := session.Start(cmd); err != nil {
		return err
	}
	if err := session.Wait(); err != nil {
		return err
	}
	return nil
}

func (s *SecureShell) Output(cmd string) (string, error) {
	session, err := s.client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	b, err := session.Output(cmd)
	return string(b), err
}

func (s *SecureShell) SecureCopyData(fileName string, data []byte, destinationPath string) error {
	session, err := s.client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()
	var buf bytes.Buffer
	length, err := buf.Write(data)
	if err != nil {
		log.Printf("Failed to load contents: %+v\n", err)
		return err
	}
	return scp.Copy(int64(length), 0755, fileName, &buf, filepath.Join(destinationPath, fileName), session)
}

func (s *SecureShell) SecureCopy(filePath string, destinationPath string) error {
	return filepath.Walk(filePath, func(path string, info os.FileInfo, err error) error {
		session, err := s.client.NewSession()
		if err != nil {
			return err
		}
		defer session.Close()
		if err != nil {
			return err
		}
		if info.IsDir() {
			log.Printf("From path: %s to path: %s\n", path, destinationPath)
			return nil
		}
		log.Printf("From path: %s to path: %s\n", path, filepath.Join(destinationPath, info.Name()))
		return scp.CopyPath(path, filepath.Join(destinationPath, info.Name()), session)
	})
}

func (s *SecureShell) CheckDir(dir string) error {
	return s.ExecuteCommand(fmt.Sprintf("mkdir -p %s", dir))
}

func (s *SecureShell) RemoveDir(dir string) error {
	return s.ExecuteCommand(fmt.Sprintf("rm -rf %s", dir))
}

func (s *SecureShell) HostSCP(sourcePath, targetPath string, reversed bool) error {
	sshHost := beego.AppConfig.String("ssh-host::host")
	sshPort := beego.AppConfig.String("ssh-host::port")
	sshUsername := beego.AppConfig.String("ssh-host::username")
	scpCommand := fmt.Sprintf("scp -P %s %s@%s:%s %s", sshPort, sshUsername, sshHost, sourcePath, targetPath)
	if reversed {
		scpCommand = fmt.Sprintf("scp -P %s %s %s@%s:%s", sshPort, sourcePath, sshUsername, sshHost, sourcePath)
	}
	beego.Debug(fmt.Sprintf("Host SCP command is: %s", scpCommand))
	return s.ExecuteCommand(scpCommand)
}
