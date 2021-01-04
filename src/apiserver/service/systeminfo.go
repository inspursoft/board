package service

import (
	"fmt"
	"github.com/inspursoft/board/src/common/dao"
	"github.com/inspursoft/board/src/common/k8sassist"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"os/exec"

	"github.com/astaxie/beego/logs"
)

type SystemResourcesInfo struct {
	TotalNumberCpuCore int        `json:"total_number_cpu_core"`
	TotalMemorySize    int        `json:"total_memory_size"`
	TotalCpuUsage      float32    `json:"total_cpu_usage"`
	TotalMemoryUsage   int        `json:"total_memory_usage"`
	NodesResources     []NodeInfo `json:"nodes_resources"`
}

func GetSystemInfo() (*model.SystemInfo, error) {
	configs, err := dao.GetAllConfigs()
	if err != nil {
		return nil, err
	}
	var systemInfo model.SystemInfo
	for _, config := range configs {
		switch config.Name {
		case "MODE":
			systemInfo.Mode = config.Value
		case "BOARD_HOST_IP":
			systemInfo.BoardHost = config.Value
		case "AUTH_MODE":
			systemInfo.AuthMode = config.Value
		case "SET_ADMIN_PASSWORD":
			systemInfo.SetAdminPassword = config.Value
		case "INIT_PROJECT_REPO":
			systemInfo.InitProjectRepo = config.Value
		case "SYNC_K8S":
			systemInfo.SyncK8s = config.Value
		case "REDIRECTION_URL":
			systemInfo.RedirectionURL = config.Value
		case "BOARD_VERSION":
			systemInfo.Version = config.Value
		case "DNS_SUFFIX":
			systemInfo.DNSSuffix = config.Value
		case "KUBERNETES_VERSION":
			systemInfo.KubernetesVersion = config.Value
		case "DEVOPS_OPT":
			systemInfo.DevOpsOpt = config.Value
		}
	}

	//Get the hareware processor arch
	cmd := exec.Command("uname", "-m")
	out, err := cmd.Output()
	if err != nil {
		fmt.Errorf("Uname failed to get info %v", err)
		systemInfo.ProcessorType = ""
	} else {
		systemInfo.ProcessorType = string(out)
	}
	return &systemInfo, nil
}

func SetSystemInfo(name string, reconfigurable bool) error {
	config, err := dao.GetConfig(name)
	if err != nil {
		return err
	}
	if config.Name == "" || reconfigurable {
		value := utils.GetStringValue(name)
		if value == "" {
			return fmt.Errorf("Has not set config %s yet", name)
		}
		_, err := dao.AddOrUpdateConfig(model.Config{Name: name, Value: value, Comment: fmt.Sprintf("Set config %s.", name)})
		return err
	}
	utils.SetConfig(name, config.Value)
	return nil
}

func GetKubernetesInfo() (*model.KubernetesInfo, error) {
	// add the pv to k8s
	k8sclient := k8sassist.NewK8sAssistClient(&k8sassist.K8sAssistConfig{
		KubeConfigPath: kubeConfigPath(),
	})
	return k8sclient.AppV1().Discovery().ServerVersion()
}

func GetSystemResourcesInfo() (systemResourcesInfo SystemResourcesInfo, err error) {
	nodes, err := GetNodes()
	if err != nil {
		logs.Error("Failed to get Node information.")
		return
	}
	for _, node := range nodes {
		systemResourcesInfo.TotalNumberCpuCore += node.NumberCPUCore
		systemResourcesInfo.TotalMemorySize += node.MemorySize
		systemResourcesInfo.TotalCpuUsage += node.CPUUsage / 100
		systemResourcesInfo.TotalMemoryUsage += int(float32(node.MemorySize) * node.MemoryUsage / 100)
		systemResourcesInfo.NodesResources = append(systemResourcesInfo.NodesResources, node)
	}

	return
}
