import { Injectable } from "@angular/core";
import { HttpClient, HttpResponse } from "@angular/common/http";
import { Observable } from "rxjs";
import { ConfigMap, ConfigMapDetail, ConfigMapList } from "./resource.types";

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
    }).map((res: HttpResponse<Object>) => ConfigMapDetail.createFromRes(res.body))
  }

  deleteConfigMap(configMapName, projectName: string): Observable<any> {
    return this.http.delete(`/api/v1/configmaps/${configMapName}`, {
      observe: "response", params: {
        project_name: projectName
      }
    })
  }

  getConfigMapList(projectName: string, pageIndex, pageSize: number): Observable<Array<ConfigMapList>> {
    return this.http.get<Array<Object>>(`/api/v1/configmaps`, {
      observe: "response", params: {
        project_name: projectName,
        configmap_list_page: pageIndex.toString(),
        configmap_list_page_size: pageSize.toString()
      }
    }).map((res: HttpResponse<Array<Object>>) => {
      let result = Array<ConfigMapList>();
      res.body.forEach((configMap: Object) => result.push(ConfigMapList.createFromRes(configMap)));
      console.log(result);
      return result;
    })
  }
}