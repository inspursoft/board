import { HttpBind, ResponseBase, RequestBase } from '../shared/shared.type';

export class ApiServer extends ResponseBase implements RequestBase {
  @HttpBind('hostname') hostname: string;
  @HttpBind('api_server_port') apiServerPort: string;
  @HttpBind('kube_http_scheme') kubeHttpScheme: string;
  @HttpBind('kube_master_ip') kubeMasterIp: string;
  @HttpBind('kube_master_port') kubeMasterPort: string;
  @HttpBind('registry_ip') registryIp: string;
  @HttpBind('registry_port') registryPort: string;
  @HttpBind('image_baseline_time') imageBaselineTime: string;

  constructor(res?: object) {
    if (res) {
      super(res);
    } else {
      this.hostname = 'reg.domain.com';
      this.apiServerPort = '8088';
      this.kubeHttpScheme = 'http';
      this.kubeMasterIp = '10.0.0.0';
      this.kubeMasterPort = '8080';
      this.registryIp = '10.0.0.0';
      this.registryPort = '5000';
      this.imageBaselineTime = '2016-01-01 09:00:00';
    }
  }

  PostBody(): object {
    return {
      hostname: this.hostname.toString(),
      api_server_port: this.apiServerPort.toString(),
      kube_http_scheme: this.kubeHttpScheme.toString(),
      kube_master_ip: this.kubeMasterIp.toString(),
      kube_master_port: this.kubeMasterPort.toString(),
      registry_ip: this.registryIp.toString(),
      registry_port: this.registryPort.toString(),
      image_baseline_time: this.imageBaselineTime.toString(),
    };
  }
}

export class Gogits extends ResponseBase implements RequestBase {
  @HttpBind('gogits_host_ip') gogitsHostIp: string;
  @HttpBind('gogits_host_port') gogitsHostPort: string;
  @HttpBind('gogits_ssh_port') gogitsSshPort: string;

  constructor(res?: object) {
    if (res) {
      super(res);
    } else {
      this.gogitsHostIp = '10.0.0.0';
      this.gogitsHostPort = '10080';
      this.gogitsSshPort = '10022';
    }
  }

  PostBody(): object {
    return {
      gogits_host_ip: this.gogitsHostIp.toString(),
      gogits_host_port: this.gogitsHostPort.toString(),
      gogits_ssh_port: this.gogitsSshPort.toString(),
    };
  }
}

export class Jenkins extends ResponseBase implements RequestBase {
  @HttpBind('jenkins_host_ip') jenkinsHostIp: string;
  @HttpBind('jenkins_host_port') jenkinsHostPort: string;
  @HttpBind('jenkins_node_ip') jenkinsNodeIp: string;
  @HttpBind('jenkins_node_ssh_port') jenkinsNodeSshPort: string;
  @HttpBind('jenkins_node_username') jenkinsNodeUsername: string;
  @HttpBind('jenkins_node_password') jenkinsNodePassword: string;
  @HttpBind('jenkins_node_volume') jenkinsNodeVolume: string;
  @HttpBind('jenkins_execution_mode') jenkinsExecutionMode: string;

  constructor(res?: object) {
    if (res) {
      super(res);
    } else {
      this.jenkinsHostIp = '10.0.0.0';
      this.jenkinsHostPort = '8888';
      this.jenkinsNodeIp = '10.0.0.0';
      this.jenkinsNodeSshPort = '22';
      this.jenkinsNodeUsername = 'root';
      this.jenkinsNodePassword = '123456a?';
      this.jenkinsNodeVolume = '/data/jenkins_node';
      this.jenkinsExecutionMode = 'single';
    }
  }

