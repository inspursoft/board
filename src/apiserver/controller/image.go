package controller

import (
	"fmt"
	"git/inspursoft/board/src/apiserver/service"
	"git/inspursoft/board/src/apiserver/service/devops/travis"
	"git/inspursoft/board/src/common/model"
	"net/http"
	"os"
	"path/filepath"

	"strings"

	"github.com/astaxie/beego/logs"
)

type ImageController struct {
	BaseController
}

// API to get image list
func (p *ImageController) GetImagesAction() {
	// Get the image list from registry v2
	query := model.Project{}
	projectList, err := service.GetProjectsByUser(query, p.currentUser.ID)
	if err != nil {
		p.internalError(err)
		return
	}
	repoList, err := service.GetRegistryCatalog()
	if err != nil {
		p.internalError(err)
		return
	}

	var repoListFiltered model.RegistryRepo
	for _, imageName := range repoList.Names {
		fromIndex := strings.LastIndex(imageName, "/")
		if fromIndex == -1 {
			continue
		}
		for _, project := range projectList {
			if imageName[:fromIndex] == project.Name {
				repoListFiltered.Names = append(repoListFiltered.Names, imageName)
				break
			}
		}
	}

	/* Interpret the message to api server */
	imageList := []model.Image{}
	for _, imageName := range repoListFiltered.Names {
		var newImage model.Image
		newImage.ImageName = imageName
		reqTagList, err := service.GetRegistryImageTags(imageName)
		if err != nil {
			p.internalError(err)
			return
		}
		if len(reqTagList.Tags) == 0 {
			logs.Debug("Image: %s has no tags.", imageName)
			continue
		}

		// Check image in DB
		dbImage, err := service.GetImageByName(imageName)
		if err != nil {
			p.customAbort(http.StatusInternalServerError, fmt.Sprintf("Checking image name in DB error: %+v", err))
			return
		}
		if dbImage != nil {
			// image already in DB, use the status in DB
			newImage.ImageID = dbImage.ImageID
			newImage.ImageComment = dbImage.ImageComment
			newImage.ImageDeleted = dbImage.ImageDeleted
		} else {
			// image not in DB, add it to DB
			imageID, err := service.CreateImage(newImage)
			if err != nil {
				p.customAbort(http.StatusInternalServerError, fmt.Sprintf("Create image to DB error: %+v", err))
				return
			}
			newImage.ImageID = imageID
		}
		if newImage.ImageDeleted == 0 {
			imageList = append(imageList, newImage)
		}
	}
	p.renderJSON(imageList)
}

// API to get tag list for a specific image
func (p *ImageController) GetImageDetailAction() {

	imageName := p.Ctx.Input.Param(":imagename")
	reqTagList, err := service.GetRegistryImageTags(imageName)
	if err != nil {
		p.internalError(err)
		return
	}

	var imageDetail []model.TagDetail
	for _, tagID := range reqTagList.Tags {
		var tagDetail model.TagDetail
		tagDetail.ImageName = reqTagList.ImageName
		tagDetail.ImageTag = tagID
		tagDetail.ImageSizeUnit = "B"
		// Get version one schema

		manifest1, err := service.GetRegistryManifest1(tagDetail.ImageName, tagDetail.ImageTag)
		if err != nil {
			p.internalError(err)
			return
		}

		tagDetail.ImageDetail = (manifest1.History[0])["v1Compatibility"]
		tagDetail.ImageAuthor = ""       //TODO: get the author by frontend simply
		tagDetail.ImageCreationTime = "" //TODO: get the time by frontend simply

		// Get version two schema
		manifest2, err := service.GetRegistryManifest2(tagDetail.ImageName, tagDetail.ImageTag)
		if err != nil {
			p.internalError(err)
			return
		}

		tagDetail.ImageId = manifest2.Config.Digest
		tagDetail.ImageSize = manifest2.Config.Size

		var layerconfig model.Manifest2Config
		for _, layerconfig = range manifest2.Layers {
			tagDetail.ImageSize += layerconfig.Size
		}
		// Add the tag detail to list
		imageDetail = append(imageDetail, tagDetail)
	}
	p.renderJSON(imageDetail)
}

