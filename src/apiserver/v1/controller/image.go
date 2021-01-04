package controller

import (
	"encoding/json"
	"fmt"
	c "github.com/inspursoft/board/src/apiserver/controllers/commons"
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"
	"github.com/inspursoft/board/src/common/utils"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"time"

	"strings"

	"github.com/astaxie/beego/logs"
)

type ImageController struct {
	c.BaseController
}

var imageBaselineTime = utils.GetConfig("IMAGE_BASELINE_TIME")

// API to get image list
func (p *ImageController) GetImagesAction() {
	// Get the image list from registry v2
	query := model.Project{}
	projectList, err := service.GetProjectsByUser(query, p.CurrentUser.ID)
	if err != nil {
		p.InternalError(err)
		return
	}
	repoList, err := service.GetRegistryCatalog()
	if err != nil {
		p.InternalError(err)
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
	imageList := model.ImageList{}
	for _, imageName := range repoListFiltered.Names {
		newImage := model.Image{
			ImageName:         imageName,
			ImageCreationTime: time.Now(),
		}

		reqTagList, err := service.GetRegistryImageTags(imageName)
		if err != nil {
			p.InternalError(err)
			return
		}
		if len(reqTagList.Tags) == 0 {
			logs.Debug("Image: %s has no tags.", imageName)
			continue
		}

		for _, imageTag := range reqTagList.Tags {
			imageManifest, err := service.GetRegistryManifest1(imageName, imageTag)
			if err != nil {
				logs.Error("Failed to get resgistry image manifest: %+v", err)
				continue
			}
			if len(imageManifest.History) > 0 {
				imageDetail := struct {
					Created time.Time `json:"created"`
				}{}
				err := json.Unmarshal([]byte((imageManifest.History[0])["v1Compatibility"]), &imageDetail)
				if err != nil {
					logs.Error("Failed to Unmarshal registry manifest: %+v", err)
					continue
				}
				if newImage.ImageCreationTime.Unix() > imageDetail.Created.Unix() {
					newImage.ImageCreationTime = imageDetail.Created
				}
				if newImage.ImageUpdateTime.Unix() < imageDetail.Created.Unix() {
					newImage.ImageUpdateTime = imageDetail.Created
				}
			}
		}
		// Check image in DB
		dbImage, err := service.GetImageByName(imageName)
		if err != nil {
			p.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprintf("Checking image name in DB error: %+v", err))
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
				p.CustomAbortAudit(http.StatusInternalServerError, fmt.Sprintf("Create image to DB error: %+v", err))
				return
			}
			newImage.ImageID = imageID
		}

		baselineTime, err := time.Parse("2006-01-02 15:04:05", imageBaselineTime())
		if err != nil {
			logs.Error("Illegal image baseline time: %s, err:%+v", imageBaselineTime(), err)
			baselineTime, _ = time.Parse("2006-01-02 15:04:05", "2017-06-06 00:00:00")
		}
		if newImage.ImageDeleted == 0 && newImage.ImageUpdateTime.Unix() > baselineTime.Unix() {
			imageList = append(imageList, &newImage)
		}
	}
	sort.Sort(imageList)
	p.RenderJSON(imageList)
}

