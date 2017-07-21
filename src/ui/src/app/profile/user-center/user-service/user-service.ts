import { Injectable } from "@angular/core";
import { Http, RequestOptions, Headers, Response } from "@angular/http";
import { user } from "app/profile/user-center/user";
import { MessageService } from "../../../shared/message-service/message.service";
import { Message } from "../../../shared/message-service/message";
import "rxjs/add/operator/toPromise";

const BASE_URL = "/api/v1";
@Injectable()
export class UserService {
  readonly defaultHeaders: Headers = new Headers({
    contentType: "application/json"
  });

  static getErrorMsg(reason: Response, statusArr: Array<number>, errorKey: string): string {
    return statusArr.indexOf(reason.status) > -1 ?
      `USER_CENTER.${errorKey}_ERR_${reason.status}` :
      `${reason.status}:${reason.statusText}`;
  }

  handleErrMsg(reason: Response, statusArr: Array<number>, errorKey: string): void {
    let m: Message = new Message();
    if (reason.status == 500) {
      m.message = "USER_CENTER.ERR_500";
      this.messageService.globalMessage(m);
    } else {
      m.message = UserService.getErrorMsg(reason, statusArr, errorKey);
      this.messageService.inlineAlertMessage(m);
    }
  }

  constructor(private http: Http,
              private messageService: MessageService) {
  }

  deleteUser(user: user): Promise<user> {
    let options = new RequestOptions({
      headers: this.defaultHeaders
    });
    return this.http.delete(`${BASE_URL}/users/${user.user_id}`, options).toPromise()
      .then(res => res.json())
      .catch(reason => {
        if (reason instanceof Response) {
          this.handleErrMsg(reason, Array.from([400, 401, 403, 404]), "DEL");
        } else if (reason instanceof Error) {
          console.error(`name:${(<Error>reason).name};message:${(<Error>reason).message}`);
        } else {
          console.error(reason);
        }
        return Promise.reject(reason);
      })
  }

  getUser(userID: number): Promise<user> {
    let options = new RequestOptions({
      headers: this.defaultHeaders
    });
    return this.http.get(`${BASE_URL}/users/${userID}`, options)
      .toPromise()
      .then(res => res.json())
      .catch(reason => {
        if (reason instanceof Response) {
          this.handleErrMsg(reason, Array.from([401, 404]), "GET");
        } else if (reason instanceof Error) {
          console.error(`name:${(<Error>reason).name};message:${(<Error>reason).message}`);
        } else {
          console.error(reason);
        }
        return Promise.reject(reason);
      });
  }

  updateUser(user: user): Promise<user> {
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
          r = `name:${(<Error>reason).name};message:${(<Error>reason).message}`;
        }
        return Promise.reject(r);
      });
  }

  newUser(userParams: user): Promise<user> {
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
          r = `name:${(<Error>reason).name};message:${(<Error>reason).message}`;
        }
        return Promise.reject(r);
      });
  }

  getUserList(username?: string,
              user_list_page: number = 0,
              user_list_page_size: number = 0): Promise<user[]> {
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
        return Array.from(res.json()) as user[];
      })
      .catch(reason => {
        if (reason instanceof Response) {
          this.handleErrMsg(reason, Array.from([400, 401, 403, 404]), "GET");
        } else if (reason instanceof Error) {
          console.error(`name:${(<Error>reason).name};message:${(<Error>reason).message}`);
        } else {
          console.error(reason);
        }
        return Promise.reject(reason);
      })
  }
}