export class ConfigMap {
  namespace = '';
  name = '';
  dataList: Array<{key: string, value: string}>;

  constructor() {
    this.dataList = Array<{key: string, value: string}>();
  }

  static createFromRes(res: Object): ConfigMap {
    let result = new ConfigMap();
    result.namespace = res['namespace'];
    result.name = res['name'];
    if (Reflect.has(res, 'datalist')) {
      Reflect.ownKeys(res['datalist']).forEach((key: string) =>
        result.dataList.push({key: key, value: res['datalist'][key]})
      );
    }
    return result;
  }

  postBody(): Object {
    let obj = Object.create({});
    this.dataList.forEach(value =>
      Object.defineProperties(obj, {[value.key]: {enumerable: true, value: value.value}})
    );
    return {
      namespace: this.namespace,
      name: this.name,
      datalist: obj
    }
  }
}

export class ConfigMapDetail {
  namespace = '';
  name = '';
  creationTime = '';
  deletionTime = '';
  labels = '';
  dataList: Array<{key: string, value: string}>;

  constructor() {
    this.dataList = Array<{key: string, value: string}>();
  }

  static createFromRes(res: Object): ConfigMapDetail {
    let result = new ConfigMapDetail();
    let metadata = res['metadata'];
    result.namespace = metadata['Namespace'];
    result.name = metadata['Name'];
    result.creationTime = metadata['CreationTimestamp'];
    result.deletionTime = metadata['DeletionTimestamp'];
    result.labels = metadata['Labels'];
    if (Reflect.has(res, 'data')) {
      Reflect.ownKeys(res['data']).forEach((key: string) =>
        result.dataList.push({key: key, value: res['data'][key]})
      );
    }
    return result;
  }
}
