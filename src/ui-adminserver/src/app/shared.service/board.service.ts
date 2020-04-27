import { Injectable } from '@angular/core';
import { HttpClient } from '@angular/common/http';
import { User } from '../account/account.model';
import { Observable } from 'rxjs';

const BASE_URL = '/v1/admin';

@Injectable()
export class BoardService {

  constructor(private http: HttpClient) { }

  applyCfg(user: User): Observable<any> {
    const token = window.sessionStorage.getItem('token');
    return this.http.post(
      `${BASE_URL}/board/applycfg?token=${token}`,
      user.PostBody()
    );
  }

  shutdown(user: User, uninstall: boolean): Observable<any> {
    const token = window.sessionStorage.getItem('token');
    return this.http.post(
      `${BASE_URL}/board/shutdown?token=${token}&uninstall=${uninstall}`,
      user.PostBody()
    );
  }

  start(user: User): Observable<any> {
    const token = window.sessionStorage.getItem('token');
    return this.http.post(
      `${BASE_URL}/board/start?token=${token}`,
      user.PostBody()
    );
  }
}
