package controller

import (
	//"encoding/json"
	//"errors"
	//"fmt"
	//"io/ioutil"
	"net/http"
	"net/http/httptest"
	//"strings"
	"path/filepath"
	"testing"

	//"git/inspursoft/board/src/common/model"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/logs"
	"github.com/stretchr/testify/assert"
)

func TestGetImageRegistryAction(t *testing.T) {
	assert := assert.New(t)
	token := signIn("admin", "123456a?")
	defer signOut("admin")
	assert.NotEmpty(token, "signIn error")

	reqURL := "/api/v1/images/registry?token=" + token
	r, _ := http.NewRequest("GET", reqURL, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code, "Get registry fail.")
	logs.Info("Tested GetImageRegistry %s pass", w.Body.String())
}

func TestGetImagesAction(t *testing.T) {
	assert := assert.New(t)
	token := signIn("admin", "123456a?")
	defer signOut("admin")
	assert.NotEmpty(token, "signIn error")

	reqURL := "/api/v1/images?token=" + token
	r, _ := http.NewRequest("GET", reqURL, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code, "Get images fail.")
	logs.Info("Tested GetImagesAction %s pass", w.Body.String())
}

var testproject = "library"
var testimage = "nginx"

//var testerrimage = "noimage"
func TestGetImageDetailAction(t *testing.T) {
	assert := assert.New(t)
	token := signIn("admin", "123456a?")
	defer signOut("admin")
	assert.NotEmpty(token, "signIn error")

	imagepath := filepath.Join(testproject, testimage)
	reqURL := "/api/v1/images/" + imagepath + "?token=" + token
	r, _ := http.NewRequest("GET", reqURL, nil)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)

	assert.Equal(http.StatusOK, w.Code, "Get image detail fail.")
	logs.Info("Tested GetImageDetailAction %s pass", w.Body.String())
}
