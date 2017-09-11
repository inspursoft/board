package service_test

import (
	"fmt"
	"testing"

	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
)

func TestSearchSource(t *testing.T) {
	var user *model.User
	user = &model.User{Username: "aaa"}
	res, err := service.SearchSource(user, "b")
	fmt.Println(res, err)
}
