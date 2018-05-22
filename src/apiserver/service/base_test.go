package service

import (
	"fmt"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"os"
	"testing"
	"time"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	"k8s.io/client-go/kubernetes"
	modelK8s "k8s.io/client-go/pkg/api/v1"
)

func connectToDB() {
	hostIP := os.Getenv("HOST_IP")
	orm.RegisterDriver("mysql", orm.DRMySQL)
	err := orm.RegisterDataBase("default", "mysql", fmt.Sprintf("root:root123@tcp(%s:3306)/board?charset=utf8", hostIP))
	if err != nil {
		logs.Error("Failed to connect to DB.")
	}

}

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

	serviceID, err := DeleteServiceByID(*serviceStatus)
	if err != nil {
		logs.Error("Failed to delete service info in DB, service ID:%d.", serviceID)
		return err
	}
	logs.Debug("Deleted service ID %d.", serviceID)
	return nil
}

func TestMain(m *testing.M) {
	utils.Initialize()
	utils.AddEnv("HOST_IP")
	utils.AddEnv("KUBE_MASTER_IP")
	utils.AddEnv("KUBE_MASTER_PORT")
	utils.AddEnv("NODE_IP")
	utils.AddEnv("REGISTRY_BASE_URI")
	utils.AddEnv("JENKINS_BASE_URL")

	utils.SetConfig("KUBE_MASTER_URL", "http://%s:%s", "KUBE_MASTER_IP", "KUBE_MASTER_PORT")
	utils.SetConfig("GOGITS_HOST_IP", "%s", "HOST_IP")
	utils.SetConfig("GOGITS_HOST_PORT", "10080")
	utils.SetConfig("GOGITS_SSH_PORT", "10022")
	utils.SetConfig("GOGITS_BASE_URL", "http://%s:%s", "GOGITS_HOST_IP", "GOGITS_HOST_PORT")
	utils.SetConfig("GOGITS_REPO_URL", "ssh://git@%s:%s", "GOGITS_HOST_IP", "GOGITS_SSH_PORT")
	utils.SetConfig("BASE_REPO_PATH", "/tmp/test-repos")
	utils.SetConfig("SSH_KEY_PATH", "/tmp/test-keys")

	connectToDB()
	os.Exit(m.Run())
}
