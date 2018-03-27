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
	var imagelist []model.Image
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

func (p *ImageController) generateRepoPathByProject(project *model.Project) string {
	if project == nil {
		p.customAbort(http.StatusBadRequest, "Failed to generate repo path since project is nil.")
	}
	return filepath.Join(baseRepoPath(), p.currentUser.Username, project.Name)
}

func (p *ImageController) generateRepoPathByProjectName(projectName string) string {
	return filepath.Join(baseRepoPath(), p.currentUser.Username, projectName)
}

func (p *ImageController) BuildImageAction() {
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
	reqImageConfig.RepoPath = p.generateRepoPathByProjectName(reqImageConfig.ProjectName)
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

	repoPath := p.generateRepoPathByProjectName(reqImageConfig.ProjectName)
	reqImageConfig.ImageDockerfilePath = filepath.Join(repoPath, imageProcess, reqImageConfig.ImageName, reqImageConfig.ImageTag)

	// Check image:tag path existing for rebuild
	//existing, err := exists(reqImageConfig.ImageDockerfilePath)
	//if err != nil {
	//	p.internalError(err)
	//	return
	//}
	//
	//if existing {
	//	logs.Error("This image:tag existing in system %s", reqImageConfig.ImageDockerfilePath)
	//	p.customAbort(http.StatusConflict, "This image:tag already existing.")
	//	return
	//}

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

	//move the upload directory
	tempPath := filepath.Join(repoPath, imageProcess, wrapStringWithSymbol(p.currentUser.Username))
	tempUploadPath := filepath.Join(tempPath, "upload")
	if _, err = os.Stat(tempUploadPath); err == nil {
		dstUploadPath := filepath.Join(repoPath, imageProcess, reqImageConfig.ImageName,
			reqImageConfig.ImageTag, "upload")
		err = os.Rename(tempUploadPath, dstUploadPath)
		if err != nil {
			logs.Error("Failed to move from %s to %s", tempUploadPath, dstUploadPath)
			p.internalError(err)
			return
		}
		err = os.RemoveAll(tempPath)
		if err != nil {
			logs.Error("Failed to remove temp path: %s", tempPath)
			p.internalError(err)
			return
		}
	}

	//push to git
	var pushobject pushObject
	pushobject.UserID = p.currentUser.ID
	pushobject.FileName = defaultDockerfilename
	pushobject.JobName = imageProcess
	pushobject.Value = filepath.Join(imageProcess, reqImageConfig.ImageName, reqImageConfig.ImageTag)
	pushobject.ProjectName = reqImageConfig.ProjectName

	pushobject.Extras = filepath.Join(reqImageConfig.ProjectName,
		reqImageConfig.ImageName) + ":" + reqImageConfig.ImageTag
	pushobject.Message = fmt.Sprintf("Build image: %s", pushobject.Extras)

	//Get file list for Jenkis git repo
	uploads, err := service.ListUploadFiles(filepath.Join(reqImageConfig.ImageDockerfilePath, "upload"))
	if err != nil {
		p.internalError(err)
		return
	}
	generateMetaConfiguration(&pushobject, repoPath)
	pushobject.Items = append(pushobject.Items, "META.cfg")
	// Add upload files
	for _, finfo := range uploads {
		filefullname := filepath.Join(pushobject.Value, "upload", finfo.FileName)
		pushobject.Items = append(pushobject.Items, filefullname)
	}
	// Add Dockerfile
	pushobject.Items = append(pushobject.Items, filepath.Join(pushobject.Value,
		defaultDockerfilename))

	ret, msg, err := InternalPushObjects(&pushobject, &(p.baseController))
	if err != nil {
		p.internalError(err)
		return
	}
	logs.Info("Internal push object: %d %s", ret, msg)
	p.ServeJSON()
}

func (p *ImageController) GetImageDockerfileAction() {
	imageName := strings.TrimSpace(p.GetString("image_name"))
	imageTag := strings.TrimSpace(p.GetString("image_tag"))
	projectName := strings.TrimSpace(p.GetString("project_name"))

	repoPath := p.generateRepoPathByProjectName(projectName)
	dockerfilePath := filepath.Join(repoPath, imageProcess, imageName, imageTag)
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
	reqImageConfig.RepoPath = p.generateRepoPathByProjectName(reqImageConfig.ProjectName)
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

	repoPath := p.generateRepoPathByProject(currentProject)
	reqImageConfig.ImageDockerfilePath = filepath.Join(repoPath, imageProcess, reqImageConfig.ImageName, reqImageConfig.ImageTag)
	err = service.BuildDockerfile(reqImageConfig, p.Ctx.ResponseWriter)
	if err != nil {
		p.internalError(err)
		return
	}
}

