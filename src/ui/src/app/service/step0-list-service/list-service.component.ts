import { Component, Input, OnInit, OnDestroy } from '@angular/core';

import { Subscription } from 'rxjs/Subscription';

import { AppInitService } from '../../app.init.service';
import { K8sService } from '../service.k8s';
import { Service } from '../service';
import { MessageService } from '../../shared/message-service/message.service';
import { MESSAGE_TARGET, BUTTON_STYLE, MESSAGE_TYPE } from '../../shared/shared.const';
import { Message } from '../../shared/message-service/message';

class ServiceData {
  id: number;
  name: string;
  status: boolean;

  constructor(id: number, name: string, status: boolean) {
    this.id = id;
    this.name = name;
    this.status = status;
  }
}

@Component({
  templateUrl: './list-service.component.html'
})
export class ListServiceComponent implements OnInit, OnDestroy {
  @Input() data: any;
  currentUser: {[key: string]: any};
  services: Service[];

  _subscription: Subscription;

  constructor(private appInitService: AppInitService,
              private k8sService: K8sService,
              private messageService: MessageService) {
    this._subscription = this.messageService.messageConfirmed$.subscribe(m => {
      let confirmationMessage = <Message>m;
      if (confirmationMessage) {
        let serviceData = <ServiceData>confirmationMessage.data;
        let m: Message = new Message();
        switch (confirmationMessage.target) {
          case MESSAGE_TARGET.DELETE_SERVICE:
            this.k8sService
              .deleteService(serviceData.id)
              .then(res => {
                m.message = 'SERVICE.SUCCESSFUL_DELETE';
                this.messageService.inlineAlertMessage(m);
                this.retrieve();
              })
              .catch(err => {
                m.message = 'SERVICE.FAILED_TO_DELETE';
                m.type = MESSAGE_TYPE.COMMON_ERROR;
                this.messageService.inlineAlertMessage(m);
              });
            break;
          case MESSAGE_TARGET.TOGGLE_SERVICE: {
            let service:ServiceData = confirmationMessage.params[0];
            this.k8sService
              .toggleService(service.id,!service.status)
              .then(res => {
                m.message = 'SERVICE.SUCCESSFUL_TOGGLE';
                this.messageService.inlineAlertMessage(m);
                this.retrieve();
              })
              .catch(err => {
                m.message = 'SERVICE.FAILED_TO_TOGGLE';
                m.type = MESSAGE_TYPE.COMMON_ERROR;
                this.messageService.inlineAlertMessage(m);
              });
              break;
          }
        }
      }
    });
  }

  ngOnInit(): void {
    this.currentUser = this.appInitService.currentUser;
    this.retrieve();
  }

  ngOnDestroy(): void {
    if (this._subscription) {
      this._subscription.unsubscribe();
    }
  }

  get createActionIsDisabled(): boolean {
    if (this.currentUser &&
      this.currentUser.hasOwnProperty("user_project_admin") &&
      this.currentUser.hasOwnProperty("user_system_admin")) {
      return this.currentUser["user_project_admin"] == 0 && this.currentUser["user_system_admin"] == 0;
    }
    return true;
  }

  createService(): void {
    this.k8sService.stepSource.next(1);
  }

  retrieve(): void {
    this.k8sService.getServices()
      .then(services => this.services = services)
      .catch(err => this.messageService.dispatchError(err));
  }

  getServiceStatus(status: number): string {
    //0: preparing 1: running 2: suspending
    switch (status) {
      case 0:
        return 'SERVICE.STATUS_PREPARING';
      case 1:
        return 'SERVICE.STATUS_RUNNING';
      case 2:
        return 'SERVICE.STATUS_STOPPED';
    }
  }

  getServicePublic(publicity: number): string {
    return publicity === 1 ? 'SERVICE.STATUS_PUBLIC' : 'SERVICE.STATUS_PRIVATE';
  }

  editService(s: Service) {

  }

  confirmToServiceAction(s: Service, action: string): void {
    let serviceData = new ServiceData(s.service_id, s.service_name, (s.service_status === 1));
    let title: string;
    let message: string;
    let target: MESSAGE_TARGET;
    let buttonStyle: BUTTON_STYLE;
    switch (action) {
      case 'DELETE':
        title = 'SERVICE.DELETE_SERVICE';
        message = 'SERVICE.CONFIRM_TO_DELETE_SERVICE';
        target = MESSAGE_TARGET.DELETE_SERVICE;
        buttonStyle = BUTTON_STYLE.DELETION;
        break;
      case 'TOGGLE':
        title = 'SERVICE.TOGGLE_SERVICE';
        message = 'SERVICE.CONFIRM_TO_TOGGLE_SERVICE';
        target = MESSAGE_TARGET.TOGGLE_SERVICE;
        buttonStyle = BUTTON_STYLE.CONFIRMATION;
        break;
    }
    let announceMessage = new Message();
    announceMessage.title = title;
    announceMessage.message = message;
    announceMessage.params = [ s.service_name ];
    announceMessage.target = target;
    announceMessage.buttons = buttonStyle;
    announceMessage.data = serviceData;
    this.messageService.announceMessage(announceMessage);
  }

  confirmToDeleteService(s: Service) {

  }
}
