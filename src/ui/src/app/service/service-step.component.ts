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

interface DeepCopy<T> {
  deepCopyBySelf(): T;
}

export class ServerServiceStep {
  public phase: ServiceStepPhase;
  public project_id?: number = 0;
  public service_name?: string = "";
  public instance?: number = 0;
  public postData?: Object;
}

export abstract class UIServiceStepBase {
  abstract uiToServer(): ServerServiceStep;

  abstract serverToUi(serverResponse: Object): UIServiceStepBase;
}

export class ImageIndex implements DeepCopy<ImageIndex> {
  image_name: string = "";
  image_tag: string = "";

  deepCopyBySelf(): ImageIndex {
    return Object.assign(new ImageIndex(), this);
  }
}

export class EnvStruct implements DeepCopy<EnvStruct> {
  dockerfile_envname: string = "";
  dockerfile_envvalue: string = "";

  deepCopyBySelf(): EnvStruct {
    return Object.assign(new EnvStruct(), this);
  }
}

export class VolumeStruct implements DeepCopy<VolumeStruct> {
  public target_storage_service: string = "";
  public target_path: string = "";
  public volume_name: string = "";
  public container_path: string = "";

  deepCopyBySelf(): VolumeStruct {
    return Object.assign(new VolumeStruct(), this);
  }
}

export class Container implements DeepCopy<Container> {
  public name: string = "";
  public working_dir: string = "";
  public volume_mount: VolumeStruct = new VolumeStruct();
  public image: ImageIndex = new ImageIndex();
  public project_name: string = "";
  public env: Array<EnvStruct> = Array<EnvStruct>();
  public container_port: Array<number> = Array();
  public command: string = "";

  deepCopyBySelf(): Container {
    let result: Container = new Container;
    result.name = this.name;
    result.working_dir = this.working_dir;
    result.volume_mount = this.volume_mount.deepCopyBySelf();
    result.image = this.image.deepCopyBySelf();
    result.project_name = this.project_name;
    this.env.forEach((env: EnvStruct) => {
      result.env.push(env.deepCopyBySelf());
    });
    result.container_port = Array.from(this.container_port);
    result.command = this.command;
    return result;
  }
}

export class NodeType implements DeepCopy<NodeType> {
  target_port: number = 0;
  node_port: number = 0;

  deepCopyBySelf(): NodeType {
    return Object.assign(new NodeType(), this);
  }
}

export class LoadBalancer implements DeepCopy<LoadBalancer> {
  external_access: string;

  deepCopyBySelf(): LoadBalancer {
    let result: LoadBalancer = new LoadBalancer();
    result.external_access = this.external_access;
    return result
  }
}

export class ExternalService implements DeepCopy<ExternalService> {
  public container_name: string = "";
  public node_config: NodeType = new NodeType();
  public load_balancer_config: LoadBalancer = new LoadBalancer();

  deepCopyBySelf(): ExternalService {
    let result: ExternalService = new ExternalService();
    result.container_name = this.container_name;
    result.node_config = this.node_config.deepCopyBySelf();
    result.load_balancer_config = this.load_balancer_config.deepCopyBySelf();
    return result;
  }
}

export class ConfigServiceStep {
  project_id: number = 0;
  service_id: number = 0;
  image_list: Array<ImageIndex> = Array<ImageIndex>();
  service_name: string = "";
  instance: number = 0;
  container_list: Array<Container> = Array<Container>();
  external_service_list: Array<ExternalService> = Array<ExternalService>();
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
    let step1 = new UIServiceStep1();
    if (serverResponse && serverResponse["project_id"]) {
      step1.projectId = serverResponse["project_id"];
    }
    return step1;
  }
}

export class UIServiceStep2 extends UIServiceStepBase {
  public imageList: Array<ImageIndex> = Array<ImageIndex>();

  uiToServer(): ServerServiceStep {
    let result = new ServerServiceStep();
    let postData: Array<ImageIndex> = Array<ImageIndex>();
    result.phase = PHASE_SELECT_IMAGES;
    this.imageList.forEach((value: ImageIndex) => {
      postData.push(value.deepCopyBySelf());
    });
    result.postData = postData;
    return result;
  }

  serverToUi(serverResponse: Object): UIServiceStep2 {
    let step2 = new UIServiceStep2();
    if (serverResponse && serverResponse["image_list"]) {
      let list: Array<ImageIndex> = serverResponse["image_list"];
      list.forEach((value: ImageIndex) => {
        step2.imageList.push(value.deepCopyBySelf())
      });
    }
    return step2;
  }
}

export class UIServiceStep3 extends UIServiceStepBase {
  public containerList: Array<Container> = Array<Container>();

  uiToServer(): ServerServiceStep {
    let result = new ServerServiceStep();
    let postData: Array<Container> = Array<Container>();
    result.phase = PHASE_CONFIG_CONTAINERS;
    this.containerList.forEach((value: Container) => {
      postData.push(value.deepCopyBySelf());
    });
    result.postData = postData;
    return result;
  }

  serverToUi(serverResponse: Object): UIServiceStep3 {
    let step3 = new UIServiceStep3();
    if (serverResponse && serverResponse["container_list"]) {
      let list: Array<Container> = serverResponse["container_list"];
      list.forEach((value: Container) => {
        step3.containerList.push(value.deepCopyBySelf())
      });
    }
    return step3;
  }
}

export class UIServiceStep4 extends UIServiceStepBase {
  public serviceName: string = "";
  public instance: number = 0;
  public externalServiceList: Array<ExternalService> = Array<ExternalService>();

  uiToServer(): ServerServiceStep {
    let result = new ServerServiceStep();
    let postData: Array<ExternalService> = Array<ExternalService>();
    result.phase = PHASE_ENTIRE_SERVICE;
    result.service_name = this.serviceName;
    result.instance = this.instance;
    this.externalServiceList.forEach((value: ExternalService) => {
      postData.push(value.deepCopyBySelf());
    });
    result.postData = postData;
    return result;
  }

  serverToUi(serverResponse: Object): UIServiceStep4 {
    let step4 = new UIServiceStep4();
    if (serverResponse && serverResponse["external_service_list"]) {
      let list: Array<ExternalService> = serverResponse["external_service_list"];
      list.forEach((value: ExternalService) => {
        step4.externalServiceList.push(value.deepCopyBySelf());
      });
    }
    if (serverResponse && serverResponse["instance"]) {
      step4.instance = serverResponse["instance"];
    }
    if (serverResponse && serverResponse["service_name"]) {
      step4.serviceName = serverResponse["service_name"];
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