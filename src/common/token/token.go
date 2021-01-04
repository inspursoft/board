package common

import (
	"errors"
	"fmt"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"
	"strings"

	"github.com/astaxie/beego/logs"
)

var ErrInvalidToken = errors.New("error for invalid token")

func SignToken(tokenServerURL string, payload map[string]interface{}) (*model.Token, error) {
	var token model.Token
	err := utils.RequestHandle(http.MethodPost, tokenServerURL, func(req *http.Request) error {
		req.Header = http.Header{
			"Content-Type": []string{"application/json"},
		}
		return nil
	}, payload, func(req *http.Request, resp *http.Response) error {
		return utils.UnmarshalToJSON(resp.Body, &token)
	})
	return &token, err
}

func VerifyToken(tokenServerURL string, tokenString string) (map[string]interface{}, error) {
	if strings.TrimSpace(tokenString) == "" {
		return nil, fmt.Errorf("no token provided")
	}
	var payload map[string]interface{}
	err := utils.RequestHandle(http.MethodGet, fmt.Sprintf("%s?token=%s", tokenServerURL, tokenString), nil, nil, func(req *http.Request, resp *http.Response) error {
		if resp.StatusCode == http.StatusUnauthorized {
			logs.Error("Invalid token due to session timeout.")
			return ErrInvalidToken
		}
		return utils.UnmarshalToJSON(resp.Body, &payload)
	})
	return payload, err
}
