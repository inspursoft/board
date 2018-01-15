import { HttpInterceptor, HttpRequest, HttpResponse } from '@angular/common/http'
import { Observable } from "rxjs/Observable";
import { HttpHandler } from "@angular/common/http/src/backend";
import { HttpEvent } from "@angular/common/http/src/response";
import { AppTokenService } from "../../app.init.service";
import { Injectable } from "@angular/core";
import "rxjs/add/operator/do";

@Injectable()
export class HttpClientInterceptor implements HttpInterceptor {

  constructor(private appTokenService: AppTokenService) {

  }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    const noNeedTokenUrls = Array<string>();
    noNeedTokenUrls.push("/api/v1/sign-in");
    let authReq: HttpRequest<any> = req.clone({
      headers: req.headers.set("contentType", 'application/json')
    });
    if (noNeedTokenUrls.indexOf(req.url) == -1 && this.appTokenService.token) {
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
      })
  }
}