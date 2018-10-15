import { ChangeDetectorRef, Component, Injector, OnInit } from '@angular/core';
import {
  Container,
  ExternalService,
  PHASE_CONFIG_CONTAINERS,
  PHASE_EXTERNAL_SERVICE,
  ServiceStepPhase,
  UIServiceStep3,
  UIServiceStep4
} from '../service-step.component';
import { ServiceStepBase } from "../service-step";
import { ValidationErrors } from "@angular/forms/forms";
import { HttpErrorResponse } from "@angular/common/http";
import { Observable } from "rxjs/Observable";

@Component({
  styleUrls: ["./config-setting.component.css"],
  templateUrl: './config-setting.component.html'
})
export class ConfigSettingComponent extends ServiceStepBase implements OnInit {
  patternServiceName: RegExp = /^[a-z]([-a-z0-9]*[a-z0-9])+$/;
  dropDownListNum: Array<number>;
  showAdvanced: boolean = true;
  showExternal: boolean = false;
  showCollaborative: boolean = false;
  showNodeSelector: boolean = false;
  uiPreData: UIServiceStep3 = new UIServiceStep3();
  collaborativeServiceList: Array<string>;
  /*Todo:Only for collaborative plus action.It must be delete after update UIServiceStep4*/
  collaborativeList:Array<Object>;
  nodeSelectorList:Array<string>;
  noPortForExtent: boolean = false;
  isActionWip: boolean = false;

  constructor(protected injector: Injector, private changeDetectorRef: ChangeDetectorRef) {
    super(injector);
    this.dropDownListNum = Array<number>();
    this.collaborativeServiceList = Array<string>();
    this.collaborativeList = Array<Object>();
    this.nodeSelectorList = Array<string>()
  }

  ngOnInit() {
    this.k8sService.getServiceConfig(PHASE_CONFIG_CONTAINERS).subscribe(res => {
      this.uiPreData = res as UIServiceStep3;
      this.noPortForExtent = this.uiPreData.containerList.every(value => !value.isHavePort())
    });
    this.k8sService.getServiceConfig(this.stepPhase).subscribe(res => {
      this.uiBaseData = res;
      this.setServiceName(this.uiData.serviceName);
      this.changeDetectorRef.detectChanges();
    });
    this.k8sService.getNodeSelectors().subscribe((res:Array<string>)=>{
      this.nodeSelectorList = res;
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

  get isCanAddContainerInfo(){
    return this.uiPreData.containerList.find(value => value.isHavePort());
  }

  get nodeSelectorDefaultText(){
    return this.uiData.nodeSelector == "" ? 'SERVICE.STEP_3_NODE_SELECTOR_COMMENT': this.uiData.nodeSelector;
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

  setNodePort(index: number, port: number) {
    this.uiData.externalServiceList[index].node_config.node_port = Number(port).valueOf();
  }

  setServiceName(serviceName: string): void {
    this.uiData.serviceName = serviceName;
    /*Todo:add reset the Collaborative service Info*/
    this.collaborativeServiceList.splice(0, this.collaborativeServiceList.length);
    this.k8sService.getCollaborativeService(serviceName, this.uiData.projectName).subscribe(
      res => this.collaborativeServiceList = res,
      (err: HttpErrorResponse) => {
        if (err.status == 404) {
          this.messageService.cleanNotification();
        }
      });
  }

  addContainerInfo() {
    if (this.isCanAddContainerInfo){
      let externalService = new ExternalService();
      if (this.uiPreData.containerList.length > 0) {
        externalService.container_name = this.uiPreData.containerList[0].name;
        let containerPorts = this.getContainerPorts(externalService.container_name);
        if (containerPorts.length > 0) {
          externalService.node_config.target_port = containerPorts[0];
        }
      }
      this.uiData.externalServiceList.push(externalService);
    }
  }

  addOneCollaborativeService(){
    if (this.collaborativeServiceList.length > 0){
      this.collaborativeList.push({});
    }
  }

  removeContainerInfo(index: number) {
    this.uiData.externalServiceList.splice(index, 1);
  }

  getContainerDropdownText(index: number): string {
    let result = this.uiData.externalServiceList[index].container_name;
    return result == "" ? "SERVICE.STEP_3_SELECT_CONTAINER" : result;
  }

  getContainerPortDropdownText(index: number): string {
    let result = this.uiData.externalServiceList[index].node_config.target_port;
    return result == 0 ? "SERVICE.STEP_3_SELECT_PORT" : result.toString();
  }

  setExternalInfo(container: Container, index: number) {
    this.uiData.externalServiceList[index].container_name = container.name;
    let containerPorts = this.getContainerPorts(container.name);
    if (containerPorts.length > 0) {
      this.uiData.externalServiceList[index].node_config.target_port = containerPorts[0];
    }
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
    if (this.verifyInputValid()) {
      this.isActionWip = true;
      this.k8sService.setServiceConfig(this.uiData.uiToServer()).subscribe(
        () => this.k8sService.stepSource.next({index: 5, isBack: false})
      );
    }
  }

  backUpStep(): void {
    this.k8sService.stepSource.next({index: 2, isBack: true});
  }
}