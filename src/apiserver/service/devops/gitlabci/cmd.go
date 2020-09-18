package gitlabci

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type GitlabCI int

type Image struct {
	Name       string   `json:"name"`
	Entrypoint []string `json:"entrypoint"`
}

type Job struct {
	Stage  string   `json:"stage"`
	Image  *Image   `json:"image,omitempty"`
	Tags   []string `json:"tags"`
	Script []string `json:"script"`
}

var GitlabCIFilename = ".gitlab-ci.yml"

func marshalToBytes(g *map[string]Job) ([]byte, error) {
	return yaml.Marshal(g)
}

func unmarshalToObject(data []byte) (*map[string]Job, error) {
	var gy map[string]Job
	err := yaml.Unmarshal(data, &gy)
	return &gy, err
}

func (g GitlabCI) WriteMultiLine(formatter string, params ...interface{}) string {
	appendCR := fmt.Sprintf("%s\n", formatter)
	return fmt.Sprintf(appendCR, params...)
}

func (g GitlabCI) GenerateGitlabCI(ci map[string]Job, targetPath string) error {
	data, err := json.Marshal(ci)
	if err != nil {
		return err
	}
	gy, err := unmarshalToObject(data)
	if err != nil {
		return err
	}
	datay, err := marshalToBytes(gy)
	if err != nil {
		return err
	}
	for _, job := range *gy {
		datay = append([]byte("- "+job.Stage+"\n"), datay...)
	}
	datay = append([]byte("stages:\n"), datay...)
	return ioutil.WriteFile(filepath.Join(targetPath, GitlabCIFilename), datay, 0755)
}
