import { Injectable } from '@angular/core';
import { Http, Headers } from '@angular/http';
import 'rxjs/add/operator/toPromise';

import { Account } from './account';

export const BASE_URL = '/api/v1';

@Injectable()
export class AccountService {
  
  defaultHeaders: Headers = new Headers({contentType: 'application/json'});
  constructor(private http: Http){}

  signIn(principal: string, password: string): Promise<any> {
    return this.http
      .post(
        BASE_URL + '/sign-in', 
        { user_name: principal, user_password: password }, 
        { headers: this.defaultHeaders })
      .toPromise()
      .then(res=>res)
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
} 