import { HttpErrorResponse, HttpInterceptor, HttpRequest, HttpResponse } from '@angular/common/http'
import { HttpHandler } from "@angular/common/http/src/backend";
import { HttpEvent } from "@angular/common/http/src/response";
import { AppTokenService } from "../../app.init.service";
import { Injectable } from "@angular/core";
import { MessageService } from "../message-service/message.service";
import { GlobalAlertType } from "../shared.types";
import { TranslateService } from "@ngx-translate/core";
import { Observable, of, throwError, TimeoutError } from "rxjs";
import { catchError, tap, timeout } from "rxjs/operators";

@Injectable()
export class HttpClientInterceptor implements HttpInterceptor {

  constructor(private appTokenService: AppTokenService,
              private messageService: MessageService,
              private translateService: TranslateService) {

  }

  intercept(req: HttpRequest<any>, next: HttpHandler): Observable<HttpEvent<any>> {
    let authReq: HttpRequest<any> = req.clone({
      headers: req.headers
    });
    if (this.appTokenService.token) {
      authReq = authReq.clone({
        headers: authReq.headers.set("token", this.appTokenService.token),
        params: authReq.params.set("Timestamp", Date.now().toString())
      });
    }
    return next.handle(authReq)
      .pipe(
        tap((event: HttpEvent<any>) => {
          if (event instanceof HttpResponse) {
            let res = event as HttpResponse<Object>;
            if (res.ok && res.headers.has("token")) {
              this.appTokenService.chainResponse(res);
            }
          }
        }), timeout(30000),
        catchError((err: HttpErrorResponse | TimeoutError) => {
          if (err instanceof HttpErrorResponse) {
            if (err.status >= 200 && err.status < 300) {
              const res = new HttpResponse({
                body: null,
                headers: err.headers,
                status: err.status,
                statusText: err.statusText,
                url: err.url
              });
              return of(res);
            } else if (err.status == 502) {
              window.location.replace('/bad-gateway-page');
            } else if (err.status == 504) {
              this.messageService.showGlobalMessage('ERROR.HTTP_504', {
                globalAlertType: GlobalAlertType.gatShowDetail,
                errorObject: err
              });
            } else if (err.status == 500) {
              this.messageService.showGlobalMessage('ERROR.HTTP_500', {
                globalAlertType: GlobalAlertType.gatShowDetail,
                errorObject: err
              });
            } else if (err.status == 400) {
              this.messageService.showGlobalMessage(`ERROR.HTTP_400`, {
                globalAlertType: GlobalAlertType.gatShowDetail,
                errorObject: err
              });
            } else if (err.status == 401 && this.appTokenService.token) {
              this.messageService.showGlobalMessage(`ERROR.HTTP_401`, {
                globalAlertType: GlobalAlertType.gatLogin,
                alertType: 'alert-warning'
              });
            } else if (err.status == 403) {
              this.messageService.showAlert(`ERROR.HTTP_403`, {alertType: 'alert-danger'});
            } else if (err.status == 404) {
              this.messageService.showAlert(`ERROR.HTTP_404`, {alertType: 'alert-danger'});
            } else if (err.status == 412) {
              this.messageService.showAlert(`ERROR.HTTP_412`, {alertType: 'alert-warning'});
            } else if (err.status == 422) {
              this.translateService.get(`ERROR.HTTP_422`).subscribe((msg: string) => {
                let alertMsg = `${msg},${err.error}`;
                this.messageService.showAlert(alertMsg, {alertType: 'alert-danger'});
              });
            } else if (this.appTokenService.token) {
              this.messageService.showGlobalMessage(`ERROR.HTTP_UNK`, {
                globalAlertType: GlobalAlertType.gatShowDetail,
                errorObject: err
              });
            }
          } else {
            this.messageService.showGlobalMessage(`ERROR.HTTP_TIME_OUT`, {
              globalAlertType: GlobalAlertType.gatShowDetail,
              errorObject: err,
              endMessage: req.url
            });
          }
          return throwError(err);
        }));
  }
}
