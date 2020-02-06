package searches

import (
	"github.com/astaxie/beego"
)

// Operation about search
type SearchController struct {
	beego.Controller
}

// @Title Search by query item.
// @Description Search by query item.
// @Param	search	query	string	false	"Search by query item."
// @Success 200 Successful got.
// @Failure 400 Bad request.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router / [get]
func (s *SearchController) Get() {

}
