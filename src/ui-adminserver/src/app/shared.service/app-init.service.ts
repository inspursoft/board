import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { InitStatus } from './app-init.type';

const BASE_URL = '/v1/admin';

@Injectable()
export class AppInitService {
  isInited = false;
  currentLang: string;

  constructor(private http: HttpClient) { }

  getVerify(): Observable<any> {
    return this.http.get(`${BASE_URL}/account/install`, {
      observe: 'response',
    });
  }

  getSystemStatus(): Observable<InitStatus>  {
    return this.http.get(
      `${BASE_URL}/boot/checksysstatus`, {
      observe: 'response',
    }).pipe(map((res: HttpResponse<InitStatus>) => {
      return res.body;
    }));
  }
}
