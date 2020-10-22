import { DragStatus } from '../shared/shared.types';
import { SERVICE_STATUS } from '../shared/shared.const';
import { HttpBase, HttpBind, HttpBindArray, HttpBindBoolean, HttpBindObject } from '../shared/ui-model/model-types';
import { ServiceType } from './service.types';

export const PHASE_SELECT_PROJECT = 'SELECT_PROJECT';
export const PHASE_CONFIG_CONTAINERS = 'CONFIG_CONTAINERS';
export const PHASE_CONFIG_INIT_CONTAINERS = 'CONFIG_INIT_CONTAINERS';
export const PHASE_EXTERNAL_SERVICE = 'EXTERNAL_SERVICE';
export const PHASE_ENTIRE_SERVICE = 'ENTIRE_SERVICE';
export type ServiceStepPhase =
  'SELECT_PROJECT'
  | 'SELECT_IMAGES'
  | 'CONFIG_CONTAINERS'
  | 'CONFIG_INIT_CONTAINERS'
  | 'EXTERNAL_SERVICE'
  | 'ENTIRE_SERVICE';

export type VolumeType = 'nfs' | 'pvc' | 'configmap';

export enum AffinityCardListView {
  aclvColumn = 'column', aclvRow = 'row'
}

export enum ContainerType {
  runContainer, initContainer
}

export abstract class ServiceStepDataBase extends HttpBase {
  abstract getParams(): { [key: string]: string };
}

export class ImageIndex extends HttpBase {
  @HttpBind('image_name') imageName = '';
  @HttpBind('image_tag') imageTag = '';
  @HttpBind('project_name') projectName = '';
}

export class EnvStruct extends HttpBase {
  @HttpBind('dockerfile_envname') dockerFileEnvName = '';
  @HttpBind('dockerfile_envvalue') dockerFileEnvValue = '';
  @HttpBind('configmap_name') configMapName = '';
  @HttpBind('configmap_key') configMapKey = '';
}

export class Volume extends HttpBase {
  @HttpBind('volume_type') volumeType: VolumeType = 'nfs';
  @HttpBind('target_storage_service') targetStorageService = '';
  @HttpBind('target_path') targetPath = '';
  @HttpBind('volume_name') volumeName = '';
  @HttpBind('container_path') containerPath = '';
  @HttpBind('container_path_flag') containerPathFlagProp = 0;
  @HttpBind('target_pvc') targetPvc = '';
  @HttpBind('target_configmap') targetConfigMap = '';
  @HttpBind('container_file') containerFile = '';
  @HttpBind('targetFile') targetFile = '';

  get containerPathFlag(): boolean {
    return this.containerPathFlagProp === 1;
  }

  set containerPathFlag(value) {
    this.containerPathFlagProp = value ? 1 : 0;
  }
}

export class Container extends HttpBase {
  @HttpBind('name') name = '';
  @HttpBind('working_dir') workingDir = '';
  @HttpBind('command') command = '';
  @HttpBind('cpu_request') cpuRequest = '';
  @HttpBind('mem_request') memRequest = '';
  @HttpBind('cpu_limit') cpuLimit = '';
  @HttpBind('mem_limit') memLimit = '';
  @HttpBind('gpu_limit') gpuLimit = '';
  @HttpBindArray('volume_mounts', Volume) volumeMounts: Array<Volume>;
  @HttpBindArray('env', EnvStruct) env: Array<EnvStruct>;
  @HttpBindObject('image', ImageIndex) image: ImageIndex;
  @HttpBind('container_port') containerPort: Array<number>;

  protected prepareInit() {
    this.image = new ImageIndex();
    this.volumeMounts = Array<Volume>();
    this.env = Array<EnvStruct>();
    this.containerPort = Array<number>();
  }

  get gpuLimitValue(): number {
    return Number(this.gpuLimit);
  }

  set gpuLimitValue(value) {
    this.gpuLimit = `${value}`;
  }
}

export class NodeType extends HttpBase {
  @HttpBind('target_port') targetPort = 0;
  @HttpBind('node_port') nodePort = 0;
}

export class LoadBalance extends HttpBase {
  @HttpBind('external_access') externalAccess = '';
}

export class ExternalService extends HttpBase {
  @HttpBind('container_name') containerName = '';
  @HttpBindObject('node_config', NodeType) nodeConfig: NodeType;
  @HttpBindObject('load_balancer_config', LoadBalance) loadBalance: LoadBalance;

  protected prepareInit() {
    this.nodeConfig = new NodeType();
    this.loadBalance = new LoadBalance();
  }
}

export class AffinityCardData {
  serviceName = '';
  serviceStatus: SERVICE_STATUS;
  status ? = DragStatus.dsReady;

  get key(): string {
    return `${this.serviceName}`;
  }
}

export class ServiceStep1Data extends ServiceStepDataBase {
  @HttpBind('project_id') projectId = -1;
  @HttpBind('project_name') projectName = '';

