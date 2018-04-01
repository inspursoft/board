import { Injectable } from "@angular/core";
import { HttpClient, HttpResponse } from "@angular/common/http";
import { User } from "../user";

const BASE_URL = "/api/v1";

@Injectable()
export class UserService {

  constructor(private http: HttpClient) {
  }

  deleteUser(user: User): Promise<any> {
    return this.http.delete(`${BASE_URL}/users/${user.user_id}`, {observe: "response"}).toPromise();
  }

  getCurrentUser(): Promise<User> {
    return this.http.get(`${BASE_URL}/users/current`, {observe: "response"})
      .toPromise()
      .then((res: HttpResponse<User>) => res.body);
  }

  getUser(userID: number): Promise<User> {
    return this.http.get(`${BASE_URL}/users/${userID}`, {observe: "response"})
      .toPromise()
      .then((res: HttpResponse<User>) => res.body);
  }

  changeUserPassword(userID: number, user_password_old: string, user_password_new: string): Promise<any> {
    let body = {
      "user_password_old": user_password_old,
      "user_password_new": user_password_new
    };
    return this.http.put(`${BASE_URL}/users/${userID}/password`, body, {observe: "response"}).toPromise();
  }

  updateUser(user: User): Promise<any> {
    return this.http.put(`${BASE_URL}/users/${user.user_id}`, user, {observe: "response"}).toPromise()
  }

  newUser(userParams: User): Promise<any> {
    return this.http.post(`${BASE_URL}/adduser`, userParams, {observe: "response"}).toPromise();
  }

  getUserList(username: string, pageIndex: number, pageSize: number,sortBy: string, isReverse: boolean): Promise<Object> {
    return this.http.get(`${BASE_URL}/users`, {
      observe: "response",
      params: {
        'username': username,
        'page_index': pageIndex.toString(),
        'page_size': pageSize.toString(),
        'order_field': sortBy,
        'order_asc': isReverse ? "0" : "1"
      }
    }).toPromise()
      .then((res: HttpResponse<Object>) => res.body);
  }

  setUserSystemAdmin(userID: number, userSystemAdmin: number): Promise<any> {
    return this.http.put(`${BASE_URL}/users/${userID}/systemadmin`, {user_system_admin: userSystemAdmin}, {observe:"response"}).toPromise();
  }

  usesChangeAccount(user: User): Promise<any> {
    return this.http.put(`${BASE_URL}/users/changeaccount`, user, {observe:"response"}).toPromise();
  }
}