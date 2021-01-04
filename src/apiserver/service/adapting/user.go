package adapting

import (
	"github.com/inspursoft/board/src/apiserver/models/vm"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/utils"
)

func GetUserByEmail(email string) (target *vm.User, err error) {
	user, err := service.GetUserByEmail(email)
	utils.Adapt(user, &target)
	return
}

func GetUserByName(username string) (target *vm.User, err error) {
	user, err := service.GetUserByName(username)
	utils.Adapt(user, &target)
	return
}

func UpdateUser(user vm.User, selectedFields ...string) (bool, error) {
	return service.UpdateUser(user.ToMO(), selectedFields...)
}

func SignUp(user vm.User) (bool, error) {
	return service.SignUp(user.ToMO())
}
