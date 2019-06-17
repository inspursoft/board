import { DragStatus, IPagination } from "../shared/shared.types";

export class PaginationJob {
  pagination: IPagination;
  job_status_list: Array<Job>;

  constructor() {
    this.job_status_list = Array<Job>();
    this.pagination = {page_count: 1, page_index: 1, page_size: 15, total_count: 0}
  }
}

export class Job {
  job_id: number;
  job_name: string;
  job_project_id: number;
  job_project_name: string;
  job_comment: string;
  job_creation_time: string;
  job_update_time: string;
  job_deleted: number;
  job_owner_id: number;
  job_owner_name: string;
  job_source: number;
  job_status: number;
  job_yaml: string;
}

export class JobVolumeMounts {
  volume_type: string;
  volume_name: string;
  container_path: string;
  container_file: string;
  container_path_flag: number;
  target_storage_service: string;
  target_path: string;
  target_file: string;
  target_pvc: string;
}

export class JobImage {
  image_name: string;
  image_tag: string;
  project_name: string;
}

export class JobEnv {
  dockerfile_envname: string;
  dockerfile_envvalue: string;
  configmap_key: string;
  configmap_name: string;
}

export class JobContainer {
  name: string;
  working_Dir: string;
  command: string;
  container_port: Array<number>;
  cpu_request: string;
  mem_request: string;
  cpu_limit: string;
  mem_limit: string;
  volume_mounts: Array<JobVolumeMounts>;
  image: JobImage;
  env: Array<JobEnv>;

  constructor() {
    this.container_port = Array<number>();
    this.volume_mounts = Array<JobVolumeMounts>();
    this.env = Array<JobEnv>();
  }
}

export class JobAffinity {
  anti_flag: number;
  job_names: Array<string>;

  constructor() {
    this.job_names = Array<string>();
  }
}

export class JobDeployment {
  project_id = 0;
  project_name: string;
  job_id: number;
  job_name: string;
  node_selector: string;
  container_list: Array<JobContainer>;
  affinity_list: Array<JobAffinity>;
  parallelism = 1;
  completions = 1;
  active_Deadline_Seconds = 1;
  backoff_Limit = 6;

  constructor() {
    this.container_list = Array<JobContainer>();
    this.affinity_list = Array<JobAffinity>();
  }
}


export class JobAffinityCardData {
  jobName = '';
  status? = DragStatus.dsReady;

  get key(): string {
    return `${this.jobName}`
  }
}

export enum JobAffinityCardListView {
  aclvColumn = 'column', aclvRow = 'row'
}

export class JobPod{
  name: string;
  project_name: string;
  spec: Array<JobContainer>;
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