// API to get tag list for a specific image
func (p *ImageController) GetImageDetailAction() {

	imageName := p.Ctx.Input.Param(":imagename")
	reqTagList, err := service.GetRegistryImageTags(imageName)
	if err != nil {
		p.InternalError(err)
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
			p.InternalError(err)
			return
		}
		if len(manifest1.History) > 0 {
			imageDetail := struct {
				Created time.Time `json:"created"`
				Arch    string    `json: "architecture"`
				Author  string    `json: "author"`
			}{}
			err := json.Unmarshal([]byte((manifest1.History[0])["v1Compatibility"]), &imageDetail)
			if err != nil {
				logs.Error("Failed to Unmarshal registry manifest: %+v", err)
			} else {
				tagDetail.ImageDetail = (manifest1.History[0])["v1Compatibility"]
				tagDetail.ImageArch = imageDetail.Arch
				tagDetail.ImageAuthor = imageDetail.Author
				tagDetail.ImageCreationTime = imageDetail.Created.Format("2006-01-02 15:04:05")
			}
		}

		// Get version two schema
		manifest2, err := service.GetRegistryManifest2(tagDetail.ImageName, tagDetail.ImageTag)
		if err != nil {
			p.InternalError(err)
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
	p.RenderJSON(imageDetail)
}

func (p *ImageController) generateBuildingImageTravis(imageURI, dockerfileName string) (yamlFileName string, err error) {
	configurations := make(map[string]string)
	configurations["user_id"] = strconv.Itoa(int(p.CurrentUser.ID))
	configurations["repo_name"] = p.RepoName
	configurations["repo_token"] = p.CurrentUser.RepoToken
	configurations["token"] = p.Token
	configurations["image_uri"] = imageURI
	configurations["dockerfile"] = dockerfileName
	configurations["repo_path"] = p.RepoPath
	return service.CurrentDevOps().CreateCIYAML(service.BuildDockerImageCIYAML, configurations)
}

func (p *ImageController) generatePushImagePackageTravis(imageURI, imagePackageName string) (yamlFileName string, err error) {
	configurations := make(map[string]string)
	configurations["user_id"] = strconv.Itoa(int(p.CurrentUser.ID))
	configurations["repo_name"] = p.RepoName
	configurations["repo_token"] = p.CurrentUser.RepoToken
	configurations["token"] = p.Token
	configurations["image_uri"] = imageURI
	configurations["image_package_name"] = imagePackageName
	configurations["repo_path"] = p.RepoPath
	return service.CurrentDevOps().CreateCIYAML(service.PushDockerImageCIYAML, configurations)
}

func (p *ImageController) BuildImageAction() {
	var reqImageConfig model.ImageConfig
	var err error
	//Check user priviledge project admin
	err = p.ResolveBody(&reqImageConfig)
	if err != nil {
		return
	}
	p.ResolveUserPrivilege(reqImageConfig.ProjectName)
	//Checking invalid parameters
	p.ResolveRepoImagePath(reqImageConfig.ProjectName)
	reqImageConfig.RepoPath = p.RepoImagePath

	if reqImageConfig.NodeSelection == "" {
		reqImageConfig.NodeSelection = "slave"
	}
	utils.SetConfig("NODE_SELECTION", reqImageConfig.NodeSelection)

	err = service.CheckDockerfileConfig(&reqImageConfig)
	if err != nil {
		p.ServeStatus(http.StatusBadRequest, err.Error())
		return
	}

	reqImageConfig.ImageDockerfilePath = reqImageConfig.RepoPath
	// Check image:tag existing in registry
	existing, err := existRegistry(reqImageConfig.ProjectName, reqImageConfig.ImageName,
		reqImageConfig.ImageTag)
	if err != nil {
		p.InternalError(err)
		return
	}

	if existing {
		logs.Error("This image:tag existing in registry %s", reqImageConfig.ImageDockerfilePath)
		p.CustomAbortAudit(http.StatusConflict, "This image:tag already existing.")
		return
	}

	err = service.BuildDockerfile(reqImageConfig)
	if err != nil {
		p.InternalError(err)
		return
	}

	projectName := reqImageConfig.ProjectName
	imageName := reqImageConfig.ImageName
	imageTag := reqImageConfig.ImageTag
	imageURI := filepath.Join(c.RegistryBaseURI(), projectName, imageName) + ":" + imageTag
	dockerfileName := service.ResolveDockerfileName(imageName, imageTag)
	yamlFileName, err := p.generateBuildingImageTravis(imageURI, dockerfileName)
	if err != nil {
		logs.Error("Failed to generate building image travis: %+v", err)
		return
	}
	p.MergeCollaborativePullRequest()
	items := []string{yamlFileName, filepath.Join("containers", dockerfileName)}
	p.PushItemsToRepo(items...)
	p.CollaborateWithPullRequest("master", "master", items...)
}

func (p *ImageController) GetImageDockerfileAction() {
	projectName := strings.TrimSpace(p.GetString("project_name"))

	p.ResolveProjectMember(projectName)
	p.ResolveRepoImagePath(projectName)

	imageName := strings.TrimSpace(p.GetString("image_name"))
	imageTag := strings.TrimSpace(p.GetString("image_tag"))

	if imageName == "" || imageTag == "" {
		logs.Error("Missing image name or tag, current image name is: %s, tag is: %s", imageName, imageTag)
		p.CustomAbortAudit(http.StatusBadRequest, "Missing image name or tag.")
		return
	}

	dockerfile, err := service.GetDockerfileInfo(p.RepoImagePath, imageName, imageTag)
	if err != nil {
		p.CustomAbortAudit(http.StatusNotFound, err.Error())
		return
	}
	p.RenderJSON(dockerfile)
}

func (p *ImageController) DockerfilePreviewAction() {
	var reqImageConfig model.ImageConfig
	err := p.ResolveBody(&reqImageConfig)
	if err != nil {
		return
	}
	p.ResolveUserPrivilege(reqImageConfig.ProjectName)
	p.ResolveRepoImagePath(reqImageConfig.ProjectName)
	reqImageConfig.RepoPath = p.RepoImagePath
	//Checking invalid parameters
	err = service.CheckDockerfileConfig(&reqImageConfig)
	if err != nil {
		p.ServeStatus(http.StatusBadRequest, err.Error())
		return
	}
	err = service.BuildDockerfile(reqImageConfig, p.Ctx.ResponseWriter)
	if err != nil {
		p.InternalError(err)
	}
}

func (p *ImageController) ConfigCleanAction() {
	projectName := strings.TrimSpace(p.GetString("project_name"))
	logs.Debug("Cleaning config to the project: %s", projectName)
	p.ResolveUserPrivilege(projectName)

	//remove uploaded directory
	uploadedPath := filepath.Join(c.BaseRepoPath(), p.CurrentUser.Username, "upload")
	err := os.RemoveAll(uploadedPath)
	if err != nil {
		logs.Error("Failed to remove uploaded path: %s, error: %+v", uploadedPath, err)
	}

	//remove attachment file
	err = os.RemoveAll(filepath.Join(c.BaseRepoPath(), p.CurrentUser.Username, attachmentFile))
	if err != nil {
		logs.Error("Failed to remove attachment file: %+v", err)
		p.InternalError(err)
	}
}

func (p *ImageController) deleteImageWithTag(imageName, imageTag string) {
	var err error
	digest, err := service.GetRegistryImageDigest(imageName, imageTag)
	if err != nil {
		p.InternalError(err)
		return
	}
	err = service.DeleteRegistryImageWithETag(imageName, imageTag, digest)
	if err != nil {
		p.InternalError(err)
	}
}

func (p *ImageController) DeleteImageAction() {

	if p.IsSysAdmin == false {
		p.CustomAbortAudit(http.StatusForbidden, "Insufficient privileges to delete image.")
		return
	}

	imageName := strings.TrimSpace(p.GetString("image_name"))
	reqTagList, err := service.GetRegistryImageTags(imageName)
	if err != nil {
		p.InternalError(err)
		return
	}
	for _, tagName := range reqTagList.Tags {
		p.deleteImageWithTag(imageName, tagName)
	}
}

func (p *ImageController) DeleteImageTagAction() {
	if p.IsSysAdmin == false {
		p.CustomAbortAudit(http.StatusForbidden, "Insufficient privileges to delete image tag.")
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
		p.CustomAbortAudit(http.StatusBadRequest, "Cannot generate Dockerfile due to image name or tag is missing.")
		return
	}
	dockerfileName = service.ResolveDockerfileName(imageName, imageTag)
	return
}

func (p *ImageController) DockerfileBuildImageAction() {
	projectName := strings.TrimSpace(p.GetString("project_name"))
	p.ResolveUserPrivilege(projectName)
	p.ResolveRepoImagePath(projectName)
	dockerfilePath := p.RepoImagePath
	if _, err := os.Stat(dockerfilePath); os.IsNotExist(err) {
		p.CustomAbortAudit(http.StatusNotFound, "Image path does not exist.")
		return
	}
	imageName := strings.TrimSpace(p.GetString("image_name"))
	imageTag := strings.TrimSpace(p.GetString("image_tag"))
	if imageName == "" || imageTag == "" {
		logs.Error("Missing image name or tag, current image name is: %s, tag is: %s", imageName, imageTag)
		p.CustomAbortAudit(http.StatusBadRequest, "Missing image name or tag.")
		return
	}
	imageURI := filepath.Join(c.RegistryBaseURI(), projectName, imageName) + ":" + imageTag
	dockerfileName := service.ResolveDockerfileName(imageName, imageTag)
	yamlFileName, err := p.generateBuildingImageTravis(imageURI, dockerfileName)
	if err != nil {
		logs.Error("Failed to generate building image travis: %+v", err)
		return
	}
	p.MergeCollaborativePullRequest()
	items := []string{yamlFileName, filepath.Join("containers", dockerfileName)}
	p.PushItemsToRepo(items...)
	p.CollaborateWithPullRequest("master", "master", items...)

}

func (p *ImageController) UploadAndPushImagePackageAction() {
	projectName := strings.TrimSpace(p.GetString("project_name"))
	p.ResolveUserPrivilege(projectName)
	p.ResolveRepoImagePath(projectName)
	err := utils.CheckFilePath(p.RepoImagePath)
	if err != nil {
		logs.Error("Failed to create directory to store image building items.")
		p.CustomAbort(http.StatusInternalServerError, "Failed to create directory to store image building items.")
		return
	}
	imageName := strings.TrimSpace(p.GetString("image_name"))
	imageTag := strings.TrimSpace(p.GetString("image_tag"))
	if imageName == "" || imageTag == "" {
		logs.Error("Missing image name or tag, current image name is: %s, tag is: %s", imageName, imageTag)
		p.CustomAbortAudit(http.StatusBadRequest, "Missing image name or tag.")
		return
	}
	imagePackageName := strings.TrimSpace(p.GetString("image_package_name"))
	imageURI := filepath.Join(c.RegistryBaseURI(), projectName, imageName) + ":" + imageTag
	yamlFileName, err := p.generatePushImagePackageTravis(imageURI, imagePackageName)
	if err != nil {
		logs.Error("Failed to generate building image travis: %+v", err)
		return
	}
	p.MergeCollaborativePullRequest()
	p.PushItemsToRepo(yamlFileName)
	p.CollaborateWithPullRequest("master", "master", yamlFileName)
}

func (p *ImageController) CheckImageTagExistingAction() {
	var err error
	projectName := strings.TrimSpace(p.GetString("project_name"))
	p.ResolveUserPrivilege(projectName)
	// check this image:tag in system

	imageName := strings.TrimSpace(p.Ctx.Input.Param(":imagename"))
	imageTag := strings.TrimSpace(p.GetString("image_tag"))

	// TODO check image imported from registry
	existing, err := existRegistry(projectName, imageName, imageTag)
	if err != nil {
		p.InternalError(err)
		return
	}

	if existing {
		p.CustomAbortAudit(http.StatusConflict, "This image:tag already existing.")
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
		f.InternalError(err)
		return
	}
	if isExistence != true {
		f.CustomAbortAudit(http.StatusBadRequest, "Project don't exist.")
		return
	}
	f.ResolveRepoImagePath(projectName)

	_, fileHeader, err := f.GetFile("upload_file")
	if err != nil {
		f.InternalError(err)
		return
	}
	if fileHeader.Filename != "Dockerfile" {
		f.CustomAbortAudit(http.StatusBadRequest, "Update file name invalid.")
		return
	}

	imageName := strings.TrimSpace(f.GetString("image_name"))
	imageTag := strings.TrimSpace(f.GetString("image_tag"))

	if imageName == "" || imageTag == "" {
		logs.Error("Missing image name or tag, current image name is: %s, tag is: %s", imageName, imageTag)
		f.CustomAbortAudit(http.StatusBadRequest, "Missing image name or tag.")
		return
	}
	dockerfileName := service.ResolveDockerfileName(imageName, imageTag)
	if _, err := os.Stat(f.RepoImagePath); os.IsNotExist(err) {
		os.MkdirAll(f.RepoImagePath, 0755)
	}
	err = f.SaveToFile("upload_file", filepath.Join(f.RepoImagePath, dockerfileName))
	if err != nil {
		f.InternalError(err)
	}
	dockerfileInfo, err := service.UpdateDockerfileCopyCommand(f.RepoImagePath, dockerfileName)
	if err != nil {
		logs.Error("Update dockerfile err: %s", err.Error())
		f.CustomAbortAudit(http.StatusBadRequest, err.Error())
		return
	}
	f.Ctx.WriteString(string(dockerfileInfo))
}

func (f *ImageController) DownloadDockerfileFileAction() {
	projectName := f.GetString("project_name")
	f.ResolveProjectMember(projectName)
	isExistence, err := service.ProjectExists(projectName)
	if err != nil {
		f.InternalError(err)
		return
	}
	if isExistence != true {
		f.CustomAbortAudit(http.StatusBadRequest, "Project name invalid.")
		return
	}

	f.ResolveRepoImagePath(projectName)
	if _, err := os.Stat(f.RepoImagePath); os.IsNotExist(err) {
		f.CustomAbortAudit(http.StatusNotFound, "Target file path does not exist.")
		return
	}

	imageName := strings.TrimSpace(f.GetString("image_name"))
	imageTag := strings.TrimSpace(f.GetString("image_tag"))

	if imageName == "" || imageTag == "" {
		logs.Error("Missing image name or tag, current image name is: %s, tag is: %s", imageName, imageTag)
		f.CustomAbortAudit(http.StatusBadRequest, "Missing image name or tag.")
		return
	}
	dockerfileName := service.ResolveDockerfileName(imageName, imageTag)
	logs.Info("User: %s download Dockerfile file under %s.", f.CurrentUser.Username, f.RepoImagePath)
	f.Ctx.Output.Download(filepath.Join(f.RepoImagePath, dockerfileName), dockerfileName)
}

// API to get image registry address
func (p *ImageController) GetImageRegistryAction() {
	registryAddr := c.RegistryBaseURI()
	logs.Info("Docker registry is %s", registryAddr)
	p.RenderJSON(registryAddr)
}

// Check an image used by services for deleting
// TODO
func (p *ImageController) GetImageUsedAction() {
	if p.IsSysAdmin == false {
		p.CustomAbortAudit(http.StatusForbidden, "Insufficient privileges to delete image tag.")
		return
	}
	imageName := strings.TrimSpace(p.Ctx.Input.Param(":imagename"))
	//imageTag := strings.TrimSpace(p.GetString("image_tag"))
	logs.Info("Image name is %s", imageName)
	return
}
