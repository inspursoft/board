import { HttpBind, ResponseBase, RequestBase } from '../shared/shared.type';

export class Board extends ResponseBase implements RequestBase {
  @HttpBind('arch_type') archType: string;
  @HttpBind('mode') mode: string;
  @HttpBind('hostname') hostname: string;
  @HttpBind('api_server_port') apiServerPort: string;
  @HttpBind('devops_opt') devopsOpt: string;
  @HttpBind('auth_mode') authMode: string;
  @HttpBind('audit_debug') auditDebug: string;

  constructor(res?: object) {
    super(res);
    if (!res) {
      this.archType = 'x86_64';
      this.mode = 'normal';
      this.hostname = 'reg.mydomain.com';
      this.apiServerPort = '8088';
      this.devopsOpt = 'legacy';
      this.authMode = 'db_auth';
      this.auditDebug = 'false';
    }
  }

  PostBody(): object {
    return {
      arch_type: this.archType.toString(),
      mode: this.mode.toString(),
      hostname: this.hostname.toString(),
      api_server_port: this.apiServerPort.toString().toString(),
      devops_opt: this.devopsOpt.toString(),
      auth_mode: this.authMode.toString(),
      audit_debug: this.auditDebug.toString().toString(),
    };
  }
}

export class K8s extends ResponseBase implements RequestBase {
  @HttpBind('kube_http_scheme') kubeHttpScheme: string;
  @HttpBind('kube_master_ip') kubeMasterIP: string;
  @HttpBind('kube_master_port') kubeMasterPort: string;
  @HttpBind('registry_ip') registryIP: string;
  @HttpBind('registry_port') registryPort: string;
  @HttpBind('image_baseline_time') imageBaselineTime: string;
  @HttpBind('tiller_port') tillerPort: string;
  @HttpBind('dns_suffix') dnsSuffix: string;

  constructor(res?: object) {
    super(res);
    if (!res) {
      this.kubeHttpScheme = 'http';
      this.kubeMasterIP = ' 10.0.0.0';
      this.kubeMasterPort = '8080';
      this.registryIP = '10.0.0.0';
      this.registryPort = '5000';
      this.imageBaselineTime = '2016-01-01 09:00:00';
      this.tillerPort = '31111';
      this.dnsSuffix = '.cluster.local';
    }
  }

  PostBody(): object {
    return {
      kube_http_scheme: this.kubeHttpScheme.toString(),
      kube_master_ip: this.kubeMasterIP.toString(),
      kube_master_port: this.kubeMasterPort.toString(),
      registry_ip: this.registryIP.toString(),
      registry_port: this.registryPort.toString(),
      image_baseline_time: this.imageBaselineTime.toString(),
      tiller_port: this.tillerPort.toString(),
      dns_suffix: this.dnsSuffix.toString(),
    };
  }
}

export class Gogs extends ResponseBase implements RequestBase {
  @HttpBind('gogits_host_ip') hostIP: string;
  @HttpBind('gogits_host_port') hostPort: string;
  @HttpBind('gogits_ssh_port') sshPort: string;

  constructor(res?: object) {
    super(res);
    if (!res) {
      this.hostIP = '10.0.0.0';
      this.hostPort = '10080';
      this.sshPort = '10022';
    }
  }

  PostBody(): object {
    return {
      gogits_host_ip: this.hostIP.toString(),
      gogits_host_port: this.hostPort.toString(),
      gogits_ssh_port: this.sshPort.toString(),
    };
  }
}

export class Gitlab extends ResponseBase implements RequestBase {
  @HttpBind('gitlab_host_ip') hostIP: string;
  @HttpBind('gitlab_host_port') hostPort: string;
  @HttpBind('gitlab_host_ssh_port') sshPort: string;
  @HttpBind('gitlab_helper_version') helperVersion: string;
  @HttpBind('gitlab_admin_token') adminToken: string;
  @HttpBind('gitlab_ssh_username') sshUsername: string;
  @HttpBind('gitlab_ssh_password') sshPassword: string;

  constructor(res?: object) {
    super(res);
    if (!res) {
      this.hostIP = '10.0.0.0';
      this.hostPort = '10088';
      this.sshPort = '10028';
      this.helperVersion = '1.1';
      this.adminToken = '1234567901234567890';
      this.sshUsername = 'root';
      this.sshPassword = '123456a?';
    }
  }

  PostBody(): object {
    return {
      gitlab_host_ip: this.hostIP.toString(),
      gitlab_host_port: this.hostPort.toString(),
      gitlab_host_ssh_port: this.sshPort.toString(),
      gitlab_helper_version: this.helperVersion.toString(),
      gitlab_admin_token: this.adminToken.toString(),
      gitlab_ssh_username: this.sshUsername.toString(),
      gitlab_ssh_password: this.sshPassword.toString(),
    };
  }
}

