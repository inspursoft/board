import { Injectable } from '@angular/core';
import { HttpClient, HttpHeaders } from '@angular/common/http';
import { AppInitService } from '../app.init.service';
import { Account } from './account';
import { CookieService } from "ngx-cookie";

export const BASE_URL = '/api/v1';

@Injectable()
export class AccountService {

  defaultHeaders: HttpHeaders = new HttpHeaders({contentType: 'application/json'});

  constructor(
    private http: HttpClient,
    private cookieService: CookieService,
    private appInitService: AppInitService) {
  }

  signIn(principal: string, password: string): Promise<any> {
    return this.http
      .post(
        BASE_URL + '/sign-in',
        {user_name: principal, user_password: password},
        {observe: 'response'})
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
        {headers: this.defaultHeaders}
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
          params: {
            'username': this.appInitService.currentUser.user_name
          }
        }
      )
      .toPromise()
      .then(res => {
        this.cookieService.remove('token');
      })
      .catch(err => Promise.reject(err));
  }

  retrieve(credential): Promise<any> {
    return this.http
      .post(
        BASE_URL + `/forgot-password?credential=${credential}`,
        {},
        {observe: "response"}
      ).toPromise()
      .then(res => res)
      .catch(err => Promise.reject(err));
  }
} 