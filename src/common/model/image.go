package model

import "time"

type RegistryRepo struct {
	Names []string `json:"repositories"`
}

type Image struct {
	ImageID           int64     `json:"-" orm:"column(id);pk"`
	ImageName         string    `json:"image_name" orm:"column(name)"`
	ImageComment      string    `json:"image_comment" orm:"column(comment)"`
	ImageDeleted      int       `json:"image_deleted" orm:"column(deleted)"`
	ImageUpdateTime   time.Time `json:"image_update_time" orm:"-"`
	ImageCreationTime time.Time `json:"image_creation_time" orm:"-"`
}

type ImageList []*Image

func (c ImageList) Len() int {
	return len(c)
}

func (c ImageList) Swap(i, j int) {
	if c[i].ImageUpdateTime.Unix() > c[j].ImageUpdateTime.Unix() {
		c[i], c[j] = c[j], c[i]
	}
}

func (c ImageList) Less(i, j int) bool {
	return c[i].ImageUpdateTime.Unix() > c[j].ImageUpdateTime.Unix()
}

type ImageTag struct {
	ImageTagID      int64  `orm:"column(id);pk"`
	ImageName       string `orm:"column(image_name)"`
	Tag             string `orm:"column(tag)"`
	ImageTagDeleted int    `orm:"column(deleted)"`
}

type RegistryTags struct {
	ImageName string   `json:"name"`
	Tags      []string `json:"tags"`
}

type Manifest2Config struct {
	MediaType string `json:"mediaType"`
	Size      int    `json:"size"`
	Digest    string `json:"digest"`
}

type RegistryManifest2 struct {
	SchemaVersion int               `json:"schemaVersion"`
	MediaType     string            `json:"mediaType"`
	Config        Manifest2Config   `json:"config"`
	Layers        []Manifest2Config `json:"layers"`
	//Layers interface{} `json:"layers"`
}

type RegistryManifest1 struct {
	SchemaVersion int                 `json:"schemaVersion"`
	ImageName     string              `json:"name"`
	ImageTag      string              `json:"tag"`
	ImageArch     string              `json:"architecture"`
	FsLayers      interface{}         `json:"fsLayers"`
	History       []map[string]string `json:"history"`
}

type TagDetail struct {
	ImageName         string `json:"image_name"`
	ImageTag          string `json:"image_tag"`
	ImageAuthor       string `json:"image_author"`
	ImageId           string `json:"image_id"`
	ImageCreationTime string `json:"image_creationtime"`
	ImageSize         int    `json:"image_size_number"`
	ImageSizeUnit     string `json:"image_size_unit"`
	ImageDetail       string `json:"-"`
	ImageArch         string `json:"image_arch"`
	ImageTagDeleted   int    `json:"image_tag_deleted"`
}

// The structure for dockerfile template
type CopyStruct struct {
	CopyFrom string `json:"dockerfile_copyfrom"`
	CopyTo   string `json:"dockerfile_copyto"`
}

type EnvStruct struct {
	EnvName  string `json:"dockerfile_envname"`
	EnvValue string `json:"dockerfile_envvalue"`
}

type Dockerfile struct {
	Base       string       `json:"image_base"`
	Author     string       `json:"image_author"`
	Volume     []string     `json:"image_volume"`
	Copy       []CopyStruct `json:"image_copy"`
	RUN        []string     `json:"image_run"`
	EntryPoint string       `json:"image_entrypoint"`
	Command    string       `json:"image_cmd"`
	EnvList    []EnvStruct  `json:"image_env"`
	ExposePort []string     `json:"image_expose"`
}

type ImageConfig struct {
	ImageName           string     `json:"image_name"`
	ImageTag            string     `json:"image_tag"`
	ProjectName         string     `json:"project_name"`
	ImageTemplate       string     `json:"image_template"`
	ImageDockerfile     Dockerfile `json:"image_dockerfile"`
	ImageDockerfilePath string     `json:"-"`
	RepoPath            string     `json:"-"`
	NodeSelection       string     `json:"node_selection"`
}

type ImageIndex struct {
	ImageName   string `json:"image_name"`
	ImageTag    string `json:"image_tag"`
	ProjectName string `json:"project_name"`
}
