import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';
import { Service } from './service';
import { AppInitService } from "../app.init.service";
import { Http, Headers, RequestOptions, Response } from "@angular/http";
import { Project } from "../project/project";
import { Image, ImageDetail } from "../image/image";
import {
  ImageDockerfile, ServiceStep2NewImageType, ServiceStep4Output,
  ServiceStep6Output
} from "./service-step.component";

@Injectable()
export class K8sService {
  stepSource: Subject<number> = new Subject<number>();
  step$: Observable<number> = this.stepSource.asObservable();
  stepData: Map<number, Object>;

  get defaultHeader(): Headers {
    let headers = new Headers();
    headers.append('Content-Type', 'application/json');
    headers.append('token', this.appInitService.token);
    return headers;
  }

  constructor(private http: Http,
              private appInitService: AppInitService) {
    this.stepData = new Map<number, Object>();
  }

  clearStepData() {
    this.stepData.clear();
  }

  setStepData(step: number, Data: Object) {
    this.stepData.set(step, Data);
    console.log(Data);
  }

  buildImage(imageData: ServiceStep2NewImageType): Promise<boolean> {
    return this.http.post(`/api/v1/images/building`, imageData, {
      headers: this.defaultHeader
    }).toPromise()
      .then(resp => {
        this.appInitService.chainResponse(resp);
        return resp.status == 200;
      })
      .catch(err => Promise.reject(err));
  }