export class Prometheus extends ResponseBase implements RequestBase {
  @HttpBind('prometheus_url') url: string;

  constructor(res?: object) {
    super(res);
    if (!res) {
      this.url = 'http://10.0.0.0:9090';
    }
  }

  PostBody(): object {
    return {
      prometheus_url: this.url.toString(),
    };
  }
}

export class Jenkins extends ResponseBase implements RequestBase {
  @HttpBind('jenkins_host_ip') hostIP: string;
  @HttpBind('jenkins_host_port') hostPort: string;
  @HttpBind('jenkins_node_ip') nodeIP: string;
  @HttpBind('jenkins_node_ssh_port') nodeSshPort: string;
  @HttpBind('jenkins_node_username') nodeUsername: string;
  @HttpBind('jenkins_node_password') nodePassword: string;
  @HttpBind('jenkins_node_volume') nodeVolume: string;
  @HttpBind('jenkins_execution_mode') executionMode = 'single';
  constructor(res?: object) {
    super(res);
    if (!res) {
      this.hostIP = '10.0.0.0';
      this.hostPort = '8888';
      this.nodeIP = '10.0.0.0';
      this.nodeSshPort = '22';
      this.nodeUsername = 'root';
      this.nodePassword = '123456a?';
      this.nodeVolume = '/data/jenkins_node';
    }
  }

  PostBody(): object {
    return {
      jenkins_host_ip: this.hostIP.toString(),
      jenkins_host_port: this.hostPort.toString(),
      jenkins_node_ip: this.nodeIP.toString(),
      jenkins_node_ssh_port: this.nodeSshPort.toString(),
      jenkins_node_username: this.nodeUsername.toString(),
      jenkins_node_password: this.nodePassword.toString(),
      jenkins_node_volume: this.nodeVolume.toString(),
    };
  }
}

export class ES extends ResponseBase implements RequestBase {
  @HttpBind('elaseticsearch_memory_in_megabytes') memoryInMegabytes: string;
  @HttpBind('elastic_password') password: string;

  constructor(res?: object) {
    super(res);
    if (!res) {
      this.memoryInMegabytes = '1024';
      this.password = 'root123';
    }
  }

  PostBody(): object {
    return {
      elaseticsearch_memory_in_megabytes: this.memoryInMegabytes.toString(),
      elastic_password: this.password.toString(),
    };
  }
}

export class DB extends ResponseBase implements RequestBase {
  @HttpBind('db_password') dbPassword: string;
  @HttpBind('db_max_connections') dbMaxConnections: string;
  @HttpBind('board_admin_password') boardAdminPassword: string;

  constructor(res?: object) {
    super(res);
    if (!res) {
      this.dbPassword = 'root123';
      this.dbMaxConnections = '1000';
      this.boardAdminPassword = '123456a?';
    }
  }

  PostBody(): object {
    return {
      db_password: this.dbPassword.toString(),
      db_max_connections: this.dbMaxConnections.toString(),
      board_admin_password: this.boardAdminPassword.toString(),
    };
  }
}

export class Indata extends ResponseBase implements RequestBase {
  @HttpBind('verification_url') verificationUrl: string;
  @HttpBind('redirection_url') redirectionUrl: string;

  constructor(res?: object) {
    super(res);
    if (!res) {
      this.verificationUrl = 'http://verification.mydomain.com';
      this.redirectionUrl = 'http://redirection.mydomain.com';
    }
  }

  PostBody(): object {
    return {
      verification_url: this.verificationUrl.toString(),
      redirection_url: this.redirectionUrl.toString(),
    };
  }
}

export class LDAP extends ResponseBase implements RequestBase {
  @HttpBind('ldap_url') url: string;
  @HttpBind('ldap_searchdn') searchdn: string;
  @HttpBind('ldap_search_pwd') searchPwd: string;
  @HttpBind('ldap_basedn') basedn: string;
  @HttpBind('ldap_filter') filter: string;
  @HttpBind('ldap_uid') uid: string;
  @HttpBind('ldap_scope') scope: string;
  @HttpBind('ldap_timeout') timeout: string;

  constructor(res?: object) {
    super(res);
    if (!res) {
      this.url = 'ldaps://ldap.mydomain.com';
      this.searchdn = 'uid=searchuser,ou=people,dc=mydomain,dc=com';
      this.searchPwd = 'password';
      this.basedn = 'ou=people,dc=mydomain,dc=com';
      this.filter = '(objectClass=person)';
      this.uid = 'uid';
      this.scope = 'LDAP_SCOPE_SUBTREE';
      this.timeout = '5';
    }
  }