func cleanGitImageTag(username, imageName, imageTag, projectName string, p *ImageController) error {

	repoPath := p.generateRepoPathByProjectName(projectName)
	configPath := filepath.Join(repoPath, imageProcess, imageName, imageTag)

	// Update git repo
	var pushobject pushObject

	pushobject.FileName = defaultDockerfilename
	pushobject.JobName = imageProcess
	pushobject.Value = filepath.Join(imageProcess, imageName, imageTag)
	pushobject.ProjectName = projectName

	pushobject.Extras = filepath.Join(projectName, imageName) + ":" + imageTag
	pushobject.Message = fmt.Sprintf("Build image: %s", pushobject.Extras)

	//Get file list for Jenkis git repo
	uploads, err := service.ListUploadFiles(filepath.Join(configPath, "upload"))
	if err != nil {
		logs.Error("Failed to list upload files")
		return err
	}
	// Add upload files
	for _, finfo := range uploads {
		filefullname := filepath.Join(pushobject.Value, "upload", finfo.FileName)
		pushobject.Items = append(pushobject.Items, filefullname)
	}
	// Add Dockerfile
	pushobject.Items = append(pushobject.Items, filepath.Join(pushobject.Value,
		defaultDockerfilename))

	ret, msg, err := InternalPushObjects(&pushobject, &(p.baseController), toBeRemoved)
	if err != nil {
		logs.Error("Failed to push object for git repo clean", msg, ret)
		return err
	}
	logs.Info("Internal push object for git repo clean: %s", msg)
	return err
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

	repoPath := p.generateRepoPathByProject(currentProject)

	//remove upload temp directory
	tempPath := filepath.Join(repoPath, imageProcess, wrapStringWithSymbol(p.currentUser.Username))
	err = os.RemoveAll(tempPath)
	if err != nil {
		logs.Error("Failed to remove temp path: %s", tempPath)
		p.internalError(err)
		return
	}

	configPath := filepath.Join(repoPath, imageProcess, strings.TrimSpace(imageName), strings.TrimSpace(imageTag))

	// Update git repo
	var pushobject pushObject

	pushobject.FileName = defaultDockerfilename
	pushobject.JobName = imageProcess
	pushobject.Value = filepath.Join(imageProcess, imageName, imageTag)
	pushobject.ProjectName = currentProject.Name

	pushobject.Extras = filepath.Join(currentProject.Name, imageName) + ":" + imageTag
	pushobject.Message = fmt.Sprintf("Build image: %s", pushobject.Extras)

	//Get file list for Jenkis git repo
	uploads, err := service.ListUploadFiles(filepath.Join(configPath, "upload"))
	if err != nil {
		p.internalError(err)
		return
	}
	// Add upload files
	for _, finfo := range uploads {
		filefullname := filepath.Join(pushobject.Value, "upload", finfo.FileName)
		pushobject.Items = append(pushobject.Items, filefullname)
	}
	// Add Dockerfile
	pushobject.Items = append(pushobject.Items, filepath.Join(pushobject.Value,
		defaultDockerfilename))

	ret, msg, err := InternalPushObjects(&pushobject, &(p.baseController), toBeRemoved)
	if err != nil {
		logs.Info("Failed to push object for git repo clean", msg, ret)
		p.internalError(err)
		return
	}
	logs.Info("Internal push object for git repo clean: %s", msg)

	//Delete the config files
	err = service.ImageConfigClean(configPath)
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

		// Clean image tag path in git
		username := p.currentUser.Username
		projectName := imageName[:strings.Index(imageName, "/")]
		realName := imageName[strings.Index(imageName, "/")+1:]
		err = cleanGitImageTag(username, realName, tag, projectName, p)
		if err != nil {
			logs.Error("failed to clean image tag git %s:%s %s", realName, tag, projectName)
			p.internalError(err)
			return
		}

		//Delete the config files
		repoPath := p.generateRepoPathByProjectName(projectName)
		configPath := filepath.Join(repoPath, imageProcess, realName, tag)
		err = service.ImageConfigClean(configPath)
		if err != nil {
			logs.Error("failed to delete config files %s", configPath)
			p.internalError(err)
			return
		}
	}

	//	var image model.Image
	//	image.ImageName = imageName
	//
	//	dbImage, err := service.GetImage(image, "name")
	//	if err != nil {
	//		p.internalError(err)
	//		return
	//	}
	//	if dbImage == nil {
	//		p.serveStatus(http.StatusNotFound, "Image name not found")
	//		return
	//	}
	//
	//	err = service.DeleteImage(*dbImage)
	//	if err != nil {
	//		p.internalError(err)
	//		return
	//	}
}

