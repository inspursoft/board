package utils

import (
	"regexp"
	"strings"

	"github.com/astaxie/beego/logs"
)

var patterns = map[string]string{
	"username": `^[0-9a-z_]{4,40}$`,
	"email":    `^(([^<>()[\]\.,;:\s@\"]+(\.[^<>()[\]\.,;:\s@\"]+)*)|(\".+\"))@((\[[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\.[0-9]{1,3}\])|(([a-zA-Z\-0-9]+\.)+[a-zA-Z]{2,}))$`,
	"project":  `^[a-z0-9]+(?:[._-][a-z0-9]+)*$`,
}

const infiniteLength = 65535

func ValidateWithPattern(target string, input string) bool {
	if pattern, exists := patterns[target]; exists {
		matched, err := regexp.MatchString(pattern, input)
		if err != nil {
			logs.Error("Error occurred while validating item: %s, err: %+v", target, err)
			return false
		}
		return matched
	}
	logs.Error("No pattern provided for validating %s currently.", target)
	return false
}

func ValidateWithLengthRange(target string, min int, max int) bool {
	target = strings.TrimSpace(target)
	return len(target) >= min && len(target) <= max
}

func ValidateWithMaxLength(target string, max int) bool {
	return ValidateWithLengthRange(target, 0, max)
}

func ValidateWithMinLength(target string, min int) bool {
	return ValidateWithLengthRange(target, min, infiniteLength)
}
