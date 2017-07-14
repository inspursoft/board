package controller

import (
	"encoding/json"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"strconv"
	"strings"
)

func createProjectAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod("POST") {
		reqData := c.resolveBody()
		if reqData != nil {
			var reqProject model.Project
			err := json.Unmarshal(reqData, &reqProject)
			if err != nil {
				c.internalError(err)
				return
			}
			if strings.TrimSpace(reqProject.Name) == "" {
				c.customAbort(http.StatusBadRequest, "Project name cannot be empty.")
				return
			}
			projectExists, err := service.ProjectExists(reqProject.Name)
			if err != nil {
				c.internalError(err)
				return
			}
			if projectExists {
				c.serveStatus(http.StatusConflict, "Project name already exists.")
				return
			}

			currentUser, err := getCurrentUser()
			if err != nil {
				c.internalError(err)
				return
			}
			if currentUser == nil {
				c.customAbort(http.StatusUnauthorized, "Need sign in first.")
				return
			}

			reqProject.OwnerID = int(currentUser.ID)
			reqProject.OwnerName = currentUser.Username

			isSuccess, err := service.CreateProject(reqProject)
			if err != nil {
				c.internalError(err)
				return
			}
			if !isSuccess {
				c.serveStatus(http.StatusBadRequest, "Project contains invalid characters.")
			}
		}
	}
}

func getProjectsAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod("GET") {

		projectName := req.FormValue("project_name")
		strPublic := req.FormValue("project_public")

		query := model.Project{Name: projectName, Public: 0}

		var err error

		public, err := strconv.Atoi(strPublic)
		if err == nil {
			query.Public = public
		}

		currentUser, err := getCurrentUser()
		if err != nil {
			c.internalError(err)
			return
		}
		if currentUser == nil {
			c.customAbort(http.StatusUnauthorized, "Need sign in first.")
			return
		}

		isSysAdmin, err := service.IsSysAdmin(currentUser.ID)
		if err != nil {
			c.internalError(err)
			return
		}

		var projects []*model.Project

		if isSysAdmin {
			projects, err = service.GetAllProjects(query)
		} else {
			projects, err = service.GetProjectsByUser(query, currentUser.ID)
		}

		if err != nil {
			c.internalError(err)
			return
		}
		c.serveJSON(projects)
	}
}

func ListAndCreateProjectAction(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		getProjectsAction(resp, req)
	case "POST":
		createProjectAction(resp, req)
	}
}

func getProjectAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod("GET") {
		projectID, err := strconv.Atoi(c.GetStringFromPath("id"))
		if err != nil {
			c.internalError(err)
			return
		}
		projectQuery := model.Project{ID: int64(projectID), Deleted: 0}
		project, err := service.GetProject(projectQuery, "id", "deleted")
		if err != nil {
			c.internalError(err)
			return
		}
		if project == nil {
			c.customAbort(http.StatusNotFound, "No project was found with provided ID.")
			return
		}
		c.serveJSON(project)
	}
}
func deleteProjectAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod("DELETE") {
		projectID, err := strconv.Atoi(c.GetStringFromPath("id"))
		if err != nil {
			c.internalError(err)
			return
		}

		isExists, err := service.ProjectExistsByID(int64(projectID))
		if err != nil {
			c.internalError(err)
			return
		}
		if !isExists {
			c.customAbort(http.StatusNotFound, "Cannot find project by ID")
			return
		}

		hasPermission, err := checkUserPermission(int64(projectID))
		if err != nil {
			c.internalError(err)
			return
		}
		if !hasPermission {
			c.customAbort(http.StatusForbidden, "Insuffient privileges.")
			return
		}

		isSuccess, err := service.DeleteProject(int64(projectID))
		if err != nil {
			c.internalError(err)
			return
		}
		if !isSuccess {
			c.customAbort(http.StatusBadRequest, "Failed to delete project.")
		}
	}
}

func GetAndDeleteProjectAction(resp http.ResponseWriter, req *http.Request) {
	switch req.Method {
	case "GET":
		getProjectAction(resp, req)
	case "DELETE":
		deleteProjectAction(resp, req)
	}
}

func ToggleProjectPublicAction(resp http.ResponseWriter, req *http.Request) {
	c := NewCustomController(resp, req)
	if c.assertMethod("PUT") {
		projectID, err := strconv.Atoi(c.GetStringFromPath("id"))
		if err != nil {
			c.internalError(err)
			return
		}
		if projectID == 0 {
			c.customAbort(http.StatusBadRequest, "Invalid project ID")
			return
		}

		isExists, err := service.ProjectExistsByID(int64(projectID))
		if err != nil {
			c.internalError(err)
			return
		}
		if !isExists {
			c.customAbort(http.StatusNotFound, "Cannot find project by ID")
			return
		}

		hasPermission, err := checkUserPermission(int64(projectID))
		if err != nil {
			c.internalError(err)
			return
		}
		if !hasPermission {
			c.customAbort(http.StatusForbidden, "Insuffient privileges.")
			return
		}

		reqData := c.resolveBody()
		if reqData != nil {
			var reqProject model.Project
			var err error
			err = json.Unmarshal(reqData, &reqProject)
			if err != nil {
				c.internalError(err)
				return
			}
			reqProject.ID = int64(projectID)
			isSuccess, err := service.UpdateProject(reqProject, "public")
			if err != nil {
				c.internalError(err)
				return
			}
			if !isSuccess {
				c.customAbort(http.StatusBadRequest, "Failed to update project public.")
			}
		}
	}
}
