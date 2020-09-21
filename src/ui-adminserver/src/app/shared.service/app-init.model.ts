import { ResponseBase, HttpBind } from '../shared/shared.type';

export enum InitStatusCode {
  InitStatusFirst = 1, InitStatusSecond, InitStatusThird
}

export class InitStatus extends ResponseBase {
  @HttpBind('status') status: InitStatusCode;
  @HttpBind('log') log = '';
}
