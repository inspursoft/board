import { Injectable } from '@angular/core';
import { Http, Headers } from '@angular/http';
import 'rxjs/add/operator/toPromise';

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
            'token': this.appInitService.token
          }
        }
      )
      .toPromise()
      .then(res=>res)
      .catch(err=>Promise.reject(err));
  }
} 