  PostBody(): object {
    return {
      jenkins_host_ip: this.jenkinsHostIp.toString(),
      jenkins_host_port: this.jenkinsHostPort.toString(),
      jenkins_node_ip: this.jenkinsNodeIp.toString(),
      jenkins_node_ssh_port: this.jenkinsNodeSshPort.toString(),
      jenkins_node_username: this.jenkinsNodeUsername.toString(),
      jenkins_node_password: this.jenkinsNodePassword.toString(),
      jenkins_node_volume: this.jenkinsNodeVolume.toString(),
      jenkins_execution_mode: this.jenkinsExecutionMode.toString(),
    };
  }
}

export class Kvm extends ResponseBase implements RequestBase {
  @HttpBind('kvm_registry_size') kvmRegistrySize: string;
  @HttpBind('kvm_registry_port') kvmRegistryPort: string;
  @HttpBind('kvm_toolkits_path') kvmToolkitsPath: string;

  constructor(res?: object) {
    if (res) {
      super(res);
    } else {
      this.kvmRegistrySize = '5';
      this.kvmRegistryPort = '8890';
      this.kvmToolkitsPath = '/root/kvm_toolkits';
    }
  }

  PostBody(): object {
    return {
      kvm_registry_size: this.kvmRegistrySize.toString(),
      kvm_registry_port: this.kvmRegistryPort.toString(),
      kvm_toolkits_path: this.kvmToolkitsPath.toString(),
    };
  }
}

export class Ldap extends ResponseBase implements RequestBase {
  @HttpBind('ldap_url') ldapUrl: string;
  @HttpBind('ldap_basedn') ldapBasedn: string;
  @HttpBind('ldap_uid') ldapUid: string;
  @HttpBind('ldap_scope') ldapScope: string;
  @HttpBind('ldap_timeout') ldapTimeout: string;

  constructor(res?: object) {
    if (res) {
      super(res);
    } else {
      this.ldapUrl = 'ldaps://ldap.mydomain.com';
      this.ldapBasedn = 'ou=people,dc=mydomain,dc=com';
      this.ldapUid = 'uid';
      this.ldapScope = 'LDAP_SCOPE_SUBTREE';
      this.ldapTimeout = '5';
    }
  }

  PostBody(): object {
    return {
      ldap_url: this.ldapUrl.toString(),
      ldap_basedn: this.ldapBasedn.toString(),
      ldap_uid: this.ldapUid.toString(),
      ldap_scope: this.ldapScope.toString(),
      ldap_timeout: this.ldapTimeout.toString(),
    };
  }
}

export class Email extends ResponseBase implements RequestBase {
  @HttpBind('email_identity') emailIdentity: string;
  @HttpBind('email_server') emailServer: string;
  @HttpBind('email_server_port') emailServerPort: string;
  @HttpBind('email_username') emailUsername: string;
  @HttpBind('email_password') emailPassword: string;
  @HttpBind('email_from') emailFrom: string;
  @HttpBind('email_ssl') emailSsl: string;

  constructor(res?: object) {
    if (res) {
      super(res);
    } else {
      this.emailIdentity = '';
      this.emailServer = 'smtp.mydomain.com';
      this.emailServerPort = '25';
      this.emailUsername = 'admin@mydomain.com';
      this.emailPassword = '123456a?';
      this.emailFrom = 'admin <admin@mydomain.com>';
      this.emailSsl = 'false';
    }
  }

  PostBody(): object {
    return {
      email_identity: this.emailIdentity === 'null' ? '' : this.emailIdentity.toString(),
      email_server: this.emailServer.toString(),
      email_server_port: this.emailServerPort.toString(),
      email_username: this.emailUsername.toString(),
      email_password: this.emailPassword.toString(),
      email_from: this.emailFrom.toString(),
      email_ssl: this.emailSsl.toString(),
    };
  }
}

