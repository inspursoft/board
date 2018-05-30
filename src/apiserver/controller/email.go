package controller

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/utils"
	"net/http"

	"github.com/astaxie/beego/logs"
)

type emailPingParam struct {
	Hostname string
	Port     int
	Username string
	Password string
	IsTLS    bool
}

type emailSendParam struct {
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
	Value  int                    `json:"value"`
}

var emailIdentity = ""

type EmailController struct {
	BaseController
}

func (e *EmailController) Prepare() {}

func (e *EmailController) Ping() {
	var pingEmail emailPingParam
	e.resolveBody(&pingEmail)
	err := utils.NewEmailHandler(emailIdentity,
		pingEmail.Username, pingEmail.Password,
		pingEmail.Hostname, pingEmail.Port).
		IsTLS(pingEmail.IsTLS).
		Ping()
	if err != nil {
		logs.Error("Failed to ping SMTP: %+v", err)
		e.customAbort(http.StatusBadRequest, "Failed to ping SMTP server.")
	}
}

func (e *EmailController) send(to []string, subject string, content string) {
	from := utils.GetStringValue("EMAIL_FROM")
	if err := service.SendMail(from, to, subject, content); err != nil {
		logs.Error("Failed to send email to error: %+v", err)
		e.internalError(err)
	}
}

func (e *EmailController) GrafanaNotification() {
	var n grafanaNotification
	e.resolveBody(&n)
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
