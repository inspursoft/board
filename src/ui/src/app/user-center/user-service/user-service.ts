import { Injectable } from "@angular/core";
import { HttpClient, HttpHeaders, HttpResponse } from "@angular/common/http";
import { AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE } from "../../shared/shared.const";
import { User } from "../../shared/shared.types";
import { Observable } from "rxjs/Observable";
import "rxjs/add/operator/map"

const BASE_URL = "/api/v1";

@Injectable()
export class UserService {

  constructor(private http: HttpClient) {
  }

  deleteUser(user: User): Observable<any> {
    return this.http.delete(`${BASE_URL}/users/${user.user_id}`, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: "response"
    });
  }

  getUser(userID: number): Observable<User> {
    return this.http.get(`${BASE_URL}/users/${userID}`, {observe: "response"})
      .map((res: HttpResponse<User>) => res.body);
  }

  changeUserPassword(userID: number, user_password_old: string, user_password_new: string): Observable<any> {
    let body = {
      "user_password_old": user_password_old,
      "user_password_new": user_password_new
    };
    return this.http.put(`${BASE_URL}/users/${userID}/password`, body, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: "response"
    });
  }

  updateUser(user: User): Observable<any> {
    return this.http.put(`${BASE_URL}/users/${user.user_id}`, user, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: "response"
    })
  }

  newUser(userParams: User): Observable<any> {
    return this.http.post(`${BASE_URL}/adduser`, userParams, {
      headers: new HttpHeaders().set(AUDIT_RECORD_HEADER_KEY, AUDIT_RECORD_HEADER_VALUE),
      observe: "response"
    });
  }

  getUserList(username: string, pageIndex: number, pageSize: number,sortBy: string, isReverse: boolean): Observable<Object> {
    return this.http.get(`${BASE_URL}/users`, {
      observe: "response",
      params: {
        'username': username,
        'page_index': pageIndex.toString(),
        'page_size': pageSize.toString(),
        'order_field': sortBy,
        'order_asc': isReverse ? "0" : "1"
      }
    }).map((res: HttpResponse<Object>) => res.body);
  }

  setUserSystemAdmin(userID: number, userSystemAdmin: number): Observable<any> {
    return this.http.put(`${BASE_URL}/users/${userID}/systemadmin`, {user_system_admin: userSystemAdmin}, {observe:"response"});
  }

  usesChangeAccount(user: User): Observable<any> {
    return this.http.put(`${BASE_URL}/users/changeaccount`, user, {observe:"response"});
  }
}