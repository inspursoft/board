package utils

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"html/template"
	"net"
	"net/smtp"
	"time"

	"github.com/astaxie/beego/logs"
)

const emailTpl = "From: {{ .From }}\r\n" +
	"To: {{range .To}}{{ . }},{{end}}\r\n" +
	"Subject: {{ .Subject }}\r\n" +
	"MIME-version: 1.0;\r\n" +
	"Content-Type: text/html; charset=\"UTF-8\"\r\n" +
	"\n" +
	"{{ .Content }}" +
	"\r\n"

const emailConnectionTimeout = time.Duration(time.Second * 60)

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
	client, err := e.newClient()
	if err != nil {
		return
	}
	logs.Debug("Successful ping SMTP server.")
	return client.Close()
}

func (e *emailHandler) Send(from string, recipients []string, subject, content string) error {
	var buf bytes.Buffer
	t, err := template.New("email").Parse(emailTpl)
	if err != nil {
		return err
	}
	err = t.Execute(&buf, struct {
		From    string
		To      []string
		Subject string
		Content string
	}{
		From:    from,
		To:      recipients,
		Subject: subject,
		Content: content,
	})
	if err != nil {
		return err
	}
	client, err := e.newClient()
	if err != nil {
		return err
	}
	defer client.Close()

	if err = client.Mail(from); err != nil {
		return err
	}
	for _, t := range recipients {
		if err = client.Rcpt(t); err != nil {
			return err
		}
	}
	w, err := client.Data()
	if err != nil {
		return err
	}
	_, err = w.Write(buf.Bytes())
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return client.Quit()
}

func (e *emailHandler) newClient() (*smtp.Client, error) {
	address := fmt.Sprintf("%s:%d", e.hostname, e.port)
	conn, err := net.DialTimeout("tcp", address, emailConnectionTimeout)
	if err != nil {
		return nil, err
	}
	if e.isTLS {
		tlsConn := tls.Client(conn, &tls.Config{
			ServerName:         e.hostname,
			InsecureSkipVerify: e.insecure,
		})
		if err = tlsConn.Handshake(); err != nil {
			return nil, err
		}
		conn = tlsConn
	}

	client, err := smtp.NewClient(conn, e.hostname)
	if err != nil {
		return nil, err
	}

	if !e.isTLS {
		logs.Debug("Detecting whether the email server support STARTTLS ...")
		if isSupportSTARTTLS, _ := client.Extension("STARTTLS"); isSupportSTARTTLS {
			err = client.StartTLS(&tls.Config{
				ServerName:         e.hostname,
				InsecureSkipVerify: e.insecure,
			})
			if err != nil {
				return nil, err
			}
			e.isTLS = true
			logs.Debug("Successful switch email client as STARTTLS...")
		} else {
			logs.Warning("The SMTP server %s does not support STARTTLS.", address)
		}
	}

	logs.Debug("Detecting whether the email server support AUTH ...")
	if isSupportAUTH, _ := client.Extension("AUTH"); isSupportAUTH {

		var auth smtp.Auth
		if e.isTLS {
			auth = smtp.PlainAuth(e.identity, e.username, e.password, e.hostname)
		} else {
			auth = smtp.CRAMMD5Auth(e.username, e.password)
		}
		err = client.Auth(auth)
		if err != nil {
			return nil, err
		}
	} else {
		logs.Warning("The SMTP server %s does not support AUTH.", address)
	}
	logs.Debug("Successful create SMTP client.")
	return client, nil
}
