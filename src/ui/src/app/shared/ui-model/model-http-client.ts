import { HttpClient, HttpHandler, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { ResponseArrayBase, ResponseBase, ResponsePaginationBase } from './model-types';

export class ModelHttpClient extends HttpClient {
  defaultHeaders = new Headers({
    'Content-Type': 'application/json'
  });

  constructor(handler: HttpHandler) {
    super(handler);
  }

  getJson<T extends ResponseBase>(url: string,
                                  returnType: new(res: object) => T,
                                  options?: {
                                    param?: { [param: string]: string },
                                    header?: HttpHeaders
                                  }): Observable<T> {
    return super.get(url, {
      observe: 'body',
      responseType: 'json',
      params: options && options.param ? options.param : null,
      headers: options && options.header ? options.header : null
    }).pipe(map((res: object) => new returnType(res)));
  }

  getPagination<T extends ResponsePaginationBase<ResponseBase>>(url: string,
                                                                paginationType: new(res: object) => T,
                                                                param?: { [param: string]: string }): Observable<T> {
    return super.get(url, {observe: 'body', responseType: 'json', params: param})
      .pipe(map((res: object) => new paginationType(res)));
  }

  getArray<T extends ResponseBase>(url: string,
                                   itemType: new(res: object) => T,
                                   options?: {
                                     param?: { [param: string]: string },
                                     header?: HttpHeaders
                                   }): Observable<Array<T>> {
    return super.get(url, {
      observe: 'body',
      responseType: 'json',
      params: options && options.param ? options.param : null,
      headers: options && options.header ? options.header : null
    }).pipe(map((res: Array<object>) => {
      const result = Array<T>();
      res.forEach(item => {
        const newItem = new itemType(item);
        result.push(newItem);
      });
      return result;
    }));
  }

  getArrayJson<T extends ResponseArrayBase<ResponseBase>>(url: string,
                                                          arrayType: new(res: object) => T,
                                                          param?: { [param: string]: string }): Observable<T> {
    return super.get(url, {observe: 'body', responseType: 'json', params: param})
      .pipe(map((res: object) => new arrayType(res)));
  }

}

export function CustomHttpFactory(handler: HttpHandler): ModelHttpClient {
  return new ModelHttpClient(handler);
}

export const CustomHttpProvider = {
  provide: ModelHttpClient,
  useFactory: CustomHttpFactory,
  deps: [HttpHandler]
};
