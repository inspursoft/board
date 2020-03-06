import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpResponse } from '@angular/common/http';
import { Observable, Subject, zip } from 'rxjs';
import { map } from 'rxjs/operators';
import { Project } from '../project/project';
import { BuildImageDockerfileData, Image, ImageDetail } from '../image/image';
import {
  ImageIndex,
  ServerServiceStep,
  ServiceStepPhase,
  UiServiceFactory,
  UIServiceStepBase
} from './service-step.component';
import { Service } from './service';
import { AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE } from '../shared/shared.const';
import { INode, INodeGroup, NodeAvailableResources, PersistentVolumeClaim, ServiceHPA } from '../shared/shared.types';

@Injectable()
export class K8sService {
  stepSource: Subject<{ index: number, isBack: boolean }> = new Subject<{ index: number, isBack: boolean }>();
  step$: Observable<{ index: number, isBack: boolean }> = this.stepSource.asObservable();

  constructor(private http: HttpClient) {
  }

  cancelBuildService(): void {
    this.deleteServiceConfig().subscribe(() => this.stepSource.next({index: 0, isBack: false}));
  }

  checkServiceExist(projectName: string, serviceName: string): Observable<any> {
    return this.http.get(`/api/v1/services/exists`, {
      observe: 'response',
      params: {project_name: projectName, service_name: serviceName}
    });
  }

  getServiceConfig(phase: ServiceStepPhase): Observable<UIServiceStepBase> {
    return this.http.get(`/api/v1/services/config`, {
      observe: 'response',
      params: {phase}
    }).pipe(map((res: HttpResponse<object>) => {
      const stepBase = UiServiceFactory.getInstance(phase);
      return stepBase.serverToUi(res.body);
    }));
  }

  setServiceConfig(config: ServerServiceStep): Observable<any> {
    return this.http.post(`/api/v1/services/config`, config.postData, {
      observe: 'response',
      params: {
        phase: config.phase,
        project_id: config.project_id.toString(),
        service_name: config.service_name,
        cluster_ip: config.cluster_ip,
        service_type: config.service_type.toString(),
        instance: config.instance.toString(),
        session_affinity_flag: config.session_affinity_flag.toString(),
        service_public: config.service_public.toString(),
        node_selector: config.node_selector
      }
    });
  }

  deleteServiceConfig(): Observable<any> {
    return this.http.delete(`/api/v1/services/config`, {observe: 'response'});
  }

  deleteDeployment(serviceId: number): Observable<any> {
    return this.http.delete(`/api/v1/services/${serviceId}/deployment`, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  serviceDeployment(): Observable<object> {
    return this.http.post(`/api/v1/services/deployment`, {}, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    }).pipe(map((res: HttpResponse<object>) => res.body));
  }

  serviceStatefulDeployment(): Observable<object> {
    return this.http.post(`/api/v1/services/statefulsets`, {}, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    }).pipe(map((res: HttpResponse<object>) => res.body));
  }

  getContainerDefaultInfo(imageName: string, imageTag: string, projectName: string):
    Observable<BuildImageDockerfileData> {
    return this.http.get<BuildImageDockerfileData>(`/api/v1/images/dockerfile`, {
      observe: 'response',
      params: {image_name: imageName, project_name: projectName, image_tag: imageTag}
    }).pipe(map((res: HttpResponse<BuildImageDockerfileData>) => res.body));
  }

  getProjects(projectName: string = ''): Observable<Array<Project>> {
    return this.http.get<Array<Project>>('/api/v1/projects', {
      observe: 'response',
      params: {project_name: projectName, member_only: '1'}
    }).pipe(map((res: HttpResponse<Array<Project>>) => res.body));
  }

  getDeployStatus(serviceId: number): Observable<any> {
    return this.http.get(`/api/v1/services/${serviceId}/status`, {observe: 'response'});
  }

  getImages(imageName?: string, imageListPage?: number, imageListPageSize?: number): Observable<Array<Image>> {
    return this.http.get('/api/v1/images', {
      observe: 'response',
      params: {
        image_name: imageName,
        image_list_page: imageListPage.toString(),
        image_list_page_size: imageListPageSize.toString()
      }
    }).pipe(map((res: HttpResponse<Image[]>) => res.body || []));
  }

  getImageDetailList(imageName: string): Observable<Array<ImageDetail>> {
    return this.http.get(`/api/v1/images/${imageName}`, {observe: 'response'})
      .pipe(map((res: HttpResponse<Array<ImageDetail>>) => res.body || []));
  }

