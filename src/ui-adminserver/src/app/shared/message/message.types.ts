import {HttpErrorResponse} from '@angular/common/http';
import {Type} from '@angular/core';
import {TimeoutError} from 'rxjs';

export const DISMISS_ALERT_INTERVAL = 4;

export enum ReturnStatus {
  rsNone, rsConfirm, rsCancel
}

export enum ExecuteStatus {
  esNotExe = 'NotExe',
  esExecuting = 'Executing',
  esSuccess = 'Success',
  esFailed = 'Failed'
}

export enum ButtonStyle {
  Confirmation = 1, Deletion, YesNo, OnlyConfirm
}

export class Message {
  title = '';
  message = '';
  data: any;
  buttonStyle: ButtonStyle = ButtonStyle.Confirmation;
  returnStatus: ReturnStatus = ReturnStatus.rsNone;
}

export type AlertType = 'success' | 'danger' | 'info' | 'warning';

export class AlertMessage {
  message = '';
  alertType: AlertType = 'success';
}

export enum GlobalAlertType {
  gatNormal, gatShowDetail, gatLogin
}

export class GlobalAlertMessage {
  type: GlobalAlertType = GlobalAlertType.gatNormal;
  message = '';
  alertType: AlertType = 'danger';
  errorObject: HttpErrorResponse | Type<Error> | TimeoutError;
  endMessage = '';
}
