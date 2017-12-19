import { Component, OnInit, AfterContentChecked, ViewChildren, QueryList, Injector } from '@angular/core';
import {
  PHASE_EXTERNAL_SERVICE,
  PHASE_CONFIG_CONTAINERS,
  Container,
  ServiceStepPhase,
  UIServiceStep4,
  UIServiceStep3,
  ExternalService
} from '../service-step.component';
import { CsInputComponent } from "../../shared/cs-components-library/cs-input/cs-input.component";
import { ServiceStepBase } from "../service-step";
import { Response } from "@angular/http";
import { ValidationErrors } from "@angular/forms/forms";

@Component({
  styleUrls: ["./config-setting.component.css"],
  templateUrl: './config-setting.component.html'
})
export class ConfigSettingComponent extends ServiceStepBase implements OnInit, AfterContentChecked {
  @ViewChildren(CsInputComponent) inputComponents: QueryList<CsInputComponent>;
  patternServiceName: RegExp = /^[a-z]([-a-z0-9]*[a-z0-9])+$/;
  dropDownListNum: Array<number>;
  showAdvanced: boolean = true;
  showExternal: boolean = false;
  showCollaborative: boolean = false;
  isInputComponentsValid = false;
  uiPreData: UIServiceStep3 = new UIServiceStep3();
  collaborativeServiceList: Array<string>;

  constructor(protected injector: Injector) {
    super(injector);
    this.dropDownListNum = Array<number>();
    this.collaborativeServiceList = Array<string>();
  }

  ngAfterContentChecked() {
    this.isInputComponentsValid = true;
    if (this.inputComponents) {
      this.inputComponents.forEach(item => {
        if (!item.valid) {
          this.isInputComponentsValid = false;
        }
      });
    }
  }

  ngOnInit() {
    this.k8sService.getServiceConfig(PHASE_CONFIG_CONTAINERS).then(res => {
      this.uiPreData = res as UIServiceStep3;
    });
    this.k8sService.getServiceConfig(this.stepPhase).then(res => {
      this.uiBaseData = res;
    });
    for (let i = 1; i <= 100; i++) {
      this.dropDownListNum.push(i)
    }
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

  checkServiceName(control: HTMLInputElement): Promise<ValidationErrors | null> {
    return this.k8sService.checkServiceExist(this.uiData.projectName, control.value)
      .then(() => null)
      .catch(err => {
        if (err && err instanceof Response && (err as Response).status == 409) {
          return {serviceExist: "SERVICE.STEP_4_SERVICE_NAME_EXIST"}
        } else {
          return null;
        }
      });
  }

  setNodePort(index: number, port: number) {
    this.uiData.externalServiceList[index].node_config.node_port = Number(port).valueOf();
  }

  setServiceName(serviceName: string): void {
    this.uiData.serviceName = serviceName;
    /*Todo:add reset the Collaborative service Info*/
    this.collaborativeServiceList.splice(0, this.collaborativeServiceList.length);
    this.k8sService.getCollaborativeService(serviceName, this.uiData.projectName)
      .then(res => {
        this.collaborativeServiceList = res;
      })
  }

  addContainerInfo() {
    this.uiData.externalServiceList.push(new ExternalService());
  }

  removeContainerInfo(index: number) {
    this.uiData.externalServiceList.splice(index, 1);
  }

  getContainerDropdownText(index: number): string {
    let result = this.uiData.externalServiceList[index].container_name;
    return result == "" ? "SERVICE.STEP_4_SELECT_CONTAINER" : result;
  }

  getContainerPortDropdownText(index: number): string {
    let result = this.uiData.externalServiceList[index].node_config.target_port;
    return result == 0 ? "SERVICE.STEP_4_SELECT_PORT" : result.toString();
  }

  setExternalInfo(container: Container, index: number) {
    this.uiData.externalServiceList[index].container_name = container.name;
  }

  getContainerPorts(containerName: string): Array<number> {
    let result: Array<number> = Array<number>();
    this.uiPreData.containerList.forEach((container: Container) => {
      if (container.name == containerName) {
        result = container.container_port;
      }
    });
    return result;
  }

  forward(): void {
    this.k8sService.setServiceConfig(this.uiData.uiToServer()).then(() => {
      this.k8sService.stepSource.next({index: 6, isBack: false});
    });
  }

  backUpStep(): void {
    this.k8sService.stepSource.next({index: 3, isBack: true});
  }
}