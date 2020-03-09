import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { User } from './account.model';
import { Observable } from 'rxjs';

const BASE_URL = '/v1/admin/account';

@Injectable()
export class AccountService {

  constructor(private http: HttpClient) { }

  postSignIn(user: User): Observable<any> {
    return this.http.post(
      `${BASE_URL}/login/`,
      user.PostBody()
    );
  }

  postSignUp(user: User): Observable<any> {
    return this.http.post(
      `${BASE_URL}/initialize/`,
      user.PostBody()
    );
  }

  getVerify(alpha: string): Observable<any> {
    return this.http.get(`${BASE_URL}/verify/?alpha=${alpha}`, {
      observe: 'response',
    });
  }
}
