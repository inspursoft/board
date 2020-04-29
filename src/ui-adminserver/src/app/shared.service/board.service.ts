import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpResponse } from '@angular/common/http';
import { User } from '../account/account.model';
import { Observable, of, TimeoutError } from 'rxjs';
import { timeout, catchError } from 'rxjs/operators';
import { MessageService } from '../shared/message/message.service';

const BASE_URL = '/v1/admin';

@Injectable()
export class BoardService {

  constructor(private http: HttpClient,
              private messageService: MessageService,) { }

  applyCfg(user: User): Observable<any> {
    const token = window.sessionStorage.getItem('token');
    return this.http.post(
      `${BASE_URL}/board/applycfg?token=${token}`,
      user.PostBody()
    ).pipe(
      catchError((err: HttpErrorResponse | TimeoutError) => {
        if (err instanceof TimeoutError) {
          this.messageService.showOnlyOkDialog('ERROR.HTTP_TIME_OUT', 'GLOBAL_ALERT.WARNING');
          return of(null);
        } else {
          const res = new HttpResponse({
            body: err.message,
            headers: err.headers,
            status: err.status,
            statusText: err.statusText,
            url: err.url
          });
          return of(res);
        }
      })
    );
  }

  shutdown(user: User, uninstall: boolean): Observable<any> {
    const token = window.sessionStorage.getItem('token');
    return this.http.post(
      `${BASE_URL}/board/shutdown?token=${token}&uninstall=${uninstall}`,
      user.PostBody()
    );
  }

  start(user: User): Observable<any> {
    const token = window.sessionStorage.getItem('token');
    return this.http.post(
      `${BASE_URL}/board/start?token=${token}`,
      user.PostBody()
    );
  }
}
