import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { User, DBInfo, UserVerify, MyToken } from './account.model';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';

const ACCOUNT_URL = '/v1/admin/account';
const BOOT_URL = '/v1/admin/boot';

@Injectable()
export class AccountService {

  constructor(private http: HttpClient) { }

  postSignIn(user: User): Observable<any> {
    return this.http.post(
      `${ACCOUNT_URL}/login/`,
      user.PostBody()
    ).pipe(map((res: HttpResponse<MyToken>) => {
      return res;
    }));
  }


  postSignUp(user: UserVerify): Observable<any> {
    return this.http.post(
      `${ACCOUNT_URL}/initialize/`,
      user.PostBody()
    );
  }

  checkInit(): Observable<any> {
    return this.http.get(`${ACCOUNT_URL}/install`, {
      observe: 'response',
    });
  }

  createUUID(): Observable<any> {
    return this.http.post(`${ACCOUNT_URL}/createUUID`, null);
  }

  validateUUID(uuid: string): Observable<any> {
    return this.http.post(
      `${ACCOUNT_URL}/ValidateUUID`,
      { UUID: uuid }
    );
  }

  initDB(dbInfo: DBInfo): Observable<any> {
    return this.http.post(
      `${BOOT_URL}/initdb`,
      dbInfo.PostBody()
    );
  }

  initSSH(user: UserVerify): Observable<any> {
    return this.http.post(
      `${BOOT_URL}/startdb`,
      user.PostBody()
    );
  }

  checkDB(): Observable<any> {
    return this.http.get(
      `${BOOT_URL}/checkdb`,
      { observe: 'response', }
    );
  }
}
