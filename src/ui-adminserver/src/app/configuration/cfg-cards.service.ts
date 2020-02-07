import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Configuration, VerifyPassword } from './cfg.models';

const BASE_URL = '/v1/admin';

@Injectable()
export class CfgCardsService {

  constructor(private http: HttpClient) { }

  getConfig(whichOne?: string): Observable<Configuration> {
    const url = whichOne ? `${BASE_URL}/configuration?which=${whichOne}` : `${BASE_URL}/configuration/`;
    return this.http.get(url, {
      observe: 'response',
    }).pipe(map((res: HttpResponse<Configuration>) => {
      return res.body;
    }));
  }

  postConfig(config: Configuration): Observable<any> {
    return this.http.post(
      `${BASE_URL}/configuration/`,
      config.PostBody()
    );
  }

  getPubKey(): Observable<string> {
    return this.http.get(`${BASE_URL}/configuration/pubkey/`, {
      observe: 'response',
    }).pipe(map((res: HttpResponse<string>) => {
      return res.body;
    }));
  }

  applyCfg(token: string): Observable<any> {
    return this.http.get(
      `${BASE_URL}/account/applycfg?token=${token}`,
      { observe: 'response', }
    );
  }

  verifyPassword(oldPwd: VerifyPassword):  Observable<any> {
    return this.http.post(
      `${BASE_URL}/account/verify/`,
      oldPwd.PostBody()
    );
  }
}
