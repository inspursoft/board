package app

import (
	"flag"
	"git/inspursoft/board/src/collector/service/collect"
	"git/inspursoft/board/src/collector/util"
	"os"
)

type ServerRunOptions struct {
	ServerDbType     string
	ServerDbIp       string
	ServerDbPort     string
	ServerDbPassword string
	ServerKubeIp     string
	ServerKubePort   string
}
type env struct {
	envDbType     string
	envDbIp       string
	envDbPort     string
	envDbPassword string
	envKubeIp     string
	envKubePort   string
}

func getOsEnv() (NewEnv env) {
	NewEnv = env{
		envDbType:     os.Getenv("DB_TYPE"),
		envDbIp:       os.Getenv("DB_IP"),
		envDbPort:     os.Getenv("DB_PORT"),
		envDbPassword: os.Getenv("DB_PASSWORD"),
		envKubeIp:     os.Getenv("KUBEMASTER_IP"),
		envKubePort:   os.Getenv("KUBEMASTER_PORT"),
	}
	return
}

func preCmdFlag(vName string, vValue string, usage string) *string {
	s := flag.String(vName, vValue, usage)
	flag.Parse()
	return s
}

func getRunFlag() map[string]*string {
	var runFlag map[string]*string
	runFlag = make(map[string]*string)
	runFlag["DbType"] = preCmdFlag("dbtype", "non", "input the Database type name")
	runFlag["DbIp"] = preCmdFlag("dbip", "non", "input the Database ip address")
	runFlag["DbPort"] = preCmdFlag("dbport", "non", "input the database port")
	runFlag["KubeIp"] = preCmdFlag("kubeip", "non", "input the KubeAPIserver ip address")
	runFlag["KubePort"] = preCmdFlag("kubeport", "non", "input the kubeAPIserver port")
	runFlag["DbPass"] = preCmdFlag("dbpass", "non", "input the sql password")
	return runFlag
}

var RunFlag ServerRunOptions

func preFlag() {
	runFlag := getRunFlag()
	osEnv := getOsEnv()
	for k, v := range runFlag {
		flag.Parse()
		switch k {
		case "DbType":
			util.Logger.SetInfo(k, *v)
			if *v == "non" {
				RunFlag.ServerDbType = osEnv.envDbType
			}
		case "DbIp":
			util.Logger.SetInfo(k, *v)
			if *v == "non" {
				RunFlag.ServerDbIp = osEnv.envDbIp
			}
		case "DbPort":
			util.Logger.SetInfo(k, *v)
			if *v == "non" {
				RunFlag.ServerDbPort = osEnv.envDbPort
			}
		case "KubeIp":
			util.Logger.SetInfo(k, *v)
			if *v == "non" {
				RunFlag.ServerKubeIp = osEnv.envKubeIp
			}
		case "KubePort":
			util.Logger.SetInfo(k, *v)
			if *v == "non" {
				RunFlag.ServerKubePort = osEnv.envKubePort
			}
		case "DbPass":
			util.Logger.SetInfo(k, *v)
			if *v == "non" {
				RunFlag.ServerDbPassword = osEnv.envDbPassword
			}
		}
	}
}
func init() {
	preFlag()
	collect.SetInitVar(RunFlag.ServerKubeIp, RunFlag.ServerKubePort)
}
