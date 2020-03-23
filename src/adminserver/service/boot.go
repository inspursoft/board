package service

import (
	"git/inspursoft/board/src/adminserver/models"
	"github.com/astaxie/beego/orm"
	"path"
	"os"
	"fmt"
	"time"
	"bytes"
	"os/exec"
	"strings"
	"git/inspursoft/board/src/adminserver/utils"
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
	cnf.WriteString(fmt.Sprintf("max_connections=%s\n", db.MaxConnections))

	o := orm.NewOrm()
	status := models.InitStatusInfo{Id: 1}
	err = o.Read(&status)
	if err == orm.ErrNoRows {
    	fmt.Println("not found")
	} else if err == orm.ErrMissPK {
		fmt.Println("pk missing")
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
	var secureShell *utils.SecureShell
	var err error

	cmd := exec.Command("sh", "-c", "ip route | awk 'NR==1 {print $3}'")
	bytes, _ := cmd.Output()
	HostIp := strings.Replace(string(bytes), "\n", "", 1)

	secureShell, err = utils.NewSecureShell(&output, HostIp, host.Username, host.Password)
	if err != nil {
		return err
	}
	return nil

}