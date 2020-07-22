import { HttpBase, HttpBind, HttpBindArray } from '../shared/ui-model/model-types';

export class ImageResult extends HttpBase {
  @HttpBind('image_name') imageName = '';
  @HttpBind('project_name') projectName = '';
}

export class NodeResult extends HttpBase {
  @HttpBind('node_ip') nodeIp = '';
  @HttpBind('node_name') nodeName = '';
}

export class UserResult extends HttpBase {
  @HttpBind('role_name') roleName = '';
  @HttpBind('user_email') userEmail = '';
  @HttpBind('user_name') userName = '';
}

export class ServiceResult extends HttpBase {
  @HttpBind('is_public') isPublic = false;
  @HttpBind('project_name') projectName = '';
  @HttpBind('service_name') serviceName = '';
}

export class ProjectResult extends HttpBase {
  @HttpBind('is_public') isPublic = false;
  @HttpBind('owner_name') ownerName = '';
  @HttpBind('project_name') projectName = '';
}

export class GlobalSearchResult extends HttpBase {
  @HttpBindArray('project_result', ProjectResult) projectResult: Array<ProjectResult>;
  @HttpBindArray('service_result', ServiceResult) serviceResult: Array<ServiceResult>;
  @HttpBindArray('user_result', UserResult) userResult: Array<UserResult>;
  @HttpBindArray('node_result', NodeResult) nodeResult: Array<NodeResult>;
  @HttpBindArray('images_name', ImageResult) imageResult: Array<ImageResult>;

  protected prepareInit() {
    this.projectResult = new Array<ProjectResult>();
    this.serviceResult = new Array<ServiceResult>();
    this.userResult = new Array<UserResult>();
    this.nodeResult = new Array<NodeResult>();
    this.imageResult = new Array<ImageResult>();
  }
}