func (p *ImageController) DeleteImageTagAction() {
	var err error

	if p.isSysAdmin == false {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to delete image tag.")
		return
	}

	imageName := strings.TrimSpace(p.Ctx.Input.Param(":imagename"))
	_imageTag := strings.TrimSpace(p.GetString("image_tag"))

	projectName := imageName[:strings.Index(imageName, "/")]
	realName := imageName[strings.Index(imageName, "/")+1:]

	var client = &http.Client{}
	URLPrefix := registryURL() + `/v2/` + imageName + `/manifests/`
	manifestsURL := URLPrefix + _imageTag
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

	// Clean image tag path in git
	username := p.currentUser.Username
	err = cleanGitImageTag(username, realName, _imageTag, projectName, p)
	if err != nil {
		logs.Error("failed to clean image tag git %s:%s %s", realName, _imageTag, projectName)
		p.internalError(err)
		return
	}

	//Delete the config files
	repoPath := p.generateRepoPathByProjectName(projectName)
	configPath := filepath.Join(repoPath, imageProcess, realName, _imageTag)
	err = service.ImageConfigClean(configPath)
	if err != nil {
		logs.Error("failed to delete config files %s", configPath)
		p.internalError(err)
		return
	}

	//	var imageTag model.ImageTag
	//	imageTag.ImageName = imageName
	//	imageTag.Tag = _imageTag
	//
	//	dbImageTag, err := service.GetImageTag(imageTag, "image_name", "tag")
	//	if err != nil {
	//		p.internalError(err)
	//		return
	//	}
	//	if dbImageTag == nil {
	//		p.serveStatus(http.StatusNotFound, "Image name or tag not found")
	//		return
	//	}
	//
	//	err = service.DeleteImageTag(*dbImageTag)
	//	if err != nil {
	//		p.internalError(err)
	//		return
	//	}
}

func (p *ImageController) DockerfileBuildImageAction() {
	imageName := strings.TrimSpace(p.GetString("image_name"))
	imageTag := strings.TrimSpace(p.GetString("image_tag"))
	projectName := strings.TrimSpace(p.GetString("project_name"))

	repoPath := p.generateRepoPathByProjectName(projectName)
	dockerfilePath := filepath.Join(repoPath, imageProcess, imageName, imageTag)
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

	// TODO check the dockerfile content in service.dockerfilecheck

	//push to git
	var pushobject pushObject

	pushobject.FileName = defaultDockerfilename
	pushobject.UserID = p.currentUser.ID
	pushobject.JobName = imageProcess
	pushobject.Value = filepath.Join(imageProcess, imageName, imageTag)
	pushobject.ProjectName = currentProject.Name

	pushobject.Extras = filepath.Join(projectName, imageName) + ":" + imageTag
	pushobject.Message = fmt.Sprintf("Build image: %s", pushobject.Extras)

	//Get file list for Jenkis git repo
	uploads, err := service.ListUploadFiles(filepath.Join(dockerfilePath, "upload"))
	if err != nil {
		p.internalError(err)
		return
	}
	generateMetaConfiguration(&pushobject, repoPath)
	pushobject.Items = append(pushobject.Items, "META.cfg")
	// Add upload files
	for _, finfo := range uploads {
		filefullname := filepath.Join(pushobject.Value, "upload", finfo.FileName)
		pushobject.Items = append(pushobject.Items, filefullname)
	}
	// Add Dockerfile
	pushobject.Items = append(pushobject.Items, filepath.Join(pushobject.Value,
		defaultDockerfilename))

	ret, msg, err := InternalPushObjects(&pushobject, &(p.baseController))
	if err != nil {
		p.internalError(err)
		return
	}
	logs.Info("Internal push object: %d %s", ret, msg)
	p.customAbort(ret, msg)
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
	repoPath := p.generateRepoPathByProjectName(projectName)
	dockerfilePath := filepath.Join(repoPath, imageProcess, imageName, imageTag)
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

	imageName := f.GetString("image_name")
	tagName := f.GetString("tag_name")

	repoPath := f.generateRepoPathByProjectName(projectName)
	logs.Debug("Repo path: %s", repoPath)
	targetFilePath := filepath.Join(repoPath, imageProcess, imageName, tagName)
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

	imageName := f.GetString("image_name")
	tagName := f.GetString("tag_name")

	repoPath := f.generateRepoPathByProjectName(projectName)
	targetFilePath := filepath.Join(repoPath, imageProcess, imageName, tagName)
	if _, err := os.Stat(targetFilePath); os.IsNotExist(err) {
		f.customAbort(http.StatusBadRequest, "image Name and  tag name are invalid.")
		return
	}

	absFileName := filepath.Join(repoPath, imageProcess, imageName, tagName, dockerfileName)
	logs.Info("User: %s download Dockerfile file from %s.", f.currentUser.Username, absFileName)

	f.Ctx.Output.Download(absFileName, dockerfileName)
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

	repoPath := p.generateRepoPathByProject(currentProject)
	tempPath := filepath.Join(repoPath, imageProcess, wrapStringWithSymbol(p.currentUser.Username))
	err = os.RemoveAll(tempPath)
	if err != nil {
		logs.Error("Failed to remove temp path: %s", tempPath)
		p.internalError(err)
		return
	}
}
