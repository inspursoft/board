import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { Observable } from 'rxjs';

const BASE_URL = '/v1/admin/account';

@Injectable()
export class AppInitService {
  isInited = false;
  currentLang: string;

  constructor(private http: HttpClient) { }

  getVerify(): Observable<any> {
    return this.http.get(`${BASE_URL}/install`, {
      observe: 'response',
    });
  }
}
