import { Component, ComponentFactoryResolver, Injector, OnDestroy, OnInit, ViewContainerRef } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Service } from '../service';
import { ClrDatagridSortOrder, ClrDatagridStateInterface } from '@clr/angular';
import { GUIDE_STEP, SERVICE_STATUS } from '../../shared/shared.const';
import { ServiceDetailComponent } from './service-detail/service-detail.component';
import { ServiceStepBase } from "../service-step";
import { Observable } from "rxjs/Observable";
import { Project } from "../../project/project";
import { ServiceControlComponent } from "./service-control/service-control.component";
import { TranslateService } from "@ngx-translate/core";
import { Message, RETURN_STATUS } from "../../shared/shared.types";
import "rxjs/add/observable/interval";

enum CreateServiceMethod{None, Wizards, YamlFile, DevOps}
@Component({
  templateUrl: './list-service.component.html',
  styleUrls: ["./list-service.component.css"]
})
export class ListServiceComponent extends ServiceStepBase implements OnInit, OnDestroy {
  services: Service[];
  isInLoading: boolean = false;
  totalRecordCount: number;
  pageIndex: number = 1;
  pageSize: number = 15;
  isBuildServiceWIP: boolean = false;
  isShowServiceCreateYaml: boolean = false;
  createServiceMethod: CreateServiceMethod = CreateServiceMethod.None;
  isActionWIP: Map<number, boolean>;
  projectList: Array<Project>;
  descSort = ClrDatagridSortOrder.DESC;
  oldStateInfo: ClrDatagridStateInterface;
  private subscriptionInterval: Subscription;

  constructor(protected injector: Injector,
              private translateService: TranslateService,
              private viewRef: ViewContainerRef,
              private factory: ComponentFactoryResolver) {
    super(injector);
    this.subscriptionInterval = Observable.interval(10000).subscribe(() => this.retrieve(true, this.oldStateInfo));
    this.isActionWIP = new Map<number, boolean>();
    this.projectList = Array<Project>();
  }

  ngOnInit(): void {
    this.k8sService.getProjects().subscribe((res: Array<Project>) => this.projectList = res);
  }

  ngOnDestroy(): void {
    this.subscriptionInterval.unsubscribe();
  }

  isServiceInStoppedStatus(s: Service): boolean {
    return s.service_status == SERVICE_STATUS.STOPPED && !this.isActionWIP.get(s.service_id);
  }

  isServiceCanChangePauseStatus(service: Service): boolean {
    return service.service_status in [SERVICE_STATUS.RUNNING, SERVICE_STATUS.WARNING]
      && !this.isActionWIP.get(service.service_id);
  }

  isDeleteDisable(service: Service): boolean{
    return service.service_status in [SERVICE_STATUS.PREPARING, SERVICE_STATUS.RUNNING]
      || this.isActionWIP.get(service.service_id)
      || service.service_is_member == 0;
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

  retrieve(isAuto: boolean, stateInfo: ClrDatagridStateInterface): void {
    if (this.isNormalStatus && stateInfo) {
      setTimeout(()=>{
        this.isInLoading = !isAuto;
        this.oldStateInfo = stateInfo;
        this.k8sService.getServices(this.pageIndex, this.pageSize, stateInfo.sort.by as string, stateInfo.sort.reverse).subscribe(
          paginatedServices => {
            this.totalRecordCount = paginatedServices["pagination"]["total_count"];
            this.services = paginatedServices["service_status_list"];
            this.isInLoading = false;
          }, () => this.isInLoading = false
        );
      });
    }
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

  toggleServicePublic(service: Service, $event:MouseEvent): void {
    let oldServicePublic = service.service_public;
    this.k8sService.toggleServicePublicity(service.service_id, service.service_public == 1 ? 0 : 1).subscribe(() => {
        service.service_public = oldServicePublic == 1 ? 0 : 1;
        this.messageService.showAlert('SERVICE.SUCCESSFUL_TOGGLE')
      }, () => ($event.srcElement as HTMLInputElement).checked = oldServicePublic == 1
    );
  }

  toggleService(service: Service){
    if (service.service_is_member == 1){
      this.translateService.get('SERVICE.CONFIRM_TO_TOGGLE_SERVICE', [service.service_name]).subscribe((msg: string) => {
        this.messageService.showConfirmationDialog(msg, 'SERVICE.TOGGLE_SERVICE').subscribe((message: Message) => {
          if (message.returnStatus == RETURN_STATUS.rsConfirm) {
            this.k8sService.toggleServiceStatus(service.service_id, service.service_status == 1 ? 0 : 1).subscribe(() => {
                this.messageService.showAlert('SERVICE.SUCCESSFUL_TOGGLE');
                this.isActionWIP.set(service.service_id, false);
                this.retrieve(false, this.oldStateInfo);
              }, () => this.isActionWIP.set(service.service_id, false)
            );
          }
        });
      });
    }
  }

  deleteService(service:Service){
    if (!this.isDeleteDisable(service)){
      this.translateService.get('SERVICE.CONFIRM_TO_DELETE_SERVICE', [service.service_name]).subscribe((msg: string) => {
        this.messageService.showDeleteDialog(msg, 'SERVICE.DELETE_SERVICE').subscribe((message: Message) => {
          if (message.returnStatus == RETURN_STATUS.rsConfirm) {
            this.k8sService.deleteService(service.service_id).subscribe(() => {
                this.messageService.showAlert('SERVICE.SUCCESSFUL_DELETE');
                this.isActionWIP.set(service.service_id, false);
                this.retrieve(false, this.oldStateInfo);
              }, () => this.isActionWIP.set(service.service_id, false)
            );
          }
        });
      });
    }
  }

  openServiceDetail(service: Service) {
    if (service.service_status != 2){
      let factory = this.factory.resolveComponentFactory(ServiceDetailComponent);
      let componentRef = this.viewRef.createComponent(factory);
      componentRef.instance.openModal(service)
        .subscribe(() => this.viewRef.remove(this.viewRef.indexOf(componentRef.hostView)));
    }
  }

  openServiceControl(service: Service) {
    if (service.service_is_member == 1 && service.service_status == SERVICE_STATUS.RUNNING){
      let factory = this.factory.resolveComponentFactory(ServiceControlComponent);
      let componentRef = this.viewRef.createComponent(factory);
      componentRef.instance.service = service;
      componentRef.instance.openModal()
        .subscribe(() => this.viewRef.remove(this.viewRef.indexOf(componentRef.hostView)));
    }
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

  isSystemAdminOrOwner(service: Service): boolean {
      return this.appInitService.currentUser.user_system_admin == 1 ||
        service.service_owner_id == this.appInitService.currentUser.user_id;
  }
}
