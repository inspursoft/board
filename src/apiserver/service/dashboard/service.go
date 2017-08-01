package dashboard

import (
	"encoding/json"
	"errors"
	"fmt"
	dao "git/inspursoft/board/src/common/dao/dashboard"
	model "git/inspursoft/board/src/common/model/dashboard"
	"log"
	"strconv"
)

type ServiceStatusLog struct {
	PodsNumber            int64 `json:"pods_number"`
	ContainerNumber       int64 `json:"container_number"`
	PodcontainerTimestamp int64 `json:"podcontainer_timestamp"`
}
type JsonOut struct {
	ServiceName       string             `json:"service_name"`
	ServiceTimeunit   string             `json:"service_timeunit"`
	ServiceCount      string             `json:"service_count"`
	ServiceStatuslogs []ServiceStatusLog `json:"service_statuslogs"`
}

func GetDashboardServiceList() []byte {
	type list struct {
		ServiceName string `json:"service_name"`
	}
	var temp list
	var serviceList []list
	s := dao.GetDashboardServiceList()
	for _, v := range s {
		temp.ServiceName = v.ServiceName
		serviceList = append(serviceList, temp)
	}
	jsonOut, _ := json.Marshal(serviceList)
	return jsonOut
}

func GetTotal(timeUnit string, timeCount string,
	timestamp string) (interface{}, error) {
	var (
		myServerList  interface{}
		sql           string
		sqlChan       chan string = make(chan string)
		timeId        []int64
		count         chan string  = make(chan string)
		doneSql       chan bool    = make(chan bool)
		doneS         chan bool    = make(chan bool)
		timestampChan chan string  = make(chan string)
		timeIdChan    chan []int64 = make(chan []int64)
	)
	go dao.PreQueryTotal(&sql, sqlChan, count, timestampChan, doneSql)
	dao.QueryTotal(timeUnit, sqlChan, timeIdChan, &myServerList, doneS)
	count <- timeCount
	timestampChan <- timestamp
	<-doneSql
	timeId = dao.QueryTime(sql)
	timeIdChan <- timeId
	<-doneS
	return myServerList, nil
}

