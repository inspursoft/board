import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';
import { HttpClient, HttpParams, HttpResponse } from "@angular/common/http";
import { Project } from "../project/project";
import { BuildImageDockerfileData, Image, ImageDetail } from "../image/image";
import { ImageIndex, ServerServiceStep, ServiceStepPhase, UiServiceFactory, UIServiceStepBase } from "./service-step.component";
import { Service } from "./service";

@Injectable()
export class K8sService {
  stepSource: Subject<{index: number, isBack: boolean}> = new Subject<{index: number, isBack: boolean}>();
  step$: Observable<{index: number, isBack: boolean}> = this.stepSource.asObservable();

  constructor(private http: HttpClient) {
  }

  cancelBuildService(): void {
    this.deleteServiceConfig()
      .then(isDelete => {
        this.stepSource.next({index: 0, isBack: false});
      })
      .catch(() => {
      });
  }

  checkServiceExist(projectName: string, serviceName: string): Promise<any> {
    return this.http.get(`/api/v1/services/exists`, {
      observe: "response",
      params: {project_name: projectName, service_name: serviceName}
    }).toPromise();
  }

  getServiceConfig(phase: ServiceStepPhase): Promise<UIServiceStepBase> {
    return this.http.get(`/api/v1/services/config`, {
      observe: "response",
      params: {phase: phase}
    }).toPromise()
      .then((res: HttpResponse<Object>) => {
        let stepBase = UiServiceFactory.getInstance(phase);
        return stepBase.serverToUi(res.body);
      })
      .catch(err => Promise.reject(err));
  }

  setServiceConfig(config: ServerServiceStep): Promise<any> {
    return this.http.post(`/api/v1/services/config`, config.postData, {
      observe: "response",
      params: {
        phase: config.phase,
        project_id: config.project_id.toString(),
        service_name: config.service_name,
        instance: config.instance.toString()
      }
    }).toPromise()
  }

  deleteServiceConfig(): Promise<any> {
    return this.http.delete(`/api/v1/services/config`, {observe: "response"}).toPromise()
  }

  deleteDeployment(serviceId: number): Promise<any> {
    return this.http.delete(`/api/v1/services/${serviceId}/deployment`, {observe: "response"}).toPromise()
  }

  serviceDeployment(): Promise<number> {
    return this.http.post(`/api/v1/services/deployment`, {}, {observe: "response"})
      .toPromise()
      .then((res: HttpResponse<Object>) => res.body["project_id"])
  }

  getContainerDefaultInfo(image_name: string, image_tag: string, project_name: string): Promise<BuildImageDockerfileData> {
    return this.http.get<BuildImageDockerfileData>(`/api/v1/images/dockerfile`, {
      observe: "response",
      params: {image_name: image_name, project_name: project_name, image_tag: image_tag}
    }).toPromise()
      .then((res: HttpResponse<BuildImageDockerfileData>) => res.body);
  }

  getProjects(projectName: string = ""): Promise<Project[]> {
    return this.http.get<Project[]>('/api/v1/projects', {
      observe: "response",
      params: {'project_name': projectName, 'member_only': "1"}
    }).toPromise()
      .then((res: HttpResponse<Project[]>) => res.body)
  }

  getDeployStatus(serviceId: number): Promise<any> {
    return this.http.get(`/api/v1/services/${serviceId}/status`, {observe: "response"}).toPromise();
  }

  getImages(image_name?: string, image_list_page?: number, image_list_page_size?: number): Promise<Image[]> {
    return this.http.get("/api/v1/images", {
      observe: "response",
      params: {
        'image_name': image_name,
        'image_list_page': image_list_page.toString(),
        'image_list_page_size': image_list_page_size.toString()
      }
    }).toPromise()
      .then((res: HttpResponse<Image[]>) => res.body || []);
  }

  getImageDetailList(image_name: string): Promise<ImageDetail[]> {
    return this.http.get(`/api/v1/images/${image_name}`, {observe: "response"})
      .toPromise()
      .then((res: HttpResponse<ImageDetail[]>) => res.body || []);
  }

  getServices(pageIndex: number, pageSize: number, sortBy: string, isReverse: boolean): Promise<Object> {
    return this.http
      .get(`/api/v1/services`, {
        observe: "response", params: {
          "page_index": pageIndex.toString(),
          "page_size": pageSize.toString(),
          "order_field": sortBy,
          "order_asc": isReverse ? "0" : "1"
        }
      })
      .toPromise()
      .then((res: HttpResponse<Object>) => res.body);
  }

