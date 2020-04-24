import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Configuration, VerifyPassword } from './cfg.models';
import { User } from '../account/account.model';

const BASE_URL = '/v1/admin';

@Injectable()
export class CfgCardsService {
  private token = '';

  constructor(private http: HttpClient) {
  }

  getConfig(whichOne?: string): Observable<Configuration> {
    this.token = window.sessionStorage.getItem('token');
    let url = `${BASE_URL}/configuration?token=${this.token}`;
    url = whichOne ? `${url}&which=${whichOne}` : url;
    return this.http.get(url, {
      observe: 'response',
    }).pipe(map((res: HttpResponse<Configuration>) => {
      return res.body;
    }));
  }

  postConfig(config: Configuration): Observable<any> {
    this.token = window.sessionStorage.getItem('token');
    return this.http.post(
      `${BASE_URL}/configuration?token=${this.token}`,
      config.PostBody()
    );
  }

  getPubKey(): Observable<any> {
    this.token = window.sessionStorage.getItem('token');
    return this.http.get(`${BASE_URL}/configuration/pubkey?token=${this.token}`, {
      observe: 'response',
    }).pipe(map((res: HttpResponse<any>) => {
      return res.body;
    }));
  }

  applyCfg(user: User): Observable<any> {
    this.token = window.sessionStorage.getItem('token');
    return this.http.post(
      `${BASE_URL}/board/applycfg?token=${this.token}`,
      user.PostBody()
    );
  }

  verifyPassword(oldPwd: VerifyPassword): Observable<any> {
    this.token = window.sessionStorage.getItem('token');
    return this.http.post(
      `${BASE_URL}/account/verify?token=${this.token}`,
      oldPwd.PostBody()
    );
  }
}
