import { HttpClient, HttpHandler } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { ResponseArrayBase, ResponseBase, ResponsePaginationBase } from './shared-model-types';

export class CsHttpClient extends HttpClient {
  defaultHeaders = new Headers({
    'Content-Type': 'application/json'
  });

  constructor(handler: HttpHandler) {
    super(handler);
  }

  getJson<T extends ResponseBase>(url: string,
                                  returnType: new(res: object) => T,
                                  param: { [param: string]: string }): Observable<T> {
    return super.get(url, {observe: 'body', responseType: 'json', params: param})
      .pipe(map((res: object) => new returnType(res)));
  }

  getPagination<T extends ResponsePaginationBase<ResponseBase>>(url: string,
                                                                paginationType: new(res: object) => T,
                                                                param?: { [param: string]: string }): Observable<T> {
    return super.get(url, {observe: 'body', responseType: 'json', params: param})
      .pipe(map((res: object) => new paginationType(res)));
  }

  getArrayJson<T extends ResponseArrayBase<ResponseBase>>(url: string,
                                                          arrayType: new(res: object) => T,
                                                          param?: { [param: string]: string }): Observable<T> {
    return super.get(url, {observe: 'body', responseType: 'json', params: param})
      .pipe(map((res: object) => new arrayType(res)));
  }

}

export function csHttpFactory(handler: HttpHandler): CsHttpClient {
  return new CsHttpClient(handler);
}

export const CsHttpProvider = {
  provide: CsHttpClient,
  useFactory: csHttpFactory,
  deps: [HttpHandler]
};
