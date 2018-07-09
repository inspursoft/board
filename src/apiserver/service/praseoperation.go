package service

import (
	"git/inspursoft/board/src/common/model"
	"net/http"
	"strings"

	"github.com/astaxie/beego/context"
)

const (
	SignType      = "sign"
	UserType      = "user"
	DashboardType = "dashboard"
	NodeGroupType = "nodegroup"
	NodeType      = "node"
	ProjectType   = "projects"
	ServiceType   = "services"
	ImageType     = "images"
	FileType      = "file"
	SystemType    = "system"
	ProfileType   = "profile"
	EmptyType     = ""
	//insert into DB
	CreateMethod = "create"
	DeleteMethod = "delete"
	UpdateMethod = "update"
	GetMethod    = "get"
	OtherMethod  = ""
)

func ParseOperationAudit(ctx *context.Context) (operation *model.Operation) {
	objectType := GetOperationObjectType(ctx)
	if objectType == DashboardType {
		return
	}
	if objectType == SignType {
		//prase ctx to get user info
		objectType = UserType
	} else {
		operation.UserName = "anonymous"
	}
	operation.Action = GetOperationAction(ctx)
	operation.Path = ctx.Input.URL()
	operation.ObjectType = objectType
	operation.Status = model.Unknown
	return nil
}

func GetOperationObjectType(ctx *context.Context) string {
	url := ctx.Input.URL()

	switch {
	case strings.Contains(url, SignType):
		return SignType
	case strings.Contains(url, UserType):
		return UserType
	case strings.Contains(url, DashboardType):
		return DashboardType
	case strings.Contains(url, NodeGroupType):
		return NodeGroupType
	case strings.Contains(url, NodeType):
		return NodeType
	case strings.Contains(url, ProjectType):
		return ProjectType
	case strings.Contains(url, ServiceType):
		return ServiceType
	case strings.Contains(url, ImageType):
		return ImageType
	case strings.Contains(url, FileType):
		return FileType
	case strings.Contains(url, SystemType):
		return SystemType
	case strings.Contains(url, ProfileType):
		return SystemType
	default:
		return EmptyType
	}
}

func GetOperationAction(ctx *context.Context) string {
	method := strings.ToUpper(ctx.Input.Method())

	switch {
	case strings.EqualFold(method, http.MethodPost):
		return CreateMethod
	case strings.EqualFold(method, http.MethodDelete):
		return DeleteMethod
	case strings.EqualFold(method, http.MethodGet):
		return GetMethod
	case strings.EqualFold(method, http.MethodPut):
		return UpdateMethod
	case strings.EqualFold(method, http.MethodPatch):
		return UpdateMethod
	default:
		return OtherMethod
	}
}
