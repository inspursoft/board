import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { ComponentStatus } from './component-status.model';

const BASE_URL = '/v1/admin';

@Injectable()
export class DashboardService {

  constructor(private http: HttpClient) { }

  restartBoard(token: string): Observable<any> {
    return this.http.get(
      `${BASE_URL}/account/restart?token=${token}`,
      { observe: 'response', }
    );
  }

  shutdownBoard(token: string): Observable<any> {
    return this.http.get(
      `${BASE_URL}/account/shutdown?token=${token}`,
      { observe: 'response', }
    );
  }

  monitorContainer(token: string): Observable<any> {
    return this.http.get(
      `${BASE_URL}/monitor?token=${token}`,
      { observe: 'response', })
      .pipe(map((res: HttpResponse<Array<ComponentStatus>>) => {
      return res.body;
    }));
  }
}
