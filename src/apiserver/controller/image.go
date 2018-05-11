package controller

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"

	"strings"

	"github.com/astaxie/beego/logs"
)

type ImageController struct {
	baseController
}

const (
	commentTemp           = "Inspur image" // TODO: get from mysql in the next release
	sizeunitTemp          = "B"
	adminID               = 1
	defaultDockerfilename = "Dockerfile"
	imageProcess          = "process-image"
)

// API to get image list
func (p *ImageController) GetImagesAction() {

	var repolist model.RegistryRepo
	var repolistFiltered model.RegistryRepo
	// Get the image list from registry v2
	httpresp, err := http.Get(registryURL() + "/v2/_catalog")
	if err != nil {
		p.internalError(err)
		return
	}

	body, err := ioutil.ReadAll(httpresp.Body)
	if err != nil {
		p.internalError(err)
		return
	}

	err = json.Unmarshal(body, &repolist)
	if err != nil {
		p.internalError(err)
		return
	}

	query := model.Project{}
	projectList, err := service.GetProjectsByUser(query, p.currentUser.ID)
	if err != nil {
		p.internalError(err)
		return
	}
	for _, imageName := range repolist.Names {
		fromIndex := strings.LastIndex(imageName, "/")
		if fromIndex == -1 {
			continue
		}
		for _, project := range projectList {
			if imageName[:fromIndex] == project.Name {
				repolistFiltered.Names = append(repolistFiltered.Names, imageName)
				break
			}
		}
	}

	logs.Info("Image list is %+v\n", repolistFiltered)

	/* Interpret the message to api server */
	imagelist := []model.Image{}
	for _, imagename := range repolistFiltered.Names {
		var newImage model.Image
		newImage.ImageName = imagename

		var reqTagList tagList
		tagListURL := registryURL() + "/v2/" + imagename + "/tags/list"
		httpresp, err := http.Get(tagListURL)
		if err != nil {
			p.internalError(err)
			return
		}
		body, err := ioutil.ReadAll(httpresp.Body)
		if err != nil {
			p.internalError(err)
			return
		}
		defer httpresp.Body.Close()

		err = json.Unmarshal(body, &reqTagList)
		if err != nil {
			p.internalError(err)
			return
		}
		if len(reqTagList.Tags) == 0 {
			continue
		}

		// Check image in DB
		dbimage, err := service.GetImage(newImage, "name")
		if err != nil {
			p.customAbort(http.StatusInternalServerError, fmt.Sprintf("Checking image name in DB error: %+v", err))
			return
		}
		if dbimage != nil {
			// image already in DB, use the status in DB
			newImage.ImageID = dbimage.ImageID
			newImage.ImageComment = dbimage.ImageComment
			newImage.ImageDeleted = dbimage.ImageDeleted
		} else {
			// image not in DB, add it to DB
			newImage.ImageComment = commentTemp
			id, err := service.CreateImage(newImage)
			if err != nil {
				p.customAbort(http.StatusInternalServerError, fmt.Sprintf("Create image in DB error: %+v", err))
				return
			}
			newImage.ImageID = id
		}

		if newImage.ImageDeleted == 0 {
			imagelist = append(imagelist, newImage)
		}
	}
	logs.Info(imagelist)
	p.Data["json"] = imagelist
	p.ServeJSON()
}

