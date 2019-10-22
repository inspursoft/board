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

export class VolumeStruct implements UiServerExchangeData<VolumeStruct> {
  public volume_type: 'nfs' | 'pvc' = 'nfs';
  public target_storage_service = '';
  public target_path = '';
  public volume_name = '';
  public container_path = '';
  public container_path_flag = 0;
  public target_pvc = '';
  public container_file = '';
  public target_file = '';

  serverToUi(serverResponse: Object): VolumeStruct {
    this.target_storage_service = serverResponse['target_storage_service'];
    this.target_path = serverResponse['target_path'];
    this.volume_name = serverResponse['volume_name'];
    this.container_path = serverResponse['container_path'];
    this.volume_type = serverResponse['volume_type'];
    this.container_path_flag = serverResponse['container_path_flag'];
    this.target_pvc = serverResponse['target_pvc'];
    this.container_file = serverResponse['container_file'];
    this.target_file = serverResponse['target_file'];
    return this;
  }

  get containerPathFlag(): boolean {
    return this.container_path_flag == 1;
  }

  set containerPathFlag(value) {
    this.container_path_flag = value ? 1 : 0;
  }

  uiToServer(): VolumeStruct {
    return this;
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
  public volume_mounts: Array<VolumeStruct>;
  public image: ImageIndex;
  public env: Array<EnvStruct>;
  public container_port: Array<number>;

  constructor() {
    this.volume_mounts = Array<VolumeStruct>();
    this.image = new ImageIndex();
    this.env = Array<EnvStruct>();
    this.container_port = Array<number>();
  }

  serverToUi(serverResponse: Object): Container {
    this.name = serverResponse["name"];
    this.working_dir = serverResponse["working_dir"];
    this.cpu_request = serverResponse["cpu_request"];
    this.cpu_limit = serverResponse["cpu_limit"];
    this.mem_request = serverResponse["mem_request"];
    this.mem_limit = serverResponse["mem_limit"];
    this.command = serverResponse["command"];
    if (serverResponse["volume_mounts"]) {
      let tempVolumeDataList: Array<Object> = serverResponse["volume_mounts"];
      tempVolumeDataList.forEach(tempVolumeData => {
        let volume = new VolumeStruct();
        this.volume_mounts.push(volume.serverToUi(tempVolumeData));
      })
    }
    this.image.serverToUi(serverResponse["image"]);
    if (serverResponse["env"]) {
      let envArr: Array<EnvStruct> = serverResponse["env"];
      envArr.forEach((env: EnvStruct) => {
        let envStruct = new EnvStruct();
        this.env.push(envStruct.serverToUi(env));
      });
    }
    if (serverResponse["container_port"]) {
      this.container_port = Array.from(serverResponse["container_port"]) as Array<number>;
    }
    return this;
  }

  uiToServer(): Container {
    return this;
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
    return `${this.serviceName}`
  }
}

export class UIServiceStep1 extends UIServiceStepBase {
  public projectId = 0;
  public projectName = '';

  uiToServer(): ServerServiceStep {
    let result = new ServerServiceStep();
    result.phase = PHASE_SELECT_PROJECT;
    result.project_id = this.projectId;
    return result;
  }

  serverToUi(serverResponse: Object): UIServiceStep1 {
    this.projectId = serverResponse["project_id"];
    this.projectName = serverResponse["project_name"];
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
    let result = new ServerServiceStep();
    let postData: Array<Container> = Array<Container>();
    result.phase = this.isInitContainers ? PHASE_CONFIG_INIT_CONTAINERS : PHASE_CONFIG_CONTAINERS;
    result.project_id = this.projectId;
    this.containerList.forEach((value: Container) => {
      postData.push(value.uiToServer());
    });
    result.postData = postData;
    return result;
  }

  serverToUi(serverResponse: Object): UIServiceStep2 {
    const containerListKey = this.isInitContainers ? 'initcontainer_list' : 'container_list';
    if (serverResponse && serverResponse[containerListKey]) {
      let list: Array<Container> = serverResponse[containerListKey];
      list.forEach((value: Container) => {
        let container = new Container();
        container.serverToUi(value);
        this.containerList.push(container);
      });
    }
    if (serverResponse && serverResponse["project_id"]) {
      this.projectId = serverResponse["project_id"];
    }
    if (serverResponse && serverResponse["project_name"]) {
      this.projectName = serverResponse["project_name"];
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
  public projectName = "";
  public serviceName = "";
  public nodeSelector = "";
  public serviceType = ServiceType.ServiceTypeNormalNodePort;
  public clusterIp = "";
  public instance = 1;
  public servicePublic = false;
  public sessionAffinityFlag = false;
  public externalServiceList: Array<ExternalService>;
  public affinityList: Array<{antiFlag: boolean, services: Array<AffinityCardData>}>;

  constructor() {
    super();
    this.affinityList = Array<{antiFlag: boolean, services: Array<AffinityCardData>}>();
    this.externalServiceList = Array<ExternalService>();
  }

  uiToServer(): ServerServiceStep {
    let result = new ServerServiceStep();
    let postAffinityData: Array<{anti_flag: number, service_names: Array<string>}> = Array<{anti_flag: number, service_names: Array<string>}>();
    result.phase = PHASE_EXTERNAL_SERVICE;
    result.service_name = this.serviceName;
    result.service_type = this.serviceType;
    result.instance = this.instance;
    result.session_affinity_flag = this.sessionAffinityFlag ? 1 : 0;
    result.cluster_ip = this.clusterIp;
    result.service_public = this.servicePublic ? 1 : 0;
    result.node_selector = this.nodeSelector;
    this.affinityList.forEach((value: {antiFlag: boolean, services: Array<AffinityCardData>}) => {
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
      let list: Array<{anti_flag: number, service_names: Array<string>}> = serverResponse["affinity_list"];
      list.forEach((value: {anti_flag: number, service_names: Array<string>}) => {
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
