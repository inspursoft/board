import { Component, ComponentFactoryResolver, Injector, OnDestroy, OnInit, ViewContainerRef } from '@angular/core';
import { ClrDatagridSortOrder, ClrDatagridStateInterface } from '@clr/angular';
import { GUIDE_STEP, SERVICE_STATUS } from '../../shared/shared.const';
import { ServiceDetailComponent } from './service-detail/service-detail.component';
import { ServiceControlComponent } from './service-control/service-control.component';
import { TranslateService } from '@ngx-translate/core';
import { Message, RETURN_STATUS } from '../../shared/shared.types';
import { interval, Observable, Subscription } from 'rxjs';
import { ServiceStepComponentBase } from '../service-step';
import { PaginationService, Service, ServiceProject, ServiceSource, ServiceType } from '../service.types';

enum CreateServiceMethod {None, Wizards, YamlFile, DevOps}

@Component({
  templateUrl: './list-service.component.html',
  styleUrls: ['./list-service.component.css']
})
export class ListServiceComponent extends ServiceStepComponentBase implements OnInit, OnDestroy {
  services: PaginationService;
  isInLoading = false;
  totalRecordCount = 0;
  pageIndex = 1;
  pageSize = 10;
  isBuildServiceWIP = false;
  isShowServiceCreateYaml = false;
  createServiceMethod: CreateServiceMethod = CreateServiceMethod.None;
  isActionWIP: Map<number, boolean>;
  projectList: Array<ServiceProject>;
  descSort = ClrDatagridSortOrder.DESC;
  oldStateInfo: ClrDatagridStateInterface;
  private subscriptionInterval: Subscription;

  constructor(protected injector: Injector,
              private translateService: TranslateService,
              private viewRef: ViewContainerRef,
              private factory: ComponentFactoryResolver) {
    super(injector);
    this.subscriptionInterval = interval(10000).subscribe(() => this.retrieve(true, this.oldStateInfo));
    this.isActionWIP = new Map<number, boolean>();
    this.projectList = Array<ServiceProject>();
  }

  ngOnInit(): void {
    this.k8sService.getProjects().subscribe((res: Array<ServiceProject>) => this.projectList = res);
  }

  ngOnDestroy(): void {
    this.subscriptionInterval.unsubscribe();
  }

  checkWithinPreparingRunning(status: number): boolean {
    return [SERVICE_STATUS.RUNNING, SERVICE_STATUS.PREPARING].indexOf(status) > -1;
  }

  isServiceCanPlay(service: Service): boolean {
    return service.serviceStatus === SERVICE_STATUS.STOPPED;
  }

  isServiceCanPause(service: Service): boolean {
    return [SERVICE_STATUS.RUNNING, SERVICE_STATUS.WARNING].indexOf(service.serviceStatus) > -1;
  }

  isServiceToggleDisabled(service: Service): boolean {
    return this.isActionWIP.get(service.serviceId)
      || service.serviceIsMember === 0
      || service.serviceSource === ServiceSource.ServiceSourceHelm;
  }

  isDeleteDisable(service: Service): boolean {
    return this.isActionWIP.get(service.serviceId)
      || this.checkWithinPreparingRunning(service.serviceStatus)
      || service.serviceIsMember === 0
      || service.serviceSource === ServiceSource.ServiceSourceHelm;
  }

  isUpdateDisable(service: Service): boolean {
    return this.isActionWIP.get(service.serviceId)
      || [SERVICE_STATUS.RUNNING, SERVICE_STATUS.WARNING].indexOf(service.serviceStatus) === -1
      || service.serviceIsMember === 0
      || service.serviceType === ServiceType.ServiceTypeStatefulSet
      || service.serviceSource === ServiceSource.ServiceSourceHelm;
  }

  serviceToggleTipMessage(service: Service): string {
    if (service.serviceIsMember === 0) {
      return 'SERVICE.STEP_0_NOT_SERVICE_MEMBER';
    } else if (service.serviceSource === ServiceSource.ServiceSourceHelm) {
      return 'SERVICE.STEP_0_SERVICE_FROM_HELM';
    }
    return '';
  }

  serviceDeleteTipMessage(service: Service): string {
    if (service.serviceIsMember === 0) {
      return 'SERVICE.STEP_0_NOT_SERVICE_MEMBER';
    } else if (service.serviceSource === ServiceSource.ServiceSourceHelm) {
      return 'SERVICE.STEP_0_SERVICE_FROM_HELM';
    } else if (this.checkWithinPreparingRunning(service.serviceStatus)) {
      return 'SERVICE.STEP_0_CAN_NOT_DELETE_MSG';
    }
    return '';
  }

  serviceUpdateTipMessage(service: Service): string {
    if (service.serviceIsMember === 0) {
      return 'SERVICE.STEP_0_NOT_SERVICE_MEMBER';
    } else if (service.serviceSource === ServiceSource.ServiceSourceHelm) {
      return 'SERVICE.STEP_0_SERVICE_FROM_HELM';
    } else if (service.serviceType === ServiceType.ServiceTypeStatefulSet) {
      return 'SERVICE.STEP_0_SERVICE_STATEFUL';
    } else if (SERVICE_STATUS.RUNNING !== service.serviceStatus) {
      return 'SERVICE.STEP_0_CAN_NOT_UPDATE_MSG';
    }
    return '';
  }