// API to get tag list for a specific image
func (p *ImageController) GetImageDetailAction() {

	var taglist model.RegistryTags

	imageName := p.Ctx.Input.Param(":imagename")

	gettagsurl := "/v2/" + imageName + "/tags/list"

	httpresp, err := http.Get(registryURL() + gettagsurl)
	if err != nil {
		logs.Debug("Get image detail URL: %s", gettagsurl)
		p.internalError(err)
		return
	}

	body, err := ioutil.ReadAll(httpresp.Body)
	if err != nil {
		p.internalError(err)
		return
	}

	err = json.Unmarshal(body, &taglist)
	if err != nil {
		p.internalError(err)
		return
	}
	//fmt.Println(taglist)

	var imagedetail []model.TagDetail
	for _, tagid := range taglist.Tags {
		var tagdetail model.TagDetail
		tagdetail.ImageName = taglist.ImageName
		tagdetail.ImageTag = tagid
		tagdetail.ImageSizeUnit = sizeunitTemp

		// Get version one schema
		getmanifesturl := "/v2/" + taglist.ImageName + "/manifests/" + tagid
		httpresp, err = http.Get(registryURL() + getmanifesturl)
		if err != nil {
			logs.Debug(getmanifesturl)
			p.internalError(err)
			return
		}

		body, err = ioutil.ReadAll(httpresp.Body)
		if err != nil {
			p.internalError(err)
			return
		}

		var manifest1 model.RegistryManifest1
		err = json.Unmarshal(body, &manifest1)
		if err != nil {
			p.internalError(err)
			return
		}

		//fmt.Println((manifest1.History[0])["v1Compatibility"])

		// Interpret it on the frontend
		tagdetail.ImageDetail = (manifest1.History[0])["v1Compatibility"]
		tagdetail.ImageAuthor = ""       //TODO: get the author by frontend simply
		tagdetail.ImageCreationTime = "" //TODO: get the time by frontend simply

		// Get version two schema
		getmanifesturl = registryURL() + getmanifesturl
		req, _ := http.NewRequest("GET", getmanifesturl, nil)
		req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
		client := http.Client{}
		httpresp, err = client.Do(req)

		body, err = ioutil.ReadAll(httpresp.Body)
		if err != nil {
			p.internalError(err)
			return
		}

		var manifest2 model.RegistryManifest2
		err = json.Unmarshal(body, &manifest2)
		if err != nil {
			p.internalError(err)
			return
		}

		tagdetail.ImageId = manifest2.Config.Digest
		tagdetail.ImageSize = manifest2.Config.Size

		var layerconfig model.Manifest2Config
		for _, layerconfig = range manifest2.Layers {
			tagdetail.ImageSize += layerconfig.Size
		}

		// Add the tag detail to list
		imagedetail = append(imagedetail, tagdetail)

	}
	logs.Debug(imagedetail)
	p.Data["json"] = imagedetail
	p.ServeJSON()

}

//  Checking the user priviledge by token
func (p *ImageController) Prepare() {
	user := p.getCurrentUser()
	if user == nil {
		p.customAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
	p.isSysAdmin = (user.SystemAdmin == 1)
}

func (p *ImageController) BuildImageAction() {
	var err error
	//Check user priviledge project admin
	reqData, err := p.resolveBody()
	if err != nil {
		p.internalError(err)
		return
	}

	var reqImageConfig model.ImageConfig
	err = json.Unmarshal(reqData, &reqImageConfig)
	if err != nil {
		p.internalError(err)
		return
	}

	//Checking invalid parameters
	p.resolveRepoPath(reqImageConfig.ProjectName)
	reqImageConfig.RepoPath = p.repoPath
	err = service.CheckDockerfileConfig(&reqImageConfig)
	if err != nil {
		p.serveStatus(http.StatusBadRequest, err.Error())
		return
	}

	currentProject, err := service.GetProject(model.Project{Name: reqImageConfig.ProjectName}, "name")
	if err != nil {
		p.internalError(err)
		return
	}
	if currentProject == nil {
		p.customAbort(http.StatusBadRequest, "Invalid project name.")
		return
	}

	isMember, err := service.IsProjectMember(currentProject.ID, p.currentUser.ID)
	if err != nil {
		p.internalError(err)
		return
	}

	if !(p.isSysAdmin || isMember) {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to build image.")
		return
	}

	reqImageConfig.ImageDockerfilePath = reqImageConfig.RepoPath

	// Check image:tag existing in registry
	existing, err := existRegistry(reqImageConfig.ProjectName, reqImageConfig.ImageName,
		reqImageConfig.ImageTag)
	if err != nil {
		p.internalError(err)
		return
	}

	if existing {
		logs.Error("This image:tag existing in registry %s", reqImageConfig.ImageDockerfilePath)
		p.customAbort(http.StatusConflict, "This image:tag already existing.")
		return
	}

	err = service.BuildDockerfile(reqImageConfig)
	if err != nil {
		p.internalError(err)
		return
	}

	repoPath := reqImageConfig.RepoPath
	projectName := reqImageConfig.ProjectName
	imageName := reqImageConfig.ImageName
	imageTag := reqImageConfig.ImageTag

	username := p.currentUser.Username
	email := p.currentUser.Email
	imageURI := filepath.Join(registryBaseURI(), projectName, imageName) + ":" + imageTag

	err = service.GenerateBuildingImageTravis(repoPath, username, email, imageURI)
	if err != nil {
		logs.Error("Failed to generate building image Travis.yml: %+v", err)
		return
	}
	if currentToken, ok := memoryCache.Get(username).(string); ok {
		service.CreateFile("key.txt", currentToken, repoPath)
	}

	items := []string{".travis.yml", "key.txt", "Dockerfile"}
	err = p.pushItemsToRepo(repoPath, items...)
	if err != nil {
		logs.Error("Failed to push to repo: %s for BuildImageAction, error: %+v", repoPath)
		p.internalError(err)
	}
	p.collaborateWithPullRequest(repoPath, "master", "master", items...)
}

func (p *ImageController) GetImageDockerfileAction() {
	projectName := strings.TrimSpace(p.GetString("project_name"))
	p.resolveRepoPath(projectName)
	dockerfilePath := p.repoPath
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		p.customAbort(http.StatusNotFound, "Image path does not exist.")
		return
	}
	dockerfile, err := service.GetDockerfileInfo(dockerfilePath)
	if err != nil {
		p.internalError(err)
		return
	}
	p.Data["json"] = dockerfile
	p.ServeJSON()
}

