package service_test

import (
	"fmt"
	"testing"

	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
)

func TestSearchSource(t *testing.T) {
	var user *model.User
	user = &model.User{Username: "aaa",ID:1,SystemAdmin:1}
	service.RegistryURL ="http://10.110.13.58:5000/v2/_catalog"
	service.NodeUrl="http://10.110.18.26:8080/api/v1/nodes"
	res, err := service.SearchSource(user, "m")
	fmt.Println(res, err)
	res, err = service.SearchSource(user, "10")
	fmt.Println(res, err)
	res, err = service.SearchSource(user, "a")
	fmt.Println(res.UserResult, err)
}
