import { Injectable } from '@angular/core';
import { HttpHeaders, HttpResponse } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE } from '../../../shared/shared.const';
import { User } from '../../../shared/shared.types';
import { ModelHttpClient } from '../../../shared/ui-model/model-http-client';
import { UserPagination } from '../../admin.types';

const BASE_URL = '/api/v1';

@Injectable()
export class UserService {

  constructor(private http: ModelHttpClient) {
  }

  deleteUser(user: User): Observable<any> {
    return this.http.delete(`${BASE_URL}/users/${user.userId}`, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  getUser(userID: number): Observable<User> {
    return this.http.getJson(`${BASE_URL}/users/${userID}`, User);
  }

  changeUserPassword(userID: number, userPasswordOld: string, userPasswordNew: string): Observable<any> {
    const body = {
      user_password_old: window.btoa(userPasswordOld),
      user_password_new: window.btoa(userPasswordNew)
    };
    return this.http.put(`${BASE_URL}/users/${userID}/password`, body, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE)
    });
  }

  updateUser(user: User): Observable<any> {
    return this.http.put(`${BASE_URL}/users/${user.userId}`, user.getPostBody(), {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  newUser(user: User): Observable<any> {
    user.userPassword = window.btoa(user.userPassword);
    user.userConfirmPassword = window.btoa(user.userConfirmPassword);
    return this.http.post(`${BASE_URL}/adduser`, user.getPostBody(), {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: 'response'
    });
  }

  getUserList(username: string, pageIndex: number, pageSize: number, sortBy: string, isReverse: boolean): Observable<Object> {
    return this.http.getPagination(`${BASE_URL}/users`, UserPagination, {
      param: {
        username,
        page_index: pageIndex.toString(),
        page_size: pageSize.toString(),
        order_field: sortBy,
        order_asc: isReverse ? '0' : '1'
      }
    });
  }

  setUserSystemAdmin(userID: number, userSystemAdmin: number): Observable<any> {
    return this.http.put(`${BASE_URL}/users/${userID}/systemadmin`, {user_system_admin: userSystemAdmin}, {observe: 'response'});
  }

  usesChangeAccount(user: User): Observable<any> {
    return this.http.put(`${BASE_URL}/users/changeaccount`, user.getPostBody(), {observe: 'response'});
  }
}
