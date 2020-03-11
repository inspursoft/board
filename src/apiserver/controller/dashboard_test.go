package controller_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"time"

	"git/inspursoft/board/src/apiserver/controller"
	"git/inspursoft/board/src/common/model/dashboard"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"github.com/stretchr/testify/assert"
)

func getBodys() ([][]byte, error) {
	bodies := []controller.DsBodyPara{
		controller.DsBodyPara{
			Service: controller.ServicePara{
				TimeUnit:      "second",
				TimeCount:     1,
				TimestampBase: time.Now().Second(),
				DurationTime:  0,
			},
			Node: controller.NodePara{
				TimeUnit:      "second",
				TimeCount:     1,
				TimestampBase: time.Now().Second(),
				DurationTime:  0,
			},
		},
		controller.DsBodyPara{
			Service: controller.ServicePara{
				TimeUnit:      "minute",
				TimeCount:     1,
				TimestampBase: time.Now().Second(),
				DurationTime:  0,
			},
			Node: controller.NodePara{
				TimeUnit:      "minute",
				TimeCount:     1,
				TimestampBase: time.Now().Second(),
				DurationTime:  0,
			},
		},
		controller.DsBodyPara{
			Service: controller.ServicePara{
				TimeUnit:      "hour",
				TimeCount:     1,
				TimestampBase: time.Now().Second(),
				DurationTime:  0,
			},
			Node: controller.NodePara{
				TimeUnit:      "hour",
				TimeCount:     1,
				TimestampBase: time.Now().Second(),
				DurationTime:  0,
			},
		},
		controller.DsBodyPara{
			Service: controller.ServicePara{
				TimeUnit:      "day",
				TimeCount:     1,
				TimestampBase: time.Now().Second(),
				DurationTime:  0,
			},
			Node: controller.NodePara{
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

func prepareNodeDataWrapper(nodeName string, testFunc func(t *testing.T)) func(t *testing.T) {
	return func(t *testing.T) {
		o := orm.NewOrm()
		logid, err := o.Insert(&dashboard.TimeListLog{RecordTime: time.Now().Unix()})
		if err != nil {
			t.Fatalf("initial node time list log data error: %v\n", err)
		}
		defer o.Delete(&dashboard.TimeListLog{Id: logid})

		nodeid, err := o.Insert(&dashboard.Node{
			NodeName:       nodeName,
			TimeListId:     logid,
			NumbersCpuCore: "1",
			NumbersGpuCore: "1",
			MemorySize:     "1024",
			PodLimit:       "1024",
			CreateTime:     "2018-01-02 05:29:54",
			InternalIp:     "127.0.0.1",
			CpuUsage:       "7.1",
			MemUsage:       "7.1",
			StorageTotal:   20,
			StorageUse:     20,
		})
		if err != nil {
			t.Fatalf("initial node data error: %v\n", err)
		}
		defer o.Delete(&dashboard.Node{Id: nodeid})

		nodeminid, err := o.Insert(&dashboard.NodeDashboardMinute{
			NodeName:       nodeName,
			TimeListId:     logid,
			NumbersCpuCore: "1",
			NumbersGpuCore: "1",
			MemorySize:     "1024",
			PodLimit:       "1024",
			CreateTime:     "2018-01-02 05:29:54",
			InternalIp:     "127.0.0.1",
			CpuUsage:       7.1,
			MemUsage:       7.1,
			StorageTotal:   20,
			StorageUse:     20,
		})
		if err != nil {
			t.Fatalf("initial node minute data error: %v\n", err)
		}
		defer o.Delete(&dashboard.NodeDashboardMinute{Id: nodeminid})

		nodehourid, err := o.Insert(&dashboard.NodeDashboardHour{
			NodeName:       nodeName,
			TimeListId:     logid,
			NumbersCpuCore: "1",
			NumbersGpuCore: "1",
			MemorySize:     "1024",
			PodLimit:       "1024",
			CreateTime:     "2018-01-02 05:29:54",
			InternalIp:     "127.0.0.1",
			CpuUsage:       7.1,
			MemUsage:       7.1,
			StorageTotal:   20,
			StorageUse:     20,
		})
		if err != nil {
			t.Fatalf("initial node hour data error: %v\n", err)
		}
		defer o.Delete(&dashboard.NodeDashboardHour{Id: nodehourid})

		nodedayid, err := o.Insert(&dashboard.NodeDashboardDay{
			NodeName:       nodeName,
			TimeListId:     logid,
			NumbersCpuCore: "1",
			NumbersGpuCore: "1",
			MemorySize:     "1024",
			PodLimit:       "1024",
			CreateTime:     "2018-01-02 05:29:54",
			InternalIp:     "127.0.0.1",
			CpuUsage:       7.1,
			MemUsage:       7.1,
			StorageTotal:   20,
			StorageUse:     20,
		})
		if err != nil {
			t.Fatalf("initial node day data error: %v\n", err)
		}
		defer o.Delete(&dashboard.NodeDashboardDay{Id: nodedayid})

		testFunc(t)
	}
}

func prepareServiceDataWrapper(serviceName string, testFunc func(t *testing.T)) func(t *testing.T) {
	return func(t *testing.T) {
		o := orm.NewOrm()
		logid, err := o.Insert(&dashboard.TimeListLog{RecordTime: time.Now().Unix()})
		if err != nil {
			t.Fatalf("initial service time list log data error: %v\n", err)
		}
		defer o.Delete(&dashboard.TimeListLog{Id: logid})

		svcid, err := o.Insert(&dashboard.ServiceDashboardSecond{
			ServiceName:     serviceName,
			TimeListId:      logid,
			PodNumber:       2,
			ContainerNumber: 4,
		})
		if err != nil {
			t.Fatalf("initial service data error: %v\n", err)
		}
		defer o.Delete(&dashboard.ServiceDashboardSecond{Id: svcid})

		svcminid, err := o.Insert(&dashboard.ServiceDashboardMinute{
			ServiceName:     serviceName,
			TimeListId:      logid,
			PodNumber:       2,
			ContainerNumber: 4,
		})
		if err != nil {
			t.Fatalf("initial service minute data error: %v\n", err)
		}
		defer o.Delete(&dashboard.ServiceDashboardMinute{Id: svcminid})

		svchourid, err := o.Insert(&dashboard.ServiceDashboardHour{
			ServiceName:     serviceName,
			TimeListId:      logid,
			PodNumber:       2,
			ContainerNumber: 4,
		})
		if err != nil {
			t.Fatalf("initial service hour data error: %v\n", err)
		}
		defer o.Delete(&dashboard.ServiceDashboardHour{Id: svchourid})

		svcdayid, err := o.Insert(&dashboard.ServiceDashboardDay{
			ServiceName:     serviceName,
			TimeListId:      logid,
			PodNumber:       2,
			ContainerNumber: 4,
		})
		if err != nil {
			t.Fatalf("initial service day data error: %v\n", err)
		}
		defer o.Delete(&dashboard.ServiceDashboardDay{Id: svcdayid})

		testFunc(t)
	}
}

func TestGetData(t *testing.T) {
	token := adminLoginTest(t)
	defer adminLogoutTest(t)

	bodies, err := getBodys()
	if err != nil {
		t.Fatal("dashboard test case initial data error")
	}
	nodeIP := os.Getenv("NODE_IP")

	testFunc := func(t *testing.T) {
		//init one assert object
		assert := assert.New(t)
		for i := range bodies {
			//case one without parameter
			r, _ := http.NewRequest("POST", "/api/v1/dashboard/data?token="+token, bytes.NewBuffer(bodies[i]))
			w := httptest.NewRecorder()
			beego.BeeApp.Handlers.ServeHTTP(w, r)

			assert.Equal(http.StatusOK, w.Code, "Get Dashboard data without parameter fail.")

			// case two with node parameter
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

	// insert meta data
	testFunc = prepareNodeDataWrapper(nodeIP, testFunc)
	testFunc = prepareServiceDataWrapper("kubernetes", testFunc)
	testFunc(t)
}
