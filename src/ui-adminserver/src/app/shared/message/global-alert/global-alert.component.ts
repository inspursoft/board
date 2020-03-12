import {Component} from '@angular/core';
import {Observable, Subject} from 'rxjs';
import {HttpErrorResponse} from '@angular/common/http';
import {GlobalAlertMessage} from '../message.types';

@Component({
  templateUrl: './global-alert.component.html',
  styleUrls: ['./global-alert.component.css']
})
export class GlobalAlertComponent {
  isOpenValue = false;
  curMessage: GlobalAlertMessage;
  onCloseEvent: Subject<any>;
  detailModalOpen = false;

  constructor() {
    this.onCloseEvent = new Subject<any>();
  }

  get isOpen(): boolean {
    return this.isOpenValue;
  }

  set isOpen(value: boolean) {
    this.isOpenValue = value;
    if (!value) {
      this.onCloseEvent.next();
    }
  }

  get errorDetailMsg(): string {
    let result = '';
    if (this.curMessage.errorObject && this.curMessage.errorObject instanceof HttpErrorResponse) {
      const err = (this.curMessage.errorObject as HttpErrorResponse).error;
      if (typeof err === 'object') {
        result = (this.curMessage.errorObject as HttpErrorResponse).error.message;
      } else {
        result = err;
      }
    } else if (this.curMessage.errorObject) {
      result = (this.curMessage.errorObject as Error).message;
    }
    return result;
  }

  public openAlert(message: GlobalAlertMessage): Observable<any> {
    this.curMessage = message;
    this.isOpen = true;
    return this.onCloseEvent.asObservable();
  }

  loginClick() {
    // TODO: add login function;2020.2.17
  }
}
