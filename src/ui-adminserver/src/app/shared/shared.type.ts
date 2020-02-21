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

  constructor(public res: object) {
    this.init();
  }
}

export abstract class RequestBase {
  abstract PostBody(): object;
}

export function HttpBind(name: string) {
  return (target: ResponseBase, propertyName: string) => {
    Reflect.defineMetadata(propertyName, name, target);
  };
}

export abstract class ResponseArrayBase<T extends ResponseBase> {
  protected data: Array<T>;

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

  get originData(): Array<T> {
    return this.data;
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
