import { DragStatus } from "../shared/shared.types";
import { SERVICE_STATUS } from "../shared/shared.const";
import { ServiceType } from "./service";

export const PHASE_SELECT_PROJECT = "SELECT_PROJECT";
export const PHASE_CONFIG_CONTAINERS = "CONFIG_CONTAINERS";
export const PHASE_CONFIG_INIT_CONTAINERS = "CONFIG_INIT_CONTAINERS";
export const PHASE_EXTERNAL_SERVICE = "EXTERNAL_SERVICE";
export const PHASE_ENTIRE_SERVICE = "ENTIRE_SERVICE";
export type ServiceStepPhase =
  "SELECT_PROJECT"
  | "SELECT_IMAGES"
  | "CONFIG_CONTAINERS"
  | "CONFIG_INIT_CONTAINERS"
  | "EXTERNAL_SERVICE"
  | "ENTIRE_SERVICE"

export type VolumeType = 'nfs' | 'pvc' | 'configmap';

export enum ContainerType {
  runContainer, initContainer
}

export interface UiServerExchangeData<T> {
  uiToServer(): Object;

  serverToUi(serverResponse: Object): T;
}

export abstract class UIServiceStepBase implements UiServerExchangeData<UIServiceStepBase> {
  abstract uiToServer(): ServerServiceStep;

  abstract serverToUi(serverResponse: Object): UIServiceStepBase;
}

export class ServerServiceStep {
  public phase: ServiceStepPhase;
  public project_id?: number = 0;
  public service_name?: string = "";
  public service_type? = ServiceType.ServiceTypeNormalNodePort;
  public node_selector?: string = "";
  public cluster_ip?: string = "";
  public instance?: number = 0;
  public postData?: Object;
  public service_public?: number = 0;
  public session_affinity_flag?: number = 0;
}

export class ImageIndex implements UiServerExchangeData<ImageIndex> {
  image_name = '';
  image_tag = '';
  project_name = '';

  serverToUi(serverResponse: Object): ImageIndex {
    this.image_tag = serverResponse['image_tag'];
    this.image_name = serverResponse['image_name'];
    this.project_name = serverResponse['project_name'];
    return this;
  }

  uiToServer(): ImageIndex {
    return this;
  }
}

export class EnvStruct implements UiServerExchangeData<EnvStruct> {
  dockerfile_envname = '';
  dockerfile_envvalue = '';
  configmap_name = '';
  configmap_key = '';

  serverToUi(serverResponse: Object): EnvStruct {
    this.dockerfile_envname = serverResponse['dockerfile_envname'];
    this.dockerfile_envvalue = serverResponse['dockerfile_envvalue'];
    this.configmap_name = serverResponse['configmap_name'];
    this.configmap_key = serverResponse['configmap_key'];
    return this;
  }

  uiToServer(): EnvStruct {
    return this;
  }
}

export class Volume implements UiServerExchangeData<Volume> {
  volumeType: VolumeType = 'nfs';
  targetStorageService = '';
  targetPath = '';
  volumeName = '';
  containerPath = '';
  containerPathFlagProp = 0;
  targetPvc = '';
  targetConfigMap = '';
  containerFile = '';
  targetFile = '';

  serverToUi(response: object): Volume {
    this.targetStorageService = Reflect.get(response, 'target_storage_service');
    this.targetPath = Reflect.get(response, 'target_path');
    this.volumeName = Reflect.get(response, 'volume_name');
    this.containerPath = Reflect.get(response, 'container_path');
    this.volumeType = Reflect.get(response, 'volume_type');
    this.containerPathFlagProp = Reflect.get(response, 'container_path_flag');
    this.targetPvc = Reflect.get(response, 'target_pvc');
    this.targetConfigMap = Reflect.get(response, 'target_configmap');
    this.containerFile = Reflect.get(response, 'container_file');
    this.targetFile = Reflect.get(response, 'target_file');
    return this;
  }

  get containerPathFlag(): boolean {
    return this.containerPathFlagProp === 1;
  }

  set containerPathFlag(value) {
    this.containerPathFlagProp = value ? 1 : 0;
  }

  uiToServer(): object {
    return {
      target_storage_service: this.targetStorageService,
      target_path: this.targetPath,
      volume_name: this.volumeName,
      container_path: this.containerPath,
      volume_type: this.volumeType,
      container_path_flag: this.containerPathFlagProp,
      target_pvc: this.targetPvc,
      target_configmap: this.targetConfigMap,
      container_file: this.containerFile,
      target_file: this.targetFile
    };
  }
}

export class Container implements UiServerExchangeData<Container> {
  public name = '';
  public working_dir = '';
  public command = '';
  public cpu_request = '';
  public mem_request = '';
  public cpu_limit = '';
  public mem_limit = '';
  public volume_mounts: Array<Volume>;
  public image: ImageIndex;
  public env: Array<EnvStruct>;
  public container_port: Array<number>;

