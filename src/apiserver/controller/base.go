package controller

import (
	"io/ioutil"
	"net/http"

	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
)

type baseController struct {
	beego.Controller
	currentUser    *model.User
	isSysAdmin     bool
	isProjectAdmin bool
}

func (b *baseController) Render() error {
	return nil
}

func (b *baseController) resolveBody() ([]byte, error) {
	data, err := ioutil.ReadAll(b.Ctx.Request.Body)
	if err != nil {
		return nil, err
	}
	return data, nil
}

type messageStatus struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

func (b *baseController) serveStatus(status int, message string) {
	ms := messageStatus{
		StatusCode: status,
		Message:    message,
	}
	b.Data["json"] = ms
	b.Ctx.ResponseWriter.WriteHeader(status)
	b.ServeJSON()
}

func (b *baseController) internalError(err error) {
	logs.Error("Error occurred: %+v", err)
	b.CustomAbort(http.StatusInternalServerError, "Unexpected error occurred.")
}

func (b *baseController) getCurrentUser() (*model.User, error) {
	return service.GetUserByID(1)
}

func (b *baseController) checkSysAdmin(user *model.User) (bool, error) {
	return service.IsSysAdmin(user.ID)
}
