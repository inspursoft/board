import { Component, OnDestroy, OnInit } from '@angular/core';
import { BUTTON_STYLE, Message, RETURN_STATUS } from '../../shared.types';
import { Observable, Subject } from 'rxjs';
import "rxjs/add/observable/fromEvent"

@Component({
  templateUrl: "./cs-dialog.component.html"
})
export class CsDialogComponent implements OnDestroy, OnInit {
  opened: boolean;
  curMessage: Message;
  protected returnSubject: Subject<Message>;
  constructor() {
    this.curMessage = new Message();
    this.returnSubject = new Subject<Message>();
  }

  ngOnDestroy(): void {
    delete this.curMessage;
    delete this.returnSubject;
  }

  ngOnInit() {
    const obsKeyPress = Observable.fromEvent(document, 'keypress');
    const unSubscribe= obsKeyPress.subscribe((event: KeyboardEvent) => {
      if (event.charCode === 13 && this.curMessage.buttonStyle == BUTTON_STYLE.ONLY_CONFIRM) {
        this.curMessage.returnStatus = RETURN_STATUS.rsCancel;
        this.returnSubject.next(this.curMessage);
        unSubscribe.unsubscribe();
      }
    });
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
