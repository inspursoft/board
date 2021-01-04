package service

import "github.com/inspursoft/board/src/common/utils"

func SendMail(from string, to []string, subject string, content string) error {
	identity := utils.GetStringValue("EMAIL_IDENTITY")
	host := utils.GetStringValue("EMAIL_HOST")
	port := utils.GetIntValue("EMAIL_PORT")
	username := utils.GetStringValue("EMAIL_USR")
	password := utils.GetStringValue("EMAIL_PWD")
	isTLS := utils.GetBoolValue("EMAIL_SSL")
	return utils.NewEmailHandler(identity, username, password, host, port).IsTLS(isTLS).Send(from, to, subject, content)
}
