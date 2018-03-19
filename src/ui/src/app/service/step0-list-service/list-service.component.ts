import { Component, Injector, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Service } from '../service';
import { BUTTON_STYLE, GUIDE_STEP, MESSAGE_TARGET, MESSAGE_TYPE, SERVICE_STATUS } from '../../shared/shared.const';
import { Message } from '../../shared/message-service/message';
import { ServiceDetailComponent } from './service-detail/service-detail.component';
import { ServiceStepBase } from "../service-step";
import { Observable } from "rxjs/Observable";
import "rxjs/add/observable/interval";

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
enum CreateServiceMethod{None, Wizards, YamlFile, DevOps}
@Component({
  templateUrl: './list-service.component.html',
  styleUrls: ["./list-service.component.css"]
})
export class ListServiceComponent extends ServiceStepBase implements OnInit, OnDestroy {
  currentUser: {[key: string]: any};
  services: Service[];
  isInLoading: boolean = false;
  isServiceControlOpen: boolean = false;
  serviceControlData: Service;
  checkboxRevertInfo: {isNeeded: boolean; value: boolean;};
  _subscription: Subscription;
  _subscriptionInterval: Subscription;

  totalRecordCount: number;
  pageIndex: number = 1;
  pageSize: number = 15;
  isBuildServiceWIP: boolean = false;
  isShowServiceCreateYaml: boolean = false;
  createServiceMethod: CreateServiceMethod = CreateServiceMethod.None;
  isActionWIP: Map<number, boolean>;

  @ViewChild(ServiceDetailComponent) serviceDetailComponent;

