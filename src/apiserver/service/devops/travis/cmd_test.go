package travis

import (
	"testing"

	"github.com/astaxie/beego/logs"
	"github.com/stretchr/testify/assert"
)

func TestGenerateCustomTravis(t *testing.T) {
	var travisCommand TravisCommand
	travisCommand.Script.Commands = []string{
		"token=`cat key.txt`",
		"status=`curl -I http://10.165.18.153:8088/api/v1/files/download?token=$token 2>/dev/null | head -n 1 | cut -d$' ' -f2`",
		"if [ $status == '200' ]; then curl -o http://10.165.18.153:8088/api/v1/files/download?token=$token && mkdir -p upload && unzip attachment.zip upload && rm -f attachment.zip; fi",
	}
	travisCommand.AfterScript.Commands = []string{
		"export PATH=/usr/bin:/bin:/usr/sbin:/sbin:/usr/local/bin",
		"docker build -t 10.110.13.136:5000/project11/myimage20180509:v2.5 .",
		"docker push 10.110.13.136:5000/project11/myimage20180509:v2.5",
	}
	err := travisCommand.GenerateCustomTravis(".")
	assert := assert.New(t)
	assert.Nilf(err, "Failed to generate custom travis: %+v", err)
}

func TestParseCustomTravis(t *testing.T) {
	var travisCommand TravisCommand
	err := travisCommand.ParseCustomTravis(".")
	assert := assert.New(t)
	assert.Nilf(err, "Failed to parse custom travis: %+v", err)
	logs.Debug("Parsed custom Travis command: %+v", travisCommand)
}
