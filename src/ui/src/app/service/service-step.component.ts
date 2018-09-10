export const PHASE_SELECT_PROJECT = "SELECT_PROJECT";
export const PHASE_SELECT_IMAGES = "SELECT_IMAGES";
export const PHASE_CONFIG_CONTAINERS = "CONFIG_CONTAINERS";
export const PHASE_EXTERNAL_SERVICE = "EXTERNAL_SERVICE";
export const PHASE_ENTIRE_SERVICE = "ENTIRE_SERVICE";
export type ServiceStepPhase =
  "SELECT_PROJECT"
  | "SELECT_IMAGES"
  | "CONFIG_CONTAINERS"
  | "EXTERNAL_SERVICE"
  | "ENTIRE_SERVICE"

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
  public node_selector?: string = "";
  public instance?: number = 0;
  public postData?: Object;
  public service_public?: number = 0;
}

export class ImageIndex implements UiServerExchangeData<ImageIndex> {
  image_name: string = "";
  image_tag: string = "";
  project_name: string = "";

  serverToUi(serverResponse: Object): ImageIndex {
    return Object.assign(this, serverResponse);
  }

  uiToServer(): ImageIndex {
    return this;
  }
}

export class EnvStruct implements UiServerExchangeData<EnvStruct> {
  dockerfile_envname: string = "";
  dockerfile_envvalue: string = "";

  serverToUi(serverResponse: Object): EnvStruct {
    return Object.assign(this, serverResponse);
  }

  uiToServer(): EnvStruct {
    return this;
  }
}

export class VolumeStruct implements UiServerExchangeData<VolumeStruct> {
  public target_storage_service: string = "";
  public target_path: string = "";
  public volume_name: string = "";
  public container_path: string = "";

  serverToUi(serverResponse: Object): VolumeStruct {
    return Object.assign(this, serverResponse);
  }

  uiToServer(): VolumeStruct {
    return this;
  }
}

export class Container implements UiServerExchangeData<Container> {
  public name: string = "";
  public working_dir: string = "";
  public volume_mount: VolumeStruct = new VolumeStruct();
  public image: ImageIndex = new ImageIndex();
  public env: Array<EnvStruct> = Array<EnvStruct>();
  public container_port: Array<number> = Array();
  public command: string = "";

  isHavePort(): boolean {
    return this.container_port.length > 0;
  }

  serverToUi(serverResponse: Object): Container {
    this.name = serverResponse["name"];
    this.working_dir = serverResponse["working_dir"];
    this.volume_mount = (new VolumeStruct()).serverToUi(serverResponse["volume_mount"]);
    this.image = (new ImageIndex()).serverToUi(serverResponse["image"]);
    if (serverResponse["env"]) {
      let envArr: Array<EnvStruct> = serverResponse["env"];
      envArr.forEach((env: EnvStruct) => {
        this.env.push((new EnvStruct()).serverToUi(env));
      });
    }
    if (serverResponse["container_port"]) {
      this.container_port = Array.from(serverResponse["container_port"]) as Array<number>;
    }
    this.command = serverResponse["command"];
    return this;
  }

  uiToServer(): Container {
    return this;
  }
}

export class NodeType implements UiServerExchangeData<NodeType> {
  target_port: number = 0;
  node_port: number = 0;

  serverToUi(serverResponse: Object): NodeType {
    return Object.assign(this, serverResponse);
  }

  uiToServer(): NodeType {
    return this;
  }
}

export class LoadBalancer implements UiServerExchangeData<LoadBalancer> {
  external_access: string;

  serverToUi(serverResponse: Object): LoadBalancer {
    return Object.assign(this, serverResponse);
  }

  uiToServer(): LoadBalancer {
    return this;
  }
}

export class ExternalService implements UiServerExchangeData<ExternalService> {
  public container_name: string = "";
  public node_config: NodeType = new NodeType();
  public load_balancer_config: LoadBalancer = new LoadBalancer();

  serverToUi(serverResponse: Object): ExternalService {
    this.container_name = serverResponse["container_name"];
    this.node_config = (new NodeType()).serverToUi(serverResponse["node_config"]);
    this.load_balancer_config = (new LoadBalancer()).serverToUi(serverResponse["load_balancer_config"]);
    return this;
  }

  uiToServer(): ExternalService {
    return this;
  }
}

export class UIServiceStep1 extends UIServiceStepBase {
  public projectId: number = 0;
  public projectName: string = "";

  uiToServer(): ServerServiceStep {
    let result = new ServerServiceStep();
    result.phase = PHASE_SELECT_PROJECT;
    result.project_id = this.projectId;
    return result;
  }

