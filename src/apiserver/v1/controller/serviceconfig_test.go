package controller_test

import (
	"bytes"
	"encoding/json"
	"github.com/inspursoft/board/src/common/model"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
)

var configServiceStep = model.ConfigServiceStep{
	ProjectID:   1,
	Instance:    1,
	ServiceName: "testservice001",
	ContainerList: []model.Container{
		model.Container{
			Name: "pod001",
			Image: model.ImageIndex{
				ImageName:   "library/mydemoshowing",
				ImageTag:    "1.0",
				ProjectName: "library",
			},
		},
	},
	ExternalServiceList: []model.ExternalService{
		model.ExternalService{
			ContainerName: "pod001",
			NodeConfig: model.NodeType{
				TargetPort: 80,
				NodePort:   32080,
			},
		},
	},
}

func TestSetConfigServiceStepAction(t *testing.T) {
	assert := assert.New(t)
	token := signIn(AdminUsername, AdminPassword)
	assert.NotEmpty(token, "Error occurred while sign in Board")
	defer func() {
		err := signOut(AdminUsername)
		assert.Nil(err, "Error occurred while sign out Board")
	}()

	body, err := json.Marshal(configServiceStep)
	assert.Nil(err, "Error occurred while testing SetConfigServiceStepAction.")

	//config service
	r, _ := http.NewRequest("POST", "/api/v1/services/config", bytes.NewReader(body))
	r.Header.Add("token", token)
	phase := url.Values{}
	phase.Add("phase", "ENTIRE_SERVICE")
	r.URL.RawQuery = phase.Encode()
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	assert.Equal(w.Code, http.StatusOK, "Error occurred while testing SetConfigServiceStepAction.")
}
