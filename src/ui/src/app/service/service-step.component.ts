export interface ServiceStepComponent {
  data: any;
}

export class ServiceEnvOutput {
  constructor(public key: string = "",
              public value: string = "") {
  }
}

export class ServiceStep1Output {
  service_id: number;
  service_name: string;

  constructor(public  project_id: number = 0,
              public project_name: string = "") {
    this.service_id = 0;
  }
}

export type ServiceStep2Output = Array<ServiceStep2Type>;
export class ServiceStep2Type {
  image_name: string;
  image_tag: string;
  project_id: number;
  project_name: string;
  image_template: string;

  constructor() {
  }
}

export class ImageDockerfile {
  image_base: string;
  image_author: string;
  image_volume?: Array<string>;
  image_copy?: Array<{dockerfile_copyfrom?: string, dockerfile_copyto?: string}>;
  image_run?: Array<string>;
  image_env?: Array<{dockerfile_envname?: string, dockerfile_envvalue?: string}>;
  image_expose?: Array<string>;
  image_entrypoint?: string;
  image_cmd?: string;

  constructor() {
    this.image_base = "";
    this.image_volume = Array<string>();
    this.image_run = Array<string>();
    this.image_expose = Array<string>();
    this.image_copy = Array<{dockerfile_copyfrom?: string, dockerfile_copyto?: string}>();
    this.image_env = Array<{dockerfile_envname?: string, dockerfile_envvalue?: string}>();
  }
}

export class ServiceStep2NewImageType extends ServiceStep2Type {
  image_dockerfile: ImageDockerfile;

  constructor() {
    super();
    this.image_dockerfile = new ImageDockerfile();
  }
}

export type ServiceStep3Output = Array<ServiceStep3Type>;
export class ServiceStep3Type {
  container_name: string;
  container_baseimage: string;
  container_workdir: string;
  container_ports: Array<number>;
  container_volumes: Array<{
    container_dir: string,
    target_storagename: string,
    target_storageServer: string,
    target_dir: string
  }>;
  container_envs: Array<{env_name: string, env_value: string}>;
  container_command: Array<string>;
  container_memory: string;
  container_cpu: string;

  constructor() {
    this.container_ports = Array<number>();
    this.container_volumes = Array<{
      container_dir: string,
      target_storagename: string,
      target_storageServer,
      target_dir: string
    }>();
    this.container_envs = Array<{env_name: string, env_value: string}>();
    this.container_command = Array<string>();
    this.container_memory = "";
    this.container_cpu = "";
  }
}

export class ServiceStep4Output {
  service_id: number;
  project_id: number;
  project_name: string;
  config_phase: string;
  deployment_yaml: {
    deployment_name: string,
    deployment_replicas: number,
    volume_list: Array<{
      volume_name: string,
      server_name: string,
      volume_path: string
    }>,
    container_list?: ServiceStep3Output
  };
  service_yaml: {
    service_name: string,
    service_external: Array<{
      service_containername: string,
      service_containerport: number,
      service_nodeport: number,
      service_externalpath: string;
    }>
    service_selectors: Array<string>
  };

  constructor() {
    this.deployment_yaml = {
      deployment_name: "",
      deployment_replicas: 1,
      volume_list: Array<{volume_name: string, server_name: string, volume_path: string}>()
    };
    this.service_yaml = {
      service_name: "",
        service_external: Array<{
        service_containername: string,
        service_containerport: number,
        service_nodeport: number,
        service_externalpath: string;
      }>(),
        service_selectors: Array<string>()
    };
    this.service_yaml.service_external.push({
      service_containername: "",
      service_externalpath: "",
      service_nodeport: 0,
      service_containerport: 0
    })
  }
}

export class ServiceStep6Output {
  service_id: number;
  service_name: string;
  project_id: number;
  project_name: string;
  service_owner: string;
  service_creationtime: string;
  service_public: number;

  constructor() {
  }
}
