export interface user {
  user_id: number;
  user_name: string;
  user_email: string;
  user_password: string;
  user_confirm_password: string;
  user_realname: string;
  user_comment: string;
  user_deleted: number;
  user_system_admin: number;
  user_project_admin: number;
  user_reset_uuid: string;
  user_salt: string;
  user_creation_time: Date;
  user_update_time: Date;
}

export class User implements user {
  public user_id: number = 0;
  public user_name: string = "";
  public user_email: string = "";
  public user_password: string = "";
  public user_confirm_password: string = "";
  public user_realname: string = "";
  public user_comment: string = "";
  public user_deleted: number = 0;
  public user_system_admin: number = 0;
  public user_project_admin: number = 0;
  public user_reset_uuid: string = "";
  public user_salt: string = "";
  public user_creation_time: Date = new Date();
  public user_update_time: Date = new Date();
  constructor() {}
}
