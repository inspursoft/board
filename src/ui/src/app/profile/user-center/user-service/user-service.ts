import { Injectable } from "@angular/core";
import { Http, RequestOptions, Headers } from "@angular/http";
import { User } from "app/profile/user-center/user";
import { AppInitService } from "../../../app.init.service";
import "rxjs/add/operator/toPromise";

const BASE_URL = "/api/v1";
@Injectable()
export class UserService {
  readonly defaultHeaders: Headers = new Headers({
    contentType: "application/json"
  });

  constructor(private http: Http,
              private appInitService: AppInitService) {
  }

  deleteUser(user: User): Promise<boolean> {
    let options = new RequestOptions({
      headers: this.defaultHeaders,
      params: {'token': this.appInitService.token}
    });
    return this.http.delete(`${BASE_URL}/users/${user.user_id}`, options).toPromise()
      .then(res => res.ok)
      .catch(err => Promise.reject(err));
  }

  getUser(userID: number): Promise<User> {
    let options = new RequestOptions({
      headers: this.defaultHeaders,
      params: {'token': this.appInitService.token}
    });
    return this.http.get(`${BASE_URL}/users/${userID}`, options)
      .toPromise()
      .then(res => res.json())
      .catch(err => Promise.reject(err));
  }

  updateUser(user: User): Promise<boolean> {
    let options = new RequestOptions({
      headers: this.defaultHeaders,
      params: {'token': this.appInitService.token}
    });
    return this.http.put(`${BASE_URL}/users/${user.user_id}`, user, options)
      .toPromise()
      .then(res => res.ok)
      .catch(err => Promise.reject(err));
  }

  newUser(userParams: User): Promise<boolean> {
    let options = new RequestOptions({
      headers: this.defaultHeaders,
      params: {'token': this.appInitService.token}
    });
    return this.http.post(`${BASE_URL}/adduser`, userParams, options).toPromise()
      .then(res => res.ok)
      .catch(err => Promise.reject(err));
  }

  getUserList(username?: string,
              user_list_page: number = 0,
              user_list_page_size: number = 0): Promise<User[]> {
    let options = new RequestOptions({
      headers: this.defaultHeaders,
      params: {
        'username': username,
        'user_list_page': user_list_page.toString(),
        'user_list_page_size': user_list_page_size.toString(),
        'token': this.appInitService.token
      }
    });
    return this.http.get(`${BASE_URL}/users`, options).toPromise()
      .then(res => {
        return Array.from(res.json()) as User[];
      })
      .catch(err => Promise.reject(err))
  }
}