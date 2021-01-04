package service

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/inspursoft/board/src/adminserver/common"
	"github.com/inspursoft/board/src/adminserver/dao"
	"github.com/inspursoft/board/src/adminserver/models"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/astaxie/beego/logs"
)

func StartBoard(host *models.Account, buf *bytes.Buffer) error {
	cmdList := []string{}
	boardComposeFile, devopsOpt, err := GetFileFromDevopsOpt()
	if err != nil {
		return err
	}
	var cmdGitlabHelper string
	cmdPrepare := fmt.Sprintf("%s", models.PrepareFile)
	cmdComposeDown := fmt.Sprintf("docker-compose -f %s down", boardComposeFile)
	cmdComposeUp := fmt.Sprintf("docker-compose -f %s up -d", boardComposeFile)

	if devopsOpt == "legacy" {
		logs.Info("starting Board in legacy mode...")
		cmdList = []string{cmdPrepare, cmdComposeDown, cmdComposeUp}
	} else {
		logs.Info("starting Gitlab-helper...")
		tag, err := common.ReadCfgItem("gitlab_helper_version", "/go/cfgfile/board.cfg")
		if err != nil {
			return err
		}
		cmdGitlabHelper = fmt.Sprintf("docker run --rm -v %s/board.cfg:/app/instance/board.cfg gitlab-helper:%s", models.MakePath, tag)
		cmdList = []string{cmdGitlabHelper, cmdPrepare, cmdComposeDown, cmdComposeUp}
	}
	go func(buf *bytes.Buffer) {
		shell, err := SSHtoHost(host, buf)
		if err != nil {
			logs.Error(err)
			return
		}
		for _, cmd := range cmdList {
			logs.Info("running cmd: %s", cmd)
			err = shell.ExecuteCommand(cmd)
			if err != nil {
				logs.Error(err)
				return
			}
			logs.Debug(buf.String())
		}
	}(buf)
	RemoveUUIDTokenCache()

	return nil
}

func CheckSysStatus() (models.InitStatus, string, error) {
	var err error
	var cfgCheck bool
	log := logBuffer.String()
	if err = CheckBoard(); err != nil {
		logs.Info("Board is down: %+v", err)
		if cfgCheck, err = CheckCfgModified(); err != nil {
			return 0, "", err
		}
		if !cfgCheck {
			return models.InitStatusFirst, log, nil
		} else {
			return models.InitStatusSecond, log, nil
		}
	}
	return models.InitStatusThird, log, nil
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
	err = utils.RequestHandle(http.MethodGet, "http://apiserver:8088/api/v1/systeminfo", nil, nil, nil)
	if err == utils.ErrNotAcceptable {
		return nil
	}
	return err
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

func GetFileFromDevopsOpt() (boardComposeFile, devopsOpt string, err error) {
	devopsOpt, err = common.ReadCfgItem("devops_opt", "/go/cfgfile/board.cfg")
	if err != nil {
		return
	}
	if devopsOpt == "legacy" {
		boardComposeFile = models.BoardcomposeLegacy
	} else {
		boardComposeFile = models.Boardcompose
	}
	return
}
