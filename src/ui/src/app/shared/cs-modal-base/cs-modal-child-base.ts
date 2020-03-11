import { AfterViewInit, Directive, HostBinding, Input, OnDestroy, Optional, Output, ViewChild, ViewContainerRef } from '@angular/core';
import { CsComponentBase } from '../cs-components-library/cs-component-base';
import { Observable, Subject } from 'rxjs';
import { MessageService } from "../../shared.service/message.service";

@Directive({selector: 'div.modal-body, .modal-title'})
export class CsModalChildBaseSelector {
  @HostBinding('tabindex') tabIndex = '-1';
  constructor(public view: ViewContainerRef) {

  }
}

export class CsModalChildBase extends CsComponentBase implements OnDestroy{
  _modalOpened: boolean = false;
  @Output() closeNotification: Subject<any>;
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

export class CsModalChildMessage extends CsModalChildBase implements OnDestroy, AfterViewInit {
  constructor(protected messageService: MessageService) {
    super();
  }

  ngAfterViewInit(): void {
    this.messageService.registerModalDialogHandle(this.alertView)
  }

  ngOnDestroy() {
    this.messageService.unregisterModalDialogHandle();
  }
}
