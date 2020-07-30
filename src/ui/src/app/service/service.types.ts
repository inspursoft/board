import { HttpBase, HttpBind, HttpBindArray, ResponsePaginationBase } from '../shared/ui-model/model-types';

export enum ServiceType {
  ServiceTypeUnknown,
  ServiceTypeNormalNodePort,
  ServiceTypeHelm,
  ServiceTypeDeploymentOnly,
  ServiceTypeClusterIP,
  ServiceTypeStatefulSet,
  ServiceTypeJob,
  ServiceTypeEdgeComputing
}

export enum ServiceSource {
  ServiceSourceBoard,
  ServiceSourceK8s,
  ServiceSourceHelm
}

export class ServiceNodeInfo extends HttpBase {
  @HttpBind('Type') type = '';
  @HttpBind('Address') address = '';
}

export class ServiceContainer extends HttpBase {
  @HttpBind('ContainerName') containerName = '';
  @HttpBind('NodeIP') nodeIp = '';
  @HttpBind('PodName') podName = '';
  @HttpBind('ServiceName') serviceName = '';
  @HttpBind('SecurityContext') securityContext = false;
  @HttpBind('InitContainer') initContainer = false;
}

export class ServiceDetailInfo extends HttpBase {
  @HttpBindArray('node_Name', ServiceNodeInfo) nodeNames: Array<ServiceNodeInfo>;
  @HttpBindArray('service_Containers', ServiceContainer) serviceContainers: Array<ServiceContainer>;
  @HttpBind('node_Port') nodePorts: Array<number>;

  protected prepareInit() {
    this.nodeNames = new Array<ServiceNodeInfo>();
    this.nodePorts = new Array<number>();
    this.serviceContainers = new Array<ServiceContainer>();
  }

  get isHasNotDetailProperty(): boolean {
    return !Reflect.has(this.res, 'detail');
  }
}

export class ServiceNode extends HttpBase {
  @HttpBind('node_name') nodeName = '';
  @HttpBind('node_ip') nodeIp = '';
  @HttpBind('status') status = 0;
}

export class ServiceNodeGroup extends HttpBase {
  @HttpBind('nodegroup_id') id = 0;
  @HttpBind('nodegroup_project') project = '';
  @HttpBind('nodegroup_name') name = '';
  @HttpBind('nodegroup_comment') comment = '';
}

export class ServiceProject extends HttpBase {
  @HttpBind('project_id') projectId = -1;
  @HttpBind('project_name') projectName = '';
  @HttpBind('publicity') publicity = false;
  @HttpBind('project_public') projectPublic = 0;
  @HttpBind('project_creation_time') projectCreationTime: Date;
  @HttpBind('project_comment') projectComment = '';
  @HttpBind('project_owner_id') projectOwnerId = '';
  @HttpBind('project_owner_name') projectOwnerName = '';
}

export class ServiceImage extends HttpBase {
  @HttpBind('image_name') imageName = '';
  @HttpBind('image_tag') imageTag = '';
  @HttpBind('image_comment') imageComment = '';
  @HttpBind('image_deleted') imageDelete = 0;
  @HttpBind('image_update_time') imageUpdateTime = '';
  @HttpBind('image_creation_time') imageCreationTime = '';
}

export class ServiceImageDetail extends HttpBase {
  @HttpBind('image_name') imageName = '';
  @HttpBind('image_tag') imageTag = '';
  @HttpBind('image_detail') imageDetail = '';
  @HttpBind('image_creationtime') imageCreationTime = '';
  @HttpBind('image_size_number') imageSizeNumber = 0;
  @HttpBind('image_size_unit') imageSizeUnit = 'MB';
}

export class ServiceDockerfileCopy extends HttpBase {
  @HttpBind('dockerfile_copyfrom') copyFrom = '';
  @HttpBind('dockerfile_copyto') copyTo = '';
}

export class ServiceDockerfileEnv extends HttpBase {
  @HttpBind('dockerfile_envname') envName = '';
  @HttpBind('dockerfile_envvalue') envValue = '';
}

export class ServiceDockerfileData extends HttpBase {
  @HttpBind('image_base') imageBase = '';
  @HttpBind('image_author') imageAuthor = '';
  @HttpBind('image_volume') imageVolume: Array<string>;
  @HttpBind('image_run') imageRun: Array<string>;
  @HttpBind('image_expose') imageExpose: Array<string>;
  @HttpBind('image_entrypoint') imageEntryPoint = '';
  @HttpBind('image_cmd') imageCmd = '';
  @HttpBindArray('image_env', ServiceDockerfileEnv) imageEnv: Array<ServiceDockerfileEnv>;
  @HttpBindArray('image_copy', ServiceDockerfileCopy) imageCopy: Array<ServiceDockerfileCopy>;

  protected prepareInit() {
    this.imageVolume = new Array<string>();
    this.imageRun = new Array<string>();
    this.imageExpose = new Array<string>();
    this.imageCopy = new Array<ServiceDockerfileCopy>();
    this.imageEnv = new Array<ServiceDockerfileEnv>();
  }
}

export class ServiceHPA extends HttpBase {
  @HttpBind('hpa_id') hpaId = 0;
  @HttpBind('hpa_name') hpaName = '';
  @HttpBind('hpa_status') hpaStatus = 0;
  @HttpBind('service_id') serviceId = -1;
  @HttpBind('min_pod') minPod = 1;
  @HttpBind('max_pod') maxPod = 1;
  @HttpBind('cpu_percent') cpuPercent = 0;
  isEdit = false;
}

export class PaginationService extends ResponsePaginationBase<Service> {
  ListKeyName(): string {
    return 'service_status_list';
  }

  CreateOneItem(res: object): Service {
    return new Service(res);
  }
}

export class Service extends HttpBase {
  @HttpBind('service_id') serviceId = -1;
  @HttpBind('service_name') serviceName = '';
  @HttpBind('service_project_id') serviceProjectId = -1;
  @HttpBind('service_project_name') serviceProjectName = '';
  @HttpBind('service_owner_id') serviceOwnerId = -1;
  @HttpBind('service_owner_name') serviceOwnerName = '';
  @HttpBind('service_creation_time') serviceCreationTime = '';
  @HttpBind('service_public') servicePublic = 0;
  @HttpBind('service_status') serviceStatus = 0;
  @HttpBind('service_source') serviceSource: ServiceSource = ServiceSource.ServiceSourceBoard;
  @HttpBind('service_is_member') serviceIsMember = -1;
  @HttpBind('service_type') serviceType: ServiceType = ServiceType.ServiceTypeUnknown;
  @HttpBind('service_comment') serviceComment = '';

  get isNotEdgeNode(): boolean {
    return this.serviceType !== ServiceType.ServiceTypeEdgeComputing;
  }

  get isNormalNode(): boolean {
    return this.serviceType === ServiceType.ServiceTypeNormalNodePort;
  }
}

export class NodeAvailableResources extends HttpBase {
  @HttpBind('node_id') id = 0;
  @HttpBind('node_name') name = '';
  @HttpBind('cpu_available') cpuAvailable = '';
  @HttpBind('mem_available') memAvailable = '';
  @HttpBind('storage_available') storageAvailable = '';
}

