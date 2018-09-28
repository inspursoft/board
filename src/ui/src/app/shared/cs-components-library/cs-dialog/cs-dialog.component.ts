import { AfterViewInit, Component, ElementRef, HostBinding, HostListener, OnDestroy } from '@angular/core';
import { BUTTON_STYLE, Message, RETURN_STATUS } from '../../shared.types';
import { Observable, Subject } from 'rxjs';
import "rxjs/add/observable/fromEvent"

@Component({
  templateUrl: "./cs-dialog.component.html"
})
export class CsDialogComponent implements OnDestroy, AfterViewInit {
  opened: boolean;
  curMessage: Message;
  private returnSubject: Subject<Message>;
  @HostBinding('tabindex') tabIndex = '-1';

  @HostListener('keypress', ['$event']) onKeypress(event: KeyboardEvent) {
    if (event.charCode === 13 && this.curMessage.buttonStyle == BUTTON_STYLE.ONLY_CONFIRM) {
      this.curMessage.returnStatus = RETURN_STATUS.rsCancel;
      this.returnSubject.next(this.curMessage);
    }
  }

  constructor(private el: ElementRef) {
    this.curMessage = new Message();
    this.returnSubject = new Subject<Message>();
  }

  ngAfterViewInit() {
    (this.el.nativeElement as HTMLElement).focus();
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
