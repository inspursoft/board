package models

import (
	"os"
	"reflect"

	"github.com/alyu/configparser"
)

type Board struct {
	ArchType       string `json:"arch_type"`
	Mode           string `json:"mode"`
	AccessProtocol string `json:"access_protocol"`
	Hostname       string `json:"hostname"`
	APIServerPort  string `json:"api_server_port"`
	DevopsOpt      string `json:"devops_opt"`
	AuthMode       string `json:"auth_mode"`
	AuditDebug     string `json:"audit_debug"`
}

type K8s struct {
	KubeHTTPScheme    string `json:"kube_http_scheme"`
	KubeMasterIP      string `json:"kube_master_ip"`
	KubeMasterPort    string `json:"kube_master_port"`
	RegistryIP        string `json:"registry_ip"`
	RegistryPort      string `json:"registry_port"`
	ImageBaselineTime string `json:"image_baseline_time"`
	TillerPort        string `json:"tiller_port"`
	DNSSuffix         string `json:"dns_suffix"`
}

type Gogs struct {
	HostIP   string `json:"gogits_host_ip"`
	HostPort string `json:"gogits_host_port"`
	SSHPort  string `json:"gogits_ssh_port"`
}

type Gitlab struct {
	HostIP        string `json:"gitlab_host_ip"`
	HostPort      string `json:"gitlab_host_port"`
	SSHPort       string `json:"gitlab_host_ssh_port"`
	AdminToken    string `json:"gitlab_admin_token"`
	SSHUsername   string `json:"gitlab_ssh_username"`
	SSHPassword   string `json:"gitlab_ssh_password"`
	HelperVersion string `json:"gitlab_helper_version"`
}

type Prometheus struct {
	URL string `json:"prometheus_url"`
}

type Jenkins struct {
	HostIP       string `json:"jenkins_host_ip"`
	HostPort     string `json:"jenkins_host_port"`
	NodeIP       string `json:"jenkins_node_ip"`
	NodeSSHPort  string `json:"jenkins_node_ssh_port"`
	NodeUsername string `json:"jenkins_node_username"`
	NodePassword string `json:"jenkins_node_password"`
	NodeVolume   string `json:"jenkins_node_volume"`
}

type Es struct {
	Memory   string `json:"elaseticsearch_memory_in_megabytes"`
	Password string `json:"elastic_password"`
}

type Db struct {
	Password           string `json:"db_password"`
	MaxConn            string `json:"db_max_connections"`
	BoardAdminPassword string `json:"board_admin_password"`
}

type Indata struct {
	VerificationURL string `json:"verification_url"`
	RedirectionURL  string `json:"redirection_url"`
}

type Ldap struct {
	URL       string `json:"ldap_url"`
	SearchDN  string `json:"ldap_searchdn"`
	SearchPWD string `json:"ldap_search_pwd"`
	Basedn    string `json:"ldap_basedn"`
	Filter    string `json:"ldap_filter"`
	UID       string `json:"ldap_uid"`
	Scope     string `json:"ldap_scope"`
	Timeout   string `json:"ldap_timeout"`
}

type Email struct {
	Identity   string `json:"email_identity"`
	Server     string `json:"email_server"`
	ServerPort string `json:"email_server_port"`
	Username   string `json:"email_username"`
	Password   string `json:"email_password"`
	From       string `json:"email_from"`
	SSL        string `json:"email_ssl"`
}

type TokenCfg struct {
	CacheExpireSeconds string `json:"token_cache_expire_seconds"`
	ExpireSeconds      string `json:"token_expire_seconds"`
}

type Configuration struct {
	Board         Board      `json:"board"`
	K8s           K8s        `json:"k8s"`
	Gogs          Gogs       `json:"gogs"`
	Gitlab        Gitlab     `json:"gitlab"`
	Prometheus    Prometheus `json:"prometheus"`
	Jenkins       Jenkins    `json:"jenkins"`
	Es            Es         `json:"es"`
	Db            Db         `json:"db"`
	Indata        Indata     `json:"indata"`
	Ldap          Ldap       `json:"ldap"`
	Email         Email      `json:"email"`
	TokenCfg      TokenCfg   `json:"token"`
	FirstTimePost bool       `json:"first_time_post"`
	TmpExist      bool       `json:"tmp_exist"`
	Current       string     `json:"current"`
}

func GetCfg(section *configparser.Section, part interface{}) interface{} {
	value := reflect.ValueOf(part).Elem()
	rtype := reflect.TypeOf(part).Elem()
	var item string
	for i := 0; i < value.NumField(); i++ {
		item = rtype.Field(i).Tag.Get("json")
		value.Field(i).SetString(section.ValueOf(item))
	}
	return part
}

func SetCfg(section *configparser.Section, part interface{}) {
	value := reflect.ValueOf(part).Elem()
	rtype := reflect.TypeOf(part).Elem()
	var item string
	for i := 0; i < value.NumField(); i++ {
		item = rtype.Field(i).Tag.Get("json")
		if value.Field(i).String() != "" {
			section.SetValueFor(item, value.Field(i).String())
		}
	}
}

func GetConfiguration(section *configparser.Section) Configuration {
	configuration := Configuration{
		Board:         *GetCfg(section, &Board{}).(*Board),
		K8s:           *GetCfg(section, &K8s{}).(*K8s),
		Gogs:          *GetCfg(section, &Gogs{}).(*Gogs),
		Gitlab:        *GetCfg(section, &Gitlab{}).(*Gitlab),
		Prometheus:    *GetCfg(section, &Prometheus{}).(*Prometheus),
		Jenkins:       *GetCfg(section, &Jenkins{}).(*Jenkins),
		Es:            *GetCfg(section, &Es{}).(*Es),
		Db:            *GetCfg(section, &Db{}).(*Db),
		Indata:        *GetCfg(section, &Indata{}).(*Indata),
		Ldap:          *GetCfg(section, &Ldap{}).(*Ldap),
		Email:         *GetCfg(section, &Email{}).(*Email),
		TokenCfg:      *GetCfg(section, &TokenCfg{}).(*TokenCfg),
		FirstTimePost: true,
		TmpExist:      false,
		Current:       "cfg"}
	return configuration
}

func UpdateConfiguration(section *configparser.Section, cfg *Configuration) {
	SetCfg(section, &(cfg.Board))
	SetCfg(section, &(cfg.K8s))
	SetCfg(section, &(cfg.Gogs))
	SetCfg(section, &(cfg.Gitlab))
	SetCfg(section, &(cfg.Prometheus))
	SetCfg(section, &(cfg.Jenkins))
	SetCfg(section, &(cfg.Es))
	SetCfg(section, &(cfg.Db))
	SetCfg(section, &(cfg.Indata))
	SetCfg(section, &(cfg.Ldap))
	SetCfg(section, &(cfg.Email))
	SetCfg(section, &(cfg.TokenCfg))
}

type Account struct {
	Id       int    `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	Port     int    `json:"port"`
}

var DBconfigdir string = "/data/board/make/config/db"
var DBcompose string = os.Getenv("DB_COMPOSE")
var Boardcompose string = os.Getenv("BOARD_COMPOSE")
var PrepareFile string = os.Getenv("PREPARE_FILE")
var MakePath string = os.Getenv("MAKE_PATH")
var BoardcomposeLegacy string = os.Getenv("BOARD_COMPOSE_LEGACY")
