package model

type RegistryRepo struct {
	Names []string `json:"repositories"`
}

type Image struct {
	ImageID      int64  `json:"-" orm:"column(id)"`
	ImageName    string `json:"image_name" orm:"column(name)"`
	ImageComment string `json:"image_comment" orm:"column(comment)"`
}

type ImageTag struct {
	ImageTagID int64  `orm:"column(id)"`
	ImageName  string `orm:"column(image_name)"`
	Tag        string `orm:"column(tag)"`
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
type CopyStruct struct {
	CopyFrom string `json:"dockerfile_copyfrom"`
	CopyTo   string `json:"dockerfile_copyto"`
}

type Dockerfile struct {
	Base       string       `json:"image_base"`
	Author     string       `json:"image_author"`
	Volume     []string     `json:"image_volume,omitempty"`
	Copy       []CopyStruct `json:"image_copy,omitempty"`
	RUN        []string     `json:"image_run,omitempty"`
	EntryPoint string       `json:"image_entrypoint"`
	Command    string       `json:"image_cmd"`
}

type ImageConfig struct {
	ImageName       string     `json:"image_name"`
	ImageTag        string     `json:"image_tag"`
	ProjectName     string     `json:"project_name"`
	ImageTemplate   string     `json:"image_template"`
	ImageDockerfile Dockerfile `json:"image_dockerfile"`
}
