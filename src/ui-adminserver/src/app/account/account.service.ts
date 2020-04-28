import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { User, MyToken } from './account.model';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

const ACCOUNT_URL = '/v1/admin/account';

@Injectable()
export class AccountService {

  constructor(private http: HttpClient) { }

  signIn(user: User): Observable<any> {
    return this.http.post(
      `${ACCOUNT_URL}/login/`,
      user.PostBody()
    ).pipe(map((res: HttpResponse<MyToken>) => {
      return res;
    }));
  }

  createUUID(): Observable<any> {
    return this.http.post(`${ACCOUNT_URL}/createUUID`, null);
  }

}
