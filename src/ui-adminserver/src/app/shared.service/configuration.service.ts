import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Configuration } from './cfg.model';

const BASE_URL = '/v1/admin/configuration';

@Injectable()
export class ConfigurationService {

  constructor(private http: HttpClient) {
  }

  getConfig(whichOne?: string): Observable<Configuration> {
    const url = whichOne ? `${BASE_URL}?which=${whichOne}` : BASE_URL;
    return this.http.get(url, {
      observe: 'response',
    }).pipe(map((res: HttpResponse<Configuration>) => {
      return res.body;
    }));
  }

  putConfig(config: Configuration): Observable<any> {
    return this.http.put(
      BASE_URL,
      config.PostBody()
    );
  }

}
