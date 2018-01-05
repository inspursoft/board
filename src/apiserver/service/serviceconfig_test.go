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

var (
	kubeMasterUrl        = "http://10.110.18.26:8080"
	kubeMasterInvalidUrl = "http://10.110.18.26:8081"
	serviceUrl           = kubeMasterUrl + "/api/v1/namespaces/default/services/kubernetes"
	noServiceUrl         = kubeMasterUrl + "/api/v1/namespaces/default/services/kubernetesinvaild"
	invalidServiceUrl    = kubeMasterInvalidUrl + "/api/v1/namespaces/default/services/kubernetes"
	nodeUrl              = kubeMasterUrl + "/api/v1/nodes/10.110.18.71"
	noNodeUrl            = kubeMasterUrl + "/api/v1/nodes/10.110.18.70"
	invalidNodeUrl       = kubeMasterInvalidUrl + "/api/v1/nodes/10.110.18.71"
	endpointUrl          = kubeMasterUrl + "/api/v1/namespaces/default/endpoints/kubernetes"
	noEndpointUrl        = kubeMasterUrl + "/api/v1/namespaces/default/endpoints/kubernetesinvaild"
	invalidEndpointUrl   = kubeMasterInvalidUrl + "/api/v1/namespaces/default/endpoints/kubernetes"
)

func TestGetServiceStatus(t *testing.T) {
	// _, err, flag := GetServiceStatus(serviceUrl)
	// if flag == false || err != nil {
	// 	t.Errorf("Error occurred while test GetServiceStatus: %+v\n", err)
	// }
	// if flag == true && err == nil {
	// 	t.Log("GetServiceStatus is ok.\n")
	// }

	// _, err, flag = GetServiceStatus(noServiceUrl)
	// if flag == false && err != nil {
	// 	t.Log("GetServiceStatus is ok.\n")
	// }
	// if flag == true || err == nil {
	// 	t.Errorf("Error occurred while test GetServiceStatus\n")
	// }

	// _, err, flag = GetServiceStatus(invalidServiceUrl)
	// if flag == false || err == nil {
	// 	t.Errorf("Error occurred while test GetServiceStatus: %+v\n", err)
	// }
	// if flag == true && err != nil {
	// 	t.Log("GetServiceStatus is ok.\n")
	// }
}

func TestGetNodesStatus(t *testing.T) {
	// _, err, flag := GetNodesStatus(nodeUrl)
	// if flag == false || err != nil {
	// 	t.Errorf("Error occurred while test GetNodesStatus: %+v\n", err)
	// }
	// if flag == true && err == nil {
	// 	t.Log("GetNodesStatus is ok.\n")
	// }

	// _, err, flag = GetNodesStatus(noNodeUrl)
	// if flag == false && err != nil {
	// 	t.Log("GetNodesStatus is ok.\n")
	// }
	// if flag == true || err == nil {
	// 	t.Errorf("Error occurred while test GetNodesStatus\n")
	// }

	// _, err, flag = GetNodesStatus(invalidNodeUrl)
	// if flag == false || err == nil {
	// 	t.Errorf("Error occurred while test GetNodesStatus: %+v\n", err)
	// }
	// if flag == true && err != nil {
	// 	t.Log("GetNodesStatus is ok.\n")
	// }
}

func TestGetEndpointStatus(t *testing.T) {
	// _, err, flag := GetEndpointStatus(endpointUrl)
	// if flag == false || err != nil {
	// 	t.Errorf("Error occurred while test GetEndpointStatus: %+v\n", err)
	// }
	// if flag == true && err == nil {
	// 	t.Log("GetEndpointStatus is ok.\n")
	// }

	// _, err, flag = GetEndpointStatus(noEndpointUrl)
	// if flag == false && err != nil {
	// 	t.Log("GetEndpointStatus is ok.\n")
	// }
	// if flag == true || err == nil {
	// 	t.Errorf("Error occurred while test GetEndpointStatus\n")
	// }

	// _, err, flag = GetEndpointStatus(invalidEndpointUrl)
	// if flag == false || err == nil {
	// 	t.Errorf("Error occurred while test GetEndpointStatus: %+v\n", err)
	// }
	// if flag == true && err != nil {
	// 	t.Log("GetEndpointStatus is ok.\n")
	// }
}

func TestSyncServiceWithK8s(t *testing.T) {
	utils.Initialize()
	utils.AddValue("KUBE_MASTER_URL", kubeMasterUrl)
	err := SyncServiceWithK8s()
	if err != nil {
		t.Errorf("Error occurred while test SyncServiceWithK8s: %+v\n", err)
	}
	if err == nil {
		t.Log("SyncServiceWithK8s is ok.\n")
	}

}

var scCreate = model.ServiceStatus{
	ProjectID:   2,
	ProjectName: "test",
	Status:      0,
}

var scUpdate = model.ServiceStatus{
	ProjectID:   2,
	ProjectName: "testproject",
	Status:      1,
	Name:        "testservice",
	OwnerID:     2,
}

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
	serviceid, err := CreateServiceConfig(scCreate)
	assert.Nil(err, "Failed, err when create service config.")
	assert.NotEqual(0, serviceid, "Failed to assign a service id")
	t.Log("clean test", serviceid)
	cleanSeviceTestByID(serviceid)
}

func TestUpdateService(t *testing.T) {
	assert := assert.New(t)
	serviceid, err := CreateServiceConfig(scCreate)
	assert.Nil(err, "Failed, err when create service config.")
	assert.NotEqual(0, serviceid, "Failed to assign a service id")
	scUpdate.ID = serviceid
	res, err := UpdateService(scUpdate, "name", "status", "owner_id")
	assert.Nil(err, "Failed, err when update service status.")
	assert.NotEqual(false, res, "Failed to update service status")
	t.Log("updated", serviceid)
	t.Log("clean test", serviceid)
	cleanSeviceTestByID(serviceid)
}
