import { Injectable } from '@angular/core';
import { Observable } from 'rxjs';
import { map, tap } from 'rxjs/operators';
import { BASE_URL, BodyData, Prometheus, QueryData } from './dashboard.types';
import { ModelHttpClient } from '../shared/ui-model/model-http-client';

@Injectable()
export class DashboardService {

  constructor(private modelHttp: ModelHttpClient) {
  }

  testGrafana(grafanaUrl: string): Observable<any> {
    return this.modelHttp.get(grafanaUrl, {observe: 'response'});
  }

  getServerTimeStamp(): Observable<number> {
    return this.modelHttp.get(`${BASE_URL}/dashboard/time`)
      .pipe(map((res: object) => Reflect.get(res, 'time_now')));
  }

  getLineData(paramData: QueryData, bodyData: BodyData): Observable<Prometheus> {
    return this.modelHttp.postJson(`${BASE_URL}/prometheus`, Prometheus, bodyData.getPostBody(),
      {param: paramData.getPostBody()}
    ).pipe(tap((prometheus: Prometheus) => {
      prometheus.analyzeData(paramData.serviceName, paramData.nodeName);
    }));
  }
}
