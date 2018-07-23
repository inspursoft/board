import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';
import { CookieService } from 'ngx-cookie';
import { GUIDE_STEP } from "./shared/shared.const";

export interface IAuditOperationData {
  operation_id?: number,
  operation_creation_time?: string,
  operation_update_time?: string,
  operation_deleted?: number,
  operation_user_id: number,
  operation_user_name: string,
  operation_project_name: string,
  operation_project_id: number,
  operation_tag?: string,
  operation_comment?: string,
  operation_object_type: string,
  operation_object_name: string,
  operation_action: string,
  operation_status: string,
  operation_path?: string
}

@Injectable()
export class AppTokenService {
  _tokenString: string | null;
  tokenMessageSource: Subject<string> = new Subject<string>();
  tokenMessage$: Observable<string> = this.tokenMessageSource.asObservable();

  constructor(private cookieService: CookieService) {
  }

  set token(t: string | null) {
    this._tokenString = t;
  }

  get token(): string | null {
    return this._tokenString;
  }

  chainResponse(r: HttpResponse<Object>): HttpResponse<Object> {
    this.token = r.headers.get('token');
    this.cookieService.put("token", this.token == null ? "" : this.token);
    this.tokenMessageSource.next(this.token);
    return r;
  }
}

@Injectable()
export class AppInitService {
  _isFirstLogin: boolean = false;
  _currentLang: string;
  grafanaViewUrl: string = "";
  cookieExpiry: Date = new Date(Date.now() + 10 * 60 * 60 * 24 * 365 * 1000);
  guideStep: GUIDE_STEP;
  currentUser: {[key: string]: any} = null;
  systemInfo: {[key: string]: any} = null;

  constructor(private cookieService: CookieService,
              private http: HttpClient,
              private tokenService: AppTokenService) {
    console.log('App initialized from current service.');
    this._isFirstLogin = this.cookieService.get("isFirstLogin") == undefined;
    if (this._isFirstLogin) {
      this.guideStep = GUIDE_STEP.PROJECT_LIST;
      this.cookieService.put("isFirstLogin", "used", {expires: this.cookieExpiry});
    }
  }

  set token(t: string) {
    this.tokenService.token = t;
  }

  get token(): string {
    return this.tokenService.token;
  }

  set currentLang(lang: string) {
    this._currentLang = lang;
  }

  get currentLang(): string {
    return this._currentLang;
  }

  get isFirstLogin(): boolean {
    return this._isFirstLogin;
  }

  getCurrentUser(tokenParam?: string): Promise<any> {
    let token = this.tokenService.token || tokenParam || this.cookieService.get("token") || '';
    return this.http
      .get('/api/v1/users/current',
        {
          observe: "response",
          params: {
            'token': token
          }
        })
      .toPromise()
      .then(res => {
        this.currentUser = res.body;
        return res.body;
      });
  }

  getSystemInfo(): Promise<any> {
    return this.http
      .get(`/api/v1/systeminfo`, {observe: "response"})
      .toPromise()
      .then(res => {
        this.systemInfo = res.body;
        return this.systemInfo
      });
  }

  setAuditLog(auditData: IAuditOperationData): Observable<any> {
    return this.http.post('/api/v1/operations', auditData, {observe: "response"})
  }

}