func genStatusLog(TimeListId int64, PodNumber int64, ContainerNumber int64) ServiceStatusLog {
	var ti model.TimeListLog
	dao.QueryTimeList(&ti, TimeListId)
	return ServiceStatusLog{
		PodsNumber:            PodNumber,
		ContainerNumber:       ContainerNumber,
		PodcontainerTimestamp: ti.RecordTime,
	}

}
func genOut(timeUnit string, timeCount string,
	ServiceName string, serviceStatusLog ServiceStatusLog, Out *JsonOut) {
	*Out = JsonOut{
		ServiceName:       ServiceName,
		ServiceCount:      timeCount,
		ServiceTimeunit:   timeUnit,
		ServiceStatuslogs: append(Out.ServiceStatuslogs, serviceStatusLog),
	}

}
func GetService(timeUnit string, timeCount string, timestamp string, serviceName string) ([]byte, error) {

	server, err := GetTotal(timeUnit, timeCount, timestamp)
	var serviceOut *JsonOut
	var totalServer JsonOut
	serviceOut = new(JsonOut)
	if err != nil {
		return nil, err
	}
	switch timeUnit {
	case "hour":
		t, ok := server.(*[]model.ServiceDashboardHour)
		if ok != true {
			return nil, errors.New("Internal error：interface assertion error")
		}
		if serviceName != "" {
			for _, v := range *t {
				statusLog := genStatusLog(v.TimeListId, v.PodNumber, v.ContainerNumber)
				if serviceName == v.ServiceName || serviceName == "" {
					genOut(timeUnit, timeCount, v.ServiceName, statusLog, serviceOut)
				}
			}

		} else {
			totalCountHour(serviceName, &totalServer, t)
			*serviceOut = totalServer
			fmt.Println(totalServer)
		}
		/*if len((*serviceOut).ServiceStatuslogs) == 0 {
			return nil, errors.New("this timeUnit has not data, try smaller timeUnit")
		}*/
		jsonOut, _ := json.Marshal(serviceOut)
		return jsonOut, nil
	case "day":
		t, ok := server.(*[]model.ServiceDashboardDay)
		if ok != true {
			return nil, errors.New("Internal error：interface assertion error")
		}
		if serviceName != "" {
			for _, v := range *t {
				statusLog := genStatusLog(v.TimeListId, v.PodNumber, v.ContainerNumber)
				if serviceName == v.ServiceName || serviceName == "" {
					genOut(timeUnit, timeCount, v.ServiceName, statusLog, serviceOut)
				}
			}
		} else {
			totalCountDay(serviceName, &totalServer, t)
			*serviceOut = totalServer
			fmt.Println(totalServer)
		}
		if len((*serviceOut).ServiceStatuslogs) == 0 {
			return nil, errors.New("this timeUnit has not data, try smaller timeUnit")
		}
		jsonOut, _ := json.Marshal(serviceOut)
		return jsonOut, err

	case "minute":
		t, ok := server.(*[]model.ServiceDashboardMinute)
		fmt.Println(ok)
		if ok != true {
			return nil, errors.New("Internal error：interface assertion error")
		}
		if serviceName != "" {
			for _, v := range *t {
				statusLog := genStatusLog(v.TimeListId, v.PodNumber, v.ContainerNumber)
				if serviceName == v.ServiceName {
					genOut(timeUnit, timeCount, v.ServiceName, statusLog, serviceOut)
				}
			}
		} else {
			totalCountMinute(serviceName, &totalServer, t)
			*serviceOut = totalServer
			fmt.Println(totalServer)
		}
		if len((*serviceOut).ServiceStatuslogs) == 0 {
			return nil, errors.New("this timeUnit has not data, try smaller timeUnit")
		}
		jsonOut, _ := json.Marshal(serviceOut)
		return jsonOut, nil

	case "second":
		t, ok := server.(*[]model.ServiceDashboardSecond)
		fmt.Println(ok, serviceName)
		if ok != true {
			return nil, errors.New("Internal error：interface assertion error")
		}
		if serviceName != "" {
			for _, v := range *t {
				statusLog := genStatusLog(v.TimeListId, v.PodNumber, v.ContainerNumber)
				if serviceName == v.ServiceName {
					genOut(timeUnit, timeCount, v.ServiceName, statusLog, serviceOut)
				}
			}
		} else {
			totalCountSecond(serviceName, &totalServer, t)
			*serviceOut = totalServer
			log.Printf("totalServer: %+v\n", totalServer)
		}
		if len((*serviceOut).ServiceStatuslogs) == 0 {
			return nil, errors.New("this timeUnit has not data, try smaller timeUnit")
		}
		jsonOut, _ := json.Marshal(serviceOut)
		return jsonOut, nil
	}
	return nil, nil
}
func totalCountSecond(serviceName string, totalServer *JsonOut, t *[]model.ServiceDashboardSecond) {
	if serviceName == "" {
		for _, v := range *t {
			var ti model.TimeListLog
			dao.QueryTimeList(&ti, v.TimeListId)
			v.TimeListId = ti.RecordTime
			(*totalServer).ServiceName = "total"
			(*totalServer).ServiceTimeunit = "Second"
			(*totalServer).ServiceCount = strconv.Itoa(len(*t))
			if len((*totalServer).ServiceStatuslogs) == 0 {
				totalServiceAppend(totalServer, v.PodNumber, v.ContainerNumber, v.TimeListId)
				continue
			}
			for k, v2 := range (*totalServer).ServiceStatuslogs {
				if v.TimeListId == v2.PodcontainerTimestamp {
					(*totalServer).ServiceStatuslogs[k] = totalCac(v2.PodsNumber, v.PodNumber, v2.ContainerNumber,
						v.ContainerNumber, v.TimeListId)
					break
				} else if k == len((*totalServer).ServiceStatuslogs)-1 && v.TimeListId != v2.PodcontainerTimestamp {
					totalServiceAppend(totalServer, v.PodNumber, v.ContainerNumber, v.TimeListId)
				}
			}
		}
	}
}

