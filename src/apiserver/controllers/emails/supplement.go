package emails

import (
	"fmt"
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/models/vm"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/apiserver/service/adapting"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"

	"github.com/astaxie/beego/logs"
)

var emailIdentity = ""

// Operations about Email service supplements
type SupplementController struct {
	c.BaseController
}

func (e *SupplementController) Prepare() {
	e.EnableXSRF = false
}

func (e *SupplementController) send(to []string, subject string, content string) {
	from := utils.GetStringValue("EMAIL_FROM")
	if err := service.SendMail(from, to, subject, content); err != nil {
		logs.Error("Failed to send email to error: %+v", err)
		e.InternalError(err)
		return
	}
}

// @Title Ping target Email service
// @Description Ping target Email service
// @Param	body	body	"vm.EmailPingParam"	true	"Email address"
// @Success 200 Successful executed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /ping [post]
func (e *SupplementController) Ping() {
	var pingEmail vm.EmailPingParam
	e.ResolveBody(&pingEmail)
	err := utils.NewEmailHandler(emailIdentity,
		pingEmail.Username, pingEmail.Password,
		pingEmail.Hostname, pingEmail.Port).
		IsTLS(pingEmail.IsTLS).
		Ping()
	if err != nil {
		logs.Error("Failed to ping SMTP: %+v", err)
		e.CustomAbortAudit(http.StatusBadRequest, "Failed to ping SMTP server.")
	}
}

// @Title Notify by Email service
// @Description Notify by Email service
// @Param	email	query	string	true	"Email address"
// @Success 200 Successful executed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /notification [post]
func (e *SupplementController) Notification() {
	var n vm.GrafanaNotification
	e.ResolveBody(&n)
	message := fmt.Sprintf(`<b>Title:</b>%s<br/>
	<b>Rule ID:</b> %d<br/>
	<b>Rule Name:</b> %s<br/>
	<b>Rule URL:</b> %s<br/>
	<b>State:</b> %s<br/>
	<b>Image URL:</b> %s<br/>
	<b>Message:</b> %s<br/>
	<b>Eval Matches:</b><br/>`, n.Title, n.RuleID, n.RuleName,
		n.RuleURL, n.State, n.ImageURL, n.Message)

	for _, m := range n.EvalMatches {
		message += fmt.Sprintf(` - Metric: %s<br> - Tags: %s<br/> - Value: %d<br/>`, m.Metric, m.Tags, m.Value)
	}
	e.send([]string{utils.GetStringValue("EMAIL_FROM")}, n.Title, message)
}

// @Title Forgot password by sending email
// @Description Forgot password for users.
// @Param	credential	query 	string	true	"View model for user changing password."
// @Success 200 Successful changed.
// @Failure 400 Bad requests.
// @Failure 401 Unauthorized.
// @Failure 403 Forbidden.
// @router /forgot [post]
func (e *SupplementController) Forgot() {
	if utils.GetBoolValue("IS_EXTERNAL_AUTH") {
		e.CustomAbortAudit(http.StatusPreconditionFailed, "Resetting password doesn't support in external auth.")
		return
	}
	credential := e.GetString("credential")
	var user *vm.User
	var err error
	if utils.ValidateWithPattern("email", credential) {
		user, err = adapting.GetUserByEmail(credential)
	} else {
		user, err = adapting.GetUserByName(credential)
	}
	if err != nil {
		logs.Error("Failed to get user with credential: %s, error: %+v", credential, err)
		e.InternalError(err)
		return
	}
	if user == nil {
		logs.Error("User not found with credential: %s", credential)
		e.CustomAbortAudit(http.StatusNotFound, "User not found")
		return
	}
	resetUUID := utils.GenerateRandomString()
	_, err = service.UpdateUserUUID(user.ID, resetUUID)
	if err != nil {
		logs.Error("Failed to update user reset UUID for user: %d, error: %+v", user.ID, err)
		e.InternalError(err)
		return
	}
	var hostIP = utils.GetStringValue("BOARD_HOST_IP")
	resetURL := fmt.Sprintf("http://%s/reset-password?reset_uuid=%s", hostIP, resetUUID)
	e.send([]string{user.Email},
		"Resetting password",
		fmt.Sprintf(`Please reset your password by clicking the URL as below:<br/><a href="%s">%s</a>`, resetURL, resetURL))
}