  getServiceDetail(serviceId: number): Promise<Object> {
    return this.http
      .get(`/api/v1/services/${serviceId}/info`, {observe: "response"})
      .toPromise()
      .then((res: HttpResponse<Object>) => res.body)
      .catch(err => Promise.reject(err));
  }

  deleteService(serviceID: number): Promise<any> {
    return this.http.delete(`/api/v1/services/${serviceID}`, {observe: "response"}).toPromise()
  }

  toggleServiceStatus(serviceID: number, isStart: 0 | 1): Promise<any> {
    return this.http.put(`/api/v1/services/${serviceID}/toggle`, {service_toggle: isStart}, {observe: "response"}).toPromise();
  }

  toggleServicePublicity(serviceID: number, service_togglable: 0 | 1): Promise<any> {
    return this.http.put(`/api/v1/services/${serviceID}/publicity`, {service_public: service_togglable}, {observe: "response"}).toPromise();
  }

  getConsole(jobName: string, buildSerialId?: string): Promise<string> {
    return this.http
      .get(`/api/v1/jenkins-job/console`, {
        observe: "response",
        responseType: "text",
        params: {
          "job_name": jobName,
          "build_serial_id": buildSerialId
        }
      })
      .toPromise()
      .then((res: HttpResponse<string>) => res.body);
  }

  getNodesList(): Promise<Array<{node_name: string, node_ip: string, status: number}>> {
    return this.http
      .get(`/api/v1/nodes`, {observe: "response"})
      .toPromise()
      .then((res: HttpResponse<Array<{node_name: string, node_ip: string, status: number}>>) => res.body || [])
  }

  addServiceRoute(serviceURL: string, serviceIdentity: string): Promise<any> {
    return this.http.post(`/api/v1/services/info`, {}, {
      observe: "response",
      params: {
        'service_url': serviceURL,
        'service_identity': serviceIdentity
      }
    }).toPromise();
  }

  setServiceScale(serviceID: number, scale: number): Promise<any> {
    return this.http.put(`/api/v1/services/${serviceID}/scale`, {service_scale: scale}, {observe: "response"}).toPromise();
  }

  getCollaborativeService(serviceName: string, projectName: string): Promise<Array<string>> {
    return this.http
      .get<Array<string>>(`/api/v1/services/selectservices`, {
        observe: "response",
        params: {
          service_name: serviceName,
          project_name: projectName
        }
      })
      .toPromise()
      .then((res: HttpResponse<Array<string>>) => res.body || Array<string>());
  }

  getServiceYamlFile(projectName: string, serviceName: string, yamlType: string): Promise<string> {
    return this.http
      .get(`/api/v1/services/yaml/download`, {
        observe: "response",
        responseType: "text",
        params: {
          service_name: serviceName,
          project_name: projectName,
          yaml_type: yamlType
        }
      })
      .toPromise()
      .then((res: HttpResponse<string>) => res.body);
  }

  getServiceImages(projectName: string, serviceName: string): Promise<Array<ImageIndex>> {
    return this.http
      .get<Array<ImageIndex>>(`/api/v1/services/rollingupdate`, {
        observe: "response",
        params: {
          service_name: serviceName,
          project_name: projectName
        }
      })
      .toPromise()
      .then((res: HttpResponse<Array<ImageIndex>>) => res.body || Array<ImageIndex>());
  }

  updateServiceImages(projectName: string, serviceName: string, postData: Array<ImageIndex>): Promise<any> {
    return this.http
      .patch(`/api/v1/services/rollingupdate`, postData, {
        observe: "response",
        params: {
          service_name: serviceName,
          project_name: projectName
        }
      }).toPromise();
  }

  uploadServiceYamlFile(projectName: string, formData: FormData,): Observable<Service> {
    return this.http
      .post(`/api/v1/services/yaml/upload`, formData, {
        observe: "response",
        params: {
          project_name: projectName
        }
      })
      .map((res: HttpResponse<Service>) => res.body)
  }

  getServiceScaleInfo(serviceId: number): Observable<Object> {
    return this.http.get(`/api/v1/services/${serviceId}/scale`, {observe: "response"})
      .map((res: HttpResponse<Object>) => res.body)
  }
}