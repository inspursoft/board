package controller_test

import (
	"bytes"
	"encoding/json"
	"git/inspursoft/board/src/apiserver/controller"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
)

func getNodeBodys() ([][]byte, error) {
	bodies := []controller.NodeBodyReqPara{
		controller.NodeBodyReqPara{
			TimeUnit:      "second",
			TimeCount:     1,
			TimestampBase: time.Now().Second(),
			DurationTime:  0,
		},
		controller.NodeBodyReqPara{
			TimeUnit:      "minute",
			TimeCount:     1,
			TimestampBase: time.Now().Second(),
			DurationTime:  0,
		},
		controller.NodeBodyReqPara{
			TimeUnit:      "hour",
			TimeCount:     1,
			TimestampBase: time.Now().Second(),
			DurationTime:  0,
		},
		controller.NodeBodyReqPara{
			TimeUnit:      "day",
			TimeCount:     1,
			TimestampBase: time.Now().Second(),
			DurationTime:  0,
		},
	}
	bodySlice := make([][]byte, len(bodies))
	for i := range bodies {
		body, err := json.Marshal(bodies[i])
		if err != nil {
			return nil, err
		}
		bodySlice[i] = body
	}
	return bodySlice, nil
}

func TestGetNodeData(t *testing.T) {
	token := adminLoginTest(t)
	defer adminLogoutTest(t)

	bodies, err := getNodeBodys()
	if err != nil {
		t.Fatal("dashboard test case initial data error")
	}

	nodeIP := os.Getenv("NODE_IP")

	testFunc := func(t *testing.T) {
		// init one assert
		assert := assert.New(t)
		for i := range bodies {
			//case one without parameter
			r, _ := http.NewRequest("POST", "/api/v1/dashboard/node?token="+token, bytes.NewBuffer(bodies[i]))
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)

			assert.Equal(http.StatusOK, w.Code, "Get Dashboard node data without parameter fail.")

			// case two with node parameter
			r, _ = http.NewRequest("POST", "/api/v1/dashboard/node?node_name="+nodeIP+"&token="+token, bytes.NewBuffer(bodies[i]))
			w = httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)

			assert.Equal(http.StatusOK, w.Code, "Get Dashboard node data with node parameter fail.")

		}
	}

	// insert meta data
	testFunc = prepareNodeDataWrapper(nodeIP, testFunc)
	testFunc(t)
}