  getParams(): { [key: string]: string } {
    return {
      phase: PHASE_SELECT_PROJECT,
      project_id: this.projectId.toString()
    };
  }
}

export class ServiceStep2Data extends ServiceStepDataBase {
  @HttpBindArray('container_list', Container) containerList: Array<Container>;
  @HttpBind('project_id') projectId = 0;
  @HttpBind('project_name') projectName = '';

  protected prepareInit() {
    super.prepareInit();
    this.containerList = Array<Container>();
  }

  getParams(): { [p: string]: string } {
    return {
      phase: PHASE_CONFIG_CONTAINERS,
      project_id: this.projectId.toString()
    };
  }

  getPortList(containerName: string): Array<number> {
    return this.containerList.find(value => value.name === containerName).containerPort;
  }

  getPostBody(): any {
    const containers = new Array<any>();
    this.containerList.forEach(value => containers.push(value.getPostBody()));
    return containers;
  }
}

export class ServiceStep2DataInit extends ServiceStep2Data {
  @HttpBindArray('initcontainer_list', Container) containerList: Array<Container>;

  getParams(): { [p: string]: string } {
    return {
      phase: PHASE_CONFIG_INIT_CONTAINERS,
      project_id: this.projectId.toString()
    };
  }

  getPostBody(): { [p: string]: any } {
    const containers = new Array<any>();
    this.containerList.forEach(value => containers.push(value.getPostBody()));
    return containers;
  }
}

export class Affinity extends HttpBase {
  services: Array<AffinityCardData>;
  @HttpBindBoolean('anti_flag', 1, 0) antiFlag = false;
  @HttpBind('service_names') serviceNames: Array<string>;

  protected prepareInit() {
    this.serviceNames = new Array<string>();
    this.services = new Array<AffinityCardData>();
  }

  protected preparePost() {
    this.serviceNames.splice(0, this.serviceNames.length);
    this.services.forEach(value => this.serviceNames.push(value.serviceName));
  }

  protected afterInit() {
    if (this.serviceNames.length > 0) {
      this.serviceNames.forEach(value => {
        const card = new AffinityCardData();
        card.serviceName = value;
        card.status = DragStatus.dsEnd;
        this.services.push(card);
      });
    }
  }
}

export class ServiceStep3Data extends ServiceStepDataBase {
  @HttpBind('service_name') serviceName = '';
  @HttpBind('project_name') projectName = '';
  @HttpBind('instance') instance = 1;
  @HttpBind('cluster_ip') clusterIp = '';
  @HttpBind('service_type') serviceType = ServiceType.ServiceTypeNormalNodePort;
  @HttpBind('node_selector') nodeSelector = '';
  @HttpBindBoolean('service_public', 1, 0) servicePublic = false;
  @HttpBindBoolean('session_affinity_flag', 1, 0) sessionAffinityFlag = false;
  @HttpBindArray('external_service_list', ExternalService) externalServiceList: Array<ExternalService>;
  @HttpBindArray('affinity_list', Affinity) affinityList: Array<Affinity>;
  edgeNodeSelectorIsNode = true;

  protected prepareInit() {
    this.externalServiceList = Array<ExternalService>();
    this.affinityList = Array<Affinity>();
  }

  getPostBody(): { [p: string]: any } {
    const externalServices = new Array<any>();
    const affinityList = new Array<any>();
    this.externalServiceList.forEach(value => externalServices.push(value.getPostBody()));
    this.affinityList.forEach(value => affinityList.push(value.getPostBody()));
    return {
      external_service_list: externalServices,
      affinity_list: affinityList
    };
  }

  getParams(): { [p: string]: string } {
    return {
      phase: PHASE_EXTERNAL_SERVICE,
      service_name: this.serviceName,
      service_type: this.serviceType.toString(),
      instance: this.instance.toString(),
      session_affinity_flag: this.sessionAffinityFlag ? '1' : '0',
      cluster_ip: this.clusterIp,
      service_public: this.servicePublic ? '1' : '0',
      node_selector: this.nodeSelector
    };
  }

  get isShowExternalConfig() {
    return this.serviceType === ServiceType.ServiceTypeNormalNodePort ||
      this.serviceType === ServiceType.ServiceTypeStatefulSet;
  }

  get isShowAdvanceConfig() {
    return this.serviceType === ServiceType.ServiceTypeNormalNodePort ||
      this.serviceType === ServiceType.ServiceTypeStatefulSet;
  }

  get isEdgeComputingType(): boolean {
    return this.serviceType === ServiceType.ServiceTypeEdgeComputing;
  }

  get isStatefulSetType() {
    return this.serviceType === ServiceType.ServiceTypeStatefulSet;
  }

  get isClusterIpType() {
    return this.serviceType === ServiceType.ServiceTypeClusterIP;
  }
}
