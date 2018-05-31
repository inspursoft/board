package utils

import (
	"crypto/tls"

	"gopkg.in/gomail.v2"

	"github.com/astaxie/beego/logs"
)

type emailHandler struct {
	identity string
	hostname string
	port     int
	username string
	password string
	insecure bool
	isTLS    bool
}

func NewEmailHandler(identity, username, password, hostname string, port int) *emailHandler {
	logs.Debug("Create email handler with username: %s, host: %s, port: %d", username, hostname, port)
	return &emailHandler{
		identity: identity,
		username: username,
		password: password,
		hostname: hostname,
		port:     port,
	}
}

func (e *emailHandler) IsTLS(isTLS bool) *emailHandler {
	e.isTLS = isTLS
	logs.Debug("Set email handler TLS as: %v", isTLS)
	return e
}

func (e *emailHandler) Ping() (err error) {
	_, err = e.newDialer().Dial()
	return
}

func (e *emailHandler) Send(from string, recipients []string, subject, content string) error {
	dialer := e.newDialer()
	message := gomail.NewMessage()
	message.SetHeader("From", from)
	message.SetHeader("To", recipients...)
	message.SetHeader("Subject", subject)
	message.SetBody("text/html;charset=utf-8", content)
	return dialer.DialAndSend(message)
}

func (e *emailHandler) newDialer() *gomail.Dialer {
	dialer := gomail.NewDialer(e.hostname, e.port, e.username, e.password)
	if e.isTLS {
		dialer.TLSConfig = &tls.Config{InsecureSkipVerify: e.isTLS}
	}
	return dialer
}
