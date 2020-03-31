package service

import (
	"git/inspursoft/board/src/adminserver/models"
	board_utils "git/inspursoft/board/src/common/utils"
	"github.com/astaxie/beego/orm"
	"github.com/astaxie/beego/logs"
	"path"
	"os"
	"fmt"
	"time"
	"bytes"
	"os/exec"
	"strings"
	"io/ioutil"

	"git/inspursoft/board/src/adminserver/tools/secureShell"
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

	o := orm.NewOrm()
	status := models.InitStatusInfo{Id: 1}
	err = o.Read(&status)
	if err != nil {
		return err
	}
	if status.Status == models.InitStatusFirst {
		status.InstallTime = time.Now().Unix()
		status.Status = models.InitStatusSecond
		o.Update(&status, "InstallTime", "Status")
	}

	return nil
}

func StartDB(host *models.Account) error {
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

	cmdDB := fmt.Sprintf("docker-compose -f %s up -d", models.DBcompose)
	err = shell.ExecuteCommand(cmdDB)
	if err != nil {
		return err
	}

	time.Sleep(time.Duration(10)*time.Second)

	b, err := ioutil.ReadFile(path.Join(models.DBconfigdir, "/env"))
	if err != nil {
		return err
	}
	DBpassword := strings.TrimPrefix(string(b), "DB_PASSWORD=")
	DBpassword = strings.Replace(DBpassword, "\n", "", 1)

	board_utils.InitializeDefaultConfig()

	if err = CheckDB(); err != nil {
		logs.Info(err)
		orm.RegisterDriver("mysql", orm.DRMySQL)
		err = orm.RegisterDataBase("mysql-db2", "mysql", fmt.Sprintf("root:%s@tcp(%s:%d)/board?charset=utf8", DBpassword, "db", 3306))
		if err != nil {
			logs.Error("error occurred on registering DB: %+v", err)
			panic(err)
		}
		logs.Info("register DB success")
	}

	o := orm.NewOrm()
	status := models.InitStatusInfo{Id: 1}
	err = o.Read(&status)
	if err != nil {
		return err
	} 
	if status.Status == models.InitStatusSecond {
		status.InstallTime = time.Now().Unix()
		status.Status = models.InitStatusThird
		o.Update(&status, "InstallTime", "Status")
	}

	return nil
}

func StartBoard(host *models.Account) error {
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

	cmdPrepare := fmt.Sprintf("%s", models.PrepareFile)
	cmdCompose := fmt.Sprintf("docker-compose -f %s up -d", models.Boardcompose)

	err = shell.ExecuteCommand(cmdPrepare)
	if err != nil {
		return err
	}

	/*
	// TODO:
	o := orm.NewOrm()
	o.Using("default")
	account := models.Account{Id: 1}
	err = o.Read(&account)
	if err != orm.ErrNoRows {
		if _, err = o.Delete(&account); err != nil {
			return err
		}
	} 
	*/
	
	time.Sleep(time.Duration(5)*time.Second)
	
	err = shell.ExecuteCommand(cmdCompose)
	if err != nil {
		return err
	}

	return nil
}

func CheckDB() error {
	/*
	// checkDB
	cmd := exec.Command("sh", "-c", "ping -q -c1 db > /dev/null 2>&1 && echo $?")
	bytes, _ := cmd.Output()
	result := strings.Replace(string(bytes), "\n", "", 1)
	return (result == "0")
	*/
	o := orm.NewOrm()
	err := o.Using("mysql-db2")
	if err != nil {
		return err
	}
	_, err = o.Raw("SELECT 1").Exec()
	return err
}


func RegisterDBWhenBooting() error {
	/*
	// checkDB
	cmd := exec.Command("sh", "-c", "ping -q -c1 db > /dev/null 2>&1 && echo $?")
	bytes, _ := cmd.Output()
	result := strings.Replace(string(bytes), "\n", "", 1)
	*/
	
	if err := CheckDB(); err == nil {
		b, err := ioutil.ReadFile(path.Join(models.DBconfigdir, "/env"))
		if err != nil {
			return err
		}
		DBpassword := strings.TrimPrefix(string(b), "DB_PASSWORD=")
		DBpassword = strings.Replace(DBpassword, "\n", "", 1)
		orm.RegisterDriver("mysql", orm.DRMySQL)
		err = orm.RegisterDataBase("mysql-db2", "mysql", fmt.Sprintf("root:%s@tcp(%s:%d)/board?charset=utf8", DBpassword, "db", 3306))
		if err != nil {
			return err
		}
		return nil
	} else {
		return nil
	}
}