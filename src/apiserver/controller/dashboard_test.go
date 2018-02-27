package controller

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
)

func getBodys() ([][]byte, error) {
	bodies := []DsBodyPara{
		DsBodyPara{
			Service: ServicePara{
				TimeUnit:      "second",
				TimeCount:     1,
				TimestampBase: time.Now().Second(),
				DurationTime:  0,
			},
			Node: NodePara{
				TimeUnit:      "second",
				TimeCount:     1,
				TimestampBase: time.Now().Second(),
				DurationTime:  0,
			},
		},
		DsBodyPara{
			Service: ServicePara{
				TimeUnit:      "minute",
				TimeCount:     1,
				TimestampBase: time.Now().Second(),
				DurationTime:  0,
			},
			Node: NodePara{
				TimeUnit:      "minute",
				TimeCount:     1,
				TimestampBase: time.Now().Second(),
				DurationTime:  0,
			},
		},
		DsBodyPara{
			Service: ServicePara{
				TimeUnit:      "hour",
				TimeCount:     1,
				TimestampBase: time.Now().Second(),
				DurationTime:  0,
			},
			Node: NodePara{
				TimeUnit:      "hour",
				TimeCount:     1,
				TimestampBase: time.Now().Second(),
				DurationTime:  0,
			},
		},
		DsBodyPara{
			Service: ServicePara{
				TimeUnit:      "day",
				TimeCount:     1,
				TimestampBase: time.Now().Second(),
				DurationTime:  0,
			},
			Node: NodePara{
				TimeUnit:      "day",
				TimeCount:     1,
				TimestampBase: time.Now().Second(),
				DurationTime:  0,
			},
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

func TestGetData(t *testing.T) {
	token := adminLoginTest(t)
	defer adminLogoutTest(t)

	bodies, err := getBodys()
	if err != nil {
		t.Fatal("dashboard test case initial data error")
	}

	//init one assert object
	assert := assert.New(t)
	for i := range bodies {
		//case one without parameter
		r, _ := http.NewRequest("POST", "/api/v1/dashboard/data?token="+token, bytes.NewBuffer(bodies[i]))
		w := httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		assert.Equal(http.StatusOK, w.Code, "Get Dashboard data without parameter fail.")

		// case two with node parameter
		nodeIP := os.Getenv("NODE_IP")
		r, _ = http.NewRequest("POST", "/api/v1/dashboard/data?node_name="+nodeIP+"&token="+token, bytes.NewBuffer(bodies[i]))
		w = httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		assert.Equal(http.StatusOK, w.Code, "Get Dashboard data with node parameter fail.")

		// case three with service parameter
		r, _ = http.NewRequest("POST", "/api/v1/dashboard/data?service_name=kubernetes&token="+token, bytes.NewBuffer(bodies[i]))
		w = httptest.NewRecorder()
		beego.BeeApp.Handlers.ServeHTTP(w, r)

		assert.Equal(http.StatusOK, w.Code, "Get Dashboard data with service parameter fail.")
	}

}
