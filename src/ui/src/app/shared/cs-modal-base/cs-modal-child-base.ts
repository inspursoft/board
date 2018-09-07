import { Directive, OnDestroy, ViewChild, ViewContainerRef } from '@angular/core';
import { CsComponentBase } from '../cs-components-library/cs-component-base';
import { Observable, Subject } from 'rxjs';

@Directive({selector: 'div.modal-body'})
export class CsModalChildBaseSelector {
  constructor(public view: ViewContainerRef) {

  }
}

export class CsModalChildBase extends CsComponentBase implements OnDestroy{
  closeNotification: Subject<any>;
  _modalOpened: boolean = false;
  @ViewChild(CsModalChildBaseSelector) alertViewSelector;

  constructor() {
    super();
    this.closeNotification = new Subject<any>();
  }

  get alertView(): ViewContainerRef {
    return this.alertViewSelector.view;
  }

  ngOnDestroy() {
    this.closeNotification.next();
  }

  set modalOpened(value: boolean) {
    this._modalOpened = value;
    if (!value) {
      this.closeNotification.next()
    }
  }

  get modalOpened(): boolean {
    return this._modalOpened;
  }

  openModal(): Observable<any> {
    this.modalOpened = true;
    return this.closeNotification.asObservable();
  }
}