import { Injectable } from '@angular/core';
import { Http, Headers, Response } from '@angular/http';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';
import { Message } from './shared/message-service/message';
import { MessageService } from './shared/message-service/message.service';

@Injectable()
export class AppInitService {
  
  tokenMessageSource: Subject<string> = new Subject<string>();
  tokenMessage$: Observable<string> = this.tokenMessageSource.asObservable()

  constructor(private http: Http) {
    console.log('App initialized from current service.');
  }

  _tokenString: string;

  currentUser: {[key: string]: any};

  set token(t: string) {
    this._tokenString = t;
  }

  get token(): string {
    return this._tokenString;
  }

  chainResponse(r: Response): Response {
    this.token = r.headers.get('token');
    this.tokenMessageSource.next(this.token);
    return r;
  }

  getCurrentUser(tokenParam?: string): Promise<any> {
    return this.http
      .get('/api/v1/users/current', 
        { headers: new Headers({
          'Content-Type': 'application/json',
          'token': this.token || tokenParam || ''
          }),
          params: {
           'token': this.token || tokenParam || ''
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
  
}