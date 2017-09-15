package service_test

import (
	"testing"
	"git/inspursoft/board/src/apiserver/service"
	"fmt"
)

func TestK8sCliFactory(t *testing.T) {
	defer func() { recover() }()
	service.MasterUrl = "http://10.110.18.26:8080"
	s, d := service.Suspend("10.110.18.71")
	fmt.Println(s, d)
}
