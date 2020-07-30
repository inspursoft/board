import { HttpBase, HttpBind, ResponsePaginationBase } from '../shared/ui-model/model-types';

export class AuditPagination extends ResponsePaginationBase<Audit> {
  CreateOneItem(res: object): Audit {
    return new Audit(res);
  }

  ListKeyName(): string {
    return 'operation_list';
  }
}

export class Audit extends HttpBase {
  @HttpBind('operation_id') id = 0;
  @HttpBind('operation_creation_time') creationTime = '';
  @HttpBind('operation_update_time') updateTime = '';
  @HttpBind('operation_deleted') deleted = 0;
  @HttpBind('operation_user_id') userId = 0;
  @HttpBind('operation_user_name') userName = '';
  @HttpBind('operation_project_name') projectName = '';
  @HttpBind('operation_project_id') projectId = 0;
  @HttpBind('operation_tag') tag = '';
  @HttpBind('operation_comment') comment = '';
  @HttpBind('operation_object_type') objectType = '';
  @HttpBind('operation_object_name') objectName = '';
  @HttpBind('operation_action') action = '';
  @HttpBind('operation_status') status = '';
  @HttpBind('operation_path') path = '';
}

export class AuditQueryData {
  pageIndex = 1;
  pageSize = 15;
  sortBy = '';
  isReverse = false;
  endTimestamp = 0;
  beginTimestamp = 0;
  status = '';
  userName = '';
  action = '';
  objectName = '';
}
