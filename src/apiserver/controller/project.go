package controller

import (
	"encoding/json"
	"fmt"
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
		projects, err := service.GetProjects("name", projectName)
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
		projectQuery := model.Project{ID: int64(projectID)}
		project, err := service.GetProject(projectQuery)
		if err != nil {
			c.internalError(err)
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
		fmt.Printf("Project ID: %d\n", projectID)
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
