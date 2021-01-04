package travis_test

import (
	"github.com/inspursoft/board/src/apiserver/service/devops/travis"
	"os"
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCustomTravis(t *testing.T) {
	var travisCommand travis.TravisCommand
	travisCommand.Script.Commands = []string{
		"token=`cat key.txt`",
		"status=`curl -I http://apiserver:8088/api/v1/files/download?token=$token 2>/dev/null | head -n 1 | cut -d$' ' -f2`",
		"if [ $status == '200' ]; then curl -o http://apiserver:8088/api/v1/files/download?token=$token && unzip attachment.zip upload && rm -f attachment.zip; fi",
	}
	travisCommand.AfterScript.Commands = []string{
		"export PATH=/usr/bin:/bin:/usr/sbin:/sbin:/usr/local/bin",
		"docker build -t registry:5000/project11/myimage20180509:v2.5 .",
		"docker push registry:5000/project11/myimage20180509:v2.5",
	}
	err := travisCommand.GenerateCustomTravis(".")
	assert := assert.New(t)
	assert.Nilf(err, "Failed to generate custom travis: %+v", err)
}

func TestParseCustomTravis(t *testing.T) {
	var travisCommand travis.TravisCommand
	err := travisCommand.ParseCustomTravis(".")
	assert := assert.New(t)
	assert.Nilf(err, "Failed to parse custom travis: %+v", err)
	logs.Debug("Parsed custom Travis command: %+v", travisCommand)
	os.Remove(travis.TravisFilename)
}
