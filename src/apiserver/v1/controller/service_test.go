package controller_test

import (
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"
	"net/http/httptest"
	"testing"

	//"io/ioutil"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/stretchr/testify/assert"
)

var configServiceStep1 = model.ConfigServiceStep{
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
	token := signIn(AdminUsername, AdminPassword)
	assert.NotEmpty(token, "Error occurred while sign in Board")
	defer func() {
		err := signOut(AdminUsername)
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
	token := signIn(AdminUsername, AdminPassword)
	assert.NotEmpty(token, "Error occurred while sign in Board")
	defer func() {
		err := signOut(AdminUsername)
		assert.Nil(err, "Error occurred while sign out Board")
	}()

	t.Log("Test KubeMaster")
	masterIP := utils.GetStringValue("KUBE_MASTER_IP")
	registryIP := utils.GetStringValue("REGISTRY_IP")
	logs.Info("KUBE_MASTER_IP %s  REGISTRY_IP %s", masterIP, registryIP)

	//get service list
	r, _ := http.NewRequest("GET", "/api/v1/services", nil)
	r.Header.Add("token", token)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	//data, _ := ioutil.ReadAll(w.Body)
	logs.Info("Response %s", w.Body.String())
	assert.Equal(w.Code, http.StatusOK, "Error occurred while testing GetServiceListAction.")
}
