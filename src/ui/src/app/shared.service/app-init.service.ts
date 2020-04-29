import { Injectable } from '@angular/core';
import { HttpClient, HttpResponse } from '@angular/common/http';
import { CookieService } from 'ngx-cookie';
import { GUIDE_STEP } from '../shared/shared.const';
import { SystemInfo, User } from '../shared/shared.types';
import { map } from 'rxjs/operators';
import { Observable, Subject } from 'rxjs';
import { AppTokenService } from './app-token.service';

export interface IAuditOperationData {
  operation_id?: number;
  operation_creation_time?: string;
  operation_update_time?: string;
  operation_deleted?: number;
  operation_user_id: number;
  operation_user_name: string;
  operation_project_name: string;
  operation_project_id: number;
  operation_tag?: string;
  operation_comment?: string;
  operation_object_type: string;
  operation_object_name: string;
  operation_action: string;
  operation_status: string;
  operation_path?: string;
}

@Injectable()
export class AppInitService {
  isFirstLogin = false;
  currentLang: string;
  cookieExpiry: Date;
  guideStep: GUIDE_STEP;
  currentUser: User;
  systemInfo: SystemInfo;

  constructor(private cookieService: CookieService,
              private http: HttpClient,
              private tokenService: AppTokenService) {
    this.systemInfo = new SystemInfo();
    this.currentUser = new User();
    this.cookieExpiry = new Date(Date.now() + 10 * 60 * 60 * 24 * 365 * 1000);
    this.isFirstLogin = this.cookieService.get('isFirstLogin') === undefined;
    if (this.isFirstLogin) {
      this.guideStep = GUIDE_STEP.PROJECT_LIST;
      this.cookieService.put('isFirstLogin', 'used', {expires: this.cookieExpiry});
    }
  }

  set token(t: string) {
    this.tokenService.token = t;
  }

  get token(): string {
    return this.tokenService.token;
  }

  get isSystemAdmin(): boolean {
    return this.currentUser && this.currentUser.user_system_admin === 1;
  }

  get isMipsSystem(): boolean {
    return this.systemInfo.processor_type &&
      this.systemInfo.processor_type.startsWith('mips64el');
  }

  get isArmSystem(): boolean {
    return this.systemInfo.processor_type &&
      this.systemInfo.processor_type.startsWith('aarch64');
  }

  get isNormalMode(): boolean {
    return this.systemInfo.mode === 'normal';
  }

  getCurrentUser(tokenParam?: string): Observable<User> {
    const token = this.tokenService.token || tokenParam;
    return this.http.get<User>('/api/v1/users/current', {observe: 'response', params: {token}})
      .pipe(map((res: HttpResponse<User>) => {
        this.currentUser = res.body;
        return res.body;
      }));
  }

  getSystemInfo(): Observable<any> {
    return this.http.get(`/api/v1/systeminfo`).pipe(map((res: SystemInfo) => {
      this.systemInfo = res;
      return this.systemInfo;
    }));
  }

  setAuditLog(auditData: IAuditOperationData): Observable<any> {
    return this.http.post('/api/v1/operations', auditData, {observe: 'response'});
  }

}
