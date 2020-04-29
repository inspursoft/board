import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Configuration } from './configuration.model';

const BASE_URL = '/v1/admin/configuration';

@Injectable()
export class ConfigurationService {

  constructor(private http: HttpClient) {
  }

  getConfig(whichOne?: string): Observable<Configuration> {
    const token = window.sessionStorage.getItem('token');
    let url = `${BASE_URL}?token=${token}`;
    url = whichOne ? `${url}&which=${whichOne}` : url;
    return this.http.get(url, {
      observe: 'response',
    }).pipe(map((res: HttpResponse<Configuration>) => {
      return res.body;
    }));
  }

  putConfig(config: Configuration): Observable<any> {
    const token = window.sessionStorage.getItem('token');
    return this.http.put(
      `${BASE_URL}?token=${token}`,
      config.PostBody()
    );
  }

}
