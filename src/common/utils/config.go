package utils

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/astaxie/beego/logs"
)

var configStorage map[string]interface{}
var boardHostIP = GetConfig("BOARD_HOST_IP")

func add(name string, value interface{}) {
	configStorage[name] = value
}

func AddEnv(name string, defaultValue ...string) {
	value := os.Getenv(name)
	if value == "" && len(defaultValue) > 0 {
		value = defaultValue[0]
	}
	add(name, value)
}

func AddValue(name string, value interface{}) {
	add(name, value)
}

func GetIntValue(name string) int {
	if value, ok := configStorage[name].(int); ok {
		return value
	}
	value, err := strconv.Atoi(GetStringValue(name))
	if err != nil {
		panic(fmt.Sprintf("Failed to get int value for key: %s", name))
	}
	return value
}

func GetBoolValue(name string) bool {
	if value, ok := configStorage[name].(bool); ok {
		return value
	}
	value, err := strconv.ParseBool(GetStringValue(name))
	if err != nil {
		panic(fmt.Sprintf("Failed to get bool value for key: %s", name))
	}
	return value
}

func GetStringValue(name string, defaultValue ...string) string {
	if s, ok := configStorage[name].(string); ok {
		return s
	}
	if defaultValue != nil && len(defaultValue) > 0 {
		return defaultValue[0]
	}
	panic(fmt.Sprintf("Failed to get value for key: %s", name))
}

func SetConfig(name, formatter string, keys ...string) {
	configStorage[name] = fmt.Sprintf(formatter,
		func() (values []interface{}) {
			for _, key := range keys {
				values = append(values, GetStringValue(key))
			}
			return
		}()...)
	return
}

func GetConfig(name string, defaultValue ...string) func() string {
	return func() string { return GetStringValue(name, defaultValue...) }
}

func Initialize() {
	configStorage = make(map[string]interface{})
}

func ShowAllConfigs() {
	logs.Info("Current configurations in storage:\n")
	for k, v := range configStorage {
		if strings.Contains(strings.ToUpper(k), "PASSWORD") || strings.Contains(strings.ToUpper(k), "PWD") {
			continue
		}
		logs.Info("\t%s: %v", k, v)
	}
}

func InitializeDefaultConfig() {
	logs.SetLogFuncCall(true)
	logs.SetLogFuncCallDepth(4)

	Initialize()

	AddEnv("MODE")
	AddEnv("BOARD_HOST_IP")
	AddEnv("API_SERVER_PORT")
	AddEnv("DB_IP")
	AddEnv("DB_PORT")
	AddEnv("DB_PASSWORD")

	AddEnv("BOARD_ADMIN_PASSWORD")

	AddEnv("TOKEN_CACHE_EXPIRE_SECONDS")

	AddEnv("KUBE_MASTER_IP")
	AddEnv("KUBE_MASTER_PORT")
	AddEnv("REGISTRY_IP")
	AddEnv("REGISTRY_PORT")
	AddEnv("IMAGE_BASELINE_TIME")

	AddEnv("AUTH_MODE")

	AddEnv("LDAP_URL")
	AddEnv("LDAP_SEARCH_DN")
	AddEnv("LDAP_SEARCH_PWD")
	AddEnv("LDAP_BASE_DN")
	AddEnv("LDAP_FILTER")
	AddEnv("LDAP_UID")
	AddEnv("LDAP_SCOPE")
	AddEnv("LDAP_TIMEOUT")
	AddEnv("FORCE_INIT_SYNC")
	AddEnv("VERIFICATION_URL")
	AddEnv("REDIRECTION_URL")

	SetConfig("BOARD_API_BASE_URL", "http://%s:%s/api/v1", "BOARD_HOST_IP", "API_SERVER_PORT")

	AddEnv("GOGITS_HOST_IP", "10.0.0.0")
	AddEnv("GOGITS_HOST_PORT", "10080")
	SetConfig("GOGITS_BASE_URL", "http://%s:%s", "GOGITS_HOST_IP", "GOGITS_HOST_PORT")

	AddEnv("GOGITS_SSH_PORT", "10022")
	SetConfig("GOGITS_SSH_URL", "ssh://git@%s:%s", "GOGITS_HOST_IP", "GOGITS_SSH_PORT")

	AddEnv("JENKINS_HOST_IP", "10.0.0.0")
	AddEnv("JENKINS_HOST_PORT", "8080")
	AddEnv("JENKINS_NODE_IP", "10.0.0.0")
	AddEnv("JENKINS_NODE_USERNAME", "root")
	AddEnv("JENKINS_NODE_PASSWORD", "123456a?")
	AddEnv("JENKINS_NODE_SSH_PORT", "22")
	AddEnv("JENKINS_NODE_VOLUME", "/data/jenkins_node")
	AddEnv("JENKINS_EXECUTION_MODE", "single")
	SetConfig("JENKINS_BASE_URL", "http://%s:%s", "JENKINS_HOST_IP", "JENKINS_HOST_PORT")

	SetConfig("JENKINSFILE_REPO_URL", "http://%s:%s/admin/devops-jenkins.git", "GOGITS_HOST_IP", "GOGITS_HOST_PORT")

	SetConfig("REGISTRY_URL", "http://%s:%s", "REGISTRY_IP", "REGISTRY_PORT")
	SetConfig("KUBE_MASTER_URL", "http://%s:%s", "KUBE_MASTER_IP", "KUBE_MASTER_PORT")
	SetConfig("KUBE_NODE_URL", "http://%s:%s/api/v1/nodes", "KUBE_MASTER_IP", "KUBE_MASTER_PORT")

	SetConfig("API_SERVER_URL", "http://%s:%s", "BOARD_HOST_IP", "API_SERVER_PORT")

	SetConfig("REGISTRY_BASE_URI", "%s:%s", "REGISTRY_IP", "REGISTRY_PORT")

	AddValue("IS_EXTERNAL_AUTH", (GetStringValue("AUTH_MODE") != "db_auth"))

	SetConfig("SYNC_K8S", "false")

	AddEnv("EMAIL_HOST")
	AddEnv("EMAIL_PORT")
	AddEnv("EMAIL_USR")
	AddEnv("EMAIL_PWD")
	AddEnv("EMAIL_SSL")
	AddEnv("EMAIL_FROM")
	AddEnv("EMAIL_IDENTITY")
	AddEnv("AUDIT_DEBUG")
	AddEnv("DNS_SUFFIX")

	AddEnv("TILLER_PORT")

	AddEnv("DEVOPS_OPT")

	AddEnv("GITLAB_HOST_IP")
	AddEnv("GITLAB_HOST_PORT")
	AddEnv("GITLAB_HOST_SSH_PORT")
	AddEnv("GITLAB_ADMIN_TOKEN")

	SetConfig("GITLAB_BASE_URL", "http://%s:%s", "GITLAB_HOST_IP", "GITLAB_HOST_PORT")

	ShowAllConfigs()
}
