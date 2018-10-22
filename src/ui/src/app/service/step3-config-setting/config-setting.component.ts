import { ChangeDetectorRef, Component, Injector, OnInit } from '@angular/core';
import { ValidationErrors } from "@angular/forms/forms";
import { HttpErrorResponse } from "@angular/common/http";
import { Observable } from "rxjs/Observable";
import {
  Container,
  ExternalService,
  PHASE_CONFIG_CONTAINERS,
  PHASE_EXTERNAL_SERVICE,
  ServiceStepPhase,
  UIServiceStep3,
  UIServiceStep4,
  UIServiceStepBase
} from '../service-step.component';
import { ServiceStepBase } from "../service-step";
import { IDropdownTag } from "../../shared/shared.types";
import { SetAffinityComponent } from "./set-affinity/set-affinity.component";

@Component({
  styleUrls: ["./config-setting.component.css"],
  templateUrl: './config-setting.component.html'
})
export class ConfigSettingComponent extends ServiceStepBase implements OnInit {
  patternServiceName: RegExp = /^[a-z]([-a-z0-9]*[a-z0-9])+$/;
  showAdvanced = false;
  showNodeSelector = false;
  isActionWip: boolean = false;
  nodeSelectorList: Array<{name: string, value: string, tag: IDropdownTag}>;
  uiPreData: UIServiceStep3;

  constructor(protected injector: Injector,
              private changeDetectorRef: ChangeDetectorRef) {
    super(injector);
    this.nodeSelectorList = Array<{name: string, value: string, tag: IDropdownTag}>();
    this.uiPreData = new UIServiceStep3();
  }

  ngOnInit() {
    let obsStepConfig = this.k8sService.getServiceConfig(this.stepPhase);
    let obsPreStepConfig = this.k8sService.getServiceConfig(PHASE_CONFIG_CONTAINERS);
    Observable.forkJoin(obsStepConfig, obsPreStepConfig).subscribe((res: [UIServiceStepBase, UIServiceStepBase]) => {
      this.uiBaseData = res[0];
      this.uiPreData = res[1] as UIServiceStep3;
      if (this.uiData.externalServiceList.length === 0 && this.uiPreData.containerHavePortList.length > 0) {
        let container = this.uiPreData.containerHavePortList[0];
        this.addNewExternalService();
        this.setExternalInfo(container, 0);
      }
      this.changeDetectorRef.detectChanges();
    });
    this.nodeSelectorList.push({name: 'SERVICE.STEP_3_NODE_DEFAULT', value: '', tag: null});
    this.k8sService.getNodeSelectors().subscribe((res: Array<{name: string, status: number}>) => {
      res.forEach((value: {name: string, status: number}) => {
        this.nodeSelectorList.push({
          name: value.name, value: value.name, tag: {
            type: value.status == 1 ? 'alert-success' : 'alert-warning',
            description: value.status == 1 ? 'SERVICE.STEP_3_NODE_STATUS_SCHEDULABLE' : 'SERVICE.STEP_3_NODE_STATUS_UNSCHEDULABLE'
          }
        })
      });
    });
  }

  get stepPhase(): ServiceStepPhase {
    return PHASE_EXTERNAL_SERVICE
  }

  get uiData(): UIServiceStep4 {
    return this.uiBaseData as UIServiceStep4;
  }

  get checkServiceNameFun() {
    return this.checkServiceName.bind(this);
  }

  getContainerDropdownText(index: number): string {
    let result = this.uiData.externalServiceList[index].container_name;
    return result == "" ? "SERVICE.STEP_3_SELECT_CONTAINER" : result;
  }

  setExternalInfo(container: Container, index: number) {
    this.uiData.externalServiceList[index].container_name = container.name;
    this.uiData.externalServiceList[index].node_config.target_port = container.container_port[0];
  }

  setNodePort(index: number, port: number) {
    this.uiData.externalServiceList[index].node_config.node_port = Number(port).valueOf();
  }

  addNewExternalService() {
    if (this.uiPreData.containerHavePortList.length > 0) {
      let externalService = new ExternalService();
      this.uiData.externalServiceList.push(externalService);
    }
  }

  removeExternalService(index: number) {
    this.uiData.externalServiceList.splice(index, 1);
  }

  setAffinity() {
    let factory = this.factoryResolver.resolveComponentFactory(SetAffinityComponent);
    let componentRef = this.selfView.createComponent(factory);
    componentRef.instance.openSetModal(this.uiData).subscribe(() => this.selfView.remove(this.selfView.indexOf(componentRef.hostView)));
  }

  checkServiceName(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.k8sService.checkServiceExist(this.uiData.projectName, control.value)
      .map(() => null)
      .catch((err:HttpErrorResponse) => {
        if (err.status == 409) {
          this.messageService.cleanNotification();
          return Observable.of({serviceExist: "SERVICE.STEP_3_SERVICE_NAME_EXIST"});
        } else if (err.status == 404) {
          this.messageService.cleanNotification();
        }
        return Observable.of(null);
      });
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
    if (this.verifyInputValid()) {
      if (this.uiData.externalServiceList.length == 0) {
        this.messageService.showAlert(`SERVICE.STEP_3_EXTERNAL_MESSAGE`, {alertType: "alert-warning"});
      } else if (this.haveRepeatNodePort()) {
        this.messageService.showAlert(`SERVICE.STEP_3_EXTERNAL_REPEAT`, {alertType: "alert-warning"});
      } else if (this.uiData.affinityList.find(value => value.services.length == 0)) {
        this.messageService.showAlert(`SERVICE.STEP_3_AFFINITY_MESSAGE`, {alertType: "alert-warning"});
      } else {
        this.isActionWip = true;
        this.k8sService.setServiceConfig(this.uiData.uiToServer()).subscribe(
          () => this.k8sService.stepSource.next({index: 5, isBack: false})
        );
      }
    }
  }

  backUpStep(): void {
    this.k8sService.stepSource.next({index: 2, isBack: true});
  }
}