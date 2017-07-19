import {Injectable} from "@angular/core";
import {Http, RequestOptions, Headers} from "@angular/http";
import {user} from "app/profile/user-center/user";
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

	constructor(private http: Http) {
	}

	deleteUser(user: user): Promise<object> {
		let options = new RequestOptions({
			headers: this.defaultHeaders
		});
		return this.http.delete(BASE_URL.concat("/users/").concat(user.user_id.toString()), options)
			.toPromise().then(
				res => {
					return Promise.resolve(res);
				},
				(reason: Response) => {
					return Promise.reject(UserService.getErrorMsg(reason, Array.from([400, 401, 403, 404]), "DEL"));
				});
	}

	getUser(userID: number): Promise<user> {
		let options = new RequestOptions({
			headers: this.defaultHeaders
		});
		return this.http.get(BASE_URL.concat("/users/").concat(userID.toString()), options)
			.toPromise()
			.then(
				res => {
					return Promise.resolve(res.json() as user);
				},
				reason => {
					return Promise.reject(UserService.getErrorMsg(reason, Array.from([401, 404]), "GET"));
				}
			);
	}

	updateUser(user: user): Promise<object> {
		let options = new RequestOptions({
			headers: this.defaultHeaders
		});
		return this.http.put(BASE_URL.concat("/users/").concat(user.user_id.toString()), user, options)
			.toPromise()
			.then(res => {
					return Promise.resolve(res);
				},
				reason => {
					return Promise.reject(UserService.getErrorMsg(reason, Array.from([400, 401, 403, 404]), "UPT"));
				});
	}

	newUser(userParams: user): Promise<object> {
		let options = new RequestOptions({
			headers: this.defaultHeaders
		});
		return this.http.post(BASE_URL.concat("/adduser"), userParams, options)
			.toPromise()
			.then(res => {
					return Promise.resolve(res);
				},
				(reason: Response) => {
					return Promise.reject(UserService.getErrorMsg(reason, Array.from([400, 403, 404, 409]), "ADD"));
				}).catch(err => {
				return Promise.reject(err);
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
		return this.http.get(BASE_URL.concat("/users"), options).toPromise().then(
			res => {
				return Promise.resolve(Array.from(res.json()) as user[]);
			},
			(reason: Response) => {
				return Promise.reject(UserService.getErrorMsg(reason, Array.from([400, 401, 403, 404]), "GET"));
			}
		)
	}
}