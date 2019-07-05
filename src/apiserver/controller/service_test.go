package controller_test

import (
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
)

var configServiceStep = model.ConfigServiceStep{
	ProjectID:   1,
	Instance:    1,
	ServiceName: "teststatefulset001",
	ServiceType: model.ServiceTypeStatefulSet,
	ClusterIP:   "None",
	ContainerList: []model.Container{
		model.Container{
			Name: "nginx",
			Image: model.ImageIndex{
				ImageName:   "library/nginx",
				ImageTag:    "1.11.5",
				ProjectName: "library",
			},
		},
	},
	ExternalServiceList: []model.ExternalService{
		model.ExternalService{
			ContainerName: "nginx",
			NodeConfig: model.NodeType{
				TargetPort: 80,
				Port:       80,
			},
		},
	},
}

func TestGetServiceListAction(t *testing.T) {
	assert := assert.New(t)
	token := signIn("admin", "123456a?")
	assert.NotEmpty(token, "Error occurred while sign in Board")
	defer func() {
		err := signOut("admin")
		assert.Nil(err, "Error occurred while sign out Board")
	}()

	//get service list
	r, _ := http.NewRequest("GET", "/api/v1/services", nil)
	r.Header.Add("token", token)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	assert.Equal(w.Code, http.StatusOK, "Error occurred while testing GetServiceListAction.")
}

func TestDeployStatefulSetAction(t *testing.T) {
	assert := assert.New(t)
	token := signIn("admin", "123456a?")
	assert.NotEmpty(token, "Error occurred while sign in Board")
	defer func() {
		err := signOut("admin")
		assert.Nil(err, "Error occurred while sign out Board")
	}()

	t.Log("Test KubeMaster")
	masterIP = utils.GetConfig("KUBE_MASTER_IP")
	registryIP = utils.GetConfig("REGISTRY_IP")
	t.Log("KUBE_MASTER_IP %s  REGISTRY_IP %s", masterIP, registryIP)

	//get service list
	r, _ := http.NewRequest("GET", "/api/v1/services", nil)
	r.Header.Add("token", token)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	assert.Equal(w.Code, http.StatusOK, "Error occurred while testing GetServiceListAction.")
}
