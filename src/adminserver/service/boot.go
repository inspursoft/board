package service

import (
	"bufio"
	"fmt"
	"git/inspursoft/board/src/adminserver/common"
	"git/inspursoft/board/src/adminserver/dao"
	"git/inspursoft/board/src/adminserver/models"
	"git/inspursoft/board/src/common/utils"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"

	"github.com/astaxie/beego/logs"
)

func StartBoard(host *models.Account, logDetail *[]string) error {
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
		if err = CheckGitlab(); err == nil {
			logs.Info("Gitlab is up")
			cmdGitlabHelper += " python action/perform.py -r true"
		}
		cmdList = []string{cmdGitlabHelper, cmdPrepare, cmdComposeDown, cmdComposeUp}
	}
	go func(logDetail *[]string) {
		shell, output, err := SSHtoHost(host)
		if err != nil {
			logs.Error(err)
			return
		}
		for _, cmd := range cmdList {
			output.Reset()
			logs.Info("running cmd: %s", cmd)
			err = shell.ExecuteCommand(cmd)
			if err != nil {
				logs.Error(err)
				return
			}
			logs.Debug(output.String())
			reader := bufio.NewReader(strings.NewReader(output.String()))
			for {
				line, err := reader.ReadString('\n')
				if err != nil {
					if err == io.EOF {
						break
					}
					logs.Error(err)
					return
				}
				line = strings.TrimSpace(line)
				*logDetail = append(*logDetail, line)
			}
		}
	}(logDetail)
	RemoveUUIDTokenCache()

	return nil
}

func CheckSysStatus() (models.InitStatus, error) {
	var err error
	var cfgCheck bool
	if err = CheckBoard(); err != nil {
		logs.Info("Board is down: %+v", err)
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

func CheckGitlab() error {
	ip, _ := common.ReadCfgItem("gitlab_host_ip", "/go/cfgfile/board.cfg")
	port, _ := common.ReadCfgItem("gitlab_host_port", "/go/cfgfile/board.cfg")
	url := fmt.Sprintf("http://%s:%s", ip, port)
	return utils.RequestHandle(http.MethodGet, url, nil, nil, nil)
}
