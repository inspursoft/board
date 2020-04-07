package models

import (
	"reflect"
	"os"
	"github.com/alyu/configparser"
	"time"
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
}

//GetOther returns data extracted from the Other part of the cfg file.
func GetOther(section *configparser.Section) Other {
	array := [...]string{"arch_type", "db_password", "token_cache_expire_seconds", "token_expire_seconds",
		"elaseticsearch_memory_in_megabytes", "tiller_port", "board_admin_password", "auth_mode",
		"verification_url", "redirection_url", "audit_debug", "dns_suffix"}
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
		"verification_url", "redirection_url", "audit_debug", "dns_suffix"}
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

//Password refers to a certain type with its value.
type Password struct {
	Which string `json:"which"`
	Value string `json:"value"`
}

//const DBconfigdir = "/data/board/Deploy/config/db"
//const DBcompose = "/data/board/Deploy/docker-compose-db.yml"
//const Boardcompose = "/data/board/Deploy/docker-compose-rest.yml"
//const PrepareFile = "/data/board/Deploy/prepare"

var DBconfigdir string = os.Getenv("DB_CONFIG_DIR")
var DBcompose string = os.Getenv("DB_COMPOSE")
var Boardcompose string = os.Getenv("BOARD_COMPOSE")
var PrepareFile string = os.Getenv("PREPARE_FILE")

type DBconf struct {
	Password   		string	`json:"db_password"`
	MaxConnections	int 	`json:"db_max_connections"`
}

type User struct {
	ID           int64     `json:"user_id" orm:"column(id)"`
	Username     string    `json:"user_name" orm:"column(username)"`
	Password     string    `json:"user_password" orm:"column(password)"`
	Email        string    `json:"user_email" orm:"column(email)"`
	Realname     string    `json:"user_realname" orm:"column(realname)"`
	Comment      string    `json:"user_comment" orm:"column(comment)"`
	Deleted      int       `json:"user_deleted" orm:"column(deleted)"`
	SystemAdmin  int       `json:"user_system_admin" orm:"column(system_admin)"`
	ResetUUID    string    `json:"user_reset_uuid" orm:"column(reset_uuid)"`
	Salt         string    `json:"user_salt" orm:"column(salt)"`
	RepoToken    string    `json:"user_token" orm:"column(repo_token)"`
	CreationTime time.Time `json:"user_creation_time" orm:"column(creation_time)"`
	UpdateTime   time.Time `json:"user_update_time" orm:"column(update_time)"`
	FailedTimes  int       `json:"user_failed_times" orm:"column(failed_times)"`
}

type Config struct {
	Name    string `json:"name" orm:"column(name);pk"`
	Value   string `json:"value" orm:"column(value)"`
	Comment string `json:"comment" orm:"column(comment)"`
}

type UUID struct {
	UUID string `json:"UUID"`
}

type Key struct {
	Key string `json:"Key"`
}