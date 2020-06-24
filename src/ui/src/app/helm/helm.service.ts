import { Injectable } from '@angular/core';
import { HttpEvent, HttpRequest, HttpResponse } from '@angular/common/http';
import { Observable, Subject } from 'rxjs';
import { map } from 'rxjs/operators';
import {
  HelmViewData,
  IChartReleaseDetail,
  IChartRelease,
  IHelmRepo,
  HelmRepoDetail,
  IChartReleasePost,
  ChartRelease
} from './helm.type';
import { Project } from '../project/project';
import { ModelHttpClient } from '../shared/ui-model/model-http-client';

@Injectable()
export class HelmService {
  viewDataList: Array<HelmViewData>;
  viewSubject: Subject<HelmViewData>;

  constructor(private http: ModelHttpClient) {
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
    const len = this.viewDataList.length - 1;
    for (let i = len; i >= 0; i--) {
      const data = this.viewDataList[i];
      if (data.type === helmViewData.type) {
        this.viewSubject.next(helmViewData);
        return;
      } else {
        this.viewDataList.splice(i, 1);
      }
    }
  }

  getRepoList(): Observable<Array<IHelmRepo>> {
    return this.http.get<Array<IHelmRepo>>('/api/v1/helm/repositories', {
      observe: 'response'
    }).pipe(map((res: HttpResponse<Array<IHelmRepo>>) => res.body || Array<IHelmRepo>()));
  }

  getRepoDetail(repoId: number, pageIndex: number = 1, pageSize: number = 1): Observable<HelmRepoDetail> {
    return this.http.getJson(`/api/v1/helm/repositories/${repoId}`, HelmRepoDetail, {
      param: {page_index: pageIndex.toString(), page_size: pageSize.toString()}
    });
  }

  uploadChart(repoId: number, formData: FormData): Observable<HttpEvent<object>> {
    const req = new HttpRequest('POST', `/api/v1/helm/repositories/${repoId}/chartupload`, formData, {
      reportProgress: true,
    });
    return this.http.request<object>(req);
  }

  deleteChartVersion(repoId: number, chartName, chartVersion: string): Observable<any> {
    return this.http.delete(`/api/v1/helm/repositories/${repoId}/charts/${chartName}/${chartVersion}`, {
      observe: 'response'
    });
  }

  deleteChartRelease(releaseId: number): Observable<any> {
    return this.http.delete(`/api/v1/helm/release/${releaseId}`, {
      observe: 'response'
    });
  }

  getProjects(): Observable<Array<Project>> {
    return this.http.get<Array<Project>>('/api/v1/projects', {
      observe: 'response',
      params: {member_only: '1'}
    }).pipe(map((res: HttpResponse<Array<Project>>) => res.body));
  }

  checkChartReleaseName(chartReleaseName: string): Observable<object> {
    return this.http.get<object>(`/api/v1/helm/release/existing`, {
      observe: 'response', params: {
        release_name: chartReleaseName
      }
    }).pipe(map((res: HttpResponse<object>) => res.body));
  }

  releaseChartVersion(postBody: IChartReleasePost): Observable<any> {
    return this.http.post(`/api/v1/helm/release`, postBody, {observe: 'response'});
  }

  getChartReleaseList(): Observable<Array<IChartRelease>> {
    return this.http.get<object>(`/api/v1/helm/release`, {observe: 'response'})
      .pipe(map((res: HttpResponse<Array<IChartRelease>>) => res.body || Array<IChartRelease>()));
  }

  getChartReleaseDetail(chartReleaseId: number): Observable<IChartReleaseDetail> {
    return this.http.get<object>(`/api/v1/helm/release/${chartReleaseId}`, {observe: 'response'})
      .pipe(map((res: HttpResponse<IChartReleaseDetail>) => res.body));
  }

  getChartRelease(repoId: number, chartName, chartVersion: string): Observable<ChartRelease> {
    return this.http.getJson(`/api/v1/helm/repositories/${repoId}/charts/${chartName}/${chartVersion}`, ChartRelease);
  }
}
