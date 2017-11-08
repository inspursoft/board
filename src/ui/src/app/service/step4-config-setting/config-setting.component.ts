import { Component, Input, OnInit, OnDestroy, AfterContentChecked, ViewChildren, QueryList } from '@angular/core';
import {
  ServiceStep1Output,
  ServiceStep2Output, ServiceStep3Output, Container, ServiceStep4Output,
  ServiceStepComponent, ContainerPort, ServicePort
} from '../service-step.component';
import { K8sService } from '../service.k8s';
import { CsInputComponent } from "../cs-input/cs-input.component";


@Component({
  styleUrls: ["./config-setting.component.css"],
  templateUrl: './config-setting.component.html'
})
export class ConfigSettingComponent implements ServiceStepComponent, OnInit, OnDestroy, AfterContentChecked {
  @Input() data: any;
  @ViewChildren(CsInputComponent) inputComponents: QueryList<CsInputComponent>;
  patternServiceName: RegExp = /^[a-z]([-a-z0-9]*[a-z0-9])+$/;
  dropDownListNum: Array<number>;
  step2Output: ServiceStep2Output;
  step3Output: ServiceStep3Output;
  step4Output: ServiceStep4Output;
  showAdvanced: boolean = true;
  showExternal: boolean = false;
  showCollaborative: boolean = false;
  isInputComponentsValid = false;

  constructor(private k8sService: K8sService) {
    this.dropDownListNum = Array<number>();
    this.step4Output = new ServiceStep4Output();
  }

  ngOnDestroy() {
    this.k8sService.setStepData(4, this.step4Output);
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
    let step1Out = this.k8sService.getStepData(1) as ServiceStep1Output;
    this.step2Output = this.k8sService.getStepData(2) as ServiceStep2Output;
    this.step3Output = this.k8sService.getStepData(3) as ServiceStep3Output;
    if (this.k8sService.getStepData(4)) {
      this.step4Output = this.k8sService.getStepData(4) as ServiceStep4Output;
    } else {
      this.step4Output.deployment_yaml.spec.template.spec.containers = this.step3Output;
      this.step4Output.projectinfo.service_id = step1Out.service_id;
      this.step4Output.projectinfo.project_id = step1Out.project_id;
      this.step4Output.projectinfo.project_name = step1Out.project_name;
      let volumeList = this.step4Output.deployment_yaml.spec.template.spec.volumes;
      this.step3Output.forEach((container: Container) => {
        container.volumeMounts.forEach(volume => {
          volumeList.push({
            name: volume.name,
            nfs: {server: volume.ui_nfs_server, path: volume.ui_nfs_path}
          });
        });
      });
    }
    for (let i = 1; i <= 100; i++) {
      this.dropDownListNum.push(i)
    }
  }

  get servicePortsArray(): Array<ServicePort> {
    return this.step4Output.service_yaml.spec.ports;
  }

  get serviceSelectorsArray(): {[key: string]: string} {
    //get selectors list api ...
    return this.step4Output.service_yaml.spec.selector;
  }

  setServiceName(serviceName: string) {
    this.step4Output.projectinfo.service_name = serviceName;
    this.step4Output.service_yaml.metadata.name = serviceName;
    this.step4Output.service_yaml.metadata.labels = {"app": serviceName};
    this.step4Output.service_yaml.spec.selector = {"app": serviceName};
    this.step4Output.deployment_yaml.metadata.name = serviceName;
    this.step4Output.deployment_yaml.spec.template.metadata.labels = {"app": serviceName};
  }

  addContainerInfo() {
    this.servicePortsArray.push({name: "", port: 0, nodePort: 0});
    this.step4Output.service_yaml.spec.type = "NodePort";
    this.step4Output.projectinfo.service_externalpath.push("");
  }

  removeContainerInfo(index: number) {
    this.servicePortsArray.splice(index, 1);
    this.step4Output.projectinfo.service_externalpath.splice(index, 1);
  }

  getContainerDropdownText(index: number): string {
    let result = this.servicePortsArray[index].name;
    return result == "" ? "SERVICE.STEP_4_SELECT_CONTAINER" : result;
  }

  getContainerPortDropdownText(index: number): string {
    let result = this.servicePortsArray[index].port;
    return result == 0 ? "SERVICE.STEP_4_SELECT_PORT" : result.toString();
  }

  setExternalInfo(container: Container, index: number) {
    this.servicePortsArray[index].name = container.name;
  }

  getContainerPorts(containerName: string):Array<ContainerPort> {
    let result:Array<ContainerPort> ;
    this.step3Output.forEach(container=>{
      if (container.name == containerName){
        result = container.ports;
      }
    });
    return result;
  }

  onRadNodeChange(index: number) {
    this.step4Output.projectinfo.service_externalpath[index] = "";
  }

  onRadExternalpathChange(index: number) {
    this.servicePortsArray[index].nodePort = 0;
  }

  setServiceNodeport(port: string, index: number) {
    this.servicePortsArray[index].nodePort = Number(port).valueOf();
  }

  getExternalPath(index: number): string {
    return this.step4Output.projectinfo.service_externalpath[index];
  }

  setExternalPath(index: number, value: string): void {
    this.step4Output.projectinfo.service_externalpath[index] = value;
  }

  forward(): void {
    this.k8sService.stepSource.next(6);
  }

  backUpStep(): void {
    this.k8sService.stepSource.next(3);
  }
}