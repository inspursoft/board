import { ChangeDetectorRef, Component, Injector, OnInit } from '@angular/core';
import { ValidationErrors } from '@angular/forms/forms';
import { HttpErrorResponse } from '@angular/common/http';
import { forkJoin, Observable, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import {
  Container,
  ExternalService,
  PHASE_CONFIG_CONTAINERS,
  PHASE_EXTERNAL_SERVICE, ServiceStep2Data, ServiceStep3Data,
  ServiceStepPhase,
} from '../service-step.component';
import { IDropdownTag } from '../../shared/shared.types';
import { SetAffinityComponent } from './set-affinity/set-affinity.component';
import { ServiceStepComponentBase } from '../service-step';
import { ServiceType } from '../service.types';

@Component({
  styleUrls: ['./config-setting.component.css'],
  templateUrl: './config-setting.component.html'
})
export class ConfigSettingComponent extends ServiceStepComponentBase implements OnInit {
  patternServiceName: RegExp = /^[a-z]([-a-z0-9]*[a-z0-9])+$/;
  patternIP: RegExp = /^((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))))$/;
  showAdvanced = false;
  showNodeSelector = false;
  isActionWip = false;
  isGetNodePortWip = false;
  nodeSelectorList: Array<{ name: string, value: string, tag: IDropdownTag }>;
  serviceStep2Data: ServiceStep2Data;
  serviceStep3Data: ServiceStep3Data;
  existingNodePorts: Array<number>;
  serviceTypes: Array<{ description: string, type: ServiceType }>;

  constructor(protected injector: Injector,
              private changeDetectorRef: ChangeDetectorRef) {
    super(injector);
    this.changeDetectorRef.detach();
    this.nodeSelectorList = Array<{ name: string, value: string, tag: IDropdownTag }>();
    this.existingNodePorts = Array<number>();
    this.serviceStep2Data = new ServiceStep2Data();
    this.serviceStep3Data = new ServiceStep3Data();
    this.serviceTypes = new Array<{ description: string, type: ServiceType }>();
  }

  ngOnInit() {
    const obsStepConfig = this.k8sService.getServiceConfig(this.stepPhase, ServiceStep3Data);
    const obsPreStepConfig = this.k8sService.getServiceConfig(PHASE_CONFIG_CONTAINERS, ServiceStep2Data);
    this.serviceTypes.push({
      description: 'SERVICE.STEP_3_SERVICE_TYPE_NORMAL',
      type: ServiceType.ServiceTypeNormalNodePort
    });
    this.serviceTypes.push({
      description: 'SERVICE.STEP_3_SERVICE_TYPE_STATEFUL',
      type: ServiceType.ServiceTypeStatefulSet
    });
    this.isGetNodePortWip = true;
    forkJoin(obsPreStepConfig, obsStepConfig).subscribe((res: [ServiceStep2Data, ServiceStep3Data]) => {
      this.serviceStep2Data = res[0];
      this.serviceStep3Data = res[1];
      if (this.serviceStep3Data.externalServiceList.length === 0) {
        const container = this.serviceStep2Data.containerList[0];
        this.addNewExternalService();
        this.setExternalInfo(container, 0);
      }
      this.k8sService.getNodePorts(this.serviceStep3Data.projectName).subscribe(
        (ports: Array<number>) => this.existingNodePorts = ports,
        () => this.isGetNodePortWip = false,
        () => this.isGetNodePortWip = false
      );
      this.changeDetectorRef.reattach();
    });
    this.nodeSelectorList.push({name: 'SERVICE.STEP_3_NODE_DEFAULT', value: '', tag: null});
    this.k8sService.getNodeSelectors().subscribe((res: Array<{ name: string, status: number }>) => {
      res.forEach((value: { name: string, status: number }) => {
        this.nodeSelectorList.push({
          name: value.name, value: value.name, tag: {
            type: value.status === 1 ? 'success' : 'warning',
            description: value.status === 1 ? 'SERVICE.STEP_3_NODE_STATUS_SCHEDULABLE' : 'SERVICE.STEP_3_NODE_STATUS_UNSCHEDULABLE'
          }
        });
      });
    });
  }

  get stepPhase(): ServiceStepPhase {
    return PHASE_EXTERNAL_SERVICE;
  }

  get checkServiceNameFun() {
    return this.checkServiceName.bind(this);
  }

  get nodeSelectorDropdownText() {
    return this.serviceStep3Data.nodeSelector === '' ? 'SERVICE.STEP_3_NODE_DEFAULT' : this.serviceStep3Data.nodeSelector;
  }

  get curNodeSelector() {
    return this.nodeSelectorList.find(value => value.name === this.serviceStep3Data.nodeSelector);
  }

  get isStatefulService() {
    return this.serviceStep3Data.serviceType === ServiceType.ServiceTypeStatefulSet;
  }

  getItemTagClass(dropdownTag: IDropdownTag) {
    return {
      'label-info': dropdownTag.type === 'success',
      'label-warning': dropdownTag.type === 'warning',
      'label-danger': dropdownTag.type === 'danger'
    };
  }

  getContainerPorts(containerName: string): Array<number> {
    return this.serviceStep2Data.getPortList(containerName);
  }

  setExternalInfo(container: Container, index: number) {
    this.serviceStep3Data.externalServiceList[index].containerName = container.name;
    this.serviceStep3Data.externalServiceList[index].nodeConfig.targetPort =
      container.containerPort.length > 0 ? container.containerPort[0] : 0;
  }

  setNodePort(index: number, port: number) {
    this.serviceStep3Data.externalServiceList[index].nodeConfig.nodePort = Number(port).valueOf();
  }

  setServiceType(value: { description: string, type: ServiceType }) {
    this.serviceStep3Data.serviceType = value.type;
    if (value.type === ServiceType.ServiceTypeStatefulSet) {
      this.serviceStep3Data.externalServiceList.forEach(
        (external: ExternalService) => external.nodeConfig.nodePort = 0
      );
    } else {
      this.serviceStep3Data.clusterIp = '';
    }
  }

  inputExternalPortEnable(containerName: string): boolean {
    return this.serviceStep2Data.getPortList(containerName).length === 0;
  }

  selectExternalPortEnable(containerName: string): boolean {
    return this.serviceStep2Data.getPortList(containerName).length > 0;
  }

  addNewExternalService() {
    if (this.serviceStep2Data.containerList.length > 0 && !this.isActionWip) {
      const externalService = new ExternalService();
      this.serviceStep3Data.externalServiceList.push(externalService);
    }
  }

  removeExternalService(index: number) {
    this.serviceStep3Data.externalServiceList.splice(index, 1);
  }

  setAffinity() {
    if (!this.isActionWip) {
      const factory = this.factoryResolver.resolveComponentFactory(SetAffinityComponent);
      const componentRef = this.selfView.createComponent(factory);
      componentRef.instance.openSetModal(this.serviceStep3Data).subscribe(
        () => this.selfView.remove(this.selfView.indexOf(componentRef.hostView))
      );
    }
  }

  setNodeSelector() {
    if (!this.isActionWip) {
      this.showNodeSelector = !this.showNodeSelector;
    }
  }

  checkServiceName(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.k8sService.checkServiceExist(this.serviceStep3Data.projectName, control.value)
      .pipe(map(() => null),
        catchError((err: HttpErrorResponse) => {
          if (err.status === 409) {
            this.messageService.cleanNotification();
            return of({serviceExist: 'SERVICE.STEP_3_SERVICE_NAME_EXIST'});
          } else if (err.status === 404) {
            this.messageService.cleanNotification();
          }
          return of(null);
        })
      );
  }

  haveRepeatNodePort(): boolean {
    let haveRepeat = false;
    this.serviceStep3Data.externalServiceList.forEach((value, index) => {
      if (this.serviceStep3Data.externalServiceList.find((value1, index1) =>
        value1.containerName === value.containerName
        && value1.nodeConfig.targetPort === value.nodeConfig.targetPort
        && index1 !== index)) {
        haveRepeat = true;
      }
    });
    return haveRepeat;
  }

  forward(): void {
    if (this.verifyInputExValid() && this.verifyDropdownExValid() && this.verifyInputNumberDropdownValid()) {
      if (this.serviceStep3Data.externalServiceList.length === 0) {
        this.messageService.showAlert(`SERVICE.STEP_3_EXTERNAL_MESSAGE`, {alertType: 'warning'});
      } else if (this.haveRepeatNodePort()) {
        this.messageService.showAlert(`SERVICE.STEP_3_EXTERNAL_REPEAT`, {alertType: 'warning'});
      } else if (this.serviceStep3Data.affinityList.find(value => value.services.length === 0)) {
        this.messageService.showAlert(`SERVICE.STEP_3_AFFINITY_MESSAGE`, {alertType: 'warning'});
      } else {
        this.isActionWip = true;
        this.k8sService.setServiceStepConfig(this.serviceStep3Data).subscribe(
          () => this.k8sService.stepSource.next({index: 5, isBack: false}),
          () => this.isActionWip = false
        );
      }
    }
  }

  backUpStep(): void {
    this.k8sService.stepSource.next({index: 2, isBack: true});
  }
}