  constructor() {
    this.volume_mounts = Array<Volume>();
    this.image = new ImageIndex();
    this.env = Array<EnvStruct>();
    this.container_port = Array<number>();
  }

  serverToUi(response: object): Container {
    this.name = Reflect.get(response, 'name');
    this.working_dir = Reflect.get(response, 'working_dir');
    this.cpu_request = Reflect.get(response, 'cpu_request');
    this.cpu_limit = Reflect.get(response, 'cpu_limit');
    this.mem_request = Reflect.get(response, 'mem_request');
    this.mem_limit = Reflect.get(response, 'mem_limit');
    this.command = Reflect.get(response, 'command');
    if (Reflect.get(response, 'volume_mounts')) {
      const volumeList: Array<object> = Reflect.get(response, 'volume_mounts');
      volumeList.forEach(data => {
        const volume = new Volume();
        this.volume_mounts.push(volume.serverToUi(data));
      });
    }
    this.image.serverToUi(Reflect.get(response, 'image'));
    if (Reflect.get(response, 'env')) {
      const envList: Array<object> = Reflect.get(response, 'env');
      envList.forEach(data => {
        const envStruct = new EnvStruct();
        this.env.push(envStruct.serverToUi(data));
      });
    }
    if (Reflect.get(response, 'container_port')) {
      this.container_port = Array.from(Reflect.get(response, 'container_port')) as Array<number>;
    }
    return this;
  }

  uiToServer(): object {
    const postVolumes = new Array<object>();
    const postEnvs = new Array<object>();
    this.volume_mounts.forEach(value => postVolumes.push(value.uiToServer()));
    this.env.forEach(value => postEnvs.push(value.uiToServer()));
    return {
      name: this.name,
      working_dir: this.working_dir,
      command: this.command,
      cpu_request: this.cpu_request,
      mem_request: this.mem_request,
      cpu_limit: this.cpu_limit,
      mem_limit: this.mem_limit,
      volume_mounts: postVolumes,
      image: this.image.uiToServer(),
      env: postEnvs,
      container_port: this.container_port
    };
  }
}

export class NodeType implements UiServerExchangeData<NodeType> {
  target_port = 0;
  node_port = 0;

  serverToUi(serverResponse: Object): NodeType {
    this.target_port = serverResponse['target_port'];
    this.node_port = serverResponse['node_port'];
    return this;
  }

  uiToServer(): NodeType {
    return this;
  }
}

export class LoadBalancer implements UiServerExchangeData<LoadBalancer> {
  external_access = '';

  serverToUi(serverResponse: Object): LoadBalancer {
    this.external_access = serverResponse['external_access'];
    return this;
  }

  uiToServer(): LoadBalancer {
    return this;
  }
}

export class ExternalService implements UiServerExchangeData<ExternalService> {
  public container_name = '';
  public node_config: NodeType;
  public load_balancer_config: LoadBalancer;

  constructor() {
    this.node_config = new NodeType();
    this.load_balancer_config = new LoadBalancer();
  }

  serverToUi(serverResponse: Object): ExternalService {
    this.container_name = serverResponse["container_name"];
    this.node_config.serverToUi(serverResponse["node_config"]);
    this.load_balancer_config.serverToUi(serverResponse["load_balancer_config"]);
    return this;
  }

  uiToServer(): ExternalService {
    return this;
  }
}

export enum AffinityCardListView {
  aclvColumn = 'column', aclvRow = 'row'
}

export class AffinityCardData {
  serviceName = '';
  serviceStatus: SERVICE_STATUS;
  status? = DragStatus.dsReady;

  get key(): string {
    return `${this.serviceName}`;
  }
}

export class UIServiceStep1 extends UIServiceStepBase {
  public projectId = 0;
  public projectName = '';

  uiToServer(): ServerServiceStep {
    const result = new ServerServiceStep();
    result.phase = PHASE_SELECT_PROJECT;
    result.project_id = this.projectId;
    return result;
  }

  serverToUi(res: object): UIServiceStep1 {
    this.projectId = Reflect.get(res, 'project_id');
    this.projectName = Reflect.get(res, 'project_name');
    return this;
  }
}

export class UIServiceStep2 extends UIServiceStepBase {
  public containerList: Array<Container>;
  public projectId = 0;
  public projectName = '';
  public isInitContainers = false;

  constructor() {
    super();
    this.containerList = Array<Container>();
  }

  uiToServer(): ServerServiceStep {
    const result = new ServerServiceStep();
    const postData: Array<object> = Array<object>();
    result.phase = this.isInitContainers ? PHASE_CONFIG_INIT_CONTAINERS : PHASE_CONFIG_CONTAINERS;
    result.project_id = this.projectId;
    this.containerList.forEach((value: Container) => {
      postData.push(value.uiToServer());
    });
    result.postData = postData;
    return result;
  }

