package controller

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"net/http"
)

type SearchSourceController struct {
	baseController
}

func (pm *SearchSourceController) Prepare() {
	user := pm.getCurrentUser()
	if user == nil {
		pm.currentUser = new(model.User)
		pm.currentUser.Username=""
		return
	}
	pm.currentUser = user
	pm.isSysAdmin = (user.SystemAdmin == 1)
	pm.isProjectAdmin = (user.ProjectAdmin == 1)
	if !pm.isProjectAdmin {
		pm.CustomAbort(http.StatusForbidden, "Insuffient privileges to for manipulating projects.")
		return
	}
}
func (pm *SearchSourceController) Search() {
	pjName := pm.GetString("search_parameter")
	res, err := service.SearchSource(pm.currentUser.Username, pjName)
	if err != nil {
		pm.CustomAbort(http.StatusInternalServerError, fmt.Sprint(err))
	}
	fmt.Println(pjName )
	pm.Data["json"] = res
	pm.ServeJSON()

}