func (p *ImageController) DockerfilePreviewAction() {
	reqData, err := p.resolveBody()
	if err != nil {
		p.internalError(err)
		return
	}

	var reqImageConfig model.ImageConfig
	err = json.Unmarshal(reqData, &reqImageConfig)
	if err != nil {
		p.internalError(err)
		return
	}

	p.resolveRepoPath(reqImageConfig.ProjectName)

	reqImageConfig.RepoPath = p.repoPath
	//Checking invalid parameters
	err = service.CheckDockerfileConfig(&reqImageConfig)
	if err != nil {
		p.serveStatus(http.StatusBadRequest, err.Error())
		return
	}

	currentProject, err := service.GetProject(model.Project{Name: reqImageConfig.ProjectName}, "name")
	if err != nil {
		p.internalError(err)
		return
	}
	if currentProject == nil {
		p.customAbort(http.StatusBadRequest, "Invalid project name.")
		return
	}

	isMember, err := service.IsProjectMember(currentProject.ID, p.currentUser.ID)
	if err != nil {
		p.internalError(err)
		return
	}

	if !(p.isSysAdmin || isMember) {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to build image.")
		return
	}

	reqImageConfig.ImageDockerfilePath = p.repoPath
	err = service.BuildDockerfile(reqImageConfig, p.Ctx.ResponseWriter)
	if err != nil {
		p.internalError(err)
		return
	}
}

func (p *ImageController) ConfigCleanAction() {
	imageName := strings.TrimSpace(p.GetString("image_name"))
	imageTag := strings.TrimSpace(p.GetString("image_tag"))
	projectName := strings.TrimSpace(p.GetString("project_name"))
	logs.Debug("clean config %s %s %s", projectName, imageName, imageTag)

	currentProject, err := service.GetProject(model.Project{Name: projectName}, "name")
	if err != nil {
		p.internalError(err)
		return
	}
	if currentProject == nil {
		p.customAbort(http.StatusBadRequest, "Invalid project name.")
		return
	}

	isMember, err := service.IsProjectMember(currentProject.ID, p.currentUser.ID)
	if err != nil {
		p.internalError(err)
		return
	}

	if !(p.isSysAdmin || isMember) {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to build image.")
		return
	}

	//remove uploaded directory
	uploadedPath := filepath.Join(baseRepoPath(), p.currentUser.Username, "upload")
	err = os.RemoveAll(uploadedPath)
	if err != nil {
		logs.Error("Failed to remove uploaded path: %s, error: %+v", uploadedPath, err)
		p.internalError(err)
		return
	}

	//remove attachment file
	err = os.RemoveAll(filepath.Join(baseRepoPath(), p.currentUser.Username, attachmentFile))
	if err != nil {
		logs.Error("Failed to remove attachment file: %+v", err)
		p.internalError(err)
		return
	}

	p.resolveRepoPath(currentProject.Name)
	//remove items to git repo
	p.removeItemsToRepo(p.repoPath)

	//Delete the config files
	err = service.ImageConfigClean(p.repoPath)
	if err != nil {
		p.internalError(err)
		return
	}
}

