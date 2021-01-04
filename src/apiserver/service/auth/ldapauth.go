package auth

import (
	"crypto/tls"
	"errors"
	"fmt"
	"strconv"
	"strings"
	"time"

	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/dao"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"

	"github.com/astaxie/beego/logs"
	goldap "gopkg.in/ldap.v2"
)

type LdapAuth struct{}

const metaChars = "&|!=~*<>()"

func (auth LdapAuth) DoAuth(principal, password string) (*model.User, error) {

	for _, c := range metaChars {
		if strings.ContainsRune(principal, c) {
			return nil, fmt.Errorf("the principal contains meta char: %q", c)
		}
	}

	var err error
	var user *model.User

	if principal != "boardadmin" {
		var ldapconf model.LdapConf
		ldapconf.LdapURL = utils.GetStringValue("LDAP_URL")
		logs.Debug("LDAP_URL %s", ldapconf.LdapURL)
		ldapconf.LdapSearchDn = utils.GetStringValue("LDAP_SEARCH_DN")
		ldapconf.LdapSearchPassword = utils.GetStringValue("LDAP_SEARCH_PWD")
		ldapconf.LdapBaseDn = utils.GetStringValue("LDAP_BASE_DN")
		ldapconf.LdapFilter = utils.GetStringValue("LDAP_FILTER")
		ldapconf.LdapUID = utils.GetStringValue("LDAP_UID")
		ldapconf.LdapScope = utils.GetStringValue("LDAP_SCOPE")
		//GetIntValue have panic issue, fix later. Here force convert to int.
		//ldapconf.LdapConnectionTimeout = utils.GetIntValue("LDAP_TIMEOUT")
		ldapconf.LdapConnectionTimeout, err = strconv.Atoi(utils.GetStringValue("LDAP_TIMEOUT"))
		if err != nil {
			logs.Debug("Failed to get LdapConnectionTimeout, %s", err)
			return nil, err
		}

		user, err = ldapAuth(principal, password, &ldapconf)

		if err != nil {
			logs.Error("Failed to auth user with LDAP: %+v\n", err)
			return nil, nil
		}
		password = "12345678AbC"
	}

	query := model.User{Username: principal, Password: password, Deleted: 0}
	user, err = dao.GetUser(query, "username", "deleted")
	if err != nil {
		logs.Error("Failed to get user in SignIn: %+v\n", err)
		return nil, err
	}
	if user == nil {
		return nil, nil
	}
	query.Password = utils.Encrypt(query.Password, user.Salt)
	return dao.GetUser(query, "username", "password")

}

// Connect checks the LDAP configuration directives, and connects to the LDAP URL
// Returns an LDAP connection
func connect(settings *model.LdapConf) (*goldap.Conn, error) {
	ldapURL := settings.LdapURL
	if ldapURL == "" {
		return nil, errors.New("can not get any available LDAP_URL")
	}
	logs.Debug("ldapURL:", ldapURL)

	// This routine keeps compability with the old format used on harbor.cfg
	splitLdapURL := strings.Split(ldapURL, "://")
	protocol, hostport := splitLdapURL[0], splitLdapURL[1]

	var host, port string

	// This tries to detect the used port, if not defined
	if strings.Contains(hostport, ":") {
		splitHostPort := strings.Split(hostport, ":")
		host, port = splitHostPort[0], splitHostPort[1]
	} else {
		host = hostport
		switch protocol {
		case "ldap":
			port = "389"
		case "ldaps":
			port = "636"
		}
	}

	// Sets a Dial Timeout for LDAP
	goldap.DefaultTimeout = time.Duration(settings.LdapConnectionTimeout) * time.Second

	var ldap *goldap.Conn
	var err error
	switch protocol {
	case "ldap":
		ldap, err = goldap.Dial("tcp", fmt.Sprintf("%s:%s", host, port))
	case "ldaps":
		ldap, err = goldap.DialTLS("tcp", fmt.Sprintf("%s:%s", host, port), &tls.Config{InsecureSkipVerify: true})
	}

	if err != nil {
		return nil, err
	}

	return ldap, nil

}

