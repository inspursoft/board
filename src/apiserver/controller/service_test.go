package controller

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/astaxie/beego"
	"github.com/stretchr/testify/assert"
)

func TestGetServiceListAction(t *testing.T) {
	assert := assert.New(t)
	token := signIn("admin", "123456a?")
	assert.NotEmpty(token, "Error occurred while sign in Board")
	defer func() {
		err := signOut("admin")
		assert.Nil(err, "Error occurred while sign out Board")
	}()

	//get service list
	r, _ := http.NewRequest("GET", "/api/v1/services", nil)
	r.Header.Add("token", token)
	w := httptest.NewRecorder()
	beego.BeeApp.Handlers.ServeHTTP(w, r)
	assert.Equal(w.Code, http.StatusOK, "Error occurred while testing GetServiceListAction.")
}
