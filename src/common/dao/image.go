package dao

import (
	"github.com/inspursoft/board/src/common/model"

	//"time"

	"github.com/astaxie/beego/orm"
)

func AddImage(image model.Image) (int64, error) {
	o := orm.NewOrm()

	imageID, err := o.Insert(&image)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return imageID, err
}

func GetImage(image model.Image, fieldNames ...string) (*model.Image, error) {
	o := orm.NewOrm()
	err := o.Read(&image, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &image, err
}

func UpdateImage(image model.Image, fieldNames ...string) (int64, error) {
	o := orm.NewOrm()

	imageID, err := o.Update(&image, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return imageID, err
}

func AddImageTag(imagetag model.ImageTag) (int64, error) {
	o := orm.NewOrm()

	imagetagID, err := o.Insert(&imagetag)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return imagetagID, err
}

func GetImageTag(imageTag model.ImageTag, fieldNames ...string) (*model.ImageTag, error) {
	o := orm.NewOrm()
	err := o.Read(&imageTag, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return &imageTag, err
}

func UpdateImageTag(imageTag model.ImageTag, fieldNames ...string) (int64, error) {
	o := orm.NewOrm()

	imageTagID, err := o.Update(&imageTag, fieldNames...)
	if err != nil {
		if err == orm.ErrNoRows {
			return 0, nil
		}
		return 0, err
	}
	return imageTagID, err
}
