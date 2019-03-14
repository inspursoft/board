import { Injectable } from "@angular/core";
import { HttpClient, HttpEvent, HttpRequest, HttpResponse } from "@angular/common/http";
import { Observable, Subject } from "rxjs";
import { HelmViewData, IChartReleaseDetail, IChartRelease, IHelmRepo } from "./helm.type";

@Injectable()
export class HelmService {
  viewDataList: Array<HelmViewData>;
  viewSubject: Subject<HelmViewData>;

  constructor(private http: HttpClient) {
    this.viewSubject = new Subject<HelmViewData>();
    this.viewDataList = Array<HelmViewData>();
  }

  pushNewView(helmViewData: HelmViewData) {
    this.viewDataList.push(helmViewData);
    this.viewSubject.next(helmViewData);
  }

  cleanViewData() {
    this.viewDataList.splice(0, this.viewDataList.length);
  }

  popToView(helmViewData: HelmViewData) {
    let len = this.viewDataList.length - 1;
    for (let i = len; i >= 0; i--) {
      let data = this.viewDataList[i];
      if (data.type == helmViewData.type) {
        this.viewSubject.next(helmViewData);
        return;
      } else {
        this.viewDataList.splice(i, 1);
      }
    }
  }

  getRepoList(): Observable<Array<IHelmRepo>> {
    return this.http.get<Array<IHelmRepo>>('/api/v1/helm/repositories', {
      observe: "response"
    }).map((res: HttpResponse<Array<IHelmRepo>>) => res.body || [])
  }

  getRepoDetail(repoId: number, pageIndex: number = 1, pageSize: number = 1): Observable<Object> {
    return this.http.get<Object>(`/api/v1/helm/repositories/${repoId}`, {
      params: {page_index: pageIndex.toString(), page_size: pageSize.toString()},
      observe: "response"
    }).map((res: HttpResponse<Object>) => res.body)
  }

  uploadChart(repoId: number, formData: FormData): Observable<HttpEvent<Object>> {
    const req = new HttpRequest('POST', `/api/v1/helm/repositories/${repoId}/chartupload`, formData, {
      reportProgress: true,
    });
    return this.http.request<Object>(req)
  }

  deleteChartVersion(repoId: number, chartName, chartVersion: string): Observable<any> {
    return this.http.delete(`/api/v1/helm/repositories/${repoId}/charts/${chartName}/${chartVersion}`, {
      observe: 'response'
    })
  }

  deleteChartRelease(releaseId: number): Observable<any> {
    return this.http.delete(`/api/v1/helm/release/${releaseId}`, {
      observe: 'response'
    })
  }

  checkChartReleaseName(chartReleaseName: string): Observable<Object> {
    return this.http.get<Object>(`/api/v1/helm/release/existing`, {
      observe: "response", params: {
        release_name: chartReleaseName
      }
    }).map((res: HttpResponse<Object>) => res.body)
  }

  releaseChartVersion(postBody: { name, chart, chartVersion: string, repoId, projectId, ownerId: number }): Observable<any> {
    return this.http.post(`/api/v1/helm/release`, {
      name: postBody.name,
      project_id: postBody.projectId,
      repository_id: postBody.repoId,
      chart: postBody.chart,
      owner_id: postBody.ownerId,
      chartversion: postBody.chartVersion
    }, {observe: "response"})
  }

  getChartReleaseList(): Observable<Array<IChartRelease>> {
    return this.http.get<Object>(`/api/v1/helm/release`, {observe: "response"})
      .map((res: HttpResponse<Array<IChartRelease>>) => res.body)
  }

  getChartReleaseDetail(chartReleaseId: number): Observable<IChartReleaseDetail> {
    return this.http.get<Object>(`/api/v1/helm/release/${chartReleaseId}`, {observe: "response"})
      .map((res: HttpResponse<IChartReleaseDetail>) => res.body)
  }

  getChartRelease(repoId: number, chartName, chartVersion: string): Observable<Object> {
    return this.http.get(`/api/v1/helm/repositories/${repoId}/charts/${chartName}/${chartVersion}`, {observe: 'response'})
      .map((res: HttpResponse<Object>) => res.body)
  }
}