type tagList struct {
	Name string   `json:"name"`
	Tags []string `json:"tags"`
}

func (p *ImageController) DeleteImageAction() {

	if p.isSysAdmin == false {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to delete image.")
		return
	}

	imageName := strings.TrimSpace(p.GetString("image_name"))
	projectName := imageName[:strings.Index(imageName, "/")]

	URLPrefix := registryURL() + `/v2/` + imageName
	tagListURL := URLPrefix + `/tags/list`
	resp, err := http.Get(tagListURL)
	if err != nil {
		p.internalError(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		p.customAbort(resp.StatusCode, "repository name not known to registry")
		return
	}

	reqBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		p.internalError(err)
		return
	}

	resp.Body.Close()

	var reqTagList tagList
	err = json.Unmarshal(reqBody, &reqTagList)
	if err != nil {
		p.internalError(err)
		return
	}

	URLPrefix += `/manifests/`
	for _, tag := range reqTagList.Tags {
		var client = &http.Client{}
		manifestsURL := URLPrefix + tag
		req, err := http.NewRequest("HEAD", manifestsURL, nil)
		req.Header.Add("Accept", `application/vnd.docker.distribution.manifest.v2+json`)
		resp, err = client.Do(req)
		if err != nil {
			p.internalError(err)
			return
		}
		resp.Body.Close()

		digest := strings.Trim(resp.Header.Get("Etag"), `"`)
		deleteURL := URLPrefix + digest
		req, err = http.NewRequest("DELETE", deleteURL, nil)
		if err != nil {
			p.internalError(err)
			return
		}

		resp, err = client.Do(req)
		if err != nil {
			p.internalError(err)
			return
		}
		if resp.StatusCode != http.StatusAccepted {
			errString := fmt.Sprintf("Remove registry image tag: %s", tag)
			p.customAbort(http.StatusInternalServerError, errString)
			return
		}
		resp.Body.Close()

		p.resolveRepoPath(projectName)

		//Remove items to git repo
		err = p.removeItemsToRepo(p.repoPath)
		if err != nil {
			logs.Error("Failed to remove items to repo: %+v", err)
			p.internalError(err)
			return
		}

		//Delete the config files
		err = service.ImageConfigClean(p.repoPath)
		if err != nil {
			logs.Error("Failed to delete config files at %s", p.repoPath)
			p.internalError(err)
		}
	}
}

