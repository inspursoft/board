package gitlabci

import (
	"encoding/json"
	"io/ioutil"

	yaml "gopkg.in/yaml.v2"
)

type GitlabCI int

type Job struct {
	Stage  []string `json:"stage"`
	Tags   []string `json:"tags"`
	Script []string `json:"script"`
}

func marshalToBytes(g *map[string]Job) ([]byte, error) {
	return yaml.Marshal(g)
}

func unmarshalToObject(data []byte) (*map[string]Job, error) {
	var gy map[string]Job
	err := yaml.Unmarshal(data, &gy)
	return &gy, err
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
	return ioutil.WriteFile(targetPath, datay, 0755)
}
