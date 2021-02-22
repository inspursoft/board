import { Injectable } from '@angular/core';
import { HttpResponse } from '@angular/common/http';
import { CookieService } from 'ngx-cookie';
import { GUIDE_STEP } from '../shared/shared.const';
import { SystemInfo, User } from '../shared/shared.types';
import { map } from 'rxjs/operators';
import { Observable } from 'rxjs';
import { AppTokenService } from './app-token.service';
import { ModelHttpClient } from '../shared/ui-model/model-http-client';

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
              private http: ModelHttpClient,
              private tokenService: AppTokenService) {
    this.systemInfo = new SystemInfo();
    this.currentUser = new User();
    this.cookieExpiry = new Date(Date.now() + 60 * 60 * 24 * 1000);
    this.isFirstLogin = localStorage.getItem('isFirstLogin') === null;
    if (this.isFirstLogin) {
      this.guideStep = GUIDE_STEP.PROJECT_LIST;
      localStorage.setItem('isFirstLogin', 'used');
    }
  }

  set token(t: string) {
    this.tokenService.token = t;
  }

  get token(): string {
    return this.tokenService.token;
  }

  get isSystemAdmin(): boolean {
    return this.currentUser && this.currentUser.userSystemAdmin === 1;
  }

  get isMipsSystem(): boolean {
    return this.systemInfo.processorType &&
      this.systemInfo.processorType.startsWith('mips64el');
  }

  get isArmSystem(): boolean {
    return this.systemInfo.processorType &&
      this.systemInfo.processorType.startsWith('aarch64');
  }

  get isNormalMode(): boolean {
    return this.systemInfo.mode === 'normal';
  }

  get isOpenBoard(): boolean {
    return this.systemInfo.boardVersion.endsWith('Openboard');
  }

  get getWebsocketPrefix(): string {
    return window.location.protocol === 'https:' ? 'wss' : 'ws';
  }

  get getHttpProtocol(): string {
    return window.location.protocol === 'https:' ? 'https' : 'http';
  }

  getCurrentUser(tokenParam?: string): Observable<User> {
    const token = this.tokenService.token || tokenParam;
    return this.http.getJson('/api/v1/users/current', User, {param: {token}})
      .pipe(map((res: User) => {
          this.currentUser = res;
          return res;
        })
      );
  }

  getSystemInfo(): Observable<any> {
    return this.http.getJson(`/api/v1/systeminfo`, SystemInfo)
      .pipe(map((res: SystemInfo) => {
          this.systemInfo = res;
          return this.systemInfo;
        })
      );
  }

  setAuditLog(auditData: IAuditOperationData): Observable<any> {
    return this.http.post('/api/v1/operations', auditData, {observe: 'response'});
  }

  getIsShowAdminServer(): Observable<any> {
    return this.http.get(`/api/v1/dashboard/adminservercheck`);
  }
}
