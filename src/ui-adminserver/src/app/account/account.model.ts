import { RequestBase, ResponseBase, HttpBind } from '../shared/shared.type';

export class User implements RequestBase {
  id = 0;
  username = '';
  password = '';
  port = 22;

  PostBody(): object {
    return {
      id: this.id,
      username: this.username,
      password: this.password,
      port: this.port,
    };
  }
}

export class UserVerify implements RequestBase {
  id = 0;
  username: string;
  password: string;
  passwordConfirm: string;

  hasEmpty(): boolean {
    if (this.username && this.password && this.passwordConfirm) {
      return true;
    }
    return false;
  }

  verifyPassword(): boolean {
    return this.password === this.passwordConfirm;
  }

  toUser(): User {
    const user = new User();
    user.id = this.id;
    user.username = this.username;
    user.password = this.password;
    return user;
  }

  PostBody(): object {
    return {
      id: this.id,
      username: this.username,
      password: this.password
    };
  }
}

export class DBInfo implements RequestBase {
  maxConnection = 1000;
  password: string;
  passwordConfirm: string;

  PostBody(): object {
    return {
      db_max_connections: this.maxConnection,
      db_password: this.password
    };
  }

  verify(): boolean {
    return this.password === this.passwordConfirm;
  }

}

export class MyToken extends ResponseBase {
  @HttpBind('token') token = '';
}
