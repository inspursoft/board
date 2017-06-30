package controller

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"
)

type controller struct {
	resp http.ResponseWriter
	req  *http.Request
}

func NewCustomController(resp http.ResponseWriter, req *http.Request) *controller {
	return &controller{resp: resp, req: req}
}

func (c *controller) GetStringFromPath(cutset string) string {
	log.Printf("%s, %s", c.req.URL.Path, cutset)
	return strings.TrimPrefix(c.req.URL.Path, cutset)
}

func (c *controller) resolveBody() []byte {
	data, err := ioutil.ReadAll(c.req.Body)
	if err != nil {
		c.customAbort(
			http.StatusInternalServerError,
			"Failed to resolve request body content",
			err)
		return nil
	}
	log.Printf("%s\n", string(data))
	return data
}

func (c *controller) assertMethod(method string) bool {
	if c.req.Method != method {
		c.customAbort(http.StatusMethodNotAllowed, "Method not allowed")
		return false
	}
	return true
}

func (c *controller) customAbort(statusCode int, message string, params ...interface{}) {
	c.resp.WriteHeader(statusCode)
	c.resp.Write([]byte(fmt.Sprintf(message, params...)))
}

func (c *controller) internalError(err error) {
	c.customAbort(http.StatusInternalServerError, "Internal server error: %+v", err)
}

func (c *controller) serveJSON(model interface{}) {
	header := c.resp.Header()
	header.Add("content-type", "application/json")
	output, err := json.Marshal(model)
	if err != nil {
		c.customAbort(http.StatusInternalServerError, "Failed to marshal object.")
		return
	}
	c.resp.Write(output)
}

type messageStatus struct {
	StatusCode int    `json:"status"`
	Message    string `json:"message"`
}

func (c *controller) serveStatus(status int, message string) {
	ms := messageStatus{
		StatusCode: status,
		Message:    message,
	}
	c.resp.WriteHeader(status)
	c.serveJSON(ms)
}
