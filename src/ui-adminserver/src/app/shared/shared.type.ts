import 'reflect-metadata';

export abstract class ResponseBase {
  protected init() {
    const metadataKeys: Array<string> = Reflect.getMetadataKeys(this);
    metadataKeys.forEach((metadataKey: string) => {
      const propertyName = Reflect.getMetadata(metadataKey, this);
      if (Reflect.has(this.res, propertyName)) {
        const value = Reflect.get(this.res, propertyName);
        Reflect.set(this, metadataKey, value)
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
