import 'reflect-metadata';
import { isArray } from 'util';

export abstract class ResponseBase {
  protected init() {
    const metadataKeys: Array<string> = Reflect.getMetadataKeys(this);
    metadataKeys.forEach((metadataKey: string) => {
      const propertyName = Reflect.getMetadata(metadataKey, this);
      if (Reflect.has(this.res, propertyName)) {
        const value = Reflect.get(this.res, propertyName);
        Reflect.set(this, metadataKey, value);
      }
    });
  }

  constructor(public res?: object) {
    if (res) {
      this.init();
    }
  }
}

export class Pagination extends ResponseBase {
  @HttpBind('page_index') pageIndex: number;
  @HttpBind('page_size') pageSize: number;
  @HttpBind('total_count') totalCount: number;
  @HttpBind('page_count') pageCount: number;
}

export abstract class RequestBase {
  abstract PostBody(): object;
}

export function HttpBind(name: string) {
  return (target: ResponseBase, propertyName: string) => {
    Reflect.defineMetadata(propertyName, name, target);
  };
}

export abstract class ResponsePaginationBase<T extends ResponseBase> {
  list: Array<T>;
  pagination: Pagination;

  abstract ListKeyName(): string;

  abstract CreateOneItem(res: object): T;

  public constructor(public res: object) {
    this.list = Array<T>();
    this.pagination = new Pagination(this.getObject('pagination'));
    const resList = this.getObject(this.ListKeyName());
    if (isArray(resList)) {
      (resList as Array<object>).forEach(item => this.list.push(this.CreateOneItem(item)));
    }
  }

  getObject(key: string): object {
    if (Reflect.has(this.res, key)) {
      return Reflect.get(this.res, key);
    } else {
      return {};
    }
  }
}

export abstract class ResponseArrayBase<T extends ResponseBase> {
  data: Array<T>;

  abstract CreateOneItem(res: object): T;

  constructor(protected res: object) {
    this.data = Array<T>();
    if (isArray(this.res)) {
      (this.res as Array<object>).forEach(item => this.data.push(this.CreateOneItem(item)));
    }
  }

  get length(): number {
    return this.data.length;
  }
}
