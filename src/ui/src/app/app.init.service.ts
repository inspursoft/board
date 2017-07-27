import { Injectable } from '@angular/core';
import { Http, Headers } from '@angular/http';

@Injectable()
export class AppInitService {
  defaultHeaders = new Headers({
    'Content-Type': 'application/json'
  });

  constructor(private http: Http) {
    console.log('App initialized from current service.');
  }

  _tokenString: string;

  set token(t: string) {
    this._tokenString = t;
  }

  get token(): string {
    return this._tokenString;
  }

  currentUser(): Promise<any> {
    return this.http
      .get('/api/v1/users/current', 
        { headers: this.defaultHeaders, 
          params: {
            'token': this.token
          }
        })
      .toPromise()
      .then(res=>res.json())
      .catch(err=>Promise.reject(err));
  }
  
}