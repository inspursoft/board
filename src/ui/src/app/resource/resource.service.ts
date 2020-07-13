import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { ConfigMap, ConfigMapDetail, ConfigMapProject } from './resource.types';
import { ModelHttpClient } from '../shared/ui-model/model-http-client';

@Injectable()
export class ResourceService {
  constructor(private http: ModelHttpClient) {

  }

  createConfigMap(configMap: ConfigMap): Observable<any> {
    return this.http.post(`/api/v1/configmaps`, configMap.getPostBody());
  }

  getConfigMapDetail(configMapName, projectName: string): Observable<ConfigMapDetail> {
    return this.http.getJson(`/api/v1/configmaps`, ConfigMapDetail, {
      param: {
        project_name: projectName,
        configmap_name: configMapName,
      }
    });
  }

  deleteConfigMap(configMapName, projectName: string): Observable<any> {
    return this.http.delete(`/api/v1/configmaps/${configMapName}`, {
      observe: 'response', params: {
        project_name: projectName
      }
    });
  }

  updateConfigMap(configMap: ConfigMap): Observable<any> {
    return this.http.put(`/api/v1/configmaps/${configMap.name}`, configMap.getPostBody());
  }

  getAllProjects(): Observable<Array<ConfigMapProject>> {
    return this.http.getArray('/api/v1/projects', ConfigMapProject);
  }

  getConfigMapList(projectName: string, pageIndex, pageSize: number): Observable<Array<ConfigMap>> {
    return this.http.getArray(`/api/v1/configmaps`, ConfigMap, {
      param: {
        project_name: projectName,
        configmap_list_page: pageIndex.toString(),
        configmap_list_page_size: pageSize.toString()
      }
    });
  }
}