  getContainerDefaultInfo(image_name: string, image_tag: string, project_name: string): Promise<ImageDockerfile> {
    return this.http.get(`/api/v1/images/dockerfile`, {
      headers: this.defaultHeader,
      params: {image_name: image_name, project_name: project_name, image_tag: image_tag}
    }).toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err => Promise.reject(err));
  }

  serviceDeployment(postData: ServiceStep4Output): Promise<ServiceStep6Output> {
    return this.http.post(`/api/v1/services/${postData.service_id}/deployment`, postData, {
      headers: this.defaultHeader
    }).toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err => Promise.reject(err));
  }

  getDockerFilePreview(imageData: ServiceStep2NewImageType): Promise<string> {
    return this.http.post(`/api/v1/images/preview`, imageData, {
      headers: this.defaultHeader
    }).toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.text();
      })
      .catch(err => Promise.reject(err));
  }

  getFileList(formData: FormData): Promise<Array<{path: string, file_name: string, size: number}>> {
    let headers = new Headers();
    headers.append('token', this.appInitService.token);
    let options = new RequestOptions({headers: headers});
    return this.http.post(`/api/v1/files/list`, formData, options).toPromise()
      .then(resp => {
        this.appInitService.chainResponse(resp);
        return resp.json();
      })
      .catch(err => Promise.reject(err));
  }

  removeFile(formData: FormData): Promise<boolean> {
    let headers = new Headers();
    headers.append('token', this.appInitService.token);
    let options = new RequestOptions({headers: headers});
    return this.http.post(`/api/v1/files/remove`, formData, options).toPromise()
      .then(resp => {
        this.appInitService.chainResponse(resp);
        return resp.status == 200;
      })
      .catch(err => Promise.reject(err));
  }

  uploadFile(formData: FormData): Promise<boolean> {
    let headers = new Headers();
    headers.append('token', this.appInitService.token);
    let options = new RequestOptions({headers: headers});
    return this.http.post(`/api/v1/files/upload`, formData, options).toPromise()
      .then(resp => {
        this.appInitService.chainResponse(resp);
        return resp.status == 200;
      })
      .catch(err => Promise.reject(err));
  }

  getServiceID(postData: {project_name: string, project_id: number}) {
    return this.http.post(`/api/v1/services`, postData, {
      headers: this.defaultHeader
    }).toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err => Promise.reject(err));
  }

  getStepData(step: number): Object {
    return this.stepData.get(step);
  }

  getProjects(projectName?: string): Promise<Project[]> {
    return this.http.get('/api/v1/projects', {
      headers: this.defaultHeader,
      params: {'project_name': projectName}
    }).toPromise()
      .then(resp => {
        this.appInitService.chainResponse(resp);
        return resp.json();
      })
      .catch(err => Promise.reject(err));
  }

  getDeployStatus(serviceName: string): Promise<Object> {
    let options = new RequestOptions({
      headers: this.defaultHeader
    });
    return this.http.get(`/api/v1/services/status/${serviceName}`, options).toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err => Promise.reject(err));
  }

  getImages(image_name?: string, image_list_page?: number, image_list_page_size?: number): Promise<Image[]> {
    let options = new RequestOptions({
      headers: this.defaultHeader,
      params: {
        'image_name': image_name,
        'image_list_page': image_list_page.toString(),
        'image_list_page_size': image_list_page_size.toString()
      }
    });
    return this.http.get("/api/v1/images", options).toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err => Promise.reject(err));
  }

  getImageDetailList(image_name: string): Promise<ImageDetail[]> {
    let options = new RequestOptions({
      headers: this.defaultHeader
    });
    return this.http.get(`/api/v1/images/${image_name}`, options)
      .toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err => Promise.reject(err));
  }

  getServices(): Promise<Service[]> {
    return this.http
      .get(`/api/v1/services`, {headers: this.defaultHeader})
      .toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return <Service[]>res.json();
      })
      .catch(err => Promise.reject(err));
  }

  getServiceDetail(serviceName: string): Promise<Object> {
    return this.http
      .get(`/api/v1/services/info/${serviceName}`, {headers: this.defaultHeader})
      .toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err => Promise.reject(err));
  }

  deleteService(serviceID: number): Promise<any> {
    return this.http
      .delete(`/api/v1/services/${serviceID}`, {headers: this.defaultHeader})
      .toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res;
      })
      .catch(err => Promise.reject(err));
  }

  toggleServiceStatus(serviceID: number, isStart: 0 | 1): Promise<any> {
    return this.http
      .put(`/api/v1/services/${serviceID}/toggle`, {service_toggle: isStart}, {headers: this.defaultHeader})
      .toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res;
      })
      .catch(err => Promise.reject(err));
  }

  toggleServicePublicity(serviceID: number, service_togglable: 0 | 1): Promise<any> {
    return this.http
      .put(`/api/v1/services/${serviceID}/publicity`, {service_public: service_togglable}, {headers: this.defaultHeader})
      .toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res;
      })
      .catch(err => Promise.reject(err));
  }

  getConsole(jobName: string, buildSerialId?: string): Promise<string> {
    return this.http
      .get(`/api/v1/jenkins-job/console`, {
        headers: this.defaultHeader,
        params: {
          "job_name": jobName,
          "build_serial_id": buildSerialId
        }
      })
      .toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.text();
      })
      .catch(err => Promise.reject(err));
  }

  getLastJobId(jobName: string): Promise<number> {
    return this.http
      .get(`/api/v1/jenkins-job/lastbuildnumber`, {
        headers: this.defaultHeader,
        params: {"job_name": jobName}
      })
      .toPromise()
      .then((res: Response) => {
        this.appInitService.chainResponse(res);
        return Number(res.text());
      })
      .catch(err => Promise.reject(err));
  }

  cancelConsole(jobName: string, buildSerialId: number): Promise<boolean> {
    return this.http
      .get(`/api/v1/jenkins-job/stop`, {
        headers: this.defaultHeader,
        params: {
          "job_name": jobName,
          "build_serial_id": buildSerialId
        }
      }).toPromise()
      .then((res: Response) => {
        this.appInitService.chainResponse(res);
        return res.status == 200;
      })
      .catch(err => Promise.reject(err));
  }

  getNodesList(): Promise<any> {
    return this.http
      .get(`/api/v1/nodes`, {headers: this.defaultHeader})
      .toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err => Promise.reject(err));
  }


}