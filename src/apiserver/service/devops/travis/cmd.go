package travis

import (
	"io/ioutil"
	"path/filepath"

	yaml "gopkg.in/yaml.v2"
)

type TravisCommand struct {
	BeforeInstall struct {
		Commands []string `yaml:"before_install,flow,omitempty"`
	} `yaml:",inline"`
	Install struct {
		Commands []string `yaml:"install,flow,omitempty"`
	} `yaml:",inline"`
	BeforeScript struct {
		Commands []string `yaml:"before_script,flow,omitempty"`
	} `yaml:",inline"`
	Script struct {
		Commands []string `yaml:"script,flow,omitempty"`
	} `yaml:",inline"`
	AfterSuccess struct {
		Commands []string `yaml:"after_success,flow,omitempty"`
	} `yaml:",inline"`
	AfterFailure struct {
		Commands []string `yaml:"after_failure,flow,omitempty"`
	} `yaml:",inline"`
	BeforeDeploy struct {
		Commands []string `yaml:"before_deploy,flow,omitempty"`
	} `yaml:",inline"`
	Deploy struct {
		Commands []string `yaml:"deploy,flow,omitempty"`
	} `yaml:",inline"`
	AfterDeploy struct {
		Commands []string `yaml:"after_deploy,flow,omitempty"`
	} `yaml:",inline"`
	AfterScript struct {
		Commands []string `yaml:"after_script,flow,omitempty"`
	} `yaml:",inline"`
}

var TravisFilename = ".travis.yml"

func unmarshal(filename string, command *TravisCommand) error {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return err
	}
	return yaml.UnmarshalStrict(data, command)
}

func marshal(filename string, command *TravisCommand) error {
	data, err := yaml.Marshal(command)
	if err != nil {
		return err
	}
	return ioutil.WriteFile(filename, data, 0755)
}

func (tc *TravisCommand) GenerateCustomTravis(targetPath string) error {
	return marshal(filepath.Join(targetPath, TravisFilename), tc)
}

func (tc *TravisCommand) ParseCustomTravis(sourcePath string) error {
	return unmarshal(filepath.Join(sourcePath, TravisFilename), tc)
}
