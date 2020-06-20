import { Injectable } from '@angular/core';
import { CookieService } from 'ngx-cookie';
import { HttpResponse } from '@angular/common/http';

@Injectable()
export class AppTokenService {
  tokenOrigin: string | null = '';

  constructor() {

  }

  get token(): string | null {
    if (this.tokenOrigin === '') {
      this.token = localStorage.getItem('token');
    }
    return this.tokenOrigin;
  }

  set token(tokenValue: string | null) {
    this.tokenOrigin = tokenValue || '';
  }

  chainResponse(r: HttpResponse<object>): HttpResponse<object> {
    this.token = r.headers.get('token');
    localStorage.setItem('token', this.token);
    return r;
  }
}
