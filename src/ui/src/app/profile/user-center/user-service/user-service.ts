import { Injectable } from "@angular/core";
import { Http, RequestOptions, Headers } from "@angular/http";
import { user } from "app/profile/user-center/user";
import "rxjs/add/operator/toPromise";

const BASE_URL = "/api/v1";
@Injectable()
export class UserService {
  readonly defaultHeaders: Headers = new Headers({
    contentType: "application/json"
  });

  constructor(private http: Http) { }

  delete(user: user): Promise<object> {
    let options = new RequestOptions({
      headers: this.defaultHeaders
    });
    return this.http.delete(BASE_URL.concat("/users/").concat(user.user_id.toString()), options)
      .toPromise()
      .then(res => {
        console.log(res);
        return Promise.resolve(res);
      },
      reason => {
        return Promise.reject(reason);
      });
  }

  editUser(user: user): Promise<object> {
    let options = new RequestOptions({
      headers: this.defaultHeaders
    });
    return this.http.put(BASE_URL.concat("/users/").concat(user.user_id.toString()), user, options)
      .toPromise()
      .then(res => {
        return Promise.resolve(res);
      },
      reason => {
        return Promise.reject(reason);
      });
  }

  newUser(userParams: user): Promise<object> {
    let options = new RequestOptions({
      headers: this.defaultHeaders
    });
    return this.http.post(BASE_URL.concat("/sign-up"), userParams, options)
      .toPromise()
      .then(res => {
        return Promise.resolve(res);
      },
      reason => {
        return Promise.reject(reason);
      });
  }

  getUserList(
    username?: string,
    user_list_page: number = 0,
    user_list_page_size: number = 0
  ): Promise<Array<user>> {
    let params: URLSearchParams = new URLSearchParams();
    params.set("username", username);
    params.set("user_list_page", user_list_page.toString());
    params.set("user_list_page_size", user_list_page_size.toString());
    let options = new RequestOptions({
      headers: this.defaultHeaders,
      params: params
    });
    return this.http.get(BASE_URL.concat("/users"), options).toPromise().then(
      res => {
        return Promise.resolve(Array.from(res.json()));
      },
      reason => {
        return Promise.reject(reason);
      }
    );
  }

}
