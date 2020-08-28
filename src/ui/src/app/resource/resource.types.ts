import { HttpBase, HttpBind } from '../shared/ui-model/model-types';

export class ConfigMapDetailMetadata extends HttpBase {
  @HttpBind('namespace') namespace = '';
  @HttpBind('name') name = '';
  @HttpBind('creation_time') creationTime = '';
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

