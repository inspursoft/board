package controller

import (
	"fmt"
	c "git/inspursoft/board/src/apiserver/controllers/commons"
	"git/inspursoft/board/src/apiserver/service"
	"net/http"
)

type SearchSourceController struct {
	c.BaseController
}

func (pm *SearchSourceController) Prepare() {
	user := pm.GetCurrentUser()
	pm.CurrentUser = user
	pm.RecordOperationAudit()
}

func (pm *SearchSourceController) Search() {
	searchCondition := pm.GetString("q")
	res, err := service.SearchSource(pm.CurrentUser, searchCondition)
	if err != nil {
		pm.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprint(err))
		return
	}
	pm.RenderJSON(res)
}
