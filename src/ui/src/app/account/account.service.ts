import { Injectable } from '@angular/core';
import { Http, Headers } from '@angular/http';
import 'rxjs/add/operator/toPromise';
import { NgXCookies } from 'ngx-cookies';

import { AppInitService } from '../app.init.service';
import { Account } from './account';

export const BASE_URL = '/api/v1';

@Injectable()
export class AccountService {
  
  defaultHeaders: Headers = new Headers({contentType: 'application/json'});
  constructor(
    private http: Http, 
    private appInitService: AppInitService){}

  signIn(principal: string, password: string): Promise<any> {
    return this.http
      .post(
        BASE_URL + '/sign-in', 
        { user_name: principal, user_password: password }, 
        { headers: this.defaultHeaders })
      .toPromise()
      .then(res=>res.json())
      .catch(err=>Promise.reject(err));
  }

  signUp(account: Account): Promise<any> {
    return this.http
      .post(
        BASE_URL + '/sign-up',
        { user_name: account.username, 
          user_email: account.email,
          user_password: account.password,
          user_realname: account.realname,
          user_comment: account.comment
        },
        { headers: this.defaultHeaders }
      )
      .toPromise()
      .then(res=>res)
      .catch(err=>Promise.reject(err));
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
      .then(res=>{
        NgXCookies.deleteCookie('token');
      })
      .catch(err=>Promise.reject(err));
  }
} 