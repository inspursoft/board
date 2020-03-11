package controller_test

import (
	"bytes"
	"encoding/json"
	"git/inspursoft/board/src/apiserver/controller"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
)

func getServiceBodys() ([][]byte, error) {
	bodies := []controller.ServiceBodyPara{
		controller.ServiceBodyPara{
			TimeUnit:      "second",
			TimeCount:     1,
			TimestampBase: time.Now().Second(),
			DurationTime:  0,
		},
		controller.ServiceBodyPara{
			TimeUnit:      "minute",
			TimeCount:     1,
			TimestampBase: time.Now().Second(),
			DurationTime:  0,
		},
		controller.ServiceBodyPara{
			TimeUnit:      "hour",
			TimeCount:     1,
			TimestampBase: time.Now().Second(),
			DurationTime:  0,
		},
		controller.ServiceBodyPara{
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

func TestGetServiceData(t *testing.T) {
	token := adminLoginTest(t)
	defer adminLogoutTest(t)

	bodies, err := getServiceBodys()
	if err != nil {
		t.FailNow()
	}

	testFunc := func(t *testing.T) {
		// init one assert
		assert := assert.New(t)
		for i := range bodies {
			//case one without parameter
			r, _ := http.NewRequest("POST", "/api/v1/dashboard/service?token="+token, bytes.NewBuffer(bodies[i]))
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)

			assert.Equal(http.StatusOK, w.Code, "Get Dashboard service data without parameter fail.")

			// case two with service parameter
			r, _ = http.NewRequest("POST", "/api/v1/dashboard/service?service_name=kubernetes"+"&token="+token, bytes.NewBuffer(bodies[i]))
			w = httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)

			assert.Equal(http.StatusOK, w.Code, "Get Dashboard service data with service parameter fail.")
		}
	}

	// insert meta data
	testFunc = prepareServiceDataWrapper("kubernetes", testFunc)
	testFunc(t)
}

func TestGetServerTime(t *testing.T) {
	token := adminLoginTest(t)
	defer adminLogoutTest(t)

	r, _ := http.NewRequest("GET", "/api/v1/dashboard/time?token="+token, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert := assert.New(t)
	assert.Equal(http.StatusOK, w.Code, "Get Server time fail.")

}