export class Others extends ResponseBase implements RequestBase {
  @HttpBind('arch_type') archType: string;
  @HttpBind('db_password') dbPassword: string;
  @HttpBind('db_max_connections') dbMaxConnections: string;
  @HttpBind('token_cache_expire_seconds') tokenCacheExpireSeconds: string;
  @HttpBind('token_expire_seconds') tokenExpireSeconds: string;
  @HttpBind('elaseticsearch_memory_in_megabytes') elaseticsearchMemoryInMegabytes: string;
  @HttpBind('tiller_port') tillerPort: string;
  @HttpBind('board_admin_password') boardAdminPassword: string;
  @HttpBind('auth_mode') authMode: string;
  @HttpBind('verification_url') verificationUrl: string;
  @HttpBind('redirection_url') redirectionUrl: string;
  @HttpBind('audit_debug') auditDebug: string;
  @HttpBind('dns_suffix') dnsSuffix: string;

  constructor(res?: object) {
    if (res) {
      super(res);
    } else {
      this.archType = 'x86_64';
      this.dbPassword = 'root123';
      this.dbMaxConnections = '1000';
      this.tokenCacheExpireSeconds = '1800';
      this.tokenExpireSeconds = '1800';
      this.elaseticsearchMemoryInMegabytes = '1024';
      this.tillerPort = '31111';
      this.boardAdminPassword = '';
      this.authMode = 'db_auth';
      this.verificationUrl = 'http://verification.mydomain.com';
      this.redirectionUrl = 'http://redirection.mydomain.com';
      this.auditDebug = 'false';
      this.dnsSuffix = '.cluster.local';
    }
  }

  PostBody(): object {
    return {
      arch_type: this.archType.toString(),
      db_password: this.dbPassword.toString(),
      db_max_connections: this.dbMaxConnections.toString(),
      token_cache_expire_seconds: this.tokenCacheExpireSeconds.toString(),
      token_expire_seconds: this.tokenExpireSeconds.toString(),
      elaseticsearch_memory_in_megabytes: this.elaseticsearchMemoryInMegabytes.toString(),
      tiller_port: this.tillerPort.toString(),
      board_admin_password: this.boardAdminPassword.toString(),
      auth_mode: this.authMode.toString(),
      verification_url: this.verificationUrl.toString(),
      redirection_url: this.redirectionUrl.toString(),
      audit_debug: this.auditDebug.toString(),
      dns_suffix: this.dnsSuffix.toString(),
    };
  }
}

export class Configuration implements RequestBase {
  apiserver: ApiServer;
  gogits: Gogits;
  jenkins: Jenkins;
  kvm: Kvm;
  others: Others;
  ldap: Ldap;
  email: Email;
  isInit: boolean;
  tmpExist: boolean;
  current: string;

  constructor(res?: object) {
    if (res) {
      this.apiserver = new ApiServer(Reflect.get(res, 'Apiserver'));
      this.gogits = new Gogits(Reflect.get(res, 'Gogitserver'));
      this.jenkins = new Jenkins(Reflect.get(res, 'Jenkinsserver'));
      this.kvm = new Kvm(Reflect.get(res, 'Kvm'));
      this.others = new Others(Reflect.get(res, 'Other'));
      this.ldap = new Ldap(Reflect.get(res, 'Ldap'));
      this.email = new Email(Reflect.get(res, 'Email'));
      this.isInit = Reflect.get(res, 'first_time_post');
      this.tmpExist = Reflect.get(res, 'tmp_exist');
      this.current = Reflect.get(res, 'current');
    } else {
      this.apiserver = new ApiServer();
      this.gogits = new Gogits();
      this.jenkins = new Jenkins();
      this.kvm = new Kvm();
      this.others = new Others();
      this.ldap = new Ldap();
      this.email = new Email();
      this.isInit = false;
      this.tmpExist = false;
      this.current = 'cfg';
    }
  }

  PostBody(): object {
    return {
      Apiserver: this.apiserver.PostBody(),
      Gogitserver: this.gogits.PostBody(),
      Jenkinsserver: this.jenkins.PostBody(),
      Kvm: this.kvm.PostBody(),
      Other: this.others.PostBody(),
      Ldap: this.ldap.PostBody(),
      Email: this.email.PostBody(),
    };
  }
}