  PostBody(): object {
    return {
      ldap_url: this.url.toString(),
      ldap_searchdn: this.searchdn.toString(),
      ldap_search_pwd: this.searchPwd.toString(),
      ldap_basedn: this.basedn.toString(),
      ldap_filter: this.filter.toString(),
      ldap_uid: this.uid.toString(),
      ldap_scope: this.scope.toString(),
      ldap_timeout: this.timeout.toString(),
    };
  }
}

export class Email extends ResponseBase implements RequestBase {
  @HttpBind('email_identity') identity: string;
  @HttpBind('email_server') server: string;
  @HttpBind('email_server_port') serverPort: string;
  @HttpBind('email_username') username: string;
  @HttpBind('email_password') password: string;
  @HttpBind('email_from') from: string;
  @HttpBind('email_ssl') ssl: string;

  constructor(res?: object) {
    super(res);
    if (!res) {
      this.identity = '';
      this.server = 'smtp.mydomain.com';
      this.serverPort = '25';
      this.username = 'admin@mydomain.com';
      this.password = '123456a?';
      this.from = 'admin <admin@mydomain.com>';
      this.ssl = 'false';
    }
  }

  PostBody(): object {
    return {
      email_identity: this.identity.toString(),
      email_server: this.server.toString(),
      email_server_port: this.serverPort.toString(),
      email_username: this.username.toString(),
      email_password: this.password.toString(),
      email_from: this.from.toString(),
      email_ssl: this.ssl.toString(),
    };
  }
}

export class Token extends ResponseBase implements RequestBase {
  @HttpBind('token_cache_expire_seconds') cacheExpireSeconds: string;
  @HttpBind('token_expire_seconds') expireSeconds: string;

  constructor(res?: object) {
    super(res);
    if (!res) {
      this.cacheExpireSeconds = '1800';
      this.expireSeconds = '1800';
    }
  }

  PostBody(): object {
    return {
      token_cache_expire_seconds: this.cacheExpireSeconds.toString(),
      token_expire_seconds: this.expireSeconds.toString(),
    };
  }
}

export class Configuration implements RequestBase {
  board: Board;
  k8s: K8s;
  gogs: Gogs;
  gitlab: Gitlab;
  prometheus: Prometheus;
  jenkins: Jenkins;
  es: ES;
  db: DB;
  indata: Indata;
  ldap: LDAP;
  email: Email;
  token: Token;
  isInit: boolean;
  tmpExist: boolean;
  current: string;

  constructor(res?: object) {
    if (res) {
      this.board = new Board(Reflect.get(res, 'board'));
      this.k8s = new K8s(Reflect.get(res, 'k8s'));
      this.gogs = new Gogs(Reflect.get(res, 'gogs'));
      this.gitlab = new Gitlab(Reflect.get(res, 'gitlab'));
      this.prometheus = new Prometheus(Reflect.get(res, 'prometheus'));
      this.jenkins = new Jenkins(Reflect.get(res, 'jenkins'));
      this.es = new ES(Reflect.get(res, 'es'));
      this.db = new DB(Reflect.get(res, 'db'));
      this.indata = new Indata(Reflect.get(res, 'indata'));
      this.ldap = new LDAP(Reflect.get(res, 'ldap'));
      this.email = new Email(Reflect.get(res, 'email'));
      this.token = new Token(Reflect.get(res, 'token'));

      this.isInit = Reflect.get(res, 'first_time_post');
      this.tmpExist = Reflect.get(res, 'tmp_exist');
      this.current = Reflect.get(res, 'current');
    } else {
      this.board = new Board();
      this.k8s = new K8s();
      this.gogs = new Gogs();
      this.gitlab = new Gitlab();
      this.prometheus = new Prometheus();
      this.jenkins = new Jenkins();
      this.es = new ES();
      this.db = new DB();
      this.indata = new Indata();
      this.ldap = new LDAP();
      this.email = new Email();
      this.token = new Token();

      this.isInit = false;
      this.tmpExist = false;
      this.current = 'cfg';
    }
  }

  PostBody(): object {
    return {
      board: this.board.PostBody(),
      k8s: this.k8s.PostBody(),
      gogs: this.gogs.PostBody(),
      gitlab: this.gitlab.PostBody(),
      prometheus: this.prometheus.PostBody(),
      jenkins: this.jenkins.PostBody(),
      es: this.es.PostBody(),
      db: this.db.PostBody(),
      indata: this.indata.PostBody(),
      ldap: this.ldap.PostBody(),
      email: this.email.PostBody(),
      token: this.token.PostBody(),
    };
  }
}
