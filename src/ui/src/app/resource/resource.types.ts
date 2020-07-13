import { HttpBase, HttpBind, HttpBindObject } from '../shared/ui-model/model-types';

export class ConfigMap extends HttpBase {
  @HttpBind('namespace') namespace = '';
  @HttpBind('name') name = '';
  @HttpBind('datalist') data: object;
  dataList: Array<{ key: string, value: string }>;

  protected prepareInit() {
    this.dataList = Array<{ key: string, value: string }>();
  }

  protected afterInit() {
    if (this.data) {
      Reflect.ownKeys(this.data).forEach((key: string) =>
        this.dataList.push({key, value: Reflect.get(this.data, key)})
      );
    }
  }

  getPostBody(): { [p: string]: any } {
    const obj = Object.create({});
    this.dataList.forEach(value =>
      Object.defineProperties(obj, {[value.key]: {enumerable: true, value: value.value}})
    );
    return {
      namespace: this.namespace,
      name: this.name,
      datalist: obj
    };
  }
}

export class ConfigMapDetailMetadata extends HttpBase {
  @HttpBind('namespace') namespace = '';
  @HttpBind('name') name = '';
  @HttpBind('creation_time') creationTime = '';
}

export class ConfigMapDetail extends HttpBase {
  @HttpBindObject('metadata', ConfigMapDetailMetadata) metadata: ConfigMapDetailMetadata;
  @HttpBind('data') data: object;
  dataList: Array<{ key: string, value: string }>;

  protected prepareInit() {
    this.metadata = new ConfigMapDetailMetadata();
    this.dataList = Array<{ key: string, value: string }>();
  }

  protected afterInit() {
    if (this.data) {
      Reflect.ownKeys(this.data).forEach((key: string) =>
        this.dataList.push({key, value: Reflect.get(this.data, key)})
      );
    }
  }
}


export class ConfigMapProject extends HttpBase {
  @HttpBind('project_id') projectId = 0;
  @HttpBind('project_name') projectName = '';
  @HttpBind('publicity') publicity = false;
  @HttpBind('project_public') projectPublic = 0;
  @HttpBind('project_creation_time') projectCreationTime = 0;
  @HttpBind('project_comment') projectComment = '';
  @HttpBind('project_owner_id') projectOwnerId = 0;
  @HttpBind('project_owner_name') projectOwnerName = '';
}

