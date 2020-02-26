import { HttpClient, HttpHandler } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { ResponseArrayBase, ResponseBase, ResponsePaginationBase } from '../../shared/shared.type';

export class CustomHttpClient extends HttpClient {
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

  postJson<T extends ResponseBase>(url: string, body: object, returnType: new(res: object) => T): Observable<T> {
    return super.post(url, body, {observe: 'body', responseType: 'json'})
      .pipe(map((res: object) => new returnType(res)));
  }

  getArrayJson<T extends ResponseArrayBase<ResponseBase>>(url: string,
                                                          arrayType: new(res: object) => T,
                                                          param?: { [param: string]: string }): Observable<T> {
    return super.get(url, {observe: 'body', responseType: 'json', params: param})
      .pipe(map((res: object) => new arrayType(res)));
  }

  getPagination<T extends ResponsePaginationBase<ResponseBase>>(url: string,
                                                                paginationType: new(res: object) => T,
                                                                param?: {[param: string]: string}): Observable<T> {
    return super.get(url, {observe: 'body', responseType: 'json', params: param})
      .pipe(map((res: object) => new paginationType(res)));
  }

}

export function CustomHttpFactory(handler: HttpHandler): CustomHttpClient {
  return new CustomHttpClient(handler);
}

export const CustomHttpProvider = {
  provide: CustomHttpClient,
  useFactory: CustomHttpFactory,
  deps: [HttpHandler]
};
