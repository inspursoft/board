import 'reflect-metadata';
import { isArray } from 'util';

export function HttpBind(name: string) {
  return (target: ResponseBase, propertyName: string) => {
    Reflect.defineMetadata(propertyName, name, target);
  };
}

export abstract class ResponseBase {
  protected init() {
    const metadataKeys: Array<string> = Reflect.getMetadataKeys(this);
    metadataKeys.forEach((metadataKey: string) => {
      const propertyName = Reflect.getMetadata(metadataKey, this);
      if (this.res && Reflect.has(this.res, propertyName)) {
        const value = Reflect.get(this.res, propertyName);
        Reflect.set(this, metadataKey, value);
      }
    });
  }

  constructor(public res: object) {
    this.init();
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
    if (isArray(this.res)) {
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
    if (isArray(resList)) {
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

export interface IPagination {
  page_index: number;
  page_size: number;
  total_count: number;
  page_count: number;
}
