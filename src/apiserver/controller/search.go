package controller

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"net/http"
)

type SearchSourceController struct {
	BaseController
}

func (pm *SearchSourceController) Prepare() {
	user := pm.getCurrentUser()
	pm.currentUser = user
}

func (pm *SearchSourceController) Search() {
	searchCondition := pm.GetString("q")
	res, err := service.SearchSource(pm.currentUser, searchCondition)
	if err != nil {
		pm.customAbort(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	pm.renderJSON(res)
}
