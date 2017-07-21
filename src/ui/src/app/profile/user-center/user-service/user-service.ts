import { Injectable } from "@angular/core";
import { Http, RequestOptions, Headers, Response } from "@angular/http";
import { User } from "app/profile/user-center/user";
import { MessageService } from "../../../shared/message-service/message.service";
import "rxjs/add/operator/toPromise";

const BASE_URL = "/api/v1";
@Injectable()
export class UserService {
  readonly defaultHeaders: Headers = new Headers({
    contentType: "application/json"
  });

  static getErrorMsg(reason: Response | Error, statusArr: Array<number>, errorKey: string): string {
    if (reason instanceof Response){
      return statusArr.indexOf(reason.status) > -1 ?
        `USER_CENTER.${errorKey}_ERR_${reason.status}` :
        `${reason.status}:${reason.statusText}`;
    }
  }

  constructor(private http: Http,
              private messageService: MessageService) {
  }

  deleteUser(user: User): Promise<User> {
    let options = new RequestOptions({
      headers: this.defaultHeaders
    });
    return this.http.delete(`${BASE_URL}/users/${user.user_id}`, options).toPromise()
      .then(res => res.json())
      .catch(reason => {
        let errMsg: string = UserService.getErrorMsg(reason, Array.from([400, 401, 403, 404]), "DEL");
        this.messageService.dispatchError(reason, errMsg);
        return Promise.reject(errMsg);
      })
  }

  getUser(userID: number): Promise<User> {
    let options = new RequestOptions({
      headers: this.defaultHeaders
    });
    return this.http.get(`${BASE_URL}/users/${userID}`, options)
      .toPromise()
      .then(res => res.json())
      .catch(reason => {
        this.messageService.dispatchError(reason, UserService.getErrorMsg(reason, Array.from([401, 404]), "GET"));
        return Promise.reject(reason);
      });
  }

  updateUser(user: User): Promise<User> {
    let options = new RequestOptions({
      headers: this.defaultHeaders
    });
    return this.http.put(`${BASE_URL}/users/${user.user_id}`, user, options)
      .toPromise()
      .then(res => res.json())
      .catch(reason => {
        let r: string = "";
        if (reason instanceof Response) {
          r = UserService.getErrorMsg(reason, Array.from([400, 401, 403, 404]), "UPT");
        } else {
          this.messageService.dispatchError(reason);
        }
        return Promise.reject(r);
      });
  }

  newUser(userParams: User): Promise<User> {
    let options = new RequestOptions({
      headers: this.defaultHeaders
    });
    return this.http.post(`${BASE_URL}/adduser`, userParams, options).toPromise()
      .then(res => res.json())
      .catch(reason => {
        let r: string = "";
        if (reason instanceof Response) {
          r = UserService.getErrorMsg(reason, Array.from([400, 403, 404, 409]), "ADD");
        } else {
          this.messageService.dispatchError(reason);
        }
        return Promise.reject(r);
      });
  }

  getUserList(username?: string,
              user_list_page: number = 0,
              user_list_page_size: number = 0): Promise<User[]> {
    let params: Map<string, string> = new Map<string, string>();
    params["username"] = username;
    params["user_list_page"] = user_list_page.toString();
    params["user_list_page_size"] = user_list_page_size.toString();
    let options = new RequestOptions({
      headers: this.defaultHeaders,
      search: params
    });
    return this.http.get(`${BASE_URL}/users`, options).toPromise()
      .then(res => {
        return Array.from(res.json()) as User[];
      })
      .catch(reason => {
        this.messageService.dispatchError(reason, UserService.getErrorMsg(reason, Array.from([400, 401, 403, 404]), "GET"));
        return Promise.reject(reason);
      })
  }
}