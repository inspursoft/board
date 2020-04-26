package service

import (
	"bufio"
	"errors"
	"fmt"
	"git/inspursoft/board/src/adminserver/dao"
	"git/inspursoft/board/src/adminserver/models"
	"os"
	"os/exec"
	"path"
	"strings"
	"time"

	"github.com/astaxie/beego/logs"
	_ "github.com/go-sql-driver/mysql"
)

func InitDB(db *models.DBconf) error {
	if _, err := os.Stat(models.DBconfigdir); os.IsNotExist(err) {
		os.MkdirAll(models.DBconfigdir, os.ModePerm)
	}
	envFile := path.Join(models.DBconfigdir, "env")
	cnfFile := path.Join(models.DBconfigdir, "my.cnf")
	if _, err := os.Stat(envFile); err == nil {
		os.Remove(envFile)
	}
	if _, err := os.Stat(cnfFile); err == nil {
		os.Remove(cnfFile)
	}

	env, err := os.Create(envFile)
	defer env.Close()
	if err != nil {
		return err
	}
	env.WriteString(fmt.Sprintf("DB_PASSWORD=%s\n", db.Password))

	cnf, err := os.Create(cnfFile)
	defer cnf.Close()
	if err != nil {
		return err
	}
	cnf.WriteString("[mysqld]\n")
	cnf.WriteString(fmt.Sprintf("max_connections=%d\n", db.MaxConnections))

	return nil
}

func StartDB(host *models.Account) error {
	shell, err := SSHtoHost(host)
	if err != nil {
		return err
	}

	cmdDB := fmt.Sprintf("docker-compose -f %s up -d", models.DBcompose)
	err = shell.ExecuteCommand(cmdDB)
	if err != nil {
		return err
	}
	time.Sleep(time.Duration(10) * time.Second)

	if err = dao.CheckDB(); err != nil {
		logs.Info(err)
		err = dao.RegisterDB()
		if err != nil {
			return err
		}
	}

	return nil
}

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
	time.Sleep(time.Duration(3) * time.Second)
	/*
		// TODO:
		delete account in sqlite
	*/

	err = shell.ExecuteCommand(cmdComposeDown)
	if err != nil {
		return err
	}
	time.Sleep(time.Duration(10) * time.Second)

	err = shell.ExecuteCommand(cmdComposeUp)
	if err != nil {
		return err
	}

	err = dao.RemoveUUIDToken()
	if err != nil {
		return err
	}
	os.Remove("/go/secrets/initialAdminPassword")

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
		return errors.New("tokenserver is down")
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
	cmd := exec.Command("sh", "-c", "ping -q -c1 tokenserver > /dev/null 2>&1 && echo $?")
	bytes, _ := cmd.Output()
	result := strings.Replace(string(bytes), "\n", "", 1)
	return (result == "0")
}
