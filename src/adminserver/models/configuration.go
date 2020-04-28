package models

import (
	"os"
	"reflect"

	"github.com/alyu/configparser"
)

//Apiserver refers to those properties controlling Apiserver parameters.
type Apiserver struct {
	Hostname          string `json:"hostname"`
	APIServerPort     string `json:"api_server_port"`
	KubeHTTPScheme    string `json:"kube_http_scheme"`
	KubeMasterIP      string `json:"kube_master_ip"`
	KubeMasterPort    string `json:"kube_master_port"`
	RegistryIP        string `json:"registry_ip"`
	RegistryPort      string `json:"registry_port"`
	ImageBaselineTime string `json:"image_baseline_time"`
}

//GetApiserver returns data extracted from the Apiserver part of the cfg file.
func GetApiserver(section *configparser.Section) Apiserver {
	array := [...]string{"hostname", "api_server_port", "kube_http_scheme", "kube_master_ip",
		"kube_master_port", "registry_ip", "registry_port", "image_baseline_time"}
	var apiserver Apiserver
	value := reflect.ValueOf(&apiserver).Elem()
	for i := 0; i < value.NumField(); i++ {
		value.Field(i).SetString(section.ValueOf(array[i]))
	}
	return apiserver
}

//UpdateApiserver sets value for properties of corresponding part of cfg file.
func UpdateApiserver(section *configparser.Section, apiserver Apiserver) {
	array := [...]string{"hostname", "api_server_port", "kube_http_scheme", "kube_master_ip",
		"kube_master_port", "registry_ip", "registry_port", "image_baseline_time"}
	value := reflect.ValueOf(&apiserver).Elem()
	for i := 0; i < value.NumField(); i++ {
		if value.Field(i).String() != "" {
			section.SetValueFor(array[i], value.Field(i).String())
		}
	}
}

//Gogitserver refers to those properties controlling Gogitserver parameters.
type Gogitserver struct {
	HostIP   string `json:"gogits_host_ip"`
	HostPort string `json:"gogits_host_port"`
	SSHPort  string `json:"gogits_ssh_port"`
}

//GetGogitserver returns data extracted from the Gogitserver part of the cfg file.
func GetGogitserver(section *configparser.Section) Gogitserver {
	array := [...]string{"gogits_host_ip", "gogits_host_port", "gogits_ssh_port"}
	var gogitserver Gogitserver
	value := reflect.ValueOf(&gogitserver).Elem()
	for i := 0; i < value.NumField(); i++ {
		value.Field(i).SetString(section.ValueOf(array[i]))
	}
	return gogitserver
}

//UpdateGogitserver sets value for properties of corresponding part of cfg file.
func UpdateGogitserver(section *configparser.Section, gogitserver Gogitserver) {
	array := [...]string{"gogits_host_ip", "gogits_host_port", "gogits_ssh_port"}
	value := reflect.ValueOf(&gogitserver).Elem()
	for i := 0; i < value.NumField(); i++ {
		if value.Field(i).String() != "" {
			section.SetValueFor(array[i], value.Field(i).String())
		}
	}
}

//Jenkinsserver refers to those properties controlling Jenkinsserver parameters.
type Jenkinsserver struct {
	HostIP        string `json:"jenkins_host_ip"`
	HostPort      string `json:"jenkins_host_port"`
	NodeIP        string `json:"jenkins_node_ip"`
	NodeSSHPort   string `json:"jenkins_node_ssh_port"`
	NodeUsername  string `json:"jenkins_node_username"`
	NodePassword  string `json:"jenkins_node_password"`
	NodeVolume    string `json:"jenkins_node_volume"`
	ExecutionMode string `json:"jenkins_execution_mode"`
}

//GetJenkinsserver returns data extracted from the Jenkinsserver part of the cfg file.
func GetJenkinsserver(section *configparser.Section) Jenkinsserver {
	array := [...]string{"jenkins_host_ip", "jenkins_host_port", "jenkins_node_ip", "jenkins_node_ssh_port",
		"jenkins_node_username", "jenkins_node_password", "jenkins_node_volume", "jenkins_execution_mode"}
	var jenkinsserver Jenkinsserver
	value := reflect.ValueOf(&jenkinsserver).Elem()
	for i := 0; i < value.NumField(); i++ {
		value.Field(i).SetString(section.ValueOf(array[i]))
	}
	return jenkinsserver
}

