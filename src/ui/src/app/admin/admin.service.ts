import { Injectable } from '@angular/core';
import { ModelHttpClient } from '../shared/ui-model/model-http-client';
import { Observable } from 'rxjs';
import { HttpHeaders } from '@angular/common/http';
import { AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE } from '../shared/shared.const';

@Injectable()
export class AdminService {

  constructor(private http: ModelHttpClient) {
  }

  getK8sProxyConfig(): Observable<{ enable: any }> {
    return this.http.get<{ enable: any }>(`/api/v1/k8sproxy`);
  }

  setK8sProxyConfig(enable: boolean): Observable<any> {
    return this.http.put(`/api/v1/k8sproxy`, {enable},
      {headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE)}
    );
  }
}
