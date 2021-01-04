package controller

import (
	"fmt"
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"

	"github.com/astaxie/beego/logs"
)

type EmailPingParam struct {
	Hostname string
	Port     int
	Username string
	Password string
	IsTLS    bool
}

type EmailSendParam struct {
	To      []string
	Subject string
	Content string
}

type grafanaNotification struct {
	Title       string      `json:"title"`
	RuleID      int         `json:"ruleId"`
	RuleName    string      `json:"ruleName"`
	RuleURL     string      `json:"ruleUrl"`
	State       string      `json:"state"`
	ImageURL    string      `json:"imageUrl"`
	Message     string      `json:"message"`
	EvalMatches []evalMatch `json:"evalMatches"`
}

type evalMatch struct {
	Metric string                 `json:"metric"`
	Tags   map[string]interface{} `json:"tags"`
	Value  interface{}            `json:"value"`
}

var emailIdentity = ""

type EmailController struct {
	c.BaseController
}

func (e *EmailController) Prepare() {
	e.EnableXSRF = false
}

func (e *EmailController) Ping() {
	var pingEmail EmailPingParam
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

func (e *EmailController) send(to []string, subject string, content string) {
	from := utils.GetStringValue("EMAIL_FROM")
	if err := service.SendMail(from, to, subject, content); err != nil {
		logs.Error("Failed to send email to error: %+v", err)
		e.InternalError(err)
		return
	}
}

func (e *EmailController) GrafanaNotification() {
	var n grafanaNotification
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

func (e *EmailController) ForgotPasswordEmail() {
	if utils.GetBoolValue("IS_EXTERNAL_AUTH") {
		e.CustomAbortAudit(http.StatusPreconditionFailed, "Resetting password doesn't support in external auth.")
		return
	}
	credential := e.GetString("credential")
	var user *model.User
	var err error
	if utils.ValidateWithPattern("email", credential) {
		user, err = service.GetUserByEmail(credential)
	} else {
		user, err = service.GetUserByName(credential)
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
	resetURL := fmt.Sprintf("http://%s/account/reset-password?reset_uuid=%s", hostIP, resetUUID)
	e.send([]string{user.Email},
		"Resetting password",
		fmt.Sprintf(`Please reset your password by clicking the URL as below:<br/><a href="%s">%s</a>`, resetURL, resetURL))
}
