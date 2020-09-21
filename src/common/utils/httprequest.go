package utils

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/astaxie/beego/logs"
)

var ErrBadRequest = errors.New("Bad request")
var ErrUnauthorized = errors.New("Unauthorized")
var ErrForbidden = errors.New("Forbidden")
var ErrNotFound = errors.New("Not found")
var ErrConflict = errors.New("Conflict")
var ErrUnprocessableEntity = errors.New("Unprocessable entity")
var ErrInternalError = errors.New("Internal server error")
var ErrBadGateway = errors.New("Bad gateway")
var ErrNotAcceptable = errors.New("Not Acceptable")

func EncodeString(content string) string {
	return base64.StdEncoding.EncodeToString([]byte(content))
}

func BasicAuthEncode(username, password string) string {
	return EncodeString(username + ":" + password)
}

func DefaultResponseHandler(req *http.Request, resp *http.Response) error {
	requestURL := req.URL.String()
	logs.Info("Requested: %s with response status code: %d", requestURL, resp.StatusCode)
	if resp.StatusCode >= http.StatusBadRequest {
		output, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			logs.Error("Failed to read response body: %+v", err)
			return err
		}
		logs.Debug("Error from response: %+s", string(output))
		switch resp.StatusCode {
		case http.StatusBadRequest:
			return ErrBadRequest
		case http.StatusUnauthorized:
			return ErrUnauthorized
		case http.StatusForbidden:
			return ErrForbidden
		case http.StatusNotFound:
			return ErrNotFound
		case http.StatusNotAcceptable:
			return ErrNotAcceptable
		case http.StatusConflict:
			return ErrConflict
		case http.StatusUnprocessableEntity:
			return ErrUnprocessableEntity
		case http.StatusInternalServerError:
			return ErrInternalError
		case http.StatusBadGateway:
			return ErrBadGateway
		default:
			return fmt.Errorf("unexpected error occurred while requesting %s with status code: %d", requestURL, resp.StatusCode)
		}
	}
	return nil
}

func RequestHandle(method string, urlStr string, callback func(req *http.Request) error, data interface{}, handler func(req *http.Request, resp *http.Response) error) error {
	var payload io.Reader
	var err error
	if data != nil {
		if content, ok := data.(string); ok {
			payload = bytes.NewBuffer([]byte(content))
		} else {
			obj, err := json.Marshal(data)
			if err != nil {
				log.Printf("Failed to marshal data: %+v\n", err)
				return err
			}
			payload = bytes.NewReader(obj)
		}
	}
	req, err := http.NewRequest(method, urlStr, payload)
	if err != nil {
		return err
	}
	logs.Debug("Requested URL: %s, with method: %s.", urlStr, method)
	if callback != nil {
		err = callback(req)
	}
	if err != nil {
		return err
	}
	client := http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if handler != nil {
		return handler(req, resp)
	}
	return DefaultResponseHandler(req, resp)
}

func DefaultRequestHandler(method, urlStr string, header http.Header, data interface{}) error {
	return RequestHandle(method, urlStr, func(req *http.Request) error {
		if header != nil {
			req.Header = header
		}
		return nil
	}, data, DefaultResponseHandler)
}

func SimpleGetRequestHandle(urlStr string) error {
	return RequestHandle(http.MethodGet, urlStr, nil, nil, DefaultResponseHandler)
}
func SimpleHeadRequestHandle(urlStr string, header http.Header) error {
	return DefaultRequestHandler(http.MethodHead, urlStr, header, nil)
}
func SimpleDeleteRequestHandle(urlStr string, header http.Header) error {
	return DefaultRequestHandler(http.MethodDelete, urlStr, header, nil)
}

func SimplePostRequestHandle(urlStr string, header http.Header, data interface{}) error {
	return DefaultRequestHandler(http.MethodPost, urlStr, header, data)
}
