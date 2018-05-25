package controller

import (
	"fmt"
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

var emailIdentity = ""

type EmailController struct {
	BaseController
}

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

func (e *EmailController) Send() {
	host := utils.GetStringValue("EMAIL_HOST")
	port := utils.GetIntValue("EMAIL_PORT")
	username := utils.GetStringValue("EMAIL_USR")
	password := utils.GetStringValue("EMAIL_PWD")
	isTLS := utils.GetBoolValue("EMAIL_SSL")
	from := utils.GetStringValue("EMAIL_FROM")
	identity := utils.GetStringValue("EMAIL_IDENTITY")

	var sendEmail emailSendParam
	e.resolveBody(&sendEmail)
	err := utils.NewEmailHandler(identity, username, password, host, port).
		IsTLS(isTLS).
		Send(from, sendEmail.To, sendEmail.Subject, sendEmail.Content)
	if err != nil {
		logs.Error("Failed to send email to addr: %s, error: %+v", fmt.Sprintf("%s:%d", host, port), err)
	}
}
