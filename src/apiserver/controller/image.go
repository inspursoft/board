package controller

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/common/model"
	"git/inspursoft/board/src/common/utils"
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

var registryURL = utils.GetConfig("REGISTRY_URL")

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
		logs.Info(string(body))
		p.internalError(err)
		return
	}

	/* Interpret the message to api server */
	var imagelist []model.Image
	for _, imagename := range repolist.Names {
		var newImage model.Image
		newImage.ImageName = imagename
		// Check image in DB
		dbimage, err := service.GetImage(newImage, "name")
		if err != nil {
			logs.Info("Checking image name in DB error")
			p.internalError(err)
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
				logs.Info("Create image in DB error")
				p.internalError(err)
				return
			}
			newImage.ImageID = id
		}

		imagelist = append(imagelist, newImage)
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
		logs.Info("url=%s", gettagsurl)
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
		logs.Info(string(body))
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
			logs.Info(getmanifesturl)
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
			logs.Info(string(body))
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
			logs.Info(string(body))
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
	logs.Info(imagedetail)
	p.Data["json"] = imagedetail
	p.ServeJSON()

}

//  Checking the user priviledge by token
func (p *ImageController) Prepare() {
	user := p.getCurrentUser()
	if user == nil {
		p.CustomAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
	p.isSysAdmin = (user.SystemAdmin == 1)
	p.isProjectAdmin = (user.ProjectAdmin == 1)
}

func (p *ImageController) BuildImageAction() {
	var err error

	//Check user priviledge project admin
	if p.isProjectAdmin == false {
		p.serveStatus(http.StatusForbidden, "Invalid user for project admin")
		return
	}

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

	reqImageConfig.ImageDockerfilePath = filepath.Join(repoPath, reqImageConfig.ProjectName,
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
	p.CustomAbort(ret, msg)
}

func (p *ImageController) GetImageDockerfileAction() {
	imageName := strings.TrimSpace(p.GetString("image_name"))
	imageTag := strings.TrimSpace(p.GetString("image_tag"))
	projectName := strings.TrimSpace(p.GetString("project_name"))

	dockerfilePath := filepath.Join(repoPath, projectName, imageName, imageTag)
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		p.CustomAbort(http.StatusNotFound, "Image path doe's not exist.")
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
	var err error

	//Check user priviledge project admin
	if p.isProjectAdmin == false {
		p.serveStatus(http.StatusForbidden, "Invalid user for project admin")
		return
	}

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

	reqImageConfig.ImageDockerfilePath = filepath.Join(repoPath, reqImageConfig.ProjectName,
		reqImageConfig.ImageName, reqImageConfig.ImageTag)
	err = service.BuildDockerfile(reqImageConfig, p.Ctx.ResponseWriter)
	if err != nil {
		p.internalError(err)
		return
	}
}

func (p *ImageController) ConfigCleanAction() {
	var err error

	if p.isProjectAdmin == false {
		p.serveStatus(http.StatusForbidden, "Invalid user for project admin")
		return
	}

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

	configPath := filepath.Join(repoPath, strings.TrimSpace(reqImageIndex.ProjectName), strings.TrimSpace(reqImageIndex.ImageName), strings.TrimSpace(reqImageIndex.ImageTag))
	err = service.ImageConfigClean(configPath)
	if err != nil {
		p.internalError(err)
		return
	}
}

func (p *ImageController) DeleteImageAction() {
	var err error

	if p.isProjectAdmin == false {
		p.serveStatus(http.StatusForbidden, "Invalid user for project admin")
		return
	}

	imageName := strings.TrimSpace(p.GetString("image_name"))

	var image model.Image
	image.ImageName = imageName

	dbImage, err := service.GetImage(image, "name")
	if err != nil {
		p.internalError(err)
		return
	}
	if dbImage == nil {
		p.serveStatus(http.StatusBadRequest, err.Error())
		return
	}

	err = service.DeleteImage(*dbImage)
	if err != nil {
		p.internalError(err)
		return
	}
}

func (p *ImageController) DeleteImageTagAction() {
	var err error

	if p.isProjectAdmin == false {
		p.serveStatus(http.StatusForbidden, "Invalid user for project admin")
		return
	}

	imageName := strings.TrimSpace(p.Ctx.Input.Param(":imagename"))
	_imageTag := strings.TrimSpace(p.GetString("image_tag"))

	var imageTag model.ImageTag
	imageTag.ImageName = imageName
	imageTag.Tag = _imageTag

	dbImageTag, err := service.GetImageTag(imageTag, "image_name", "tag")
	if err != nil {
		p.internalError(err)
		return
	}
	if dbImageTag == nil {
		p.serveStatus(http.StatusBadRequest, err.Error())
		return
	}

	err = service.DeleteImageTag(*dbImageTag)
	if err != nil {
		p.internalError(err)
		return
	}
}
