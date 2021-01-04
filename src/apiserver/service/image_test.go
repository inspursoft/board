package service_test

import (
	"github.com/inspursoft/board/src/apiserver/service"
	"github.com/inspursoft/board/src/common/model"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"github.com/stretchr/testify/assert"
)

var Copy = model.CopyStruct{
	CopyFrom: "from",
	CopyTo:   "to",
}

var Env = model.EnvStruct{
	EnvName:  "key",
	EnvValue: "value",
}

var Dockerfile model.Dockerfile
var imageConfig model.ImageConfig

func TestCheckDockerfileConfig(t *testing.T) {
	err := service.CheckDockerfileConfig(&imageConfig)
	if err == nil {
		t.Errorf("Check nil dockerfile config should be failed")
	} else {
		t.Log("Check nil config fail as expect.")
	}

	imageConfig.ImageDockerfile.Base = "Test:1.0"
	err = service.CheckDockerfileConfig(&imageConfig)
	if err == nil {
		t.Errorf("Check dockerfile with upper charactor in base should be failed")
	} else {
		t.Log("Check dockerfile with upper charactor in base successfully.")
	}

	imageConfig.ImageDockerfile.Base = "test:1.0"
	imageConfig.ImageDockerfile.EntryPoint = "a\nb"
	err = service.CheckDockerfileConfig(&imageConfig)
	if err == nil {
		t.Errorf("Check dockerfile with enter in Entrypoint should be failed")
	} else {
		t.Log("Check dockerfile with enter in Entrypoint successfully.")
	}

	imageConfig.ImageDockerfile.EntryPoint = "ab"
	imageConfig.ImageDockerfile.ExposePort = append(imageConfig.ImageDockerfile.ExposePort, "0s")
	err = service.CheckDockerfileConfig(&imageConfig)
	if err == nil {
		t.Errorf("Check dockerfile port should be failed")
	} else {
		t.Log("Check dockerfile port successfully.")
	}

	imageConfig.ImageDockerfile.ExposePort = nil
	imageConfig.ImageDockerfile.ExposePort = append(imageConfig.ImageDockerfile.ExposePort, "8888")
	imageConfig.ImageDockerfile.Volume = append(imageConfig.ImageDockerfile.Volume, "volume")
	//imageConfig.ImageDockerfile.Copy = append(imageConfig.ImageDockerfile.Copy, Copy)
	imageConfig.ImageDockerfile.RUN = append(imageConfig.ImageDockerfile.RUN, "run")
	imageConfig.ImageDockerfile.EnvList = append(imageConfig.ImageDockerfile.EnvList, Env)
	err = service.CheckDockerfileConfig(&imageConfig)
	if err != nil {
		t.Errorf("Check dockerfile error: %+v", err)
	} else {
		t.Log("Check dockerfile successfully.")
	}
}

func TestBuildDockerfile(t *testing.T) {
	imageConfig.ImageDockerfile.Base = "test:1.0"
	imageConfig.ImageDockerfile.Copy = append(imageConfig.ImageDockerfile.Copy, Copy)
	imageConfig.ImageDockerfilePath = "path"
	err := service.BuildDockerfile(imageConfig)
	if err != nil {
		t.Errorf("Build dockerfile fail: %+v", err)
	} else {
		t.Log("Build dockerfile successfully.")
	}
}

func TestGetDockerfileInfo(t *testing.T) {
	dockerfile, err := service.GetDockerfileInfo("path", imageConfig.ImageName, imageConfig.ImageTag)
	if err != nil {
		t.Errorf("Get dockerfile info error: %+v", err)
	}
	if dockerfile.Base == imageConfig.ImageDockerfile.Base &&
		dockerfile.EntryPoint == imageConfig.ImageDockerfile.EntryPoint {
		t.Log("Get dockerfile info successfully.")
	}
}

func TestImageConfigClean(t *testing.T) {
	err := service.ImageConfigClean("path")
	if err != nil {
		t.Errorf("Clean config error: %+v", err)
	} else {
		t.Log("Clean config successfully.")
	}
}

func cleanImageTestByID(imageid int64) {
	o := orm.NewOrm()
	rs := o.Raw("delete from image where id = ?", imageid)
	r, err := rs.Exec()
	if err != nil {
		logs.Error("Error occurred while deleting image: %+v", err)
	}
	affected, err := r.RowsAffected()
	if err != nil {
		logs.Error("Error occurred while deleting image: %+v", err)
	}
	if affected == 0 {
		logs.Error("Failed to delete image", imageid)
	} else {
		logs.Info("Successful cleared up.", imageid)
	}
}

