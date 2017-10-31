package auth

import (
	"git/inspursoft/board/src/common/model"
)

type LdapAuth struct{}

func (auth LdapAuth) DoAuth(principal, password string) (*model.User, error) {
	return nil, nil
}

func init() {
	registerAuth("ldap_auth", LdapAuth{})
}