func (p *ImageController) DeleteImageTagAction() {
	var err error

	if p.isSysAdmin == false {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to delete image tag.")
		return
	}

	imageName := strings.TrimSpace(p.Ctx.Input.Param(":imagename"))
	imageTag := strings.TrimSpace(p.GetString("image_tag"))

	projectName := imageName[:strings.Index(imageName, "/")]

	var client = &http.Client{}
	URLPrefix := registryURL() + `/v2/` + imageName + `/manifests/`
	manifestsURL := URLPrefix + imageTag
	req, err := http.NewRequest("HEAD", manifestsURL, nil)
	req.Header.Add("Accept", `application/vnd.docker.distribution.manifest.v2+json`)
	resp, err := client.Do(req)
	if err != nil {
		p.internalError(err)
		return
	}
	if resp.StatusCode != http.StatusOK {
		p.customAbort(resp.StatusCode, "Repository name or tag not to known to registry")
		return
	}

	resp.Body.Close()

	digest := strings.Trim(resp.Header.Get("Etag"), `"`)
	deleteURL := URLPrefix + digest
	req, err = http.NewRequest("DELETE", deleteURL, nil)
	if err != nil {
		p.internalError(err)
		return
	}

	resp, err = client.Do(req)
	if err != nil {
		p.internalError(err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusAccepted {
		p.customAbort(http.StatusInternalServerError, "Remove registry image error")
		return
	}

	p.resolveRepoPath(projectName)
	//Delete the config files
	err = service.ImageConfigClean(p.repoPath)
	if err != nil {
		logs.Error("Failed to delete config files %s", p.repoPath)
		p.internalError(err)
		return
	}

	// Clean repo to git repo
	err = p.removeItemsToRepo(p.repoPath)
	if err != nil {
		logs.Error("Failed to remove items to repo: %+v", err)
	}
}

func (p *ImageController) DockerfileBuildImageAction() {
	imageName := strings.TrimSpace(p.GetString("image_name"))
	imageTag := strings.TrimSpace(p.GetString("image_tag"))
	projectName := strings.TrimSpace(p.GetString("project_name"))

	p.resolveRepoPath(projectName)
	dockerfilePath := filepath.Join(p.repoPath, imageProcess, imageName, imageTag)
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		p.customAbort(http.StatusNotFound, "Image path does not exist.")
		return
	}

	currentProject, err := service.GetProject(model.Project{Name: projectName}, "name")
	if err != nil {
		p.internalError(err)
		return
	}
	if currentProject == nil {
		p.customAbort(http.StatusBadRequest, "Invalid project name.")
		return
	}

	isMember, err := service.IsProjectMember(currentProject.ID, p.currentUser.ID)
	if err != nil {
		p.internalError(err)
		return
	}

	if !(p.isSysAdmin || isMember) {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to build image.")
		return
	}

	username := p.currentUser.Username
	email := p.currentUser.Email
	imageURI := filepath.Join(registryBaseURI(), projectName, imageName) + ":" + imageTag
	err = service.GenerateBuildingImageTravis(p.repoPath, username, email, imageURI)
	if err != nil {
		logs.Error("Failed to generate building image Travis.yml: %+v", err)
		return
	}

	items := []string{".travis.yml", "Dockerfile"}
	err = p.pushItemsToRepo(p.repoPath, items...)
	if err != nil {
		logs.Error("Failed to push items to repo: %+v", err)
		p.internalError(err)
	}
	p.collaborateWithPullRequest(p.repoPath, "master", "master", items...)
}

func exists(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return true, nil
	}
	if os.IsNotExist(err) {
		return false, nil
	}
	return true, err
}

func (p *ImageController) CheckImageTagExistingAction() {
	var err error

	imageName := strings.TrimSpace(p.Ctx.Input.Param(":imagename"))
	imageTag := strings.TrimSpace(p.GetString("image_tag"))
	projectName := strings.TrimSpace(p.GetString("project_name"))

	currentProject, err := service.GetProject(model.Project{Name: projectName}, "name")
	if err != nil {
		p.internalError(err)
		return
	}
	if currentProject == nil {
		p.customAbort(http.StatusBadRequest, "Invalid project name.")
		return
	}

	isMember, err := service.IsProjectMember(currentProject.ID, p.currentUser.ID)
	if err != nil {
		p.internalError(err)
		return
	}

	if !(p.isSysAdmin || isMember) {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to build image.")
		return
	}

	// check this image:tag in system
	p.resolveRepoPath(projectName)
	dockerfilePath := filepath.Join(p.repoPath, imageProcess, imageName, imageTag)
	existing, err := exists(dockerfilePath)
	if err != nil {
		p.internalError(err)
		return
	}

	if existing {
		logs.Info("This image:tag existing in system %s", dockerfilePath)
		p.customAbort(http.StatusConflict, "This image:tag already existing.")
		return
	}

	// TODO check image imported from registry
	existing, err = existRegistry(projectName, imageName, imageTag)
	if err != nil {
		p.internalError(err)
		return
	}

	if existing {
		logs.Info("This image:tag existing in system %s", dockerfilePath)
		p.customAbort(http.StatusConflict, "This image:tag already existing.")
		return
	}

	logs.Debug("checking image:tag result %t", existing)
	p.ServeJSON()
	return
}

