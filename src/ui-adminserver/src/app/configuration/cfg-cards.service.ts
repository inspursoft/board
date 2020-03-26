import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { Configuration, VerifyPassword } from './cfg.models';
import { UserVerify, User } from '../account/account.model';

const BASE_URL = '/v1/admin';

@Injectable()
export class CfgCardsService {

  constructor(private http: HttpClient,
              private token: string) {
    this.token = window.sessionStorage.getItem('token');
  }

  getConfig(whichOne?: string): Observable<Configuration> {
    let url = `${BASE_URL}/configuration?token=${this.token}`;
    url = whichOne ? `${url}&which=${whichOne}` : url;
    return this.http.get(url, {
      observe: 'response',
    }).pipe(map((res: HttpResponse<Configuration>) => {
      return res.body;
    }));
  }

  postConfig(config: Configuration): Observable<any> {
    return this.http.post(
      `${BASE_URL}/configuration?token=${this.token}`,
      config.PostBody()
    );
  }

  getPubKey(): Observable<string> {
    return this.http.get(`${BASE_URL}/configuration/pubkey?token=${this.token}`, {
      observe: 'response',
    }).pipe(map((res: HttpResponse<string>) => {
      return res.body;
    }));
  }

  applyCfg(user: User): Observable<any> {
    return this.http.post(
      `${BASE_URL}/board/applycfg?token=${this.token}`,
      user.PostBody()
    );
  }

  verifyPassword(oldPwd: VerifyPassword): Observable<any> {
    return this.http.post(
      `${BASE_URL}/account/verify?token=${this.token}`,
      oldPwd.PostBody()
    );
  }
}