// Authenticate checks user's credential against LDAP based on basedn template and LDAP URL,
// if the check is successful a dummy record will be inserted into DB, such that this user can
// be associated to other entities in the system.
func ldapAuth(principal string, password string, settings *model.LdapConf) (*model.User, error) {

	for _, c := range metaChars {
		if strings.ContainsRune(principal, c) {
			return nil, fmt.Errorf("the principal contains meta char: %q", c)
		}
	}

	ldap, err := connect(settings)
	if err != nil {
		return nil, err
	}

	ldapBaseDn := settings.LdapBaseDn
	if ldapBaseDn == "" {
		return nil, errors.New("can not get any available LDAP_BASE_DN")
	}
	logs.Debug("baseDn:", ldapBaseDn)

	ldapSearchDn := settings.LdapSearchDn
	if ldapSearchDn != "" {
		logs.Debug("Search DN: ", ldapSearchDn)
		ldapSearchPwd := settings.LdapSearchPassword
		err = ldap.Bind(ldapSearchDn, ldapSearchPwd)
		if err != nil {
			logs.Debug("Bind search dn error", err)
			return nil, err
		}
	}

	attrName := settings.LdapUID
	filter := settings.LdapFilter
	if filter != "" {
		filter = "(&" + filter + "(" + attrName + "=" + principal + "))"
	} else {
		filter = "(" + attrName + "=" + principal + ")"
	}
	logs.Debug("one or more filter", filter)

	ldapScope := settings.LdapScope
	var scope int
	if ldapScope == "LDAP_SCOPE_BASE" {
		scope = goldap.ScopeBaseObject
	} else if ldapScope == "LDAP_SCOPE_ONELEVEL" {
		scope = goldap.ScopeSingleLevel
	} else {
		scope = goldap.ScopeWholeSubtree
	}
	attributes := []string{"uid", "cn", "mail", "email"}

	searchRequest := goldap.NewSearchRequest(
		ldapBaseDn,
		scope,
		goldap.NeverDerefAliases,
		0,     // Unlimited results. TODO: Limit this (as we expect only one result)?
		0,     // Search Timeout. TODO: Limit this (check what is the unit of timeout) and make configurable
		false, // Types Only
		filter,
		attributes,
		nil,
	)

	result, err := ldap.Search(searchRequest)
	if err != nil {
		return nil, err
	}

	if len(result.Entries) == 0 {
		logs.Warn("Not found an entry.")
		return nil, nil
	} else if len(result.Entries) != 1 {
		logs.Warn("Found more than one entry.")
		return nil, nil
	}

	entry := result.Entries[0]
	bindDN := entry.DN
	logs.Debug("found entry:", bindDN)
	err = ldap.Bind(bindDN, password)
	if err != nil {
		logs.Debug("Bind user error", err)
		return nil, err
	}
	defer ldap.Close()

	u := model.User{}

	for _, attr := range entry.Attributes {
		val := attr.Values[0]
		switch attr.Name {
		case "uid":
			u.Realname = val
		case "cn":
			u.Realname = val
		case "mail":
			u.Email = val
		case "email":
			u.Email = val
		}
	}

	u.Username = principal
	logs.Debug("username:", u.Username, ",email:", u.Email)

	exist, err := service.UserExists("username", u.Username, 0)

	if err != nil {
		return nil, err
	}

	currentUsers := []*model.User{}

	if exist {
		currentUsers, err = service.GetUsers("username", u.Username)
		if err != nil {
			return nil, err
		}
	} else {
		u.Realname = principal
		u.Password = "12345678AbC"
		u.Comment = "registered from LDAP."
		if u.Email == "" {
			u.Email = u.Username + "@placeholder.com"
		}
		boolFlag, err := service.SignUp(u)
		if err != nil {
			return nil, err
		}
		if boolFlag {
			logs.Debug("ldap add user sucessful: username= ", u.Username)
			currentUsers = append(currentUsers, &u)
		} else {
			logs.Debug("ldap add user fail: username= ", u.Username)
			return nil, nil
		}
	}
	return currentUsers[0], nil
}

func init() {
	registerAuth("ldap_auth", LdapAuth{})
}
