package service

import (
	"strings"

	"github.com/astaxie/beego/context"
)

const (
	SignType      = "sgin"
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
		return ProfileType
	default:
		return EmptyType
	}
	return EmptyType
}