  constructor(protected injector: Injector) {
    super(injector);
    this._subscriptionInterval = Observable.interval(10000).subscribe(() => this.retrieve(true));
    this.isActionWIP = new Map<number, boolean>();
    this._subscription = this.messageService.messageConfirmed$.subscribe(m => {
      let confirmationMessage = <Message>m;
      if (confirmationMessage) {
        let serviceData = <ServiceData>confirmationMessage.data;
        let m: Message = new Message();
        this.isActionWIP.set(serviceData.id, true);
        switch (confirmationMessage.target) {
          case MESSAGE_TARGET.DELETE_SERVICE:
            this.k8sService
              .deleteService(serviceData.id)
              .then(() => {
                m.message = 'SERVICE.SUCCESSFUL_DELETE';
                this.messageService.inlineAlertMessage(m);
                this.retrieve();
                this.isActionWIP.set(serviceData.id, false);
              })
              .catch(err => {
                this.isActionWIP.set(serviceData.id, false);
                m.message = 'SERVICE.FAILED_TO_DELETE_SERVICE';
                m.type = MESSAGE_TYPE.COMMON_ERROR;
                this.messageService.inlineAlertMessage(m);
              });
            break;
          case MESSAGE_TARGET.TOGGLE_SERVICE: {
            let service: ServiceData = confirmationMessage.data;
            this.k8sService
              .toggleServiceStatus(service.id, service.status ? 0 : 1)
              .then(() => {
                m.message = 'SERVICE.SUCCESSFUL_TOGGLE';
                this.messageService.inlineAlertMessage(m);
                this.retrieve();
                this.isActionWIP.set(service.id, false);
              })
              .catch(err => {
                this.isActionWIP.set(service.id, false);
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
  }

  ngOnDestroy(): void {
    this._subscriptionInterval.unsubscribe();
    if (this._subscription) {
      this._subscription.unsubscribe();
    }
  }

  serviceInStoppedStatus(s: Service): boolean {
    return s.service_status == SERVICE_STATUS.STOPPED && !this.isActionWIP.get(s.service_id);
  }

  serviceCanChangePauseStatus(s: Service): boolean {
    return s.service_status in [SERVICE_STATUS.RUNNING, SERVICE_STATUS.WARNING] && !this.isActionWIP.get(s.service_id);
  }

  serviceDeleteStatusDisabled(s: Service): boolean {
    return s.service_status in [SERVICE_STATUS.PREPARING, SERVICE_STATUS.RUNNING]
      || this.isActionWIP.get(s.service_id);
  }

  createService(): void {
    if (this.createServiceMethod == CreateServiceMethod.Wizards) {
      this.k8sService.stepSource.next({index: 1, isBack: false});
    } else if (this.createServiceMethod == CreateServiceMethod.YamlFile) {
      this.isShowServiceCreateYaml = true;
    }
  }

  get isNormalStatus(): boolean {
    return !this.isBuildServiceWIP && !this.isShowServiceCreateYaml;
  }

  retrieve(isAuto: boolean = false): void {
    setTimeout(() => {
      if (this.isNormalStatus) {
        this.isInLoading = !isAuto;
        this.k8sService.getServices(this.pageIndex, this.pageSize)
          .then(paginatedServices => {
            this.totalRecordCount = paginatedServices["pagination"]["total_count"];
            this.services = paginatedServices["service_status_list"];
            this.isInLoading = false;
          })
          .catch(err => {
            this.messageService.dispatchError(err);
            this.isInLoading = false;
          });
      }
    });
  }

  getServiceStatus(status: SERVICE_STATUS): string {
    switch (status) {
      case SERVICE_STATUS.PREPARING:
        return 'SERVICE.STATUS_PREPARING';
      case SERVICE_STATUS.RUNNING:
        return 'SERVICE.STATUS_RUNNING';
      case SERVICE_STATUS.STOPPED:
        return 'SERVICE.STATUS_STOPPED';
      case SERVICE_STATUS.WARNING:
        return 'SERVICE.STATUS_UNCOMPLETED';
    }
  }

  getStatusClass(status: SERVICE_STATUS) {
    return {
      'running': status == SERVICE_STATUS.RUNNING,
      'stopped': status == SERVICE_STATUS.STOPPED,
      'warning': status == SERVICE_STATUS.WARNING
    }
  }

  toggleServicePublic(s: Service): void {
    let toggleMessage = new Message();
    let oldServicePublic = s.service_public;
    this.k8sService
      .toggleServicePublicity(s.service_id, s.service_public ? 0 : 1)
      .then(() => {
        s.service_public = !oldServicePublic;
        toggleMessage.message = 'SERVICE.SUCCESSFUL_TOGGLE';
        this.messageService.inlineAlertMessage(toggleMessage);
      })
      .catch(err => {
        this.messageService.dispatchError(err, '');
        this.checkboxRevertInfo = {isNeeded: true, value: oldServicePublic};
      });
  }

  confirmToServiceAction(s: Service, action: string): void {
    if (action == 'DELETE' &&
      (s.service_status != SERVICE_STATUS.STOPPED) &&
      (s.service_status != SERVICE_STATUS.WARNING)) return;
    let serviceData = new ServiceData(s.service_id, s.service_name, (s.service_status === SERVICE_STATUS.RUNNING));
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
    announceMessage.params = [s.service_name];
    announceMessage.target = target;
    announceMessage.buttons = buttonStyle;
    announceMessage.data = serviceData;
    this.messageService.announceMessage(announceMessage);
  }

  openServiceDetail(s: Service) {
    this.serviceDetailComponent.openModal(s);
  }

  openServiceControl(service: Service) {
    this.serviceControlData = service;
    this.isServiceControlOpen = true;
  }

  get isFirstLogin(): boolean {
    return this.appInitService.isFirstLogin;
  }

  get guideStep(): GUIDE_STEP {
    return this.appInitService.guideStep;
  }

  guideNextStep(step: GUIDE_STEP) {
    this.isBuildServiceWIP = true;
    this.setGuideNoneStep();
  }

  setGuideNoneStep() {
    this.appInitService.guideStep = GUIDE_STEP.NONE_STEP;
  }

  setCreateServiceMethod(method: CreateServiceMethod): void {
    this.createServiceMethod = method;
  }

  cancelCreateService() {
    this.createServiceMethod = CreateServiceMethod.None;
    this.isBuildServiceWIP = false;
    this.isShowServiceCreateYaml = false;
  }
}
