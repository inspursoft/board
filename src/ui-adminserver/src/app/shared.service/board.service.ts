import { Injectable } from '@angular/core';
import { HttpClient, HttpErrorResponse, HttpResponse } from '@angular/common/http';
import { User } from '../account/account.model';
import { Observable, of, TimeoutError, throwError } from 'rxjs';
import { timeout, catchError } from 'rxjs/operators';
import { MessageService } from '../shared/message/message.service';

const BASE_URL = '/v1/admin';

@Injectable()
export class BoardService {

  constructor(private http: HttpClient,
              private messageService: MessageService,) { }

  applyCfg(user: User): Observable<any> {
    return this.http.post(
      `${BASE_URL}/board/applycfg`,
      user.PostBody()
    ).pipe(
      catchError((err: HttpErrorResponse | TimeoutError) => {
        if (err instanceof TimeoutError) {
          this.messageService.showOnlyOkDialog('ERROR.HTTP_TIME_OUT', 'GLOBAL_ALERT.WARNING');
        }
        return throwError(err);
      })
    );
  }

  shutdown(user: User, uninstall: boolean): Observable<any> {
    return this.http.post(
      `${BASE_URL}/board/shutdown?uninstall=${uninstall}`,
      user.PostBody()
    );
  }

  start(user: User): Observable<any> {
    return this.http.post(
      `${BASE_URL}/board/start`,
      user.PostBody()
    );
  }
}
