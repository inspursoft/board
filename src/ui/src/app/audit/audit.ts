export class Audit {
  public operation_id: number = 0;
  public operation_creation_time: string = "";
  public operation_update_time: string = "";
  public operation_deleted: number = 0;
  public operation_user_id: number = 0;
  public operation_user_name: string = "";
  public operation_project_name: string = "";
  public operation_project_id: number = 0;
  public operation_tag: string = "";
  public operation_comment: string = "";
  public operation_object_type: string = "";
  public operation_object_name: string = "";
  public operation_action: string = "";
  public operation_status: string = "";
  public operation_path: string = "";

  constructor() {
  }
}

export class AuditQueryData {
  public pageIndex: number = 1;
  public pageSize: number = 15;
  public sortBy: string = "";
  public isReverse: boolean = false;
  public endDateTamp: number = 0;
  public beginDateTamp: number = 0;
  public status: string = "";
  public user_name: string = "";
  public action: string = "";
  public object_name: string = "";

  constructor() {
  }

}