  getServices(pageIndex: number, pageSize: number, sortBy: string, isReverse: boolean): Observable<object> {
    return this.http.get(`/api/v1/services`, {
      observe: 'response', params: {
        page_index: pageIndex.toString(),
        page_size: pageSize.toString(),
        order_field: sortBy,
        order_asc: isReverse ? '0' : '1'
      }
    }).pipe(map((res: HttpResponse<object>) => res.body));
  }

  getServiceDetail(serviceId: number): Observable<object> {
    return this.http.get(`/api/v1/services/${serviceId}/info`, {observe: 'response'})
      .pipe(map((res: HttpResponse<object>) => res.body));
  }

  deleteService(serviceID: number): Observable<any> {
    return this.http.delete(`/api/v1/services/${serviceID}`, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  deleteStatefulService(serviceID: number): Observable<any> {
    return this.http.delete(`/api/v1/services/${serviceID}/statefulsets`, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  toggleServiceStatus(serviceID: number, isStart: 0 | 1): Observable<any> {
    return this.http.put(`/api/v1/services/${serviceID}/toggle`, {service_toggle: isStart}, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  toggleServicePublicity(serviceID: number, servicePublic: 0 | 1): Observable<any> {
    return this.http.put(`/api/v1/services/${serviceID}/publicity`, {service_public: servicePublic}, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  getConsole(jobName: string, buildSerialId?: string): Observable<string> {
    return this.http
      .get(`/api/v1/jenkins-job/console`, {
        observe: 'response',
        responseType: 'text',
        params: {
          job_name: jobName,
          build_serial_id: buildSerialId
        }
      }).pipe(map((res: HttpResponse<string>) => res.body));
  }

  getNodesList(param?: {}): Observable<Array<{ node_name: string, node_ip: string, status: number }>> {
    const queryParam = param || {};
    return this.http
      .get(`/api/v1/nodes`, {observe: 'response', params: queryParam})
      .pipe(map((res: HttpResponse<Array<{ node_name: string, node_ip: string, status: number }>>) => res.body || []));
  }

  addServiceRoute(serviceURL: string, serviceIdentity: string): Observable<any> {
    return this.http.post(`/api/v1/services/info`, {}, {
      observe: 'response',
      params: {
        service_url: serviceURL,
        service_identity: serviceIdentity
      }
    });
  }

  setServiceScale(serviceID: number, scale: number): Observable<any> {
    return this.http.put(`/api/v1/services/${serviceID}/scale`, {service_scale: scale}, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  getCollaborativeService(serviceName: string, projectName: string): Observable<Array<Service>> {
    return this.http.get<Array<Service>>(`/api/v1/services/selectservices`, {
      observe: 'response',
      params: {
        service_name: serviceName,
        project_name: projectName
      }
    }).pipe(map((res: HttpResponse<Array<Service>>) => res.body || Array<Service>()));
  }

  getServiceYamlFile(projectName: string, serviceName: string, yamlType: string): Observable<string> {
    return this.http
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

  getServiceImages(projectName: string, serviceName: string): Observable<Array<ImageIndex>> {
    return this.http.get<Array<ImageIndex>>(`/api/v1/services/rollingupdate/image`, {
      observe: 'response',
      params: {
        service_name: serviceName,
        project_name: projectName
      }
    }).pipe(map((res: HttpResponse<Array<ImageIndex>>) => res.body || Array<ImageIndex>()));
  }

  updateServiceImages(projectName: string, serviceName: string, postData: Array<ImageIndex>): Observable<any> {
    return this.http
      .patch(`/api/v1/services/rollingupdate/image`, postData, {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: 'response',
        params: {
          service_name: serviceName,
          project_name: projectName
        }
      });
  }

  uploadServiceYamlFile(projectName: string, formData: FormData): Observable<Service> {
    return this.http
      .post(`/api/v1/services/yaml/upload`, formData, {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: 'response',
        params: {
          project_name: projectName
        }
      }).pipe(map((res: HttpResponse<Service>) => res.body));
  }

  getServiceScaleInfo(serviceId: number): Observable<object> {
    return this.http.get(`/api/v1/services/${serviceId}/scale`, {observe: 'response'})
      .pipe(map((res: HttpResponse<object>) => res.body));
  }

  getNodePorts(projectName: string): Observable<Array<number>> {
    return this.http.get(`/api/v1/services/nodeports`, {observe: 'response'})
      .pipe(map((res: HttpResponse<Array<number>>) => res.body));
  }

  getNodeSelectors(): Observable<Array<{ name: string, status: number }>> {
    const obsNodeList = this.http
      .get(`/api/v1/nodes`, {observe: 'response'})
      .pipe(
        map((res: HttpResponse<Array<INode>>) => {
          const result = Array<{ name: string, status: number }>();
          res.body.forEach((iNode: INode) => result.push({name: String(iNode.node_name).trim(), status: iNode.status}));
          return result;
        }));
    const obsNodeGroupList = this.http
      .get(`/api/v1/nodegroup`, {observe: 'response', params: {is_valid_node_group: '1'}})
      .pipe(
        map((res: HttpResponse<Array<INodeGroup>>) => {
          const result = Array<{ name: string, status: number }>();
          res.body.forEach((iNodeGroup: INodeGroup) => result.push({
            name: String(iNodeGroup.nodegroup_name).trim(),
            status: 1
          }));
          return result;
        }));
    return zip(obsNodeList, obsNodeGroupList).pipe(
      map(value => value[0].concat(value[1]))
    );
  }

  getLocate(projectName: string, serviceName: string): Observable<string> {
    return this.http.get(`/api/v1/services/rollingupdate/nodegroup`, {
      observe: 'response',
      params: {project_name: projectName, service_name: serviceName}
    }).pipe(map((res: HttpResponse<string>) => res.body));
  }

  setLocate(nodeSelector: string, projectName: string, serviceName: string): Observable<object> {
    return this.http.patch(`/api/v1/services/rollingupdate/nodegroup`, null,
      {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: 'response',
        params: {project_name: projectName, service_name: serviceName, node_selector: nodeSelector}
      }).pipe(map((res: HttpResponse<Array<object>>) => res.body));
  }

  getNodesAvailableSources(): Observable<Array<NodeAvailableResources>> {
    return this.http.get(`/api/v1/nodes/availableresources`, {
      observe: 'response'
    }).pipe(map((res: HttpResponse<Array<NodeAvailableResources>>) => res.body));
  }

  setAutoScaleConfig(serviceId: number, hpa: ServiceHPA): Observable<any> {
    return this.http.post(`/api/v1/services/${serviceId}/autoscale`, hpa, {
      observe: 'response'
    });
  }

  modifyAutoScaleConfig(serviceId: number, hpa: ServiceHPA): Observable<any> {
    return this.http.put(`/api/v1/services/${serviceId}/autoscale/${hpa.hpa_id}`, hpa, {
      observe: 'response'
    });
  }

  deleteAutoScaleConfig(serviceId: number, hpa: ServiceHPA): Observable<any> {
    return this.http.delete(`/api/v1/services/${serviceId}/autoscale/${hpa.hpa_id}`, {
      observe: 'response',
      params: {hpa_name: hpa.hpa_name}
    });
  }

  getAutoScaleConfig(serviceId: number): Observable<Array<ServiceHPA>> {
    return this.http.get(`/api/v1/services/${serviceId}/autoscale`, {
      observe: 'response'
    }).pipe(
      map((res: HttpResponse<Array<ServiceHPA>>) => res.body),
      map((res: Array<ServiceHPA>) => {
        res.forEach(config => config.isEdit = true);
        return res;
      }));
  }

  getSessionAffinityFlag(serviceName: string, projectName: string): Observable<boolean> {
    return this.http.get(`/api/v1/services/rollingupdate/session`, {
      observe: 'response', params: {
        project_name: projectName,
        service_name: serviceName
      }
    }).pipe(map((res: HttpResponse<object>) => Reflect.get(res.body, 'SessionAffinityFlag') === 1));
  }

  setSessionAffinityFlag(serviceName: string, projectName: string, flag: boolean): Observable<any> {
    return this.http.patch(`/api/v1/services/rollingupdate/session`, null, {
      observe: 'response', params: {
        project_name: projectName,
        service_name: serviceName,
        session_affinity_flag: flag ? '1' : '0'
      }
    });
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

  getConfigMapNames(projectName: string): Observable<Array<string>> {
    return this.http.get(`/api/v1/configmaps`, {
      params: {project_name: projectName}
    }).pipe(map((res: Array<object>) => {
      const result = new Array<string>();
      res.forEach(config => result.push(Reflect.get(config, 'name')));
      return result;
    }));
  }
}
