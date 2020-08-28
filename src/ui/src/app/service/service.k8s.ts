import { Injectable, Type } from '@angular/core';
import { HttpHeaders, HttpParams, HttpRequest, HttpResponse, HttpClient } from '@angular/common/http';
import { Observable, Subject, zip } from 'rxjs';
import { map } from 'rxjs/operators';
import { ServiceStepDataBase, ServiceStepPhase } from './service-step.component';
import { AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE } from '../shared/shared.const';
import { PersistentVolumeClaim } from '../shared/shared.types';
import { ModelHttpClient } from '../shared/ui-model/model-http-client';
import { HttpBase } from '../shared/ui-model/model-types';
import {
  NodeAvailableResources,
  PaginationService,
  Service,
  ServiceDetailInfo,
  ServiceDockerfileData,
  ServiceHPA,
  ServiceImage,
  ServiceImageDetail, ServiceNode, ServiceNodeGroup,
  ServiceProject
} from './service.types';

@Injectable()
export class K8sService {
  stepSource: Subject<{ index: number, isBack: boolean }> = new Subject<{ index: number, isBack: boolean }>();
  step$: Observable<{ index: number, isBack: boolean }> = this.stepSource.asObservable();

  constructor(private httpModel: ModelHttpClient) {
  }

  cancelBuildService(): void {
    this.deleteServiceConfig().subscribe(() => this.stepSource.next({index: 0, isBack: false}));
  }

  checkServiceExist(projectName: string, serviceName: string): Observable<any> {
    return this.httpModel.get(`/api/v1/services/exists`, {
      observe: 'response',
      params: {project_name: projectName, service_name: serviceName}
    });
  }

  getServiceConfig(phase: ServiceStepPhase, returnType: Type<HttpBase>): Observable<any> {
    return this.httpModel.getJson(`/api/v1/services/config`, returnType, {
      param: {phase}
    });
  }

  setServiceStepConfig(stepData: ServiceStepDataBase): Observable<any> {
    return this.httpModel.post(`/api/v1/services/config`, stepData.getPostBody(), {
      params: stepData.getParams()
    });
  }

  deleteServiceConfig(): Observable<any> {
    return this.httpModel.delete(`/api/v1/services/config`, {observe: 'response'});
  }

  deleteDeployment(serviceId: number): Observable<any> {
    return this.httpModel.delete(`/api/v1/services/${serviceId}/deployment`, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  serviceDeployment(): Observable<object> {
    return this.httpModel.post(`/api/v1/services/deployment`, {}, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    }).pipe(map((res: HttpResponse<object>) => res.body));
  }

  serviceStatefulDeployment(): Observable<object> {
    return this.httpModel.post(`/api/v1/services/statefulsets`, {}, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    }).pipe(map((res: HttpResponse<object>) => res.body));
  }

  getContainerDefaultInfo(imageName: string,
                          imageTag: string,
                          projectName: string): Observable<ServiceDockerfileData> {
    return this.httpModel.getJson(`/api/v1/images/dockerfile`, ServiceDockerfileData, {
      param: {image_name: imageName, project_name: projectName, image_tag: imageTag}
    });
  }

  getProjects(projectName: string = ''): Observable<Array<ServiceProject>> {
    return this.httpModel.getArray('/api/v1/projects', ServiceProject, {
      param: {project_name: projectName, member_only: '1'}
    });
  }

  getOneProject(projectName: string): Observable<Array<ServiceProject>> {
    return this.httpModel.getArray('/api/v1/projects', ServiceProject, {
      param: {project_name: projectName}
    });
  }

  getDeployStatus(serviceId: number): Observable<any> {
    return this.httpModel.get(`/api/v1/services/${serviceId}/status`, {observe: 'response'});
  }

  getImages(imageName?: string, imageListPage?: number, imageListPageSize?: number): Observable<Array<ServiceImage>> {
    return this.httpModel.getArray('/api/v1/images', ServiceImage, {
      param: {
        image_name: imageName,
        image_list_page: imageListPage.toString(),
        image_list_page_size: imageListPageSize.toString()
      }
    });
  }

  getImageDetailList(imageName: string): Observable<Array<ServiceImageDetail>> {
    return this.httpModel.getArray(`/api/v1/images/${imageName}`, ServiceImageDetail);
  }

  getServices(pageIndex: number, pageSize: number, sortBy: string, isReverse: boolean): Observable<PaginationService> {
    return this.httpModel.getPagination(`/api/v1/services`, PaginationService, {
      param: {
        page_index: pageIndex.toString(),
        page_size: pageSize.toString(),
        order_field: sortBy,
        order_asc: isReverse ? '0' : '1'
      }
    });
  }

  getServiceDetail(serviceId: number): Observable<ServiceDetailInfo> {
    return this.httpModel.getJson(`/api/v1/services/${serviceId}/info`, ServiceDetailInfo);
  }

