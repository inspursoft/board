package service

import (
	"bufio"
	"fmt"
	"git/inspursoft/board/src/adminserver/common"
	"git/inspursoft/board/src/adminserver/dao"
	"git/inspursoft/board/src/adminserver/models"
	"os"
	"os/exec"
	"strings"

	"github.com/astaxie/beego/logs"
)

func StartBoard(host *models.Account) error {
	shell, err := SSHtoHost(host)
	if err != nil {
		return err
	}

	cmdPrepare := fmt.Sprintf("%s", models.PrepareFile)
	cmdComposeDown := fmt.Sprintf("docker-compose -f %s down", models.Boardcompose)
	cmdComposeUp := fmt.Sprintf("docker-compose -f %s up -d", models.Boardcompose)

	err = shell.ExecuteCommand(cmdPrepare)
	if err != nil {
		return err
	}
	//time.Sleep(time.Duration(3) * time.Second)
	/*
		// TODO:
		delete account in sqlite
	*/

	err = shell.ExecuteCommand(cmdComposeDown)
	if err != nil {
		return err
	}
	//time.Sleep(time.Duration(10) * time.Second)

	err = shell.ExecuteCommand(cmdComposeUp)
	if err != nil {
		return err
	}

	UUIDpath := "/go/secrets/initialAdminPassword"
	if _, err = os.Stat(UUIDpath); !os.IsNotExist(err) {
		if err = dao.RemoveUUIDToken(); err != nil {
			return err
		}
		os.Remove(UUIDpath)
	}

	return nil
}

func CheckSysStatus() (models.InitStatus, error) {
	var err error
	var cfgCheck bool
	if err = CheckBoard(); err != nil {
		logs.Info("Board is down")
		if cfgCheck, err = CheckCfgModified(); err != nil {
			return 0, err
		}
		if !cfgCheck {
			return models.InitStatusFirst, nil
		} else {
			return models.InitStatusSecond, nil
		}
	}
	return models.InitStatusThird, nil
}

func CheckBoard() error {
	var err error

	if err = dao.CheckDB(); err != nil {
		if err = dao.RegisterDB(); err != nil {
			return err
		}
	}
	if tokenserver := CheckTokenserver(); !tokenserver {
		return common.ErrTokenServer
	}
	return nil
}

func CheckCfgModified() (bool, error) {
	cfgPath := "/go/cfgfile/board.cfg"
	f, err := os.Open(cfgPath)
	if err != nil {
		return false, err
	}
	defer f.Close()
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		if strings.Contains(scanner.Text(), "hostname = reg.mydomain.com") {
			return false, nil
		}
	}
	return true, nil
}

func CheckTokenserver() bool {
	cmd := exec.Command("sh", "-c", "ping -q -c1 tokenserver > /dev/null 2>&1")
	cmd.Run()
	return (cmd.ProcessState.ExitCode() == 0)
}
