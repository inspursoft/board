package service

import (
	"log"
	"os"
	"os/exec"
	"path"
	"git/inspursoft/board/src/adminserver/models"
	"git/inspursoft/board/src/adminserver/tools/secureShell"
	"bytes"
	"strings"
	"fmt"

)

//Restart Board without loading cfg.
func Restart(host *models.Account) error {
	var output bytes.Buffer
	var shell *secureShell.SecureShell
	var err error

	cmd := exec.Command("sh", "-c", "ip route | awk 'NR==1 {print $3}'")
	bytes, _ := cmd.Output()
	HostIp := strings.Replace(string(bytes), "\n", "", 1)

	shell, err = secureShell.NewSecureShell(&output, HostIp, host.Username, host.Password)
	if err != nil {
		return err
	}

	cmdComposeDown := fmt.Sprintf("docker-compose -f %s down", models.Boardcompose)
	err = shell.ExecuteCommand(cmdComposeDown)
	if err != nil {
		return err
	}

	cmdComposeUp := fmt.Sprintf("docker-compose -f %s up -d", models.Boardcompose)
	err = shell.ExecuteCommand(cmdComposeUp)
	if err != nil {
		return err
	}

	return nil
}

//Applycfg restarts Board with applying of cfg.
func Applycfg(host *models.Account) error {

	cfgPath := path.Join("/go", "/cfgfile/board.cfg")
	err := os.Rename(cfgPath, cfgPath+".bak1")
	if err != nil {
		if !os.IsNotExist(err) { // fine if the file does not exists
			return err
		}
	}
	err = os.Rename(cfgPath+".tmp", cfgPath)
	if err != nil {
		if !os.IsNotExist(err) { // fine if the file does not exists
			return err
		}
	}

	err = Execute(fmt.Sprintf("cp %s %s.tmp", cfgPath, cfgPath))
	if err != nil {
		return err
	}

	//if err = Shutdown(host); err != nil {
	//	return err
	//}
	if err = StartBoard(host); err != nil {
		return err
	}
	
	if err = os.Remove(cfgPath+".tmp"); err != nil {
		return err
	} 

	return nil
}

//Shutdown Board.
func Shutdown(host *models.Account) error {
	var output bytes.Buffer
	var shell *secureShell.SecureShell
	var err error

	cmd := exec.Command("sh", "-c", "ip route | awk 'NR==1 {print $3}'")
	bytes, _ := cmd.Output()
	HostIp := strings.Replace(string(bytes), "\n", "", 1)

	shell, err = secureShell.NewSecureShell(&output, HostIp, host.Username, host.Password)
	if err != nil {
		return err
	}

	cmdCompose := fmt.Sprintf("docker-compose -f %s down", models.Boardcompose)
	err = shell.ExecuteCommand(cmdCompose)
	if err != nil {
		return err
	}

	return nil
}

//Execute command.
func Execute(command string) error {
	cmd := exec.Command("sh", "-c", command)
	bytes, err := cmd.Output()
	log.Println(string(bytes))
	return err
}
