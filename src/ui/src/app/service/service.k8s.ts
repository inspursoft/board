import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';
import { AppInitService } from "../app.init.service";
import { Http, Headers, RequestOptions, Response, RequestOptionsArgs } from "@angular/http";
import { Project } from "../project/project";
import { BuildImageDockerfileData, Image, ImageDetail } from "../image/image";
import { ServerServiceStep, ServiceStepPhase, UiServiceFactory, UIServiceStepBase } from "./service-step.component";

@Injectable()
export class K8sService {
  stepSource: Subject<{index: number, isBack: boolean}> = new Subject<{index: number, isBack: boolean}>();
  step$: Observable<{index: number, isBack: boolean}> = this.stepSource.asObservable();

  get defaultHeader(): Headers {
    let headers = new Headers();
    headers.append('Content-Type', 'application/json');
    headers.append('token', this.appInitService.token);
    return headers;
  }

  constructor(private http: Http,
              private appInitService: AppInitService) {
  }

  cancelBuildService(): void {
    this.deleteServiceConfig()
      .then(isDelete => {
        this.stepSource.next({index: 0, isBack: false});
      })
      .catch(() => {
      });
  }

  checkServiceExist(projectName: string, serviceName: string): Promise<boolean> {
    return this.http.get(`/api/v1/services/exists`, {
      headers: this.defaultHeader,
      params: {project_name: projectName, service_name: serviceName}
    }).toPromise()
      .then((res: Response) => {
        this.appInitService.chainResponse(res);
        return res.text().toLocaleLowerCase() == "true";
      })
      .catch(err => Promise.reject(err));
  }

  getServiceConfig(phase: ServiceStepPhase): Promise<UIServiceStepBase> {
    return this.http.get(`/api/v1/services/config`, {
      headers: this.defaultHeader,
      params: {phase: phase}
    }).toPromise()
      .then((res: Response) => {
        this.appInitService.chainResponse(res);
        let stepBase = UiServiceFactory.getInstance(phase);
        return stepBase.serverToUi(res.json());
      })
      .catch(err => Promise.reject(err));
  }

  setServiceConfig(config: ServerServiceStep): Promise<any> {
    let option: RequestOptionsArgs = {
      headers: this.defaultHeader,
      params: {
        phase: config.phase,
        project_id: config.project_id,
        service_name: config.service_name,
        instance: config.instance
      }
    };
    return this.http.post(`/api/v1/services/config`, config.postData, option)
      .toPromise()
      .then((res: Response) => this.appInitService.chainResponse(res))
      .catch(err => Promise.reject(err));
  }

  deleteServiceConfig(): Promise<any> {
    return this.http
      .delete(`/api/v1/services/config`, {headers: this.defaultHeader})
      .toPromise()
      .then((res: Response) => {
        this.appInitService.chainResponse(res);
        console.log(res);
      })
      .catch(err => Promise.reject(err));
  }

  deleteDeployment(serviceId: number): Promise<any> {
    return this.http
      .delete(`/api/v1/services/${serviceId}/deployment`, {headers: this.defaultHeader})
      .toPromise()
      .then((res: Response) => this.appInitService.chainResponse(res))
      .catch(err => Promise.reject(err));
  }

  serviceDeployment(): Promise<number> {
    return this.http.post(`/api/v1/services/deployment`, {}, {
      headers: this.defaultHeader
    }).toPromise()
      .then((res: Response) => {
        this.appInitService.chainResponse(res);
        let resJson = res.json();
        return resJson["project_id"]
      })
      .catch(err => Promise.reject(err));
  }

  getContainerDefaultInfo(image_name: string, image_tag: string, project_name: string): Promise<BuildImageDockerfileData> {
    return this.http.get(`/api/v1/images/dockerfile`, {
      headers: this.defaultHeader,
      params: {image_name: image_name, project_name: project_name, image_tag: image_tag}
    }).toPromise()
      .then((res: Response) => {
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err => Promise.reject(err));
  }

  getProjects(projectName?: string): Promise<Project[]> {
    return this.http.get('/api/v1/projects', {
      headers: this.defaultHeader,
      params: {'project_name': projectName}
    }).toPromise()
      .then((res: Response) => {
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err => Promise.reject(err));
  }

  getDeployStatus(serviceName: string): Promise<Object> {
    let options = new RequestOptions({
      headers: this.defaultHeader
    });
    return this.http.get(`/api/v1/services/status/${serviceName}`, options).toPromise()
      .then((res: Response) => {
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
      .then((res: Response) => {
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
      .then((res: Response) => {
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err => Promise.reject(err));
  }

  getServices(pageIndex?: number, pageSize?: number): Promise<any> {
    return this.http
      .get(`/api/v1/services`, {
        headers: this.defaultHeader,
        params: {
          'page_index': pageIndex,
          'page_size': pageSize
        }
      })
      .toPromise()
      .then((res: Response) => {
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err => Promise.reject(err));
  }

  getServiceDetail(serviceName: string): Promise<any> {
    return this.http
      .get(`/api/v1/services/info/${serviceName}`, {headers: this.defaultHeader})
      .toPromise()
      .then((res: Response) => {
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err => Promise.reject(err));
  }

  deleteService(serviceID: number): Promise<any> {
    return this.http
      .delete(`/api/v1/services/${serviceID}`, {headers: this.defaultHeader})
      .toPromise()
      .then(res => this.appInitService.chainResponse(res))
      .catch(err => Promise.reject(err));
  }

  toggleServiceStatus(serviceID: number, isStart: 0 | 1): Promise<any> {
    return this.http
      .put(`/api/v1/services/${serviceID}/toggle`, {service_toggle: isStart}, {headers: this.defaultHeader})
      .toPromise()
      .then(res => this.appInitService.chainResponse(res))
      .catch(err => Promise.reject(err));
  }

  toggleServicePublicity(serviceID: number, service_togglable: 0 | 1): Promise<any> {
    return this.http
      .put(`/api/v1/services/${serviceID}/publicity`, {service_public: service_togglable}, {headers: this.defaultHeader})
      .toPromise()
      .then(res => this.appInitService.chainResponse(res))
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
      .then((res: Response) => {
        this.appInitService.chainResponse(res);
        return res.text();
      })
      .catch(err => Promise.reject(err));
  }

  getNodesList(): Promise<any> {
    return this.http
      .get(`/api/v1/nodes`, {headers: this.defaultHeader})
      .toPromise()
      .then((res: Response) => {
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err => Promise.reject(err));
  }

  addServiceRoute(serviceURL: string, serviceIdentity: string): Promise<any> {
    return this.http
      .post(`/api/v1/services/info`, {}, {
        headers: this.defaultHeader,
        params: {
          'service_url': serviceURL,
          'service_identity': serviceIdentity
        }
      })
      .toPromise()
      .then(res => this.appInitService.chainResponse(res))
      .catch(err => Promise.reject(err));
  }

}