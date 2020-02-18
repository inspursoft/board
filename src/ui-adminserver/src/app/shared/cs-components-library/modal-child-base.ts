import { AfterViewInit, OnDestroy, Output, ViewChild, ViewContainerRef } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { MessageService } from '../message/message.service';
import { CsComponentBase } from './cs-component-base';

export abstract class ModalChildBase extends CsComponentBase implements OnDestroy {
  abstract get alertView(): ViewContainerRef;

  modalOpenedValue = false;
  @Output() closeNotification: Subject<any>;

  protected constructor() {
    super();
    this.closeNotification = new Subject<any>();
  }

  ngOnDestroy() {
    this.closeNotification.next();
    delete this.closeNotification;
  }

  set modalOpened(value: boolean) {
    this.modalOpenedValue = value;
    if (!value) {
      this.closeNotification.next();
    }
  }

  get modalOpened(): boolean {
    return this.modalOpenedValue;
  }

  openModal(): Observable<any> {
    this.modalOpened = true;
    return this.closeNotification.asObservable();
  }
}

export abstract class ModalChildMessage extends ModalChildBase implements OnDestroy, AfterViewInit {
  protected constructor(protected messageService: MessageService) {
    super();
  }

  ngAfterViewInit(): void {
    this.messageService.registerModalDialogHandle(this.alertView);
  }

  ngOnDestroy() {
    this.messageService.unregisterModalDialogHandle();
  }
}
