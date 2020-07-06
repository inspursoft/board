export class Project {
  project_id: number;
  project_name: string;
  publicity: boolean;
  project_public: number;
  project_creation_time: Date;
  project_comment: string;
  project_owner_id: number;
  project_owner_name: string;
}

export class Member {
  project_member_id?: number;
  project_member_user_id: number;
  project_member_username?: string;
  project_member_role_id: number;
  isMember?: boolean;
}

export class Role {
  role_id: number;
  role_name: string;
}

export class CreateProject {
  projectName = '';
  publicity = false;
  comment = '';
}
