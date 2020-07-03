import { HttpClient, HttpHandler, HttpHeaders } from '@angular/common/http';
import { Observable } from 'rxjs';
import { map } from 'rxjs/operators';
import { HttpBase, ResponsePaginationBase } from './model-types';
import { Type } from '@angular/core';

export class ModelHttpClient extends HttpClient {
  defaultHeaders = new Headers({
    'Content-Type': 'application/json'
  });

  constructor(handler: HttpHandler) {
    super(handler);
  }

  getJson(url: string, returnType: Type<HttpBase>, options?: {
    param?: { [param: string]: string },
    header?: HttpHeaders
  }): Observable<any> {
    return super.get(url, {
      observe: 'body',
      responseType: 'json',
      params: options && options.param ? options.param : null,
      headers: options && options.header ? options.header : null
    }).pipe(map((res: object) => {
      const returnItem = new returnType(res);
      returnItem.initFromRes();
      return returnItem;
    }));
  }

  postJson(url: string, returnType: Type<HttpBase>, body: any | null, options?: {
    param?: { [param: string]: string },
    header?: HttpHeaders
  }): Observable<any> {
    return super.post(url, body, {
      observe: 'body',
      responseType: 'json',
      params: options && options.param ? options.param : null,
      headers: options && options.header ? options.header : null
    }).pipe(map((res: object) => {
      const returnItem = new returnType(res);
      returnItem.initFromRes();
      return returnItem;
    }));
  }

  getPagination(url: string, paginationType: Type<ResponsePaginationBase<HttpBase>>, options?: {
    param?: { [param: string]: string },
    header?: HttpHeaders
  }): Observable<any> {
    return super.get(url, {
        observe: 'body', responseType: 'json',
        params: options && options.param ? options.param : null,
        headers: options && options.header ? options.header : null
      }
    ).pipe(map((res: object) => new paginationType(res)));
  }

  getArray(url: string, itemType: Type<HttpBase>, options?: {
    param?: { [param: string]: string },
    header?: HttpHeaders
  }): Observable<any> {
    return super.get(url, {
      observe: 'body',
      responseType: 'json',
      params: options && options.param ? options.param : null,
      headers: options && options.header ? options.header : null
    }).pipe(map((res: Array<object>) => {
      const result = Array<HttpBase>();
      res.forEach(item => {
        const newItem = new itemType(item);
        newItem.initFromRes();
        result.push(newItem);
      });
      return result;
    }));
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