  createService(): void {
    if (this.createServiceMethod === CreateServiceMethod.Wizards) {
      this.k8sService.stepSource.next({index: 1, isBack: false});
    } else if (this.createServiceMethod === CreateServiceMethod.YamlFile) {
      this.isShowServiceCreateYaml = true;
    }
  }

  get isNormalStatus(): boolean {
    return !this.isBuildServiceWIP && !this.isShowServiceCreateYaml;
  }

  retrieve(isAuto: boolean, stateInfo: ClrDatagridStateInterface): void {
    if (this.isNormalStatus && stateInfo) {
      setTimeout(() => {
        this.isInLoading = !isAuto;
        this.oldStateInfo = stateInfo;
        this.k8sService.getServices(this.pageIndex, this.pageSize, stateInfo.sort.by as string, stateInfo.sort.reverse).subscribe(
          res => {
            this.services = res;
            this.totalRecordCount = this.services.pagination.TotalCount;
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
      case SERVICE_STATUS.UnKnown:
        return 'SERVICE.STATUS_UNKNOWN';
      case SERVICE_STATUS.AutonomousOffline:
        return 'SERVICE.STATUS_EDGE';
      case SERVICE_STATUS.STOPPING:
        return 'SERVICE.STATUS_STOPPING';
    }
  }

  getStatusClass(status: SERVICE_STATUS) {
    return {
      running: status === SERVICE_STATUS.RUNNING,
      stopped: status === SERVICE_STATUS.STOPPED,
      warning: status === SERVICE_STATUS.WARNING || status === SERVICE_STATUS.AutonomousOffline
    };
  }

  toggleServicePublic(service: Service, $event: MouseEvent): void {
    const oldServicePublic = service.servicePublic;
    this.k8sService.toggleServicePublicity(service.serviceId, service.servicePublic === 1 ? 0 : 1).subscribe(() => {
        service.servicePublic = oldServicePublic === 1 ? 0 : 1;
        this.messageService.showAlert('SERVICE.SUCCESSFUL_TOGGLE');
      }, () => ($event.srcElement as HTMLInputElement).checked = oldServicePublic === 1
    );
  }

  toggleService(service: Service) {
    if (!this.isServiceToggleDisabled(service)) {
      this.translateService.get('SERVICE.CONFIRM_TO_TOGGLE_SERVICE', [service.serviceName]).subscribe((msg: string) => {
        this.messageService.showConfirmationDialog(msg, 'SERVICE.TOGGLE_SERVICE').subscribe((message: Message) => {
          if (message.returnStatus === RETURN_STATUS.rsConfirm) {
            const toggleValue =
              service.serviceStatus === SERVICE_STATUS.RUNNING ? 0 :
                service.serviceStatus === SERVICE_STATUS.WARNING ? 0 : 1;
            this.k8sService.toggleServiceStatus(service.serviceId, toggleValue).subscribe(
              () => {
                this.messageService.showAlert('SERVICE.SUCCESSFUL_TOGGLE');
                this.isActionWIP.set(service.serviceId, false);
                this.retrieve(false, this.oldStateInfo);
              },
              () => {
                this.isActionWIP.set(service.serviceId, false);
                this.messageService.showAlert('SERVICE.SERVICE_NOT_SUPPORT_TOGGLE', {alertType: 'warning'});
              });
          }
        });
      });
    }
  }

  deleteService(service: Service) {
    if (!this.isDeleteDisable(service)) {
      this.translateService.get('SERVICE.CONFIRM_TO_DELETE_SERVICE', [service.serviceName]).subscribe((msg: string) => {
        this.messageService.showDeleteDialog(msg, 'SERVICE.DELETE_SERVICE').subscribe((message: Message) => {
          if (message.returnStatus === RETURN_STATUS.rsConfirm) {
            const obsDelete: Observable<any> = service.serviceType === ServiceType.ServiceTypeStatefulSet ?
              this.k8sService.deleteStatefulService(service.serviceId) :
              this.k8sService.deleteService(service.serviceId);
            obsDelete.subscribe(() => {
                this.messageService.showAlert('SERVICE.SUCCESSFUL_DELETE');
                this.isActionWIP.set(service.serviceId, false);
                this.retrieve(false, this.oldStateInfo);
              }, () => this.isActionWIP.set(service.serviceId, false)
            );
          }
        });
      });
    }
  }

  openServiceDetail(service: Service) {
    if (service.serviceStatus !== SERVICE_STATUS.STOPPED) {
      const factory = this.factory.resolveComponentFactory(ServiceDetailComponent);
      const componentRef = this.viewRef.createComponent(factory);
      componentRef.instance.openModal(service)
        .subscribe(() => this.viewRef.remove(this.viewRef.indexOf(componentRef.hostView)));
    }
  }

  openServiceControl(service: Service) {
    if (!this.isUpdateDisable(service)) {
      const factory = this.factory.resolveComponentFactory(ServiceControlComponent);
      const componentRef = this.viewRef.createComponent(factory);
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
    return this.appInitService.currentUser.userSystemAdmin === 1 ||
      service.serviceOwnerId === this.appInitService.currentUser.userId;
  }
}
