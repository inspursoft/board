package service

import (
	"git/inspursoft/board/src/common/dao"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"strings"

	"github.com/astaxie/beego/context"
	"github.com/astaxie/beego/logs"
)

var objectType = map[string]string{
	"sign-up":    "user",
	"users":      "user",
	"adduser":    "user",
	"search":     "system",
	"nodes":      "node",
	"profile":    "system",
	"systeminfo": "system",
}

var methodType = map[string]string{
	http.MethodPost:   "create",
	http.MethodDelete: "delete",
	http.MethodPut:    "update",
	http.MethodPatch:  "update",
	http.MethodGet:    "get",
}

func GetPaginatedOperationList(query model.OperationParam, pageIndex int, pageSize int, orderField string, orderAsc int) (*model.PaginatedOperations, error) {
	paginatedOperations, err := dao.GetPaginatedOperations(query, pageIndex, pageSize, orderField, orderAsc)
	if err != nil {
		return nil, err
	}
	return paginatedOperations, nil
}

func ParseOperationAudit(ctx *context.Context) (operation model.Operation) {
	operation.UserName = "anonymous"
	operation.Action = func(url string) string {
		if _, ok := methodType[url]; !ok {
			return "n/a"
		}
		return methodType[url]
	}(ctx.Input.Method())
	operation.Path = ctx.Input.URL()
	operation.ObjectType = func(url string) string {
		parts := strings.Split(url, "/")
		if len(parts) < 4 {
			logs.Error("URL is invalid: %s", url)
			return "n/a"
		}
		inputType := parts[3]
		if _, ok := objectType[inputType]; !ok {
			return inputType
		}
		return objectType[inputType]
	}(operation.Path)
	operation.Status = model.Unknown
	return
}

func CreateOperationAudit(operation *model.Operation) error {
	operationID, err := dao.AddOperation(*operation)
	if err != nil {
		return err
	}
	operation.ID = operationID
	return nil
}

func UpdateOperationAuditStatus(operationID int64, status int, project *model.Project, user *model.User) error {
	//Update operation result in Mysql
	operation := model.Operation{ID: operationID}

	var operationStatus string
	if status < 400 {
		operationStatus = model.Success
	} else if status < 500 {
		operationStatus = model.Failed
	} else {
		operationStatus = model.Error
	}
	operation.Status = operationStatus
	param := []string{"status"}
	if project != nil {
		operation.ProjectID = project.ID
		operation.ProjectName = project.Name
		param = append(param, "project_id", "project_name")
	}
	if user != nil {
		operation.UserID = user.ID
		operation.UserName = user.Username
		param = append(param, "user_id", "user_name")
	}

	_, err := dao.UpdateOperation(operation, param...)
	if err != nil {
		return err
	}
	return nil
}
