import { Injectable } from '@angular/core';
import { Http, Headers, Response } from '@angular/http';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';
import { NgXCookies } from 'ngx-cookies';
import { Message } from './shared/message-service/message';
import { MessageService } from './shared/message-service/message.service';

@Injectable()
export class AppInitService {
  
  tokenMessageSource: Subject<string> = new Subject<string>();
  tokenMessage$: Observable<string> = this.tokenMessageSource.asObservable()

  constructor(
    private http: Http
  ) {
    console.log('App initialized from current service.');
  }

  _tokenString: string;
  _currentLang: string;

  currentUser: {[key: string]: any} = null;
  systemInfo: {[key: string]: any} = null;
  
  set token(t: string) {
    this._tokenString = t;
  }

  get token(): string {
    return this._tokenString;
  }

  set currentLang(lang: string) {
    this._currentLang = lang;
  }

  get currentLang(): string {
    return this._currentLang;
  }

  chainResponse(r: Response): Response {
    this.token = r.headers.get('token');
    NgXCookies.setCookie("token", this.token);
    this.tokenMessageSource.next(this.token);
    return r;
  }

  getCurrentUser(tokenParam?: string): Promise<any> {
    let token = this.token || tokenParam || NgXCookies.getCookie("token") || '';
    return this.http
      .get('/api/v1/users/current', 
        { headers: new Headers({
          'Content-Type': 'application/json',
          'token': token
          }),
          params: {
           'token': token
          }
        })
      .toPromise()
      .then(res=>{
        this.chainResponse(res);
        this.currentUser = res.json();
        Promise.resolve(this.currentUser);
      })
      .catch(err=>Promise.reject(err));
  }

  getSystemInfo(): Promise<any> {
    return this.http
      .get(`/api/v1/systeminfo`)
      .toPromise()
      .then(res=>res.json())
      .catch(err=>Promise.reject(err));
  }
  
}