func (p *ImageController) generateBuildingImageTravis(imageURI, dockerfileName string) error {
	userID := p.currentUser.ID
	var travisCommand travis.TravisCommand
	travisCommand.BeforeDeploy.Commands = []string{
		fmt.Sprintf("curl \"%s/jenkins-job/%d/$BUILD_NUMBER\"", boardAPIBaseURL(), userID),
		"if [ -f key.txt ]; then token=`cat key.txt`; fi",
		fmt.Sprintf("status=`curl -I \"%s/files/download?token=$token\" 2>/dev/null | head -n 1 | cut -d$' ' -f2`", boardAPIBaseURL()),
		fmt.Sprintf("if [ $status == '200' ]; then curl -o attachment.zip \"%s/files/download?token=$token\" && mkdir -p upload && unzip attachment.zip -d upload; fi", boardAPIBaseURL()),
	}
	travisCommand.Deploy.Commands = []string{
		"export PATH=/usr/bin:/bin:/usr/sbin:/sbin:/usr/local/bin",
		fmt.Sprintf("docker build -t %s -f containers/%s .", imageURI, dockerfileName),
		fmt.Sprintf("docker push %s", imageURI),
	}
	return travisCommand.GenerateCustomTravis(p.repoPath)
}

func (p *ImageController) BuildImageAction() {
	var reqImageConfig model.ImageConfig
	var err error
	//Check user priviledge project admin
	p.resolveBody(&reqImageConfig)
	p.resolveUserPrivilege(reqImageConfig.ProjectName)
	//Checking invalid parameters
	p.resolveRepoImagePath(reqImageConfig.ProjectName)
	reqImageConfig.RepoPath = p.repoImagePath
	err = service.CheckDockerfileConfig(&reqImageConfig)
	if err != nil {
		p.serveStatus(http.StatusBadRequest, err.Error())
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

	projectName := reqImageConfig.ProjectName
	imageName := reqImageConfig.ImageName
	imageTag := reqImageConfig.ImageTag
	imageURI := filepath.Join(registryBaseURI(), projectName, imageName) + ":" + imageTag

	if currentToken, ok := memoryCache.Get(p.currentUser.Username).(string); ok {
		service.CreateFile("key.txt", currentToken, p.repoPath)
	}

	dockerfileName := service.ResolveDockerfileName(imageName, imageTag)
	err = p.generateBuildingImageTravis(imageURI, dockerfileName)
	if err != nil {
		logs.Error("Failed to generate building image travis: %+v", err)
		return
	}

	items := []string{".travis.yml", "key.txt", filepath.Join("containers", dockerfileName)}
	p.pushItemsToRepo(items...)
	p.collaborateWithPullRequest("master", "master", items...)
}

func (p *ImageController) GetImageDockerfileAction() {
	projectName := strings.TrimSpace(p.GetString("project_name"))

	p.resolveProjectMember(projectName)
	p.resolveRepoImagePath(projectName)

	imageName := strings.TrimSpace(p.GetString("image_name"))
	imageTag := strings.TrimSpace(p.GetString("image_tag"))

	if imageName == "" || imageTag == "" {
		logs.Error("Missing image name or tag, current image name is: %s, tag is: %s", imageName, imageTag)
		p.customAbort(http.StatusBadRequest, "Missing image name or tag.")
		return
	}

	dockerfile, err := service.GetDockerfileInfo(p.repoImagePath, imageName, imageTag)
	if err != nil {
		p.internalError(err)
		return
	}
	p.renderJSON(dockerfile)
}

func (p *ImageController) DockerfilePreviewAction() {
	var reqImageConfig model.ImageConfig
	p.resolveBody(&reqImageConfig)
	p.resolveUserPrivilege(reqImageConfig.ProjectName)
	p.resolveRepoImagePath(reqImageConfig.ProjectName)
	reqImageConfig.RepoPath = p.repoImagePath
	//Checking invalid parameters
	err := service.CheckDockerfileConfig(&reqImageConfig)
	if err != nil {
		p.serveStatus(http.StatusBadRequest, err.Error())
		return
	}
	err = service.BuildDockerfile(reqImageConfig, p.Ctx.ResponseWriter)
	if err != nil {
		p.internalError(err)
	}
}

func (p *ImageController) ConfigCleanAction() {
	projectName := strings.TrimSpace(p.GetString("project_name"))
	logs.Debug("Cleaning config to the project: %s", projectName)
	p.resolveUserPrivilege(projectName)

	//remove uploaded directory
	uploadedPath := filepath.Join(baseRepoPath(), p.currentUser.Username, "upload")
	err := os.RemoveAll(uploadedPath)
	if err != nil {
		logs.Error("Failed to remove uploaded path: %s, error: %+v", uploadedPath, err)
	}

	//remove attachment file
	err = os.RemoveAll(filepath.Join(baseRepoPath(), p.currentUser.Username, attachmentFile))
	if err != nil {
		logs.Error("Failed to remove attachment file: %+v", err)
		p.internalError(err)
	}
}

func (p *ImageController) deleteImageWithTag(imageName, imageTag string) {
	var err error
	digest, err := service.GetRegistryImageDigest(imageName, imageTag)
	if err != nil {
		p.internalError(err)
		return
	}
	err = service.DeleteRegistryImageWithETag(imageName, imageTag, digest)
	if err != nil {
		p.internalError(err)
	}
}

func (p *ImageController) DeleteImageAction() {

	if p.isSysAdmin == false {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to delete image.")
		return
	}

	imageName := strings.TrimSpace(p.GetString("image_name"))
	reqTagList, err := service.GetRegistryImageTags(imageName)
	if err != nil {
		p.internalError(err)
		return
	}
	for _, tagName := range reqTagList.Tags {
		p.deleteImageWithTag(imageName, tagName)
	}
}

func (p *ImageController) DeleteImageTagAction() {
	if p.isSysAdmin == false {
		p.customAbort(http.StatusForbidden, "Insufficient privileges to delete image tag.")
		return
	}
	imageName := strings.TrimSpace(p.Ctx.Input.Param(":imagename"))
	imageTag := strings.TrimSpace(p.GetString("image_tag"))
	p.deleteImageWithTag(imageName, imageTag)
}

func (p *ImageController) resolveDockerfileName() (dockerfileName string) {
	imageName := strings.TrimSpace(p.GetString("image_name"))
	imageTag := strings.TrimSpace(p.GetString("image_tag"))

	if imageName == "" || imageTag == "" {
		logs.Error("Missing image name or tag, current image name is: %s, tag is: %s", imageName, imageTag)
		p.customAbort(http.StatusBadRequest, "Cannot generate Dockerfile due to image name or tag is missing.")
		return
	}
	dockerfileName = service.ResolveDockerfileName(imageName, imageTag)
	return
}

func (p *ImageController) DockerfileBuildImageAction() {

	projectName := strings.TrimSpace(p.GetString("project_name"))

	p.resolveUserPrivilege(projectName)
	p.resolveRepoImagePath(projectName)
	dockerfilePath := p.repoImagePath
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		p.customAbort(http.StatusNotFound, "Image path does not exist.")
		return
	}
	imageName := strings.TrimSpace(p.GetString("image_name"))
	imageTag := strings.TrimSpace(p.GetString("image_tag"))
	if imageName == "" || imageTag == "" {
		logs.Error("Missing image name or tag, current image name is: %s, tag is: %s", imageName, imageTag)
		p.customAbort(http.StatusBadRequest, "Missing image name or tag.")
		return
	}
	imageURI := filepath.Join(registryBaseURI(), projectName, imageName) + ":" + imageTag
	dockerfileName := service.ResolveDockerfileName(imageName, imageTag)
	err := p.generateBuildingImageTravis(imageURI, dockerfileName)
	if err != nil {
		logs.Error("Failed to generate building image travis: %+v", err)
		return
	}

	items := []string{".travis.yml", filepath.Join("containers", dockerfileName)}
	p.pushItemsToRepo(items...)
	p.collaborateWithPullRequest("master", "master", items...)
}

