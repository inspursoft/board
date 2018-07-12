import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders, HttpParams } from '@angular/common/http';
import { AppInitService } from '../app.init.service';
import { Account } from './account';
import { CookieService } from "ngx-cookie";
import { AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE } from "../shared/shared.const";

export const BASE_URL = '/api/v1';

@Injectable()
export class AccountService {

  constructor(private http: HttpClient,
              private cookieService: CookieService,
              private appInitService: AppInitService) {
  }

  signIn(principal: string, password: string): Promise<any> {
    return this.http
      .post(
        BASE_URL + '/sign-in',
        {user_name: principal, user_password: password},
        {headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE), observe: 'response'})
      .toPromise()
      .then(res => res.body)
      .catch(err => Promise.reject(err));
  }

  signUp(account: Account): Promise<any> {
    return this.http
      .post(
        BASE_URL + '/sign-up',
        {
          user_name: account.username,
          user_email: account.email,
          user_password: account.password,
          user_realname: account.realname,
          user_comment: account.comment
        },
        {headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE)}
      )
      .toPromise()
      .then(res => res)
      .catch(err => Promise.reject(err));
  }

  signOut(): Promise<any> {
    return this.http
      .get(
        BASE_URL + '/log-out',
        {
          headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
          params: {
            'username': this.appInitService.currentUser.user_name
          }
        }
      )
      .toPromise()
      .then(() => this.cookieService.remove('token'))
      .catch(err => Promise.reject(err));
  }

  postEmail(credential: string): Promise<any> {
    return this.http.post(BASE_URL + `/forgot-password?credential=${credential}`, null,
      {
        headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
        observe: "response"
      }).toPromise()
  }

  resetPassword(password, resetUuid): Promise<any> {
    let httpParams = new HttpParams().append('password', password).append('reset_uuid', resetUuid);
    return this.http.post(BASE_URL + '/reset-password', null, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: "response",
      params: httpParams
    }).toPromise()
  }
} 