//UpdateJenkinsserver sets value for properties of corresponding part of cfg file.
func UpdateJenkinsserver(section *configparser.Section, jenkinsserver Jenkinsserver) {
	array := [...]string{"jenkins_host_ip", "jenkins_host_port", "jenkins_node_ip", "jenkins_node_ssh_port",
		"jenkins_node_username", "jenkins_node_password", "jenkins_node_volume", "jenkins_execution_mode"}
	value := reflect.ValueOf(&jenkinsserver).Elem()
	for i := 0; i < value.NumField(); i++ {
		if value.Field(i).String() != "" {
			section.SetValueFor(array[i], value.Field(i).String())
		}
	}
}

//Kvm refers to those properties controlling Kvm parameters.
type Kvm struct {
	RegistrySize string `json:"kvm_registry_size"`
	RegistryPort string `json:"kvm_registry_port"`
	ToolkitsPath string `json:"kvm_toolkits_path"`
}

//GetKvm returns data extracted from the Kvm part of the cfg file.
func GetKvm(section *configparser.Section) Kvm {
	array := [...]string{"kvm_registry_size", "kvm_registry_port", "kvm_toolkits_path"}
	var kvm Kvm
	value := reflect.ValueOf(&kvm).Elem()
	for i := 0; i < value.NumField(); i++ {
		value.Field(i).SetString(section.ValueOf(array[i]))
	}
	return kvm
}

//UpdateKvm sets value for properties of corresponding part of cfg file.
func UpdateKvm(section *configparser.Section, kvm Kvm) {
	array := [...]string{"kvm_registry_size", "kvm_registry_port", "kvm_toolkits_path"}
	value := reflect.ValueOf(&kvm).Elem()
	for i := 0; i < value.NumField(); i++ {
		if value.Field(i).String() != "" {
			section.SetValueFor(array[i], value.Field(i).String())
		}
	}
}

//Other includes properties involving database, security, etc.
type Other struct {
	ArchType                        string `json:"arch_type"`
	DBPassword                      string `json:"db_password"`
	TokenCacheExpireSeconds         string `json:"token_cache_expire_seconds"`
	TokenExpireSeconds              string `json:"token_expire_seconds"`
	ElaseticsearchMemoryInMegabytes string `json:"elaseticsearch_memory_in_megabytes"`
	TillerPort                      string `json:"tiller_port"`
	BoardAdminPassword              string `json:"board_admin_password"`
	AuthMode                        string `json:"auth_mode"`
	VerificationURL                 string `json:"verification_url"`
	RedirectionURL                  string `json:"redirection_url"`
	AuditDebug                      string `json:"audit_debug"`
	DNSSuffix                       string `json:"dns_suffix"`
	DBMaxConnections                string `json:"db_max_connections"`
	Mode                            string `json:"mode"`
}

//GetOther returns data extracted from the Other part of the cfg file.
func GetOther(section *configparser.Section) Other {
	array := [...]string{"arch_type", "db_password", "token_cache_expire_seconds", "token_expire_seconds",
		"elaseticsearch_memory_in_megabytes", "tiller_port", "board_admin_password", "auth_mode",
		"verification_url", "redirection_url", "audit_debug", "dns_suffix", "db_max_connections", "mode"}
	var other Other
	value := reflect.ValueOf(&other).Elem()
	for i := 0; i < value.NumField(); i++ {
		value.Field(i).SetString(section.ValueOf(array[i]))
	}
	return other
}

//UpdateOther sets value for properties of corresponding part of cfg file.
func UpdateOther(section *configparser.Section, other Other) {
	array := [...]string{"arch_type", "db_password", "token_cache_expire_seconds", "token_expire_seconds",
		"elaseticsearch_memory_in_megabytes", "tiller_port", "board_admin_password", "auth_mode",
		"verification_url", "redirection_url", "audit_debug", "dns_suffix", "db_max_connections", "mode"}
	value := reflect.ValueOf(&other).Elem()
	for i := 0; i < value.NumField(); i++ {
		if value.Field(i).String() != "" {
			section.SetValueFor(array[i], value.Field(i).String())
		}
	}
}

//Ldap refers to those properties controlling Ldap parameters.
type Ldap struct {
	URL     string `json:"ldap_url"`
	Basedn  string `json:"ldap_basedn"`
	UID     string `json:"ldap_uid"`
	Scope   string `json:"ldap_scope"`
	Timeout string `json:"ldap_timeout"`
}

