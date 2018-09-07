import { HttpErrorResponse } from "@angular/common/http";
import { Type } from "@angular/core";

export interface ICsMenuItemData {
  caption: string,
  icon: string,
  url: string,
  visible: boolean
}

export enum RETURN_STATUS {
  rsNone, rsConfirm, rsCancel
}

export enum BUTTON_STYLE {
  CONFIRMATION, DELETION, YES_NO, ONLY_CONFIRM
}

export class Message {
  title: string = '';
  message: string = '';
  data: any;
  buttonStyle: BUTTON_STYLE = BUTTON_STYLE.CONFIRMATION;
  returnStatus: RETURN_STATUS = RETURN_STATUS.rsNone;
}

export type AlertType = 'alert-success' | 'alert-danger' | 'alert-info' | 'alert-warning';

export class AlertMessage {
  message: string = '';
  alertType: AlertType = 'alert-success';
}

export enum GlobalAlertType {
  gatNormal, gatShowDetail, gatLogin
}

export class GlobalAlertMessage {
  type: GlobalAlertType = GlobalAlertType.gatNormal;
  message: string = '';
  alertType: AlertType = 'alert-danger';
  errorObject: HttpErrorResponse | Type<Error>;
}

export class SignUp {
  username: string;
  email: string;
  password: string;
  confirmPassword: string;
  realname: string;
  comment: string;
}
