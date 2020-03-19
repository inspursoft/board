package models

import (
	"reflect"
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
