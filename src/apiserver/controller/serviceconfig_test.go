package controller

import (
	"bytes"
	"encoding/json"
	"errors"
	"git/inspursoft/board/src/common/model"
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

func setConfigServiceStepAction() error {
	token := signIn("admin", "123456a?")
	if token == "" {
		return errors.New("ERR_SIGNIN")
	}

	body, err := json.Marshal(configServiceStep)
	if err != nil {
		return errors.New("ERR_MARSHAL_SERVICE_CONFIG")
	}

	//config service
	r, _ := http.NewRequest("POST", "/api/v1/services/config", bytes.NewReader(body))
	r.Header.Add("token", token)
	phase := url.Values{}
	phase.Add("phase", "ENTIRE_SERVICE")
	r.URL.RawQuery = phase.Encode()
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	if w.Code != http.StatusOK {
		return errors.New("ERR_SET_SERVICE_CONFIG_ACTION")
	}

	err = signOut("admin")
	if err != nil {
		return errors.New("ERR_SIGNOUT")
	}
	return nil
}
func TestSetConfigServiceStepAction(t *testing.T) {
	assert := assert.New(t)
	err := setConfigServiceStepAction()
	assert.Empty(err, "SetConfigServiceStepAction error.")
}
