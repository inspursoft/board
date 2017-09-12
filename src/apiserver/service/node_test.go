package service_test

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"testing"
)

func TestGetNode(t *testing.T) {
	node, err := service.GetNode("10.110.18.71")
	fmt.Println(node, err)
}
