import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { Observable, zip } from 'rxjs';
import { map } from 'rxjs/operators';
import { Project } from '../project/project';
import {
  Job,
  JobDeployment,
  JobImageDetailInfo,
  JobImageInfo, JobNode,
  JobNodeAvailableResources, JobNodeGroup,
  JobPod,
  LogsSearchConfig,
  PaginationJob
} from './job.type';
import { PersistentVolumeClaim } from '../shared/shared.types';
import { ModelHttpClient } from '../shared/ui-model/model-http-client';

@Injectable()
export class JobService {
  constructor(private http: HttpClient,
              private modeHttp: ModelHttpClient) {

  }

  checkJobNameExists(projectName, jobName: string): Observable<any> {
    return this.modeHttp.get(`/api/v1/jobs/exists`, {params: {project_name: projectName, job_name: jobName}});
  }

  deleteJob(job: Job): Observable<any> {
    return this.modeHttp.delete(`/api/v1/jobs/${job.jobId}`);
  }

  getJobPods(job: Job): Observable<Array<JobPod>> {
    return this.modeHttp.getArray(`/api/v1/jobs/${job.jobId}/pods`, JobPod);
  }

  getJobLogs(job: Job, pod: JobPod, query?: LogsSearchConfig): Observable<any> {
    return this.http.get(`/api/v1/jobs/${job.jobId}/logs/${pod.name}`, {
      responseType: 'text',
      params: {
        container: query && query.container ? query.container : '',
        follow: query && query.follow ? 'true' : 'false',
        previous: query && query.previous ? 'true' : 'false',
        since_time: query && query.sinceTime ? query.sinceTime : '',
        timestamps: query && query.timestamps ? 'true' : 'false',
        limit_bytes: query && query.limitBytes ? query.limitBytes.toString() : '0'
      }
    });
  }

  getJobList(pageIndex: number, pageSize: number): Observable<PaginationJob> {
    return this.modeHttp.getPagination(`/api/v1/jobs`, PaginationJob, {
        param: {
          job_name: '',
          page_index: pageIndex.toString(),
          page_size: pageSize.toString(),
          order_field: '',
          order_asc: '0'
        }
      }
    );
  }

  getJobStatus(jobId: number): Observable<any> {
    return this.http.get(`/api/v1/jobs/${jobId}/status`);
  }

  getCollaborativeJobs(projectName: string): Observable<Array<Job>> {
    return this.modeHttp.getArray(`/api/v1/jobs/selectjobs`, Job, {
      param: {
        project_name: projectName
      }
    });
  }

  getProjectList(): Observable<Array<Project>> {
    return this.http.get<Array<Project>>('/api/v1/projects');
  }

  getOneProject(projectName: string): Observable<Project> {
    return this.http.get<Array<Project>>('/api/v1/projects', {
      params: {project_name: projectName}
    }).pipe(map(res => res && res.length > 0 ? res[0] : null));
  }

  getImageList(): Observable<Array<JobImageInfo>> {
    return this.modeHttp.getArray('/api/v1/images', JobImageInfo);
  }

  getImageDetailList(imageName: string): Observable<Array<JobImageDetailInfo>> {
    return this.modeHttp.getArray(`/api/v1/images/${imageName}`, JobImageDetailInfo);
  }

  getNodesAvailableSources(): Observable<Array<JobNodeAvailableResources>> {
    return this.modeHttp.getArray(`/api/v1/nodes/availableresources`, JobNodeAvailableResources);
  }

  getPvcNameList(): Observable<Array<PersistentVolumeClaim>> {
    return this.http.get(`/api/v1/pvclaims`, {observe: 'response'})
      .pipe(map((res: HttpResponse<Array<object>>) => {
        const result: Array<PersistentVolumeClaim> = Array<PersistentVolumeClaim>();
        res.body.forEach(resObject => {
          const persistentVolume = new PersistentVolumeClaim();
          persistentVolume.id = Reflect.get(resObject, 'pvc_id');
          persistentVolume.name = Reflect.get(resObject, 'pvc_name');
          result.push(persistentVolume);
        });
        return result;
      }));
  }

  getNodeSelectors(): Observable<Array<{ name: string, status: number }>> {
    const obsNodeList = this.modeHttp
      .getArray(`/api/v1/nodes`, JobNode)
      .pipe(map((res: Array<JobNode>) => {
        const result = Array<{ name: string, status: number }>();
        res.forEach((jobNode: JobNode) => result.push(
          {name: String(jobNode.nodeName).trim(), status: jobNode.status})
        );
        return result;
      }));
    const obsNodeGroupList = this.modeHttp
      .getArray(`/api/v1/nodegroup`, JobNodeGroup, {param: {is_valid_node_group: '1'}})
      .pipe(map((res: Array<JobNodeGroup>) => {
        const result = Array<{ name: string, status: number }>();
        res.forEach((jobNodeGroup: JobNodeGroup) => result.push(
          {name: String(jobNodeGroup.nodeGroupName).trim(), status: 1})
        );
        return result;
      }));
    return zip(obsNodeList, obsNodeGroupList).pipe(
      map(value => value[0].concat(value[1]))
    );
  }

  deploymentJob(jobDeployment: JobDeployment): Observable<any> {
    return this.http.post(`/api/v1/jobs/deployment`, jobDeployment.getPostBody());
  }

  getJobConfig(jobId: number): Observable<JobDeployment> {
    return this.modeHttp.getJson(`/api/v1/jobs/${jobId}/config`, JobDeployment);
  }
}
