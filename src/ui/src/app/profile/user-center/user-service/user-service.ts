import { Injectable } from "@angular/core";
import { Http, RequestOptions, Headers } from "@angular/http";
import { User } from "app/profile/user-center/user";
import { AppInitService } from "../../../app.init.service";
import "rxjs/add/operator/toPromise";

const BASE_URL = "/api/v1";
@Injectable()
export class UserService {
  get defaultHeader(): Headers {
    let headers = new Headers();
    headers.append('Content-Type', 'application/json');
    headers.append('token', this.appInitService.token);
    return headers;
  }

  constructor(private http: Http,
              private appInitService: AppInitService) {
  }

  deleteUser(user: User): Promise<boolean> {
    let options = new RequestOptions({
      headers: this.defaultHeader
    });
    return this.http.delete(`${BASE_URL}/users/${user.user_id}`, options).toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.ok;
      })
      .catch(err => Promise.reject(err));
  }

  getUser(userID: number): Promise<User> {
    let options = new RequestOptions({
      headers: this.defaultHeader
    });
    return this.http.get(`${BASE_URL}/users/${userID}`, options)
      .toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.json();
      })
      .catch(err => Promise.reject(err));
  }

  changeUserPassword(userID: number, user_password_old: string, user_password_new: string): Promise<boolean> {
    let options = new RequestOptions({
      headers: this.defaultHeader
    });
    let body = {
      "user_password_old": user_password_old,
      "user_password_new": user_password_new
    };
    return this.http.put(`${BASE_URL}/users/${userID}/password`, body, options).toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.ok;
      })
      .catch(err => Promise.reject(err));
  }

  updateUser(user: User): Promise<boolean> {
    let options = new RequestOptions({
      headers: this.defaultHeader
    });
    return this.http.put(`${BASE_URL}/users/${user.user_id}`, user, options)
      .toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.ok;
      })
      .catch(err => Promise.reject(err));
  }

  newUser(userParams: User): Promise<boolean> {
    let options = new RequestOptions({
      headers: this.defaultHeader
    });
    return this.http.post(`${BASE_URL}/adduser`, userParams, options).toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.ok;
      })
      .catch(err => Promise.reject(err));
  }

  getUserList(username?: string,
              user_list_page: number = 0,
              user_list_page_size: number = 0): Promise<User[]> {
    let options = new RequestOptions({
      headers: this.defaultHeader,
      params: {
        'username': username,
        'user_list_page': user_list_page.toString(),
        'user_list_page_size': user_list_page_size.toString()
      }
    });
    return this.http.get(`${BASE_URL}/users`, options).toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return Array.from(res.json()) as User[];
      })
      .catch(err => Promise.reject(err))
  }

  setUserSystemAdmin(userID: number, userSystemAdmin: number): Promise<boolean> {
    let options = new RequestOptions({
      headers: this.defaultHeader
    });
    return this.http.put(`${BASE_URL}/users/${userID}/systemadmin`, {user_system_admin: userSystemAdmin}, options).toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.status == 200;
      })
      .catch(err => Promise.reject(err));
  }

  setUserProjectAdmin(userID: number, userProjectAdmin: number): Promise<boolean> {
    let options = new RequestOptions({
      headers: this.defaultHeader
    });
    return this.http.put(`${BASE_URL}/users/${userID}/projectadmin`, {user_project_admin: userProjectAdmin}, options).toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.status == 200;
      })
      .catch(err => Promise.reject(err));
  }

  usesChangeAccount(user: User): Promise<boolean> {
    let options = new RequestOptions({
      headers: this.defaultHeader
    });
    return this.http.put(`${BASE_URL}/users/changeaccount`, user, options).toPromise()
      .then(res => {
        this.appInitService.chainResponse(res);
        return res.status == 200;
      })
      .catch(err => Promise.reject(err));
  }
}