func totalCountMinute(serviceName string, totalServer *JsonOut, t *[]model.ServiceDashboardMinute) {
	if serviceName == "" {
		for _, v := range *t {
			var ti model.TimeListLog
			dao.QueryTimeList(&ti, v.TimeListId)
			v.TimeListId = ti.RecordTime
			(*totalServer).ServiceName = "total"
			(*totalServer).ServiceTimeunit = "Minute"
			(*totalServer).ServiceCount=strconv.Itoa(len(*t))
			if len((*totalServer).ServiceStatuslogs) == 0 {
				totalServiceAppend(totalServer, v.PodNumber, v.ContainerNumber, v.TimeListId)
				continue
			}
			for k, v2 := range (*totalServer).ServiceStatuslogs {
				if v.TimeListId == v2.PodcontainerTimestamp {
					(*totalServer).ServiceStatuslogs[k] = totalCac(v2.PodsNumber, v.PodNumber, v2.ContainerNumber,
						v.ContainerNumber, v.TimeListId)
					break
				} else if k == len((*totalServer).ServiceStatuslogs)-1 && v.TimeListId != v2.PodcontainerTimestamp {
					totalServiceAppend(totalServer, v.PodNumber, v.ContainerNumber, v.TimeListId)
				}
			}
		}
	}
}
func totalCountHour(serviceName string, totalServer *JsonOut, t *[]model.ServiceDashboardHour) {
	if serviceName == "" {
		for _, v := range *t {
			var ti model.TimeListLog
			dao.QueryTimeList(&ti, v.TimeListId)
			v.TimeListId = ti.RecordTime
			(*totalServer).ServiceName = "total"
			(*totalServer).ServiceTimeunit = "Hour"
			if len((*totalServer).ServiceStatuslogs) == 0 {
				totalServiceAppend(totalServer, v.PodNumber, v.ContainerNumber, v.TimeListId)
				continue
			}
			for k, v2 := range (*totalServer).ServiceStatuslogs {
				if v.TimeListId == v2.PodcontainerTimestamp {
					(*totalServer).ServiceStatuslogs[k] = totalCac(v2.PodsNumber, v.PodNumber, v2.ContainerNumber,
						v.ContainerNumber, v.TimeListId)
					break
				} else if k == len((*totalServer).ServiceStatuslogs)-1 && v.TimeListId != v2.PodcontainerTimestamp {
					totalServiceAppend(totalServer, v.PodNumber, v.ContainerNumber, v.TimeListId)
				}
			}
		}
	}
}
func totalCountDay(serviceName string, totalServer *JsonOut, t *[]model.ServiceDashboardDay) {
	if serviceName == "" {
		for _, v := range *t {
			var ti model.TimeListLog
			dao.QueryTimeList(&ti, v.TimeListId)
			v.TimeListId = ti.RecordTime
			(*totalServer).ServiceName = "total"
			(*totalServer).ServiceTimeunit = "Day"
			if len((*totalServer).ServiceStatuslogs) == 0 {
				totalServiceAppend(totalServer, v.PodNumber, v.ContainerNumber, v.TimeListId)
				continue
			}
			for k, v2 := range (*totalServer).ServiceStatuslogs {
				if v.TimeListId == v2.PodcontainerTimestamp {
					(*totalServer).ServiceStatuslogs[k] = totalCac(v2.PodsNumber, v.PodNumber, v2.ContainerNumber,
						v.ContainerNumber, v.TimeListId)
					break
				} else if k == len((*totalServer).ServiceStatuslogs)-1 && v.TimeListId != v2.PodcontainerTimestamp {
					totalServiceAppend(totalServer, v.PodNumber, v.ContainerNumber, v.TimeListId)
				}
			}
		}
	}
}
func totalCac(PodsNumberA int64, PodsNumberB int64, ContainerNumberA int64, ContainerNumberB int64,
	TimeListId int64) ServiceStatusLog {
	return ServiceStatusLog{
		PodsNumber:            PodsNumberA + PodsNumberB,
		ContainerNumber:       ContainerNumberA + ContainerNumberB,
		PodcontainerTimestamp: TimeListId}
}

func totalServiceAppend(totalServer *JsonOut, PodNumber int64, ContainerNumber int64, TimeListId int64) {
	(*totalServer).ServiceStatuslogs = append(totalServer.ServiceStatuslogs, ServiceStatusLog{
		PodsNumber:            PodNumber,
		ContainerNumber:       ContainerNumber,
		PodcontainerTimestamp: TimeListId,
	})
}
