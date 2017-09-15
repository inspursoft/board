package service_test

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"testing"
)

func TestGetNode(t *testing.T) {
	service.NodeUrl="http://10.110.18.26:8080/api/v1/nodes"
	service.MasterUrl="http://10.110.18.26:8080"
	node, err := service.GetNode("10.110.18.71")
	fmt.Println(node, err)
	a,b:=service.SuspendNode("10.110.18.71")
	fmt.Println(a,b)
	a,b=service.ResumeNode("10.110.18.71")
	fmt.Println(a,b)
}
