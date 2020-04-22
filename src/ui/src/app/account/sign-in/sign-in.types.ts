export class ReqSignIn {
  username = '';
  password = '';
  captchaId = '';
  captcha = '';
}

export enum ResSignInType {
  normal = 1, overThreeTimes, temporarilyBlocked
}

export class ResSignIn {
  retries = 0;
  type = ResSignInType.normal;
  description = '';
  value = 0;

  assignFromRes(res: IResSignIn): ResSignIn {
    this.retries = Reflect.get(res, 'resolve_sign_in_retries');
    this.type = Reflect.get(res, 'resolve_sign_in_type');
    this.description = Reflect.get(res, 'resolve_sign_in_description');
    this.value = Reflect.get(res, 'resolve_sign_in_value');
    return this;
  }
}

export interface IResSignIn {
  resolve_sign_in_retries: number;
  resolve_sign_in_type: number;
  resolve_sign_in_description: string;
  resolve_sign_in_value: number;
}
