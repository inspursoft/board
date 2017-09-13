package service_test

import (
	"fmt"
	"testing"

	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
)

func TestSearchSource(t *testing.T) {
	var user *model.User
	user = &model.User{Username: "aaa",ID:1}
	service.RegistryIp="http://10.110.13.58:5000/v2/_catalog"
	res, err := service.SearchSource(user, "m")
	fmt.Println(res, err)
}