func cleanImageTestByName(imagename string) {
	o := orm.NewOrm()
	rs := o.Raw("delete from image where name = ?", imagename)
	r, err := rs.Exec()
	if err != nil {
		logs.Error("Error occurred while deleting image: %+v", err)
	}
	affected, err := r.RowsAffected()
	if err != nil {
		logs.Error("Error occurred while deleting image: %+v", err)
	}
	if affected == 0 {
		logs.Error("Failed to delete image", imagename)
	} else {
		logs.Info("Successful cleared up.", imagename)
	}
}

var testimage = model.Image{
	ImageName:    "testimage1",
	ImageComment: "testimage1",
	ImageDeleted: 0,
}

var testImageid int64

func TestCreateImage(t *testing.T) {
	assert := assert.New(t)
	id, err := service.CreateImage(testimage)
	assert.Nil(err, "Failed, err when create test image.")
	assert.NotEqual(0, id, "Failed to assign a image id")
	testImageid = id
	t.Log(testImageid)
}

func TestUpdateImage(t *testing.T) {
	assert := assert.New(t)
	testimage.ImageDeleted = 1
	testimage.ImageID = testImageid
	ret, err := service.UpdateImage(testimage, "deleted")
	assert.Nil(err, "Failed, err when update test image.")
	assert.Equal(true, ret, "Failed to update test image.")
}

func TestGetImage(t *testing.T) {
	assert := assert.New(t)
	testimage.ImageID = testImageid
	retimage, err := service.GetImage(testimage, "id")
	assert.Nil(err, "Failed, err when get test image.")
	assert.Equal("testimage1", retimage.ImageName, "Failed to get image name.")
	t.Log(retimage)
}

func TestClean(t *testing.T) {
	t.Log("Clean test image", testImageid)
	cleanImageTestByID(testImageid)
}

var testITag = model.ImageTag{
	ImageName:       "testimage1",
	Tag:             "1.0",
	ImageTagDeleted: 0,
}

var testITagid int64

func TestCreateImageTag(t *testing.T) {
	assert := assert.New(t)
	id, err := service.CreateImageTag(testITag)
	assert.Nil(err, "Failed, err when create test image tag.")
	assert.NotEqual(0, id, "Failed to assign a image tag id")
	testITagid = id
	t.Log(testITagid)
}

func TestUpdateImageTag(t *testing.T) {
	assert := assert.New(t)
	testITag.Tag = "1.1"
	testITag.ImageTagID = testITagid
	ret, err := service.UpdateImageTag(testITag, "tag")
	assert.Nil(err, "Failed, err when update test image tag.")
	assert.Equal(true, ret, "Failed to update test image tag.")
	t.Log(ret)
}

func TestDeleteImageTag(t *testing.T) {
	assert := assert.New(t)
	testITag.ImageTagID = testITagid
	err := service.DeleteImageTag(testITag)
	assert.Nil(err, "Failed, err when delete test image tag.")
}

func TestGetImageTag(t *testing.T) {
	assert := assert.New(t)
	testITag.ImageTagID = testITagid
	retimageTag, err := service.GetImageTag(testITag, "id")
	assert.Nil(err, "Failed, err when get test image tag.")
	assert.Equal("testimage1", retimageTag.ImageName, "Failed to get image tag.")
	t.Log(retimageTag)
}

func TestCleanTag(t *testing.T) {
	t.Log("Clean test image tag", testITagid)
	cleanImageTagTestByID(testITagid)
}

func cleanImageTagTestByID(imageTagid int64) {
	o := orm.NewOrm()
	rs := o.Raw("delete from image_tag where id = ?", imageTagid)
	r, err := rs.Exec()
	if err != nil {
		logs.Error("Error occurred while deleting image tag: %+v", err)
	}
	affected, err := r.RowsAffected()
	if err != nil {
		logs.Error("Error occurred while deleting image: %+v", err)
	}
	if affected == 0 {
		logs.Error("Failed to delete image tag", imageTagid)
	} else {
		logs.Info("Successfully clean up image tag.", imageTagid)
	}
}
