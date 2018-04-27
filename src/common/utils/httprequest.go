package utils

import (
	"encoding/base64"
	"io"
	"net/http"

	"github.com/astaxie/beego/logs"
)

func EncodeString(content string) string {
	return base64.StdEncoding.EncodeToString([]byte(content))
}

func BasicAuthEncode(username, password string) string {
	return EncodeString(username + ":" + password)
}

func RequestHandle(method string, urlStr string, callback func(req *http.Request) error, data io.Reader) (*http.Response, error) {
	client := http.Client{}
	req, err := http.NewRequest(method, urlStr, data)
	if err != nil {
		return nil, err
	}
	logs.Debug("Requested URL: %s, with method: %s.", urlStr, method)
	if callback != nil {
		err = callback(req)
	}
	if err != nil {
		return nil, err
	}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