  serverToUi(serverResponse: Object): UIServiceStep1 {
    if (serverResponse && serverResponse["project_id"]) {
      this.projectId = serverResponse["project_id"];
    }
    return this;
  }
}

export class UIServiceStep2 extends UIServiceStepBase {
  public imageList: Array<ImageIndex> = Array<ImageIndex>();
  public projectId: number = 0;
  public projectName: string = "";

  uiToServer(): ServerServiceStep {
    let result = new ServerServiceStep();
    let postData: Array<ImageIndex> = Array<ImageIndex>();
    result.phase = PHASE_SELECT_IMAGES;
    this.imageList.forEach((value: ImageIndex) => {
      postData.push(value.uiToServer());
    });
    result.postData = postData;
    return result;
  }

  serverToUi(serverResponse: Object): UIServiceStep2 {
    if (serverResponse && serverResponse["project_id"]) {
      this.projectId = serverResponse["project_id"];
    }
    if (serverResponse && serverResponse["project_name"]) {
      this.projectName = serverResponse["project_name"];
    }
    if (serverResponse && serverResponse["image_list"]) {
      let list: Array<ImageIndex> = serverResponse["image_list"];
      list.forEach((value: ImageIndex) => {
        this.imageList.push((new ImageIndex()).serverToUi(value))
      });
    }
    return this;
  }
}

export class UIServiceStep3 extends UIServiceStepBase {
  public containerList: Array<Container> = Array<Container>();

  uiToServer(): ServerServiceStep {
    let result = new ServerServiceStep();
    let postData: Array<Container> = Array<Container>();
    result.phase = PHASE_CONFIG_CONTAINERS;
    this.containerList.forEach((value: Container) => {
      postData.push(value.uiToServer());
    });
    result.postData = postData;
    return result;
  }

  serverToUi(serverResponse: Object): UIServiceStep3 {
    if (serverResponse && serverResponse["container_list"]) {
      let list: Array<Container> = serverResponse["container_list"];
      list.forEach((value: Container) => {
        this.containerList.push((new Container()).serverToUi(value))
      });
    }
    return this;
  }
}

export class UIServiceStep4 extends UIServiceStepBase {
  public projectName: string = "";
  public serviceName: string = "";
  public nodeSelector: string = "";
  public instance: number = 1;
  public servicePublic: boolean;
  public externalServiceList: Array<ExternalService> = Array<ExternalService>();

  uiToServer(): ServerServiceStep {
    let result = new ServerServiceStep();
    let postData: Array<ExternalService> = Array<ExternalService>();
    result.phase = PHASE_EXTERNAL_SERVICE;
    result.service_name = this.serviceName;
    result.instance = this.instance;
    result.service_public = this.servicePublic ? 1 : 0;
    result.node_selector = this.nodeSelector;
    this.externalServiceList.forEach((value: ExternalService) => {
      postData.push(value.uiToServer());
    });
    result.postData = postData;
    return result;
  }

  serverToUi(serverResponse: Object): UIServiceStep4 {
    let step4 = new UIServiceStep4();
    if (serverResponse && serverResponse["external_service_list"]) {
      let list: Array<ExternalService> = serverResponse["external_service_list"];
      list.forEach((value: ExternalService) => {
        step4.externalServiceList.push((new ExternalService()).serverToUi(value));
      });
    }
    if (serverResponse && serverResponse["instance"]) {
      step4.instance = serverResponse["instance"];
    }
    if (serverResponse && serverResponse["service_name"]) {
      step4.serviceName = serverResponse["service_name"];
    }
    if (serverResponse && serverResponse["project_name"]) {
      step4.projectName = serverResponse["project_name"];
    }
    if (serverResponse && serverResponse["service_public"]) {
      step4.servicePublic = serverResponse["service_public"] == 1;
    }
    if (serverResponse && serverResponse["node_selector"]) {
      step4.nodeSelector = serverResponse["node_selector"] ;
    }
    return step4;
  }
}

export class UiServiceFactory {
  static getInstance(phase: ServiceStepPhase): UIServiceStepBase {
    switch (phase) {
      case PHASE_SELECT_PROJECT:
        return new UIServiceStep1();
      case PHASE_SELECT_IMAGES:
        return new UIServiceStep2();
      case PHASE_CONFIG_CONTAINERS:
        return new UIServiceStep3();
      case PHASE_EXTERNAL_SERVICE:
        return new UIServiceStep4();
    }
  }
}