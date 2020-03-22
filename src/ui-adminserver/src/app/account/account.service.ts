import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { User } from './account.model';
import { Observable } from 'rxjs';

const BASE_URL = '/v1/admin/account';
const BOOT_URL = '/v1/admin/boot';

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

  checkInit(): Observable<any> {
    return this.http.get(`${BASE_URL}/install`, {
      observe: 'response',
    });
  }

  createUUID(): Observable<any> {
    return this.http.post(`${BASE_URL}/createUUID`, null);
  }

  validateUUID(uuid: string): Observable<any> {
    return this.http.post(
      `${BASE_URL}/ValidateUUID`,
      uuid
    );
  }

  initDB(dbPwd: string): Observable<any> {
    return this.http.post(
      `${BOOT_URL}/initdb`,
      dbPwd
    );
  }

  initSSH(user: User): Observable<any> {
    return this.http.post(
      `${BOOT_URL}/startdb`,
      user.PostBody()
    );
  }
}
