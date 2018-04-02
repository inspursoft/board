import { HttpErrorResponse, HttpInterceptor, HttpRequest, HttpResponse } from '@angular/common/http'
import { Observable } from "rxjs/Observable";
import { HttpHandler } from "@angular/common/http/src/backend";
import { HttpEvent } from "@angular/common/http/src/response";
import { AppTokenService } from "../../app.init.service";
import { Injectable } from "@angular/core";
import "rxjs/add/operator/do";
import "rxjs/add/observable/throw"

@Injectable()
export class HttpClientInterceptor implements HttpInterceptor {

  constructor(private appTokenService: AppTokenService) {

  }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    let authReq: HttpRequest<any> = req.clone({
      headers: req.headers.set("contentType", 'application/json')
    });
    if (this.appTokenService.token) {
      authReq = authReq.clone({
        headers: authReq.headers.set("token", this.appTokenService.token)
      });
    }
    return next.handle(authReq)
      .do((event: HttpEvent<any>) => {
        if (event instanceof HttpResponse) {
          let res = event as HttpResponse<Object>;
          if (res.ok && res.headers.has("token")){
            this.appTokenService.chainResponse(res);
          }
        }
      }).catch((err: HttpErrorResponse) => {
        if (err.status >= 200 && err.status < 300) {
          const res = new HttpResponse({
            body: null,
            headers: err.headers,
            status: err.status,
            statusText: err.statusText,
            url: err.url
          });
          return Observable.of(res);
        } else {
          return Observable.throw(err);
        }
      });
  }
}