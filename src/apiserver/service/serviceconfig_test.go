package service

import (
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

func TestGetServiceStatus(t *testing.T) {
	kubeMasterURL := utils.GetConfig("KUBE_MASTER_URL")
	serviceURL := kubeMasterURL() + "/api/v1/namespaces/default/services/kubernetes"
	assert := assert.New(t)
	service, err := GetServiceStatus(serviceURL)
	assert.Nil(err, "Error occurred while testing GetServiceStatus.")
	assert.NotEmpty(service, "Error occurred while testing GetServiceStatus.")
}

func TestGetNodesStatus(t *testing.T) {
	kubeMasterURL := utils.GetConfig("KUBE_MASTER_URL")
	nodeIP := utils.GetConfig("NODE_IP")
	nodeURL := kubeMasterURL() + "/api/v1/nodes/" + nodeIP()
	assert := assert.New(t)
	node, err := GetNodesStatus(nodeURL)
	assert.Nil(err, "Error occurred while testing GetNodesStatus.")
	assert.NotEmpty(node, "Error occurred while testing GetNodesStatus.")
}

func TestGetEndpointStatus(t *testing.T) {
	kubeMasterURL := utils.GetConfig("KUBE_MASTER_URL")
	endpointURL := kubeMasterURL() + "/api/v1/namespaces/default/endpoints/kubernetes"
	assert := assert.New(t)
	endpoint, err := GetEndpointStatus(endpointURL)
	assert.Nil(err, "Error occurred while testing GetEndpointStatus.")
	assert.NotEmpty(endpoint, "Error occurred while testing GetEndpointStatus.")
}

func TestSyncServiceWithK8s(t *testing.T) {
	err := SyncServiceWithK8s()
	if err != nil {
		t.Errorf("Error occurred while test SyncServiceWithK8s: %+v\n", err)
	}
	if err == nil {
		t.Log("SyncServiceWithK8s is ok.\n")
	}

}

var scCreate = model.ServiceStatus{
	ProjectID:   1,
	ProjectName: "library",
	Status:      0,
	Name:        "testservice",
}

var scUpdate = model.ServiceStatus{
	ProjectID:   1,
	ProjectName: "library",
	Status:      1,
	Name:        "testservice",
	OwnerID:     2,
}

var testService = serviceName //ToDo Should be changed to common
var testProject = "library"   //ToDo Should be changed to common

//var scID int64

func cleanSeviceTestByID(scid int64) {
	o := orm.NewOrm()
	rs := o.Raw("delete from service_status where id = ?", scid)
	r, err := rs.Exec()
	if err != nil {
		logs.Error("Error occurred while deleting service: %+v", err)
	}
	affected, err := r.RowsAffected()
	if err != nil {
		logs.Error("Error occurred while deleting service: %+v", err)
	}
	if affected == 0 {
		logs.Error("Failed to delete service: %d", scid)
	} else {
		logs.Info("Successful cleared up service: %d", scid)
	}
}

func TestCreateServiceConfig(t *testing.T) {
	assert := assert.New(t)
	serviceInfo, err := CreateServiceConfig(scCreate)
	assert.Nil(err, "Error occurred while testing creating service config.")
	assert.NotEqual(0, serviceInfo.ID, "Error occurred while assigning a service id")
	t.Log("clean test", serviceInfo.ID)
	cleanSeviceTestByID(serviceInfo.ID)
}

func TestUpdateService(t *testing.T) {
	assert := assert.New(t)
	serviceInfo, err := CreateServiceConfig(scCreate)
	assert.Nil(err, "Error occurred while creating service config.")
	serviceInfo.Status = 1
	serviceInfo.OwnerID = 2
	res, err := UpdateService(*serviceInfo, "status", "owner_id")
	assert.Nil(err, "Error occurred while updating service status.")
	assert.NotEqual(false, res, "Error occurred while updating service status")
	t.Log("updated", serviceInfo.ID)
	t.Log("clean test", serviceInfo.ID)
	cleanSeviceTestByID(serviceInfo.ID)
}

func TestDeleteServiceByID(t *testing.T) {
	assert := assert.New(t)
	serviceInfo, err := CreateServiceConfig(scCreate)
	assert.Nil(err, "Error occurred while creating service config.")

	retnum, err := DeleteServiceByID(serviceInfo.ID)
	assert.Nil(err, "Error occurred while deleting service status.")
	assert.NotEqual(0, retnum, "Error occurred while deleting service status")
	if err != nil {
		// try clean again
		cleanSeviceTestByID(serviceInfo.ID)
	}
	t.Log("deleted", serviceInfo.ID)
}

func TestGetSelectableServices(t *testing.T) {
	assert := assert.New(t)
	serviceInfo, err := CreateServiceConfig(scUpdate)
	assert.Nil(err, "Error occurred while testing creating in GetSelectableServices.")
	assert.NotEqual(0, serviceInfo.ID, "Error occurred while assigning a service id")
	serviceList, err := GetSelectableServices(serviceInfo.ProjectName, serviceInfo.Name)
	assert.Nil(err, "Error occurred while testing GetSelectableServices.")
	for _, serviceName := range serviceList {
		assert.NotEqual(serviceName, serviceInfo.Name, "Error in selectable services")
	}
	t.Log("clean test", serviceInfo.ID)
	cleanSeviceTestByID(serviceInfo.ID)
}

func TestGetK8sService(t *testing.T) {
	assert := assert.New(t)
	service, err := GetK8sService("default", "kubernetes")
	assert.Nil(err, "Error occurred while testing GetK8sService.")
	assert.Equal("kubernetes", service.Name, "Error service while testing GetK8sService.")
	t.Log("Get kubernetes service pass")
}

func TestGetDeployment(t *testing.T) {
	assert := assert.New(t)
	deployment, err := GetDeployment(testProject, testService)
	assert.Nil(err, "Error occurred while testing GetDeployment.")
	assert.Equal(deployment.Name, testService, "Error deployment while testing GetDeployment.")
	t.Log("Get kubernetes deployment pass")
}
