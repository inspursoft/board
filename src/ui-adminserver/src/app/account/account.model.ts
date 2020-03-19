import { RequestBase } from '../shared/shared.type';

export class User implements RequestBase {
  username: string;
  password: string;

  PostBody(): object {
    return {
      username: this.username,
      password: this.password
    };
  }
}

export class UserVerify {
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
    user.username = this.username;
    user.password = this.password;
    return user;
  }
}