func existRegistry(projectName string, imageName string, imageTag string) (bool, error) {
	var repolist model.RegistryRepo
	realName := filepath.Join(projectName, imageName)

	//check image
	httpresp, err := http.Get(registryURL() + "/v2/_catalog")
	if err != nil {
		logs.Error("Get image URL: %s", registryURL())
		return true, err
	}

	body, err := ioutil.ReadAll(httpresp.Body)
	if err != nil {
		logs.Error("Failed to read image body %+v", err)
		return true, err
	}

	err = json.Unmarshal(body, &repolist)
	if err != nil {
		logs.Error("Failed to unmarshal repolist body %+v", err)
		return true, err
	}
	for _, imageRegistry := range repolist.Names {
		if imageRegistry == realName {
			//check tag
			var taglist model.RegistryTags
			gettagsurl := "/v2/" + realName + "/tags/list"

			httpresp, err := http.Get(registryURL() + gettagsurl)
			if err != nil {
				logs.Error("Get image detail URL: %s", gettagsurl)
				return true, err
			}

			body, err := ioutil.ReadAll(httpresp.Body)
			if err != nil {
				logs.Error("Failed to read body %+v", err)
				return true, err
			}

			err = json.Unmarshal(body, &taglist)
			if err != nil {
				logs.Error("Failed to unmarshal body %+v", err)
				return true, err
			}

			for _, tagid := range taglist.Tags {
				if imageTag == tagid {
					logs.Info("Image tag existing %s:%s", realName, tagid)
					return true, nil
				}
			}
		}
	}
	return false, err
}

func (f *ImageController) UploadDockerfileFileAction() {
	projectName := f.GetString("project_name")
	isExistence, err := service.ProjectExists(projectName)
	if err != nil {
		f.internalError(err)
		return
	}
	if isExistence != true {
		f.customAbort(http.StatusBadRequest, "Project don't exist.")
		return
	}
	f.resolveRepoPath(projectName)
	targetFilePath := f.repoPath
	err = os.MkdirAll(targetFilePath, 0755)
	if err != nil {
		f.internalError(err)
		return
	}
	logs.Info("User: %s uploaded Dockerfile file to %s.", f.currentUser.Username, targetFilePath)

	_, fileHeader, err := f.GetFile("upload_file")
	if err != nil {
		f.internalError(err)
	}
	if fileHeader.Filename != dockerfileName {
		f.customAbort(http.StatusBadRequest, "Update file name invalid.")
		return
	}
	err = f.SaveToFile("upload_file", filepath.Join(targetFilePath, dockerfileName))
	if err != nil {
		f.internalError(err)
	}

}

func (f *ImageController) DownloadDockerfileFileAction() {
	projectName := f.GetString("project_name")
	isExistence, err := service.ProjectExists(projectName)
	if err != nil {
		f.internalError(err)
		return
	}
	if isExistence != true {
		f.customAbort(http.StatusBadRequest, "Project name invalid.")
		return
	}

	f.resolveRepoPath(projectName)
	targetFilePath := f.repoPath
	if _, err := os.Stat(targetFilePath); os.IsNotExist(err) {
		f.customAbort(http.StatusBadRequest, "image Name and  tag name are invalid.")
		return
	}
	logs.Info("User: %s download Dockerfile file from %s.", f.currentUser.Username, targetFilePath)
	f.Ctx.Output.Download(targetFilePath, dockerfileName)
}

// API to get image registry address
func (p *ImageController) GetImageRegistryAction() {
	registryAddr := registryBaseURI()
	logs.Info("The image registry is %s", registryAddr)
	p.Data["json"] = registryAddr
	p.ServeJSON()
}

// API to reset build image temp
func (p *ImageController) ResetBuildImageTempAction() {
	projectName := p.GetString("project_name")

	currentProject, err := service.GetProject(model.Project{Name: projectName}, "name")
	if err != nil {
		p.internalError(err)
		return
	}
	if currentProject == nil {
		p.customAbort(http.StatusBadRequest, "Invalid project name.")
		return
	}

	isMember, err := service.IsProjectMember(currentProject.ID, p.currentUser.ID)
	if err != nil {
		p.internalError(err)
		return
	}

	if !(p.isSysAdmin || isMember) {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to build image.")
		return
	}

	uploadedPath := filepath.Join(baseRepoPath(), p.currentUser.Username, "upload")
	err = os.RemoveAll(uploadedPath)
	if err != nil {
		logs.Error("Failed to remove uploaded path: %s", uploadedPath)
		p.internalError(err)
		return
	}
	//remove attachment file
	err = os.Remove(filepath.Join(baseRepoPath(), p.currentUser.Username, attachmentFile))
	if err != nil {
		logs.Error("Failed to remove attachment file: %+v", err)
		p.internalError(err)
	}
}
