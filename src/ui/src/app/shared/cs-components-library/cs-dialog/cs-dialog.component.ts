import { Component, OnDestroy } from '@angular/core';
import { Message, RETURN_STATUS } from '../../shared.types';
import { Observable, Subject } from 'rxjs';

@Component({
  templateUrl: "./cs-dialog.component.html"
})
export class CsDialogComponent implements OnDestroy {
  opened: boolean;
  curMessage: Message;
  private returnSubject:Subject<Message>;
  constructor() {
    this.curMessage = new Message();
    this.returnSubject = new Subject<Message>();
  }

  ngOnDestroy(): void {

  }

  public openDialog(message:Message):Observable<Message>{
    this.curMessage = message;
    this.opened = true;
    return this.returnSubject.asObservable();
  }

  confirm(): void {
    this.opened = false;
    this.curMessage.returnStatus = RETURN_STATUS.rsConfirm;
    this.returnSubject.next(this.curMessage);
  }

  cancel(): void {
    this.opened = false;
    this.curMessage.returnStatus = RETURN_STATUS.rsCancel;
    this.returnSubject.next(this.curMessage);
  }
}
