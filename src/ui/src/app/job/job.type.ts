import { DragStatus } from '../shared/shared.types';
import { HttpBase, HttpBind, HttpBindArray, HttpBindObject, ResponsePaginationBase } from '../shared/ui-model/model-types';

export enum CreateMethod {
  byDefault, byExistsJob
}

export class PaginationJob extends ResponsePaginationBase<Job> {

  ListKeyName(): string {
    return 'job_status_list';
  }

  CreateOneItem(res: object): Job {
    return new Job(res);
  }
}

export class Job extends HttpBase {
  @HttpBind('job_id') jobId: number;
  @HttpBind('job_name') jobName = '';
  @HttpBind('job_project_id') jobProjectId: number;
  @HttpBind('job_project_name') jobProjectName: string;
  @HttpBind('job_comment') jobComment: string;
  @HttpBind('job_creation_time') jobCreationTime: string;
  @HttpBind('job_update_time') jobUpdateTime: string;
  @HttpBind('job_deleted') jobDeleted: number;
  @HttpBind('job_owner_id') jobOwnerId: number;
  @HttpBind('job_owner_name') jobOwnerName: string;
  @HttpBind('job_source') jobSource: number;
  @HttpBind('job_status') jobStatus: number;
  @HttpBind('job_yaml') jobYaml: string;
}

export class JobVolumeMounts extends HttpBase {
  @HttpBind('volume_type') volumeType = '';
  @HttpBind('volume_name') volumeName = '';
  @HttpBind('container_path') containerPath = '';
  @HttpBind('container_file') containerFile = '';
  @HttpBind('container_path_flag') containerPathFlag = 0;
  @HttpBind('target_storage_service') targetStorageService = '';
  @HttpBind('target_path') targetPath = '';
  @HttpBind('target_file') targetFile = '';
  @HttpBind('target_pvc') targetPvc = '';
}

export class JobImage extends HttpBase {
  @HttpBind('image_name') imageName = '';
  @HttpBind('image_tag') imageTag = '';
  @HttpBind('project_name') projectName = '';
}

export class JobEnv extends HttpBase {
  @HttpBind('dockerfile_envname') dockerfileEnvName = '';
  @HttpBind('dockerfile_envvalue') dockerfileEnvValue = '';
  @HttpBind('configmap_key') configMapKey = '';
  @HttpBind('configmap_name') configMapName = '';
}

export class JobContainer extends HttpBase {
  @HttpBind('name') name = '';
  @HttpBind('working_dir') workingDir = '';
  @HttpBind('command') command = '';
  @HttpBind('container_port') containerPort: Array<number>;
  @HttpBind('cpu_request') cpuRequest = '';
  @HttpBind('mem_request') memRequest = '';
  @HttpBind('cpu_limit') cpuLimit = '';
  @HttpBind('mem_limit') memLimit = '';
  @HttpBind('gpu_limit') gpuLimit = '';
  @HttpBindObject('image', JobImage) image: JobImage;
  @HttpBindArray('volume_mounts', JobVolumeMounts) volumeMounts: Array<JobVolumeMounts>;
  @HttpBindArray('env', JobEnv) env: Array<JobEnv>;

  protected prepareInit() {
    super.prepareInit();
    this.containerPort = Array<number>();
    this.volumeMounts = Array<JobVolumeMounts>();
    this.env = Array<JobEnv>();
  }

  get gpuLimitValue(): number {
    return Number(this.gpuLimit);
  }

  set gpuLimitValue(value) {
    this.gpuLimit = `${value}`;
  }
}

export class JobAffinity extends HttpBase {
  @HttpBind('anti_flag') antiFlag = 0;
  @HttpBind('job_names') jobNames: Array<string>;

  protected prepareInit() {
    super.prepareInit();
    this.jobNames = Array<string>();
  }
}

export class JobDeployment extends HttpBase {
  @HttpBind('project_id') projectId = 0;
  @HttpBind('project_name') projectName = '';
  @HttpBind('job_id') jobId = 0;
  @HttpBind('job_name') jobName = '';
  @HttpBind('node_selector') nodeSelector = '';
  @HttpBindArray('container_list', JobContainer) containerList: Array<JobContainer>;
  @HttpBindArray('affinity_list', JobAffinity) affinityList: Array<JobAffinity>;
  @HttpBind('parallelism') parallelism = 1;
  @HttpBind('completions') completions = 1;
  @HttpBind('active_Deadline_Seconds') activeDeadlineSeconds: number;
  @HttpBind('backoff_Limit') backOffLimit = 6;

  protected prepareInit() {
    super.prepareInit();
    this.containerList = Array<JobContainer>();
    this.affinityList = Array<JobAffinity>();
  }
}


export class JobAffinityCardData {
  jobName = '';
  status ? = DragStatus.dsReady;

  get key(): string {
    return `${this.jobName}`;
  }
}

export enum JobAffinityCardListView {
  aclvColumn = 'column', aclvRow = 'row'
}

export class JobPod extends HttpBase {
  @HttpBind('name') name: string;
  @HttpBind('project_name') projectName: string;
}

export class LogsSearchConfig {
  container?: string;
  follow?: boolean;
  previous?: boolean;
  sinceSeconds?: number;
  sinceTime?: string;
  timestamps?: boolean;
  tailLines?: number;
  limitBytes?: number;
}

export class JobImageInfo extends HttpBase {
  @HttpBind('image_name') imageName: string;
  @HttpBind('image_comment') imageComment: string;
  @HttpBind('image_deleted') imageDeleted: number;
  @HttpBind('image_update_time') imageUpdateTime: string;
  @HttpBind('image_creation_time') imageCreationTime: string;
}

export class JobImageDetailInfo extends HttpBase {
  @HttpBind('image_name') imageName: string;
  @HttpBind('image_tag') imageTag: string;
  @HttpBind('image_detail') imageDetail: string;
  @HttpBind('image_creation_time') imageCreationTime: string;
  @HttpBind('image_size_number') imageSizeNumber: number;
  @HttpBind('image_size_unit') imageSizeUnit: string;
}

export class JobNodeAvailableResources extends HttpBase {
  @HttpBind('node_id') nodeId = 0;
  @HttpBind('node_name') nodeName = '';
  @HttpBind('cpu_available') cpuAvailable = '';
  @HttpBind('mem_available') memAvailable = '';
  @HttpBind('storage_available') storageAvailable = '';
}

export class JobNode extends HttpBase {
  @HttpBind('node_name') nodeName = '';
  @HttpBind('node_ip') nodeIp = '';
  @HttpBind('status') status = -1;
}

export class JobNodeGroup extends HttpBase {
  @HttpBind('nodegroup_id') nodeGroupId = -1;
  @HttpBind('nodegroup_project') nodeGroupProject = '';
  @HttpBind('nodegroup_name') nodeGroupName = '';
  @HttpBind('nodegroup_comment') nodeGroupComment = '';
}