  deleteService(serviceID: number): Observable<any> {
    return this.httpModel.delete(`/api/v1/services/${serviceID}`, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  deleteStatefulService(serviceID: number): Observable<any> {
    return this.httpModel.delete(`/api/v1/services/${serviceID}/statefulsets`, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  toggleServiceStatus(serviceID: number, isStart: 0 | 1): Observable<any> {
    return this.httpModel.put(`/api/v1/services/${serviceID}/toggle`, {service_toggle: isStart}, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  toggleServicePublicity(serviceID: number, servicePublic: 0 | 1): Observable<any> {
    return this.httpModel.put(`/api/v1/services/${serviceID}/publicity`, {service_public: servicePublic}, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  getConsole(jobName: string, buildSerialId?: string): Observable<string> {
    return this.httpModel
      .get(`/api/v1/jenkins-job/console`, {
        observe: 'response',
        responseType: 'text',
        params: {
          job_name: jobName,
          build_serial_id: buildSerialId
        }
      }).pipe(map((res: HttpResponse<string>) => res.body));
  }

  getNodesList(param?: {}): Observable<Array<ServiceNode>> {
    const queryParam = param || {};
    return this.httpModel.getArray(`/api/v1/nodes`, ServiceNode, {param: queryParam});
  }

  addServiceRoute(serviceURL: string, serviceIdentity: string): Observable<any> {
    return this.httpModel.post(`/api/v1/services/info`, {}, {
      observe: 'response',
      params: {
        service_url: serviceURL,
        service_identity: serviceIdentity
      }
    });
  }

  setServiceScale(serviceID: number, scale: number): Observable<any> {
    return this.httpModel.put(`/api/v1/services/${serviceID}/scale`, {service_scale: scale}, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  getCollaborativeService(serviceName: string, projectName: string): Observable<Array<Service>> {
    return this.httpModel.getArray(`/api/v1/services/selectservices`, Service, {
      param: {
        service_name: serviceName,
        project_name: projectName
      }
    });
  }

  getServiceYamlFile(projectName: string, serviceName: string, yamlType: string): Observable<string> {
    return this.httpModel
      .get(`/api/v1/services/yaml/download`, {
        observe: 'response',
        responseType: 'text',
        params: {
          service_name: serviceName,
          project_name: projectName,
          yaml_type: yamlType
        }
      }).pipe(map((res: HttpResponse<string>) => res.body));
  }

  getServiceImages(projectName: string, serviceName: string): Observable<Array<ServiceImage>> {
    return this.httpModel.getArray(`/api/v1/services/rollingupdate/image`, ServiceImage, {
      param: {
        service_name: serviceName,
        project_name: projectName
      }
    });
  }

  updateServiceImages(projectName: string, serviceName: string, postData: Array<{ [key: string]: string }>): Observable<any> {
    return this.httpModel
      .patch(`/api/v1/services/rollingupdate/image`, postData, {
          headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
          observe: 'response',
          params: {
            service_name: serviceName,
            project_name: projectName
          }
        }
      );
  }

  uploadServiceYamlFile(projectName: string, formData: FormData): Observable<Service> {
    return this.httpModel
      .post(`/api/v1/services/yaml/upload`, formData, {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        params: {
          project_name: projectName
        }
      }).pipe(map(res => {
        const service = new Service(res);
        service.initFromRes();
        return service;
      }));
  }

  getServiceScaleInfo(serviceId: number): Observable<object> {
    return this.httpModel.get(`/api/v1/services/${serviceId}/scale`, {observe: 'response'})
      .pipe(map((res: HttpResponse<object>) => res.body));
  }

  getNodePorts(projectName: string): Observable<Array<number>> {
    return this.httpModel.get(`/api/v1/services/nodeports`).pipe(map((res: Array<number>) => res || new Array<number>()));
  }

  getNodeSelectors(): Observable<Array<{ name: string, status: number }>> {
    const obsNodeList = this.httpModel.getArray(`/api/v1/nodes`, ServiceNode)
      .pipe(
        map((res: Array<ServiceNode>) => {
          const result = Array<{ name: string, status: number }>();
          res.forEach((node: ServiceNode) => {
            if (node.isNormalNode) {
              result.push({name: String(node.nodeName).trim(), status: node.status});
            }
          }
          );
          return result;
        }));
    const obsNodeGroupList = this.httpModel
      .getArray(`/api/v1/nodegroup`, ServiceNodeGroup, {param: {is_valid_node_group: '1'}})
      .pipe(
        map((res: Array<ServiceNodeGroup>) => {
          const result = Array<{ name: string, status: number }>();
          res.forEach((group: ServiceNodeGroup) => result.push({
            name: String(group.name).trim(),
            status: 1
          }));
          return result;
        }));
    return zip(obsNodeList, obsNodeGroupList).pipe(
      map(value => value[0].concat(value[1]))
    );
  }

  getLocate(projectName: string, serviceName: string): Observable<string> {
    return this.httpModel.get(`/api/v1/services/rollingupdate/nodegroup`, {
      observe: 'response',
      params: {project_name: projectName, service_name: serviceName}
    }).pipe(map((res: HttpResponse<string>) => res.body));
  }

  setLocate(nodeSelector: string, projectName: string, serviceName: string): Observable<object> {
    return this.httpModel.patch(`/api/v1/services/rollingupdate/nodegroup`, null,
      {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: 'response',
        params: {project_name: projectName, service_name: serviceName, node_selector: nodeSelector}
      }).pipe(map((res: HttpResponse<Array<object>>) => res.body));
  }

  getNodesAvailableSources(): Observable<Array<NodeAvailableResources>> {
    return this.httpModel.getArray(`/api/v1/nodes/availableresources`, NodeAvailableResources);
  }

  setAutoScaleConfig(serviceId: number, hpa: ServiceHPA): Observable<any> {
    return this.httpModel.post(`/api/v1/services/${serviceId}/autoscale`, hpa.getPostBody(), {
      observe: 'response'
    });
  }

  modifyAutoScaleConfig(serviceId: number, hpa: ServiceHPA): Observable<any> {
    return this.httpModel.put(`/api/v1/services/${serviceId}/autoscale/${hpa.hpaId}`, hpa.getPostBody(), {
      observe: 'response'
    });
  }

  deleteAutoScaleConfig(serviceId: number, hpa: ServiceHPA): Observable<any> {
    return this.httpModel.delete(`/api/v1/services/${serviceId}/autoscale/${hpa.hpaId}`, {
      observe: 'response',
      params: {hpa_name: hpa.hpaName}
    });
  }

  getAutoScaleConfig(serviceId: number): Observable<Array<ServiceHPA>> {
    return this.httpModel.getArray(`/api/v1/services/${serviceId}/autoscale`, ServiceHPA).pipe(
      map((res: Array<ServiceHPA>) => {
        res.forEach(serviceHPA => serviceHPA.isEdit = true);
        return res;
      })
    );
  }

  getSessionAffinityFlag(serviceName: string, projectName: string): Observable<boolean> {
    return this.httpModel.get(`/api/v1/services/rollingupdate/session`, {
      observe: 'response', params: {
        project_name: projectName,
        service_name: serviceName
      }
    }).pipe(map((res: HttpResponse<object>) => Reflect.get(res.body, 'SessionAffinityFlag') === 1));
  }

  setSessionAffinityFlag(serviceName: string, projectName: string, flag: boolean): Observable<any> {
    return this.httpModel.patch(`/api/v1/services/rollingupdate/session`, null, {
      observe: 'response', params: {
        project_name: projectName,
        service_name: serviceName,
        session_affinity_flag: flag ? '1' : '0'
      }
    });
  }

  getPvcNameList(): Observable<Array<PersistentVolumeClaim>> {
    return this.httpModel.get(`/api/v1/pvclaims`, {observe: 'response'})
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

  getConfigMapNames(projectName: string): Observable<Array<string>> {
    return this.httpModel.get(`/api/v1/configmaps`, {
      params: {project_name: projectName}
    }).pipe(map((res: Array<object>) => {
      const result = new Array<string>();
      res.forEach(config => result.push(Reflect.get(config, 'name')));
      return result;
    }));
  }

  downloadFile(projectId: number, podName, containerName, srcPath: string): Observable<any> {
    let httpParams = new HttpParams();
    httpParams = httpParams.set('container', containerName);
    httpParams = httpParams.set('src', srcPath);
    const req = new HttpRequest('GET', `/api/v1/pods/${projectId}/${podName}/download`, {
      reportProgress: true,
      responseType: 'blob',
      params: httpParams
    });
    return this.httpModel.request<any>(req);
  }

  uploadFile(projectId: number, podName, containerName, destPath: string, formData: FormData): Observable<any> {
    let httpParams = new HttpParams();
    httpParams = httpParams.set('container', containerName);
    httpParams = httpParams.set('dest', destPath);
    const req = new HttpRequest('POST', `/api/v1/pods/${projectId}/${podName}/upload`, formData, {
      reportProgress: true,
      params: httpParams
    });
    return this.httpModel.request<any>(req);
  }

  getEdgeNodes(): Observable<Array<{ description: string }>> {
    return this.httpModel.get(`/api/v1/edgenodes`).pipe(
      map((res: Array<string>) => {
        const result = new Array<{ description: string }>();
        if (res && res.length > 0) {
          res.forEach(value => result.push({description: value}));
        }
        return result;
      })
    );
  }

  getNodeGroups(): Observable<Array<{ description: string }>> {
    return this.httpModel.getArray(`/api/v1/nodegroup`, ServiceNodeGroup,
      {param: {is_valid_node_group: '1'}}
    ).pipe(map((res: Array<ServiceNodeGroup>) => {
        const result = new Array<{ description: string }>();
        res.forEach(group => result.push({description: group.name}));
        return result;
      })
    );
  }
}
