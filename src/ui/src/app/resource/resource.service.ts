import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { ConfigMapProject } from './resource.types';
import { ModelHttpClient } from '../shared/ui-model/model-http-client';
import { SharedConfigMap, SharedConfigMapDetail } from '../shared/shared.types';

@Injectable()
export class ResourceService {
  constructor(private http: ModelHttpClient) {

  }

  createConfigMap(configMap: SharedConfigMap): Observable<any> {
    return this.http.post(`/api/v1/configmaps`, configMap.getPostBody());
  }

  getConfigMapDetail(configMapName, projectName: string): Observable<SharedConfigMapDetail> {
    return this.http.getJson(`/api/v1/configmaps`, SharedConfigMapDetail, {
      param: {
        project_name: projectName,
        configmap_name: configMapName,
      }
    });
  }

  deleteConfigMap(configMapName, projectName: string): Observable<any> {
    return this.http.delete(`/api/v1/configmaps`, {
      observe: 'response', params: {
        project_name: projectName,
        configmap_name: configMapName
      }
    });
  }

  updateConfigMap(configMap: SharedConfigMap): Observable<any> {
    return this.http.put(`/api/v1/configmaps`, configMap.getPostBody(), {
      params: {project_name: configMap.namespace, configmap_name: configMap.name}
    });
  }

  getAllProjects(): Observable<Array<ConfigMapProject>> {
    return this.http.getArray('/api/v1/projects', ConfigMapProject);
  }

  getConfigMapList(projectName: string, pageIndex, pageSize: number): Observable<Array<SharedConfigMap>> {
    return this.http.getArray(`/api/v1/configmaps`, SharedConfigMap, {
      param: {
        project_name: projectName,
        configmap_list_page: pageIndex.toString(),
        configmap_list_page_size: pageSize.toString()
      }
    });
  }
}
