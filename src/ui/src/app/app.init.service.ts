import { HostListener, Injectable } from '@angular/core';
import { Http, Headers, Response } from '@angular/http';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';
import { CookieService } from 'ngx-cookie';
import { GUIDE_STEP } from "./shared/shared.const";

@Injectable()
export class AppInitService {
  
  tokenMessageSource: Subject<string> = new Subject<string>();
  tokenMessage$: Observable<string> = this.tokenMessageSource.asObservable();
  cookieExpiry: Date = new Date(Date.now() + 10 * 60 * 60 * 24 * 365 * 1000);
  guideStep:GUIDE_STEP;
  _isFirstLogin: boolean = false;

  constructor(
    private cookieService:CookieService,
    private http: Http
  ) {
    console.log('App initialized from current service.');
    this._isFirstLogin = this.cookieService.get("isFirstLogin") == undefined;
    if (this._isFirstLogin){
      this.guideStep = GUIDE_STEP.PROJECT_LIST;
      this.cookieService.put("isFirstLogin","used",{expires: this.cookieExpiry});
    }
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

  get isFirstLogin(): boolean{
    return this._isFirstLogin;
  }

  chainResponse(r: Response): Response {
    this.token = r.headers.get('token');
    this.cookieService.put("token", this.token);
    this.tokenMessageSource.next(this.token);
    return r;
  }

  getCurrentUser(tokenParam?: string): Promise<any> {
    let token = this.token || tokenParam || this.cookieService.get("token") || '';
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