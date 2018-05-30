package service

import (
	"fmt"
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"os"
	"testing"
	"time"

	"github.com/astaxie/beego/logs"
	"k8s.io/client-go/kubernetes"
	modelK8s "k8s.io/client-go/pkg/api/v1"
)

var boardHostIP = utils.GetConfig("BOARD_HOST_IP")

func connectToK8S() (*kubernetes.Clientset, error) {
	cli, err := K8sCliFactory("", kubeMasterURL(), "v1")
	cliSet, err := kubernetes.NewForConfig(cli)
	if err != nil {
		logs.Error("Failed to connect to K8s.")
		return nil, err
	}
	return cliSet, nil
}

func createService(cliSet *kubernetes.Clientset, serviceConfig Service, deploymentConfig Deployment) (*model.ServiceStatus, error) {
	serviceInfo, err := cliSet.CoreV1().Services(serviceConfig.Namespace).Create(&serviceConfig.Service)
	if err != nil {
		logs.Error("Created service failed.\n")
		return nil, err
	}
	logs.Debug("Created service, serviceInfo:%s.\n", serviceInfo)

	deploymentInfo, err := cliSet.Deployments(deploymentConfig.Namespace).Create(&deploymentConfig.Deployment)
	if err != nil {
		logs.Error("Created deployment failed.\n")
		return nil, err
	}
	logs.Debug("Created deployment, deploymentInfo:%s.\n", deploymentInfo)

	serviceStatus, err := CreateServiceConfig(model.ServiceStatus{
		Name:        serviceConfig.Name,
		ProjectName: serviceConfig.Namespace,
		Status:      defaultStatus,
	})
	if err != nil {
		logs.Error("Created Service info in DB failed.\n")
		return nil, err
	}
	logs.Debug("Service info in DB:%+v\n", serviceStatus)

	return serviceStatus, nil
}

func deleteService(cliSet *kubernetes.Clientset, serviceConfig Service, deploymentConfig Deployment, serviceStatus *model.ServiceStatus) error {
	err := cliSet.CoreV1().Services(serviceConfig.Namespace).Delete(serviceConfig.Name, nil)
	if err != nil {
		return err
	}
	logs.Debug("Delete service %s.\n", serviceConfig.Name)
	replicas = 0
	cliSetDeployment := cliSet.Deployments(deploymentConfig.Namespace)
	_, err = cliSetDeployment.Update(&deploymentConfig.Deployment)
	if err != nil {
		return err
	}
	time.Sleep(2)
	err = cliSetDeployment.Delete(deploymentConfig.Name, nil)
	if err != nil {
		return err
	}

	var opt modelK8s.ListOptions
	opt.LabelSelector = fmt.Sprintf("app=%s", deploymentConfig.Name)
	cliSetRS := cliSet.ReplicaSets(deploymentConfig.Namespace)
	RSList, err := cliSetRS.List(opt)
	if err != nil {
		logs.Error("Failed to get RS list")
		return err
	}

	for _, rs := range RSList.Items {
		err = cliSetRS.Delete(rs.Name, nil)
		if err != nil {
			logs.Error("Failed to delete RS:%s", rs.Name)
			return err
		}
		logs.Debug("Deleted RS:%s", rs.Name)
	}

	service, err := GetService(*serviceStatus, "name")
	if err != nil {
		logs.Error("Failed to delete service info in DB, service ID:%d.", service.ID)
		return err
	}
	logs.Debug("Deleted service ID %d.", service.ID)
	return nil
}

func TestMain(m *testing.M) {
	utils.InitializeDefaultConfig()
	utils.SetConfig("BASE_REPO_PATH", "/tmp/test-repos")
	utils.SetConfig("SSH_KEY_PATH", "/tmp/test-keys")

	dao.InitDB()
	os.Exit(m.Run())
}
