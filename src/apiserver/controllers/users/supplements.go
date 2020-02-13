package users

import "github.com/astaxie/beego"

type SupplementController struct {
	beego.Controller
}

// @Title Supplement checking user existing.
// @Description Supplement for user existing.
// @Param	key	query	string 	true	"Request for probe key."
// @Param	val query	string	true	"Request for probe val."
// @Success 200 Successful checked.
// @Failure 400 Bad request.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /existing [get]
func (p *SupplementController) Exists() {

}
