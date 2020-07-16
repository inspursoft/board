import 'reflect-metadata';
import { Type } from '@angular/core';

export interface IBindType {
  serverPropertyTypeName: string;
  serverPropertyName: string;
  itemType?: Type<HttpBase> | null;
  booleanTrueValue?: any;
  booleanFalseValue?: any;
}

export function HttpBind(serverPropertyName: string) {
  return (target: HttpBase, propertyName: string) => {
    const resValue: IBindType = {serverPropertyTypeName: 'normal', serverPropertyName};
    Reflect.defineMetadata(propertyName, resValue, target);
  };
}

export function HttpBindObject(serverPropertyName: string, itemType: Type<HttpBase>) {
  return (target: HttpBase, propertyName: string) => {
    const resValue: IBindType = {serverPropertyTypeName: 'object', serverPropertyName, itemType};
    Reflect.defineMetadata(propertyName, resValue, target);
  };
}

export function HttpBindArray(serverPropertyName: string, itemType: Type<HttpBase>) {
  return (target: HttpBase, propertyName: string) => {
    const resValue: IBindType = {serverPropertyTypeName: 'array', serverPropertyName, itemType};
    Reflect.defineMetadata(propertyName, resValue, target);
  };
}

export function HttpBindBoolean(serverPropertyName: string, trueValue: any, falseValue: any) {
  return (target: HttpBase, propertyName: string) => {
    const resValue: IBindType = {
      serverPropertyTypeName: 'boolean',
      serverPropertyName,
      booleanTrueValue: trueValue,
      booleanFalseValue: falseValue
    };
    Reflect.defineMetadata(propertyName, resValue, target);
  };
}

export abstract class HttpBase {

  constructor(public res?: object) {
    this.prepareInit();
  }

  protected prepareInit() {

  }

  protected afterInit() {

  }

  protected preparePost() {

  }

  getPostBody(): { [key: string]: any } {
    this.preparePost();
    const postBody = {};
    const metadataKeys: Array<string> = Reflect.getMetadataKeys(this);
    metadataKeys.forEach((propertyKey: string) => {
      const property = Reflect.get(this, propertyKey);
      const metadataValue: IBindType = Reflect.getMetadata(propertyKey, this);
      if (metadataValue.serverPropertyTypeName === 'array') {
        const reqArray = property as Array<HttpBase>;
        const reqPropertyValue = new Array<{ [key: string]: string }>();
        reqArray.forEach(reqItem => reqPropertyValue.push(reqItem.getPostBody()));
        Reflect.set(postBody, metadataValue.serverPropertyName, reqPropertyValue);
      } else if (metadataValue.serverPropertyTypeName === 'object') {
        const reqObject = property as HttpBase;
        Reflect.set(postBody, metadataValue.serverPropertyName, reqObject.getPostBody());
      } else if (metadataValue.serverPropertyTypeName === 'boolean') {
        const booleanValue = property ? metadataValue.booleanTrueValue : metadataValue.booleanFalseValue;
        Reflect.set(postBody, metadataValue.serverPropertyName, booleanValue);
      } else {
        Reflect.set(postBody, metadataValue.serverPropertyName, property);
      }
    });
    return postBody;
  }

  initFromRes() {
    if (!this.res || typeof this.res !== 'object') {
      return;
    }
    const metadataKeys: Array<string> = Reflect.getMetadataKeys(this);
    metadataKeys.forEach((propertyKey: string) => {
      const property = Reflect.get(this, propertyKey);
      const metadataValue: IBindType = Reflect.getMetadata(propertyKey, this);
      if (metadataValue.serverPropertyTypeName === 'array') {
        if (Reflect.has(this.res, metadataValue.serverPropertyName)) {
          const resArray = Reflect.get(this.res, metadataValue.serverPropertyName) as Array<object>;
          const propertyArray = property as Array<HttpBase>;
          if (resArray && resArray.length > 0) {
            resArray.forEach(resItem => {
              const item = new (metadataValue.itemType as Type<HttpBase>)(resItem);
              item.initFromRes();
              propertyArray.push(item);
            });
          }
        }
      } else if (metadataValue.serverPropertyTypeName === 'object') {
        if (Reflect.has(this.res, metadataValue.serverPropertyName)) {
          const resValue = Reflect.get(this.res, metadataValue.serverPropertyName);
          const item = new (metadataValue.itemType as Type<HttpBase>)(resValue);
          item.initFromRes();
          Reflect.set(this, propertyKey, item);
        }
      } else if (metadataValue.serverPropertyTypeName === 'boolean') {
        if (Reflect.has(this.res, metadataValue.serverPropertyName)) {
          const resValue = Reflect.get(this.res, metadataValue.serverPropertyName);
          Reflect.set(this, propertyKey, resValue === metadataValue.booleanTrueValue);
        }
      } else {
        if (Reflect.has(this.res, metadataValue.serverPropertyName) &&
          Reflect.get(this.res, metadataValue.serverPropertyName)) {
          const resValue = Reflect.get(this.res, metadataValue.serverPropertyName);
          Reflect.set(this, propertyKey, resValue);
        }
      }
    });
    this.afterInit();
  }
}

export class Pagination extends HttpBase {
  @HttpBind('page_index') PageIndex = 1;
  @HttpBind('page_size') PageSize = 1;
  @HttpBind('total_count') TotalCount = 0;
  @HttpBind('page_count') PageCount = 0;
}

export abstract class ResponsePaginationBase<T extends HttpBase> {
  list: Array<T>;
  pagination: Pagination;

  abstract ListKeyName(): string;

  abstract CreateOneItem(res: object): T;

  constructor(public res: object) {
    this.list = Array<T>();
    this.pagination = new Pagination(this.getObject('pagination'));
    this.pagination.initFromRes();
    const resList = this.getObject(this.ListKeyName());
    if (Array.isArray(resList)) {
      (resList as Array<object>).forEach(response => {
          const item = this.CreateOneItem(response);
          item.initFromRes();
          this.list.push(item);
        }
      );
    }
  }

  [Symbol.iterator]() {
    let index = 0;
    const self = this;
    return {
      next() {
        if (index < self.list.length) {
          return {value: self.list[index++], done: false};
        } else {
          return {value: undefined, done: true};
        }
      }
    };
  }

  getObject(key: string): object {
    if (Reflect.has(this.res, key)) {
      return Reflect.get(this.res, key);
    } else {
      return {};
    }
  }
}

