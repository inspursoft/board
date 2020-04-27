package common

import "errors"

var ErrAdminLogin = errors.New("another admin user has signed in other place")
var ErrForbidden = errors.New("Forbidden")
var ErrWrongPassword = errors.New("Wrong password")
var ErrTokenServer = errors.New("tokenserver is down")
