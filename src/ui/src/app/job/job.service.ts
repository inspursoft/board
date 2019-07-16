import { Injectable } from "@angular/core";
import { HttpClient, HttpResponse } from "@angular/common/http";
import { Job, JobDeployment, JobPod, LogsSearchConfig, PaginationJob } from "./job.type";
import { Observable, zip } from "rxjs";
import { Project } from "../project/project";
import { map } from "rxjs/operators";
import { Image, ImageDetail } from "../image/image";
import { INode, INodeGroup, NodeAvailableResources, PersistentVolumeClaim } from "../shared/shared.types";

@Injectable()
export class JobService {
  constructor(private http: HttpClient) {

  }

  checkJobNameExists(projectName, jobName: string): Observable<any> {
    return this.http.get(`/api/v1/jobs/exists`, {params: {project_name: projectName, job_name: jobName}})
  }

  deleteJob(job: Job): Observable<any> {
    return this.http.delete(`/api/v1/jobs/${job.job_id}`);
  }

  getJobPods(job: Job): Observable<Array<JobPod>> {
    return this.http.get<Array<JobPod>>(`/api/v1/jobs/${job.job_id}/pods`);
  }

  getJobLogs(job: Job, pod: JobPod, query?: LogsSearchConfig): Observable<any> {
    return this.http.get(`/api/v1/jobs/${job.job_id}/logs/${pod.name}`, {
      responseType: 'text',
      params: {
        container: query && query.container ? query.container : '',
        follow: query && query.follow ? 'true' : 'false',
        previous: query && query.previous ? 'true' : 'false',
        since_seconds: query && query.sinceSeconds ? query.sinceSeconds.toString() : '0',
        since_time: query && query.sinceTime ? query.sinceTime : '',
        timestamps: query && query.timestamps ? 'true' : 'false',
        tail_lines: query && query.tailLines ? query.tailLines.toString() : '0',
        limit_bytes: query && query.limitBytes ? query.limitBytes.toString() : '0'
      }
    });
  }

  getJobList(pageIndex: number, pageSize: number): Observable<PaginationJob> {
    return this.http.get<PaginationJob>(`/api/v1/jobs`, {
      params: {
        job_name: '',
        page_index: pageIndex.toString(),
        page_size: pageSize.toString(),
        order_field: '',
        order_asc: '0'
      }
    });
  }

  getJobStatus(jobId: number): Observable<any> {
    return this.http.get(`/api/v1/jobs/${jobId}/status`);
  }

  getCollaborativeJobs(projectName: string): Observable<Array<Job>> {
    return this.http.get<Array<Job>>(`/api/v1/jobs/selectjobs`, {
      params: {
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

  getImageList(): Observable<Array<Image>> {
    return this.http.get<Array<Image>>("/api/v1/images");
  }

  getImageDetailList(image_name: string): Observable<Array<ImageDetail>> {
    return this.http.get<Array<ImageDetail>>(`/api/v1/images/${image_name}`);
  }

  getNodesAvailableSources(): Observable<Array<NodeAvailableResources>> {
    return this.http.get(`/api/v1/nodes/availableresources`, {
      observe: "response"
    }).pipe(map((res: HttpResponse<Array<NodeAvailableResources>>) => res.body));
  }

  getPvcNameList(): Observable<Array<PersistentVolumeClaim>> {
    return this.http.get(`/api/v1/pvclaims`, {observe: "response"})
      .pipe(map((res: HttpResponse<Array<Object>>) => {
        let result: Array<PersistentVolumeClaim> = Array<PersistentVolumeClaim>();
        res.body.forEach(resObject => {
          let persistentVolume = new PersistentVolumeClaim();
          persistentVolume.id = Reflect.get(resObject, 'pvc_id');
          persistentVolume.name = Reflect.get(resObject, 'pvc_name');
          result.push(persistentVolume);
        });
        return result;
      }));
  }

  getNodeSelectors(): Observable<Array<{name: string, status: number}>> {
    let obsNodeList = this.http
      .get(`/api/v1/nodes`, {observe: "response"})
      .pipe(map((res: HttpResponse<Array<INode>>) => {
        let result = Array<{name: string, status: number}>();
        res.body.forEach((iNode: INode) => result.push({name: String(iNode.node_name).trim(), status: iNode.status}));
        return result;
      }));
    let obsNodeGroupList = this.http
      .get(`/api/v1/nodegroup`, {observe: "response", params: {is_valid_node_group: '1'}})
      .pipe(map((res: HttpResponse<Array<INodeGroup>>) => {
        let result = Array<{name: string, status: number}>();
        res.body.forEach((iNodeGroup: INodeGroup) => result.push({name: String(iNodeGroup.nodegroup_name).trim(), status: 1}));
        return result;
      }));
    return zip(obsNodeList, obsNodeGroupList).pipe(
      map(value => value[0].concat(value[1]))
    );
  }

  deploymentJob(jobDeployment: JobDeployment): Observable<any> {
    return this.http.post(`/api/v1/jobs/deployment`, jobDeployment)
  }
}