func (p *ImageController) CheckImageTagExistingAction() {
	var err error
	projectName := strings.TrimSpace(p.GetString("project_name"))
	p.resolveUserPrivilege(projectName)
	// check this image:tag in system

	imageName := strings.TrimSpace(p.Ctx.Input.Param(":imagename"))
	imageTag := strings.TrimSpace(p.GetString("image_tag"))

	// TODO check image imported from registry
	existing, err := existRegistry(projectName, imageName, imageTag)
	if err != nil {
		p.internalError(err)
		return
	}

	if existing {
		p.customAbort(http.StatusConflict, "This image:tag already existing.")
		return
	}
	logs.Debug("checking image:tag result %t", existing)
}

func existRegistry(projectName string, imageName string, imageTag string) (bool, error) {
	currentName := filepath.Join(projectName, imageName)
	//check image
	repoList, err := service.GetRegistryCatalog()
	if err != nil {
		logs.Error("Failed to unmarshal repoList body %+v", err)
		return false, err
	}
	for _, imageRegistry := range repoList.Names {
		if imageRegistry == currentName {
			//check tag
			tagList, err := service.GetRegistryImageTags(currentName)
			if err != nil {
				logs.Error("Failed to unmarshal body %+v", err)
				return false, err
			}
			for _, tagID := range tagList.Tags {
				if imageTag == tagID {
					logs.Info("Image tag existing %s:%s", currentName, tagID)
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
	f.resolveRepoImagePath(projectName)

	_, fileHeader, err := f.GetFile("upload_file")
	if err != nil {
		f.internalError(err)
		return
	}
	if fileHeader.Filename != "Dockerfile" {
		f.customAbort(http.StatusBadRequest, "Update file name invalid.")
		return
	}

	imageName := strings.TrimSpace(f.GetString("image_name"))
	imageTag := strings.TrimSpace(f.GetString("image_tag"))

	if imageName == "" || imageTag == "" {
		logs.Error("Missing image name or tag, current image name is: %s, tag is: %s", imageName, imageTag)
		f.customAbort(http.StatusBadRequest, "Missing image name or tag.")
		return
	}
	dockerfileName := service.ResolveDockerfileName(imageName, imageTag)
	if _, err := os.Stat(f.repoImagePath); os.IsNotExist(err) {
		os.MkdirAll(f.repoImagePath, 0755)
	}
	err = f.SaveToFile("upload_file", filepath.Join(f.repoImagePath, dockerfileName))
	if err != nil {
		f.internalError(err)
	}
}

func (f *ImageController) DownloadDockerfileFileAction() {
	projectName := f.GetString("project_name")
	f.resolveProjectMember(projectName)
	isExistence, err := service.ProjectExists(projectName)
	if err != nil {
		f.internalError(err)
		return
	}
	if isExistence != true {
		f.customAbort(http.StatusBadRequest, "Project name invalid.")
		return
	}

	f.resolveRepoImagePath(projectName)
	if _, err := os.Stat(f.repoImagePath); os.IsNotExist(err) {
		f.customAbort(http.StatusNotFound, "Target file path does not exist.")
		return
	}

	imageName := strings.TrimSpace(f.GetString("image_name"))
	imageTag := strings.TrimSpace(f.GetString("image_tag"))

	if imageName == "" || imageTag == "" {
		logs.Error("Missing image name or tag, current image name is: %s, tag is: %s", imageName, imageTag)
		f.customAbort(http.StatusBadRequest, "Missing image name or tag.")
		return
	}
	dockerfileName := service.ResolveDockerfileName(imageName, imageTag)
	logs.Info("User: %s download Dockerfile file under %s.", f.currentUser.Username, f.repoImagePath)
	f.Ctx.Output.Download(filepath.Join(f.repoImagePath, dockerfileName), dockerfileName)
}

// API to get image registry address
func (p *ImageController) GetImageRegistryAction() {
	registryAddr := registryBaseURI()
	logs.Info("Docker registry is %s", registryAddr)
	p.renderJSON(registryAddr)
}
