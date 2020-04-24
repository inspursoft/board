package models

import (
	"reflect"
	"time"
	"github.com/astaxie/beego/orm"
	"os"
)

//Boardinfo contains information output by docker ps and docker stats commands.
type Boardinfo struct {
	ID        string `json:"id"`
	Image     string `json:"image"`
	CreatedAt string `json:"created_at"`
	Status    string `json:"status"`
	Ports     string `json:"ports"`
	Name      string `json:"name"`
	CPUPerc   string `json:"cpu_perc"`
	MemUsage  string `json:"mem_usage"`
	NetIO     string `json:"net_io"`
	BlockIO   string `json:"block_io"`
	MemPerc   string `json:"mem_perc"`
	PIDs      string `json:"pids"`
}

//GetBoardinfo transfers string array into struct.
func GetBoardinfo(container []string) Boardinfo {
	var boardinfo Boardinfo
	value := reflect.ValueOf(&boardinfo).Elem()
	for i := 0; i < value.NumField(); i++ {
		value.Field(i).SetString(container[i])
	}
	return boardinfo
}

type InitStatus int

const (
	InitStatusTrue		InitStatus = 0
	InitStatusFirst		InitStatus = 1	
	InitStatusSecond	InitStatus = 2
	InitStatusThird		InitStatus = 3
	InitStatusFalse		InitStatus = 4
)

//InitStatus saves the status indicating if the adminserver is first-time installed. 
type InitStatusInfo struct {
	Id          int        	`json:"id"`
	InstallTime	int64		`json:"install_time"`
	Status		InitStatus	`json:"status"`
}

type Token struct {
	Id		int		`json:"id"`
	Token	string	`json:"token"`
	Time	int64	`json:"time"`
}


func InitInstallationStatus() error {
	o := orm.NewOrm()
	status := &InitStatusInfo{Id: 1}
	err := o.Read(status,"Id")
	if err == orm.ErrNoRows {
		initStatus := InitStatusInfo{InstallTime: time.Now().Unix(), Status: InitStatusTrue}
		_, err := o.Insert(&initStatus)
		if err != nil {
			return err
		}	
	} 
	return nil
}

var ImagePrefix string = os.Getenv("IMAGE_PREFIX")
var ContainerPrefix string = os.Getenv("CONTAINER_PREFIX")