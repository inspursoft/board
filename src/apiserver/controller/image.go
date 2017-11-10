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
	commentTemp  = "Inspur image" // TODO: get from mysql in the next release
	sizeunitTemp = "B"

	defaultDockerfilename = "Dockerfile"
	imageProcess          = "process_image"
)

// API to get image list
func (p *ImageController) GetImagesAction() {

	var repolist model.RegistryRepo
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

	/* Interpret the message to api server */
	var imagelist []model.Image
	for _, imagename := range repolist.Names {
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

	reqImageConfig.ImageDockerfilePath = filepath.Join(repoPath(), reqImageConfig.ProjectName,
		reqImageConfig.ImageName, reqImageConfig.ImageTag)
	err = service.BuildDockerfile(reqImageConfig)
	if err != nil {
		p.internalError(err)
		return
	}

	//push to git
	var pushobject pushObject

	pushobject.FileName = defaultDockerfilename
	pushobject.JobName = imageProcess
	pushobject.Value = filepath.Join(reqImageConfig.ProjectName,
		reqImageConfig.ImageName, reqImageConfig.ImageTag)
	pushobject.Extras = filepath.Join(reqImageConfig.ProjectName,
		reqImageConfig.ImageName) + ":" + reqImageConfig.ImageTag
	pushobject.Message = fmt.Sprintf("Build image: %s", pushobject.Extras)

	//Get file list for Jenkis git repo
	uploads, err := service.ListUploadFiles(filepath.Join(reqImageConfig.ImageDockerfilePath, "upload"))
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

	ret, msg, err := InternalPushObjects(&pushobject, &(p.baseController))
	if err != nil {
		p.internalError(err)
		return
	}
	logs.Info("Internal push object: %d %s", ret, msg)
	p.customAbort(ret, msg)
}

func (p *ImageController) GetImageDockerfileAction() {
	imageName := strings.TrimSpace(p.GetString("image_name"))
	imageTag := strings.TrimSpace(p.GetString("image_tag"))
	projectName := strings.TrimSpace(p.GetString("project_name"))

	dockerfilePath := filepath.Join(repoPath(), projectName, imageName, imageTag)
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

	reqImageConfig.ImageDockerfilePath = filepath.Join(repoPath(), reqImageConfig.ProjectName,
		reqImageConfig.ImageName, reqImageConfig.ImageTag)
	err = service.BuildDockerfile(reqImageConfig, p.Ctx.ResponseWriter)
	if err != nil {
		p.internalError(err)
		return
	}
}

func (p *ImageController) ConfigCleanAction() {
	reqData, err := p.resolveBody()
	if err != nil {
		p.internalError(err)
		return
	}

	var reqImageIndex model.ImageIndex
	err = json.Unmarshal(reqData, &reqImageIndex)
	if err != nil {
		p.internalError(err)
		return
	}

	currentProject, err := service.GetProject(model.Project{Name: reqImageIndex.ProjectName}, "name")
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

	configPath := filepath.Join(repoPath(), strings.TrimSpace(reqImageIndex.ProjectName),
		strings.TrimSpace(reqImageIndex.ImageName), strings.TrimSpace(reqImageIndex.ImageTag))

	// Update git repo
	var pushobject pushObject

	pushobject.FileName = defaultDockerfilename
	pushobject.JobName = imageProcess
	pushobject.Value = filepath.Join(reqImageIndex.ProjectName,
		reqImageIndex.ImageName, reqImageIndex.ImageTag)
	pushobject.Extras = filepath.Join(reqImageIndex.ProjectName,
		reqImageIndex.ImageName) + ":" + reqImageIndex.ImageTag
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

	ret, msg, err := InternalCleanObjects(&pushobject, &(p.baseController))
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
