package service

import (
	"strings"

	"github.com/astaxie/beego/context"
)

const (
	SignType        = "sign"
	UserType        = "user"
	DashboardType   = "dashboard"
	NodeGroupType   = "nodegroup"
	NodeType        = "node"
	ProjectType     = "projects"
	ServiceType     = "services"
	ImageType       = "images"
	FileType        = "file"
	SystemType      = "system"
	ProfileType     = "profile" 
	EmptyType       = ""
	//request method
	PostOperation = "POST" 
	DeleteOperation = "DELETE"
	GetOperation    = "GET"
	PutOperation = "PUT"
	PatchOperation  = "PATCH"
	//insert into DB
	CreateMethod 	= "create"
	DeleteMethod	= "delete"
	UpdateMethod	= "update"
	GetMethod		= "get"
	OtherMethod		= ""
)

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
		case strings.EqualFold(method, PostOperation):
			return CreateMethod
		case strings.EqualFold(method, DeleteOperation):
			return DeleteMethod 
		case strings.EqualFold(method, GetOperation):
			return GetMethod 
		case strings.EqualFold(method, PutOperation):
			return UpdateMethod 
		case strings.EqualFold(method, PatchOperation):
			return UpdateMethod 
		default:
			return OtherMethod
	}
}