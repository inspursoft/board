package controller

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"net/http"
)

type SearchSourceController struct {
	baseController
}

func (pm *SearchSourceController) Prepare() {
	user := pm.getCurrentUser()
	pm.currentUser = user

}
func (pm *SearchSourceController) Search() {
	projectName := pm.GetString("search_parameter")
	res, err := service.SearchSource(pm.currentUser, projectName)
	if err != nil {
		pm.CustomAbort(http.StatusInternalServerError, fmt.Sprint(err))
	}
	pm.Data["json"] = res
	pm.ServeJSON()

}
