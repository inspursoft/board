package service

import (
	"git/inspursoft/board/src/common/utils"
	"testing"
)

var (
	kubeMasterUrl        = "http://10.110.18.26:8080"
	kubeMasterInvalidUrl = "http://10.110.18.26:8081"
	serviceUrl           = kubeMasterUrl + "/api/v1/namespaces/default/services/demoshow2"
	noServiceUrl         = kubeMasterUrl + "/api/v1/namespaces/default/services/demoshow1"
	invalidServiceUrl    = kubeMasterInvalidUrl + "/api/v1/namespaces/default/services/demoshow1"
	nodeUrl              = kubeMasterUrl + "/api/v1/nodes/10.110.18.71"
	noNodeUrl            = kubeMasterUrl + "/api/v1/nodes/10.110.18.70"
	invalidNodeUrl       = kubeMasterInvalidUrl + "/api/v1/nodes/10.110.18.71"
	endpointUrl          = kubeMasterUrl + "/api/v1/namespaces/default/endpoints/demoshow2"
	noEndpointUrl        = kubeMasterUrl + "/api/v1/namespaces/default/endpoints/demoshow1"
	invalidEndpointUrl   = kubeMasterInvalidUrl + "/api/v1/namespaces/default/endpoints/demoshow2"
)

func TestGetServiceStatus(t *testing.T) {
	_, err, flag := GetServiceStatus(serviceUrl)
	if flag == false || err != nil {
		t.Errorf("Error occurred while test GetServiceStatus: %+v\n", err)
	}
	if flag == true && err == nil {
		t.Log("GetServiceStatus is ok.\n")
	}

	_, err, flag = GetServiceStatus(noServiceUrl)
	if flag == false && err != nil {
		t.Log("GetServiceStatus is ok.\n")
	}
	if flag == true || err == nil {
		t.Errorf("Error occurred while test GetServiceStatus\n")
	}

	_, err, flag = GetServiceStatus(invalidServiceUrl)
	if flag == false || err == nil {
		t.Errorf("Error occurred while test GetServiceStatus: %+v\n", err)
	}
	if flag == true && err != nil {
		t.Log("GetServiceStatus is ok.\n")
	}
}

func TestGetNodesStatus(t *testing.T) {
	_, err, flag := GetNodesStatus(nodeUrl)
	if flag == false || err != nil {
		t.Errorf("Error occurred while test GetNodesStatus: %+v\n", err)
	}
	if flag == true && err == nil {
		t.Log("GetNodesStatus is ok.\n")
	}

	_, err, flag = GetNodesStatus(noNodeUrl)
	if flag == false && err != nil {
		t.Log("GetNodesStatus is ok.\n")
	}
	if flag == true || err == nil {
		t.Errorf("Error occurred while test GetNodesStatus\n")
	}

	_, err, flag = GetNodesStatus(invalidNodeUrl)
	if flag == false || err == nil {
		t.Errorf("Error occurred while test GetNodesStatus: %+v\n", err)
	}
	if flag == true && err != nil {
		t.Log("GetNodesStatus is ok.\n")
	}
}

func TestGetEndpointStatus(t *testing.T) {
	_, err, flag := GetEndpointStatus(endpointUrl)
	if flag == false || err != nil {
		t.Errorf("Error occurred while test GetEndpointStatus: %+v\n", err)
	}
	if flag == true && err == nil {
		t.Log("GetEndpointStatus is ok.\n")
	}

	_, err, flag = GetEndpointStatus(noEndpointUrl)
	if flag == false && err != nil {
		t.Log("GetEndpointStatus is ok.\n")
	}
	if flag == true || err == nil {
		t.Errorf("Error occurred while test GetEndpointStatus\n")
	}

	_, err, flag = GetEndpointStatus(invalidEndpointUrl)
	if flag == false || err == nil {
		t.Errorf("Error occurred while test GetEndpointStatus: %+v\n", err)
	}
	if flag == true && err != nil {
		t.Log("GetEndpointStatus is ok.\n")
	}
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
