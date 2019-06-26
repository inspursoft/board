import { ChangeDetectorRef, Component, Injector, OnInit } from '@angular/core';
import { ValidationErrors } from "@angular/forms/forms";
import { HttpErrorResponse } from "@angular/common/http";
import {
  Container,
  ExternalService,
  PHASE_CONFIG_CONTAINERS,
  PHASE_EXTERNAL_SERVICE,
  ServiceStepPhase,
  UIServiceStep2,
  UIServiceStep3,
  UIServiceStepBase
} from '../service-step.component';
import { ServiceStepBase } from "../service-step";
import { IDropdownTag } from "../../shared/shared.types";
import { SetAffinityComponent } from "./set-affinity/set-affinity.component";
import { forkJoin, Observable, of } from "rxjs";
import { catchError, map } from "rxjs/operators";
import { ServiceType } from "../service";

@Component({
  styleUrls: ["./config-setting.component.css"],
  templateUrl: './config-setting.component.html'
})
export class ConfigSettingComponent extends ServiceStepBase implements OnInit {
  patternServiceName: RegExp = /^[a-z]([-a-z0-9]*[a-z0-9])+$/;
  patternIP: RegExp = /^((?:(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d)))\.){3}(?:25[0-5]|2[0-4]\d|((1\d{2})|([1-9]?\d))))$/;
  showAdvanced = false;
  showNodeSelector = false;
  isActionWip = false;
  isGetNodePortWip = false;
  nodeSelectorList: Array<{name: string, value: string, tag: IDropdownTag}>;
  uiPreData: UIServiceStep2;
  existingNodePorts: Array<number>;
  serviceTypes: Array<{description: string, type: ServiceType}>;

  constructor(protected injector: Injector,
              private changeDetectorRef: ChangeDetectorRef) {
    super(injector);
    this.changeDetectorRef.detach();
    this.nodeSelectorList = Array<{name: string, value: string, tag: IDropdownTag}>();
    this.existingNodePorts = Array<number>();
    this.uiPreData = new UIServiceStep2();
    this.serviceTypes = new Array<{description: string, type: ServiceType}>();
  }

  ngOnInit() {
    let obsStepConfig = this.k8sService.getServiceConfig(this.stepPhase);
    let obsPreStepConfig = this.k8sService.getServiceConfig(PHASE_CONFIG_CONTAINERS);
    this.serviceTypes.push({description: 'SERVICE.STEP_3_SERVICE_TYPE_NORMAL', type: ServiceType.ServiceTypeNormalNodePort});
    this.serviceTypes.push({description: 'SERVICE.STEP_3_SERVICE_TYPE_STATEFUL', type: ServiceType.ServiceTypeStatefulSet});
    this.isGetNodePortWip = true;
    forkJoin(obsStepConfig, obsPreStepConfig).subscribe((res: [UIServiceStepBase, UIServiceStepBase]) => {
      this.uiBaseData = res[0];
      this.uiPreData = res[1] as UIServiceStep2;
      if (this.uiData.externalServiceList.length === 0) {
        let container = this.uiPreData.containerList[0];
        this.addNewExternalService();
        this.setExternalInfo(container, 0);
      }
      this.k8sService.getNodePorts(this.uiData.projectName).subscribe(
        (res: Array<number>) => this.existingNodePorts = res,
        () => this.isGetNodePortWip = false,
        () => this.isGetNodePortWip = false
      );
      this.changeDetectorRef.reattach();
    });
    this.nodeSelectorList.push({name: 'SERVICE.STEP_3_NODE_DEFAULT', value: '', tag: null});
    this.k8sService.getNodeSelectors().subscribe((res: Array<{name: string, status: number}>) => {
      res.forEach((value: {name: string, status: number}) => {
        this.nodeSelectorList.push({
          name: value.name, value: value.name, tag: {
            type: value.status == 1 ? 'success' : 'warning',
            description: value.status == 1 ? 'SERVICE.STEP_3_NODE_STATUS_SCHEDULABLE' : 'SERVICE.STEP_3_NODE_STATUS_UNSCHEDULABLE'
          }
        })
      });
    });
  }

  get stepPhase(): ServiceStepPhase {
    return PHASE_EXTERNAL_SERVICE
  }

  get uiData(): UIServiceStep3 {
    return this.uiBaseData as UIServiceStep3;
  }

  get checkServiceNameFun() {
    return this.checkServiceName.bind(this);
  }

  get nodeSelectorDropdownText() {
    return this.uiData.nodeSelector === '' ? 'SERVICE.STEP_3_NODE_DEFAULT' : this.uiData.nodeSelector;
  }

  get curNodeSelector() {
    return this.nodeSelectorList.find(value => value.name === this.uiData.nodeSelector);
  }

  get serviceTypeDescription() {
    return this.serviceTypes.find(value => value.type == this.uiData.serviceType).description;
  }

  get isStatefulService() {
    return this.uiData.serviceType == ServiceType.ServiceTypeStatefulSet;
  }

  getContainerDropdownText(index: number): string {
    let result = this.uiData.externalServiceList[index].container_name;
    return result == "" ? "SERVICE.STEP_3_SELECT_CONTAINER" : result;
  }

  getContainerPorts(containerName: string): Array<number> {
    return this.uiPreData.getPortList(containerName);
  }

  setExternalInfo(container: Container, index: number) {
    this.uiData.externalServiceList[index].container_name = container.name;
    this.uiData.externalServiceList[index].node_config.target_port = container.container_port.length > 0 ? container.container_port[0] : 0;
  }

  setNodePort(index: number, port: number) {
    this.uiData.externalServiceList[index].node_config.node_port = Number(port).valueOf();
  }

  setServiceType(value: {description: string, type: ServiceType}) {
    this.uiData.serviceType = value.type;
    if (value.type === ServiceType.ServiceTypeStatefulSet) {
      this.uiData.externalServiceList.forEach((external:ExternalService) => external.node_config.node_port = 0);
    }
  }

  inputExternalPortEnable(containerName: string): boolean {
    return this.uiPreData.getPortList(containerName).length == 0;
  }

  selectExternalPortEnable(containerName: string): boolean {
    return this.uiPreData.getPortList(containerName).length > 0;
  }

  addNewExternalService() {
    if (this.uiPreData.containerList.length > 0 && !this.isActionWip) {
      let externalService = new ExternalService();
      this.uiData.externalServiceList.push(externalService);
    }
  }

  removeExternalService(index: number) {
    this.uiData.externalServiceList.splice(index, 1);
  }

  setAffinity() {
    if (!this.isActionWip) {
      let factory = this.factoryResolver.resolveComponentFactory(SetAffinityComponent);
      let componentRef = this.selfView.createComponent(factory);
      componentRef.instance.openSetModal(this.uiData).subscribe(() => this.selfView.remove(this.selfView.indexOf(componentRef.hostView)));
    }
  }

  setNodeSelector() {
    if (!this.isActionWip) {
      this.showNodeSelector = !this.showNodeSelector;
    }
  }

  checkServiceName(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.k8sService.checkServiceExist(this.uiData.projectName, control.value)
      .pipe(
        map(() => null),
        catchError((err: HttpErrorResponse) => {
          if (err.status == 409) {
            this.messageService.cleanNotification();
            return of({serviceExist: "SERVICE.STEP_3_SERVICE_NAME_EXIST"});
          } else if (err.status == 404) {
            this.messageService.cleanNotification();
          }
          return of(null);
        }));
  }

  haveRepeatNodePort(): boolean {
    let haveRepeat = false;
    this.uiData.externalServiceList.forEach((value, index) => {
      if (this.uiData.externalServiceList.find((value1, index1) =>
        value1.container_name === value.container_name
        && value1.node_config.target_port === value.node_config.target_port
        && index1 !== index)) {
        haveRepeat = true
      }
    });
    return haveRepeat;
  }

  forward(): void {
    if (this.verifyInputValid() && this.verifyInputDropdownValid()) {
      if (this.uiData.externalServiceList.length == 0) {
        this.messageService.showAlert(`SERVICE.STEP_3_EXTERNAL_MESSAGE`, {alertType: "warning"});
      } else if (this.haveRepeatNodePort()) {
        this.messageService.showAlert(`SERVICE.STEP_3_EXTERNAL_REPEAT`, {alertType: "warning"});
      } else if (this.uiData.affinityList.find(value => value.services.length == 0)) {
        this.messageService.showAlert(`SERVICE.STEP_3_AFFINITY_MESSAGE`, {alertType: "warning"});
      } else {
        this.isActionWip = true;
        this.k8sService.setServiceConfig(this.uiData.uiToServer()).subscribe(
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
