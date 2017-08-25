package model

type RegistryRepo struct {
	Names []string `json:"repositories"`
}

type BoardImage struct {
	ImageName    string `json:"image_name"`
	ImageComment string `json:"image_comment"`
}

type RegistryTags struct {
	ImageName string   `json:"name"`
	Tags      []string `json:"tags"`
}

type Menifest2Config struct {
	MediaType string `json:"mediaType"`
	Size      int    `json:"size"`
	Digest    string `json:"digest"`
}

type RegistryMenifest2 struct {
	SchemaVersion int               `json:"schemaVersion"`
	MediaType     string            `json:"mediaType"`
	Config        Menifest2Config   `json:"config"`
	Layers        []Menifest2Config `json:"layers"`
	//Layers interface{} `json:"layers"`
}

type RegistryMenifest1 struct {
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
	ImageDetail       string `json:"image_detail"`
}

// The structure for dockerfile template
type Dockerfile struct {
	Base       string `json:"image_base"`
	Author     string `json:"image_author"`
	Volume     string `json:"image_volume"`
	CopyFrom   string `json:"image_copyfrom"`
	CopyTo     string `json:"image_copyto"`
	RUN        string `json:"image_run"`
	EntryPoint string `json:"image_entrypoint"`
	Command    string `json:"image_cmd"`
}

type ImageConfig struct {
	ImageName       string     `json:"image_name"`
	ImageTag        string     `json:"image_tag"`
	ProjectID       int        `json:"project_id"`
	ProjectName     string     `json:"project_name"`
	ImageTemplate   string     `json:"image_template"`
	ImageDockerfile Dockerfile `json:"image_dockerfile"`
}

// For Image info in Database
type ImageDB struct {
	ImageID      int64  `orm:"column(id)"`
	ImageName    string `orm:"column(name)"`
	ImageComment string `orm:"column(comment)"`
}