  serverToUi(res: object): UIServiceStep2 {
    const containerListKey = this.isInitContainers ? 'initcontainer_list' : 'container_list';
    if (res && Reflect.get(res, containerListKey)) {
      const list: Array<object> = Reflect.get(res, containerListKey);
      list.forEach((value: object) => {
        const container = new Container();
        container.serverToUi(value);
        this.containerList.push(container);
      });
    }
    if (res && Reflect.get(res, 'project_id')) {
      this.projectId = Reflect.get(res, 'project_id');
    }
    if (res && Reflect.get(res, 'project_name')) {
      this.projectName = Reflect.get(res, 'project_name');
    }
    return this;
  }

  getPortList(containerName: string): Array<number> {
    return this.containerList.find(value => value.name === containerName).container_port;
  }
}

export class UIServiceStep2InitContainer extends UIServiceStep2 {
  constructor() {
    super();
    this.isInitContainers = true;
  }
}

export class UIServiceStep3 extends UIServiceStepBase {
  public projectName = '';
  public serviceName = "";
  public nodeSelector = "";
  public serviceType = ServiceType.ServiceTypeNormalNodePort;
  public clusterIp = "";
  public instance = 1;
  public servicePublic = false;
  public sessionAffinityFlag = false;
  public externalServiceList: Array<ExternalService>;
  public affinityList: Array<{ antiFlag: boolean, services: Array<AffinityCardData> }>;

  constructor() {
    super();
    this.affinityList = Array<{ antiFlag: boolean, services: Array<AffinityCardData> }>();
    this.externalServiceList = Array<ExternalService>();
  }

  uiToServer(): ServerServiceStep {
    let result = new ServerServiceStep();
    let postAffinityData: Array<{ anti_flag: number, service_names: Array<string> }> =
      Array<{ anti_flag: number, service_names: Array<string> }>();
    result.phase = PHASE_EXTERNAL_SERVICE;
    result.service_name = this.serviceName;
    result.service_type = this.serviceType;
    result.instance = this.instance;
    result.session_affinity_flag = this.sessionAffinityFlag ? 1 : 0;
    result.cluster_ip = this.clusterIp;
    result.service_public = this.servicePublic ? 1 : 0;
    result.node_selector = this.nodeSelector;
    this.affinityList.forEach((value: { antiFlag: boolean, services: Array<AffinityCardData> }) => {
      let serviceNames = Array<string>();
      value.services.forEach((card: AffinityCardData) => serviceNames.push(card.serviceName));
      postAffinityData.push({anti_flag: value.antiFlag ? 1 : 0, service_names: serviceNames})
    });
    result.postData = {external_service_list: this.externalServiceList, affinity_list: postAffinityData};
    return result;
  }

  serverToUi(serverResponse: Object): UIServiceStep3 {
    let step3 = new UIServiceStep3();
    if (serverResponse && serverResponse["affinity_list"]) {
      let list: Array<{ anti_flag: number, service_names: Array<string> }> = serverResponse["affinity_list"];
      list.forEach((value: { anti_flag: number, service_names: Array<string> }) => {
        let services = Array<AffinityCardData>();
        if (value.service_names && value.service_names.length > 0) {
          value.service_names.forEach((serviceName: string) => {
            let card = new AffinityCardData();
            card.serviceName = serviceName;
            card.status = DragStatus.dsEnd;
            services.push(card);
          });
        }
        step3.affinityList.push({antiFlag: value.anti_flag == 1, services: services});
      });
    }
    if (serverResponse && serverResponse["external_service_list"]) {
      step3.externalServiceList = serverResponse["external_service_list"];
    }
    if (serverResponse && serverResponse["instance"]) {
      step3.instance = serverResponse["instance"];
    }
    if (serverResponse && serverResponse["service_name"]) {
      step3.serviceName = serverResponse["service_name"];
    }
    if (serverResponse && serverResponse["project_name"]) {
      step3.projectName = serverResponse["project_name"];
    }
    if (serverResponse && serverResponse["service_type"]) {
      step3.serviceType = serverResponse["service_type"];
    }
    if (serverResponse && serverResponse["service_public"]) {
      step3.servicePublic = serverResponse["service_public"] == 1;
    }
    if (serverResponse && serverResponse["node_selector"]) {
      step3.nodeSelector = serverResponse["node_selector"];
    }
    if (serverResponse && serverResponse["cluster_ip"]) {
      step3.clusterIp = serverResponse["cluster_ip"];
    }
    if (serverResponse && serverResponse["session_affinity_flag"]) {
      step3.sessionAffinityFlag = serverResponse["session_affinity_flag"] == 1;
    }
    return step3;
  }
}

export class UiServiceFactory {
  static getInstance(phase: ServiceStepPhase): UIServiceStepBase {
    switch (phase) {
      case PHASE_SELECT_PROJECT:
        return new UIServiceStep1();
      case PHASE_CONFIG_CONTAINERS:
        return new UIServiceStep2();
      case PHASE_EXTERNAL_SERVICE:
        return new UIServiceStep3();
      case PHASE_CONFIG_INIT_CONTAINERS:
        return new UIServiceStep2InitContainer();
      default:
        return null;
    }
  }
}
