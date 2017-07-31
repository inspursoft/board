package controller

import (
	"encoding/json"
	"fmt"
	"git/inspursoft/board/src/common/model"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"time"
)

type ImageController struct {
	baseController
}

var RegistryIp string = "http://10.110.13.58:5000" // TODO: get from cfg env
var RegistryStatus bool
var commentTemp = "Inspur image" // TODO: get from mysql in the next release
var sizeunitTemp = "B"

func init() {
	_, err := http.Get(RegistryIp + "/v2/")
	if err != nil {
		RegistryStatus = false
	} else {
		RegistryStatus = true
	}
	log.Printf("%s\t%s\t%s\t", "RegistryStatus status is ", strconv.FormatBool(RegistryStatus), time.Now())
}

// API to get image list
func (p *ImageController) GetImagesAction() {

	var repolist model.RegistryRepo

	// Get the image list from registry v2
	httpresp, err := http.Get(RegistryIp + "/v2/_catalog")
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
		fmt.Println(body)
		p.internalError(err)
		return
	}

	// fmt.Println(repolist)
	/* Interpret the message to api server */
	var imagelist []model.BoardImage
	for _, imagename := range repolist.Names {
		var newImage model.BoardImage
		newImage.ImageName = imagename
		newImage.ImageComment = commentTemp
		//fmt.Println(newImage)
		imagelist = append(imagelist, newImage)
	}
	fmt.Println(imagelist)
	p.Data["json"] = imagelist
	p.ServeJSON()
}

// API to get tag list for a specific image
func (p *ImageController) GetImageDetailAction() {

	var taglist model.RegistryTags

	imageName := p.Ctx.Input.Param(":imagename")

	gettagsurl := "/v2/" + imageName + "/tags/list"

	httpresp, err := http.Get(RegistryIp + gettagsurl)
	if err != nil {
		fmt.Println("url=", gettagsurl)
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
		fmt.Println(string(body))
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
		getmenifesturl := "/v2/" + taglist.ImageName + "/manifests/" + tagid
		httpresp, err = http.Get(RegistryIp + getmenifesturl)
		if err != nil {
			fmt.Println(getmenifesturl)
			p.internalError(err)
			return
		}

		body, err = ioutil.ReadAll(httpresp.Body)
		if err != nil {
			p.internalError(err)
			return
		}

		var menifest1 model.RegistryMenifest1
		err = json.Unmarshal(body, &menifest1)
		if err != nil {
			fmt.Println(string(body))
			p.internalError(err)
			return
		}

		//fmt.Println((menifest1.History[0])["v1Compatibility"])

		// Interpret it on the frontend
		tagdetail.ImageAuthor = (menifest1.History[0])["v1Compatibility"]
		tagdetail.ImageCreationTime = (menifest1.History[0])["v1Compatibility"]

		// Get version two schema
		getmenifesturl = RegistryIp + getmenifesturl
		req, _ := http.NewRequest("GET", getmenifesturl, nil)
		req.Header.Set("Accept", "application/vnd.docker.distribution.manifest.v2+json")
		client := http.Client{}
		httpresp, err = client.Do(req)

		body, err = ioutil.ReadAll(httpresp.Body)
		if err != nil {
			p.internalError(err)
			return
		}

		var menifest2 model.RegistryMenifest2
		err = json.Unmarshal(body, &menifest2)
		if err != nil {
			fmt.Println(string(body))
			p.internalError(err)
			return
		}

		tagdetail.ImageId = menifest2.Config.Digest
		tagdetail.ImageSize = menifest2.Config.Size

		var layerconfig model.Menifest2Config
		for _, layerconfig = range menifest2.Layers {
			tagdetail.ImageSize += layerconfig.Size
		}

		// Add the tag detail to list
		imagedetail = append(imagedetail, tagdetail)

	}
	fmt.Println(imagedetail)
	p.Data["json"] = imagedetail
	p.ServeJSON()

}

/*  Checking the user priviledge
func (p *ProjectController) Prepare() {
	user, err := p.getCurrentUser()
	if err != nil {
		p.internalError(err)
		return
	}
	if user == nil {
		p.CustomAbort(http.StatusUnauthorized, "Need to login first.")
		return
	}
	p.currentUser = user
	p.isSysAdmin = (user.SystemAdmin == 1)
	p.isProjectAdmin = (user.ProjectAdmin == 1)
}
*/
