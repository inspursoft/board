package router

import (
	"git/inspursoft/board/src/apiserver/controller"

	"github.com/astaxie/beego"
)

func init() {
	ns := beego.NewNamespace("/api",
		beego.NSNamespace("/v1",
			beego.NSRouter("/sign-in",
				&controller.AuthController{},
				"post:SignInAction"),
			beego.NSRouter("/sign-up",
				&controller.AuthController{},
				"post:SignUpAction"),
			beego.NSRouter("/log-out",
				&controller.AuthController{},
				"get:LogOutAction"),
			beego.NSRouter("/users/current",
				&controller.AuthController{},
				"get:CurrentUserAction"),
			beego.NSRouter("/users/:id([0-9]+)/password",
				&controller.UserController{},
				"put:ChangePasswordAction"),
			beego.NSRouter("/adduser",
				&controller.SystemAdminController{},
				"post:AddUserAction"),
			beego.NSRouter("/users/:id([0-9]+)",
				&controller.SystemAdminController{},
				"get:GetUserAction;put:UpdateUserAction;delete:DeleteUserAction"),
			beego.NSRouter("/users",
				&controller.SystemAdminController{},
				"get:GetUsersAction"),
			beego.NSRouter("/users/:id([0-9]+)/systemadmin",
				&controller.SystemAdminController{},
				"put:ToggleSystemAdminAction"),
			beego.NSRouter("/users/:id([0-9]+)/projectadmin",
				&controller.SystemAdminController{},
				"put:ToggleProjectAdminAction"),
			beego.NSRouter("/projects",
				&controller.ProjectController{},
				"get:GetProjectsAction;post:CreateProjectAction"),
			beego.NSRouter("/projects/:id([0-9]+)/publicity",
				&controller.ProjectController{},
				"put:ToggleProjectPublicAction"),
			beego.NSRouter("/projects/:id([0-9]+)",
				&controller.ProjectController{},
				"get:GetProjectAction;delete:DeleteProjectAction"),
			beego.NSRouter("/projects/:id([0-9]+)/members",
				&controller.ProjectMemberController{},
				"get:GetProjectMembersAction;post:AddOrUpdateProjectMemberAction"),
			beego.NSRouter("/projects/:projectId([0-9]+)/members/:userId([0-9]+)",
				&controller.ProjectMemberController{},
				"delete:DeleteProjectMemberAction"),
			beego.NSRouter("/images",
				&controller.ImageController{},
				"get:GetImagesAction"),
			beego.NSRouter("/images/:imagename(.*)",
				&controller.ImageController{},
				"get:GetImageDetailAction"),
		),
	)
	beego.AddNamespace(ns)
	beego.SetStaticPath("/swagger", "swagger")
}
