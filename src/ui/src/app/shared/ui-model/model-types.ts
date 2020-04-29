import 'reflect-metadata';
import { Type } from '@angular/core';

export interface IBindType {
  resPropertyTypeName: string;
  resPropertyName: string;
  arrayItemType?: Type<ResponseBase> | null;
}

export function HttpBind(resPropertyName: string) {
  return (target: ResponseBase, propertyName: string) => {
    const resValue: IBindType = {resPropertyTypeName: 'normal', resPropertyName};
    Reflect.defineMetadata(propertyName, resValue, target);
  };
}

export function HttpBindObject(resPropertyName: string) {
  return (target: ResponseBase, propertyName: string) => {
    const resValue: IBindType = {resPropertyTypeName: 'object', resPropertyName};
    Reflect.defineMetadata(propertyName, resValue, target);
  };
}

export function HttpBindArray(resPropertyName: string, arrayItemType: Type<ResponseBase>) {
  return (target: ResponseBase, propertyName: string) => {
    const resValue: IBindType = {resPropertyTypeName: 'array', resPropertyName, arrayItemType};
    Reflect.defineMetadata(propertyName, resValue, target);
  };
}

export abstract class ResponseBase {
  protected init() {
    const metadataKeys: Array<string> = Reflect.getMetadataKeys(this);
    metadataKeys.forEach((propertyKey: string) => {
      const property = Reflect.get(this, propertyKey);
      const metadataValue: IBindType = Reflect.getMetadata(propertyKey, this);
      if (metadataValue.resPropertyTypeName === 'array') {
        if (Reflect.has(this.res, metadataValue.resPropertyName)) {
          const resArray = Reflect.get(this.res, metadataValue.resPropertyName) as Array<object>;
          resArray.forEach(resItem => {
            const item = new metadataValue.arrayItemType(resItem);
            const propertyArray = property as Array<ResponseBase>;
            propertyArray.push(item);
          });
        }
      } else if (metadataValue.resPropertyTypeName === 'object') {
        if (Reflect.has(this.res, metadataValue.resPropertyName)) {
          const resValue = Reflect.get(this.res, metadataValue.resPropertyName);
          const v = {};
          Object.assign(v, resValue);
          Reflect.set(this, propertyKey, v);
        }
      } else {
        if (Reflect.has(this.res, metadataValue.resPropertyName)) {
          const resValue = Reflect.get(this.res, metadataValue.resPropertyName);
          Reflect.set(this, propertyKey, resValue);
        }
      }
    });
  }

  constructor(public res: object) {
    this.prepareInit();
    this.init();
  }

  protected prepareInit() {

  }
}

export class Pagination extends ResponseBase {
  @HttpBind('page_index') PageIndex: number;
  @HttpBind('page_size') PageSize: number;
  @HttpBind('total_count') TotalCount: number;
  @HttpBind('page_count') PageCount: number;
}

export abstract class ResponseArrayBase<T extends ResponseBase> {
  protected data: Array<T>;

  abstract CreateOneItem(res: object): T;

  protected constructor(protected res: object) {
    this.data = Array<T>();
    if (Array.isArray(this.res)) {
      (this.res as Array<object>).forEach(item => this.data.push(this.CreateOneItem(item)));
    }
  }

  get length() {
    return this.data.length;
  }

  [Symbol.iterator]() {
    let index = 0;
    const self = this;
    return {
      next() {
        if (index < self.data.length) {
          return {value: self.data[index++], done: false};
        } else {
          return {value: undefined, done: true};
        }
      }
    };
  }
}

export abstract class ResponsePaginationBase<T extends ResponseBase> {
  list: Array<T>;
  pagination: Pagination;

  abstract ListKeyName(): string;

  abstract CreateOneItem(res: object): T;

  protected constructor(public res: object) {
    this.list = Array<T>();
    this.pagination = new Pagination(this.getObject('pagination'));
    const resList = this.getObject(this.ListKeyName());
    if (Array.isArray(resList)) {
      (resList as Array<object>).forEach(item => this.list.push(this.CreateOneItem(item)));
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
