
export class User {
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

  constructor() {
  }
}
