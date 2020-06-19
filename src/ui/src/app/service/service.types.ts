export class ServiceNodeInfo {
  Type: string;
  Address: string;
}

export class ServiceContainer {
  ContainerName: string;
  NodeIP: string;
  PodName: string;
  ServiceName: string;
  SecurityContext: boolean;
  InitContainer: boolean;
}

export class ServiceDetailInfo {
  node_Name: Array<ServiceNodeInfo>;
  node_Port: Array<number>;
  service_Containers: Array<ServiceContainer>;

  constructor() {
    this.node_Name = Array<ServiceNodeInfo>();
    this.node_Port = Array<number>();
    this.service_Containers = Array<ServiceContainer>();
  }
}
