package utils

import (
	"bytes"
	"html/template"
)

func GenerateURL(rawTemplate string, data interface{}) (string, error) {
	t := template.Must(template.New("").Parse(rawTemplate))
	var targetURL bytes.Buffer
	err := t.Execute(&targetURL, data)
	if err != nil {
		return "", err
	}
	return targetURL.String(), nil
}
