package service

import (
	"bytes"
	"fmt"
	"git/inspursoft/board/src/adminserver/common"
	"git/inspursoft/board/src/adminserver/models"
	"git/inspursoft/board/src/adminserver/tools/secureShell"
	"os"
	"os/exec"
	"path"

	"github.com/astaxie/beego/logs"
)

var logBuffer bytes.Buffer

//Start Board without loading cfg.
func Start(host *models.Account) error {
	var buf bytes.Buffer
	shell, err := SSHtoHost(host, &buf)
	if err != nil {
		return err
	}
	boardComposeFile, _, err := GetFileFromDevopsOpt()
	if err != nil {
		return err
	}
	cmdComposeUp := fmt.Sprintf("docker-compose -f %s up -d", boardComposeFile)
	err = shell.ExecuteCommand(cmdComposeUp)
	if err != nil {
		return err
	}
	logs.Debug(buf.String())
	RemoveUUIDTokenCache()

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
	_, err = Execute(fmt.Sprintf("cp %s %s.tmp", cfgPath, cfgPath))
	if err != nil {
		return err
	}

	if err = StartBoard(host, &logBuffer); err != nil {
		return err
	}

	if err = os.Remove(cfgPath + ".tmp"); err != nil {
		return err
	}

	return nil
}

//Shutdown Board.
func Shutdown(host *models.Account, uninstall bool) error {
	var buf bytes.Buffer
	shell, err := SSHtoHost(host, &buf)
	if err != nil {
		return err
	}
	boardComposeFile, _, err := GetFileFromDevopsOpt()
	if err != nil {
		return err
	}
	cmdCompose := fmt.Sprintf("docker-compose -f %s down", boardComposeFile)
	err = shell.ExecuteCommand(cmdCompose)
	if err != nil {
		return err
	}
	logs.Debug(buf.String())
	RemoveUUIDTokenCache()

	if uninstall {
		//check existing files under repo.
		cmdCheck := `ls /data/board/ -lR | grep "^-" | wc -l | xargs echo -n`
		output, err := shell.Output(cmdCheck)
		if err != nil {
			return err
		}
		if output == "0" {
			return common.ErrNoData
		}
		cmdRm := fmt.Sprintf("rm -rf /data/board/* %s/board.cfg* && cp %s/adminserver/board.cfg %s/.", models.MakePath, models.MakePath, models.MakePath)
		err = shell.ExecuteCommand(cmdRm)
		if err != nil {
			return err
		}
	}
	return nil
}

func SSHtoHost(host *models.Account, buf *bytes.Buffer) (*secureShell.SecureShell, error) {
	HostIP, err := Execute("ip route | awk 'NR==1 {print $3}'|xargs echo -n")
	if err != nil {
		return nil, err
	}
	shell, err := secureShell.NewSecureShell(buf, HostIP, host.Username, host.Password, host.Port)
	if err != nil {
		return nil, err
	}
	return shell, nil
}

//Execute command in container.
func Execute(command string) (string, error) {
	cmd := exec.Command("sh", "-c", command)
	bytes, err := cmd.Output()
	if err != nil {
		return "", err
	}
	return string(bytes), nil
}
