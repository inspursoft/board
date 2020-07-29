import { Injectable } from "@angular/core";
import { HttpClient, HttpResponse } from "@angular/common/http";
import { Observable } from "rxjs";
import { ConfigMap, ConfigMapDetail } from "./resource.types";
import { map } from "rxjs/operators";
import { decode, encode } from 'punycode';
import { encodeUriQuery } from '@angular/router/src/url_tree';
import { safeDecodeURIComponent } from 'ngx-cookie';

@Injectable()
export class ResourceService {
  constructor(private http: HttpClient) {

  }

  createConfigMap(configMap: ConfigMap): Observable<any> {
    return this.http.post(`/api/v1/configmaps`, configMap.postBody(), {observe: "response"})
  }

  getConfigMapDetail(configMapName, projectName: string): Observable<ConfigMapDetail> {
    return this.http.get(`/api/v1/configmaps/${configMapName}`, {
      observe: "response", params: {
        project_name: projectName
      }
    }).pipe(map((res: HttpResponse<Object>) => ConfigMapDetail.createFromRes(res.body)));
  }

  deleteConfigMap(configMapName, projectName: string): Observable<any> {
    return this.http.delete(`/api/v1/configmaps/${configMapName}`, {
      observe: "response", params: {
        project_name: projectName
      }
    })
  }

  updateConfigMap(configMap: ConfigMap): Observable<any> {
    return this.http.put(`/api/v1/configmaps/${configMap.name}`, configMap.postBody(), {observe: "response"})
  }

  getConfigMapList(projectName: string, pageIndex, pageSize: number): Observable<Array<ConfigMap>> {
    return this.http.get<Array<Object>>(`/api/v1/configmaps`, {
      observe: "response", params: {
        project_name: projectName,
        configmap_list_page: pageIndex.toString(),
        configmap_list_page_size: pageSize.toString()
      }
    }).pipe(map((res: HttpResponse<Array<Object>>) => {
      let result = Array<ConfigMap>();
      res.body.forEach((configMap: Object) => result.push(ConfigMap.createFromRes(configMap)));
      return result;
    }));
  }
}