//GetLdap returns data extracted from the Ldap part of the cfg file.
func GetLdap(section *configparser.Section) Ldap {
	array := [...]string{"ldap_url", "ldap_basedn", "ldap_uid", "ldap_scope", "ldap_timeout"}
	var ldap Ldap
	value := reflect.ValueOf(&ldap).Elem()
	for i := 0; i < value.NumField(); i++ {
		value.Field(i).SetString(section.ValueOf(array[i]))
	}
	return ldap
}

//UpdateLdap sets value for properties of corresponding part of cfg file.
func UpdateLdap(section *configparser.Section, ldap Ldap) {
	array := [...]string{"ldap_url", "ldap_basedn", "ldap_uid", "ldap_scope", "ldap_timeout"}
	value := reflect.ValueOf(&ldap).Elem()
	for i := 0; i < value.NumField(); i++ {
		if value.Field(i).String() != "" {
			section.SetValueFor(array[i], value.Field(i).String())
		}
	}
}

//Email refers to those properties controlling Email parameters.
type Email struct {
	Identity   string `json:"email_identity"`
	Server     string `json:"email_server"`
	ServerPort string `json:"email_server_port"`
	Username   string `json:"email_username"`
	Password   string `json:"email_password"`
	From       string `json:"email_from"`
	SSL        string `json:"email_ssl"`
}

//GetEmail returns data extracted from the Email part of the cfg file.
func GetEmail(section *configparser.Section) Email {
	array := [...]string{"email_identity", "email_server", "email_server_port", "email_username",
		"email_password", "email_from", "email_ssl"}
	var email Email
	value := reflect.ValueOf(&email).Elem()
	for i := 0; i < value.NumField(); i++ {
		value.Field(i).SetString(section.ValueOf(array[i]))
	}
	return email
}

//UpdateEmail sets value for properties of corresponding part of cfg file.
func UpdateEmail(section *configparser.Section, email Email) {
	array := [...]string{"email_identity", "email_server", "email_server_port", "email_username",
		"email_password", "email_from", "email_ssl"}
	value := reflect.ValueOf(&email).Elem()
	for i := 0; i < value.NumField(); i++ {
		if value.Field(i).String() != "" {
			section.SetValueFor(array[i], value.Field(i).String())
		}
	}
}

//Configuration combines all the sections above together, referring to the whole cfg file.
type Configuration struct {
	Apiserver     Apiserver
	Gogitserver   Gogitserver
	Jenkinsserver Jenkinsserver
	Kvm           Kvm
	Other         Other
	Ldap          Ldap
	Email         Email
	FirstTimePost bool   `json:"first_time_post"`
	TmpExist      bool   `json:"tmp_exist"`
	Current       string `json:"current"`
}

//GetConfiguration returns data extracted from the whole cfg file.
func GetConfiguration(section *configparser.Section) Configuration {
	configuration := Configuration{
		Apiserver:     GetApiserver(section),
		Gogitserver:   GetGogitserver(section),
		Jenkinsserver: GetJenkinsserver(section),
		Kvm:           GetKvm(section),
		Other:         GetOther(section),
		Ldap:          GetLdap(section),
		Email:         GetEmail(section),
		FirstTimePost: true,
		TmpExist:      false,
		Current:       "cfg"}
	return configuration
}

//UpdateConfiguration sets value for properties for the cfg file.
func UpdateConfiguration(section *configparser.Section, cfg *Configuration) {
	UpdateApiserver(section, cfg.Apiserver)
	UpdateGogitserver(section, cfg.Gogitserver)
	UpdateJenkinsserver(section, cfg.Jenkinsserver)
	UpdateKvm(section, cfg.Kvm)
	UpdateOther(section, cfg.Other)
	UpdateLdap(section, cfg.Ldap)
	UpdateEmail(section, cfg.Email)
}

//Account refers to a username with its password.
type Account struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

var DBconfigdir string = "/data/board/make/config/db"
var DBcompose string = os.Getenv("DB_COMPOSE")
var Boardcompose string = os.Getenv("BOARD_COMPOSE")
var PrepareFile string = os.Getenv("PREPARE_FILE")
var MakePath string = os.Getenv("MAKE_PATH")
