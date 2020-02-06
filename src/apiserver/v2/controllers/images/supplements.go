package images

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"path/filepath"
	"strings"

	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/apiserver/v2/models/images/vm"

	"github.com/astaxie/beego"
)

// Operations about images supplements
type SupplementController struct {
	beego.Controller
}

// @Title Preview dockerfile
// @Description Preview dockerfile.
// @Param	project_id	path	int	true	"ID of project"
// @Param 	body	body	models.images.vm.Image	true	"Request models for image."
// @Success 200 Successful executed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/preview [post]
func (c *SupplementController) Preview() {

}

// @Title Clean config about building image
// @Description Clean config about building image
// @Param	project_id	path	int	true	"ID of project"
// @Param 	body	body	models.images.vm.Image	true	"Request models for image."
// @Success 200 Successful executed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /:project_id/clean_config [post]
func (c *SupplementController) CleanConfig() {

}

// @Title Check existing about building image
// @Description Check existing about building image
// @Param	project_id	path	string	true	"Name of project"
// @Param 	body	body	vm.Image	true	"Request models for image."
// @Success 200 Image not existing.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @Failure 409 Image existing.
// @router /:project_name/exists [post]
func (c *SupplementController) Exists() {
	//Get project name
	projectName := strings.TrimSpace(c.Ctx.Input.Param(":project_name"))

	//TODO
	//Check user privilege
	//service.IsProjectMemberByName(projectName, c.CurrentUser.ID)

	//Get image name and tag
	data, err := ioutil.ReadAll(c.Ctx.Request.Body)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	body := vm.Image{}
	err = json.Unmarshal(data, &body)
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	if body.Name == "" || body.Tag == "" {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusBadRequest)
		return
	}

	//Check registry
	currentName := filepath.Join(projectName, body.Name)
	repoList, err := service.GetRegistryCatalog()
	if err != nil {
		c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
		return
	}
	for _, imageRegistry := range repoList.Names {
		if imageRegistry == currentName {
			//check tag
			tagList, err := service.GetRegistryImageTags(currentName)
			if err != nil {
				c.Ctx.ResponseWriter.WriteHeader(http.StatusInternalServerError)
				return
			}
			for _, tagID := range tagList.Tags {
				if body.Tag == tagID {
					c.Ctx.ResponseWriter.WriteHeader(http.StatusConflict)
					return
				}
			}
		}
	}
	c.Ctx.ResponseWriter.WriteHeader(http.StatusOK)
	return
}
