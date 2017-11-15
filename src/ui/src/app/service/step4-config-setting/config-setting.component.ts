import { Component, OnInit, AfterContentChecked, ViewChildren, QueryList, Injector } from '@angular/core';
import { Container, ContainerPort, ServicePort } from '../service-step.component';
import { CsInputComponent } from "../cs-input/cs-input.component";
import { ServiceStepBase } from "../service-step";

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

  constructor(protected injector: Injector) {
    super(injector);
    this.dropDownListNum = Array<number>();
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
    this.k8sService.getServiceConfig(this.newServiceId, this.outputData).then(res => {
      this.outputData = res;
      this.containerList.forEach((container: Container) => {
        container.volumeMounts.forEach(volume => {
          this.deployVolumes.push({
            name: volume.name,
            nfs: {server: volume.ui_nfs_server, path: volume.ui_nfs_path}
          });
        });
      });
    });
    for (let i = 1; i <= 100; i++) {
      this.dropDownListNum.push(i)
    }
  }

  get servicePortsArray(): Array<ServicePort> {
    return this.outputData.service_yaml.spec.ports;
  }

  get serviceSelectorsArray(): {[key: string]: string} {
    //get selectors list api ...
    return this.outputData.service_yaml.spec.selector;
  }

  setServiceName(serviceName: string) {
    this.outputData.service_yaml.metadata.name = serviceName;
    this.outputData.service_yaml.metadata.labels = {"app": serviceName};
    this.outputData.service_yaml.spec.selector = {"app": serviceName};
    this.outputData.deployment_yaml.metadata.name = serviceName;
    this.outputData.deployment_yaml.spec.template.metadata.labels = {"app": serviceName};
  }

  addContainerInfo() {
    this.servicePortsArray.push({name: "", port: 0, nodePort: 0});
    this.outputData.service_yaml.spec.type = "NodePort";
    this.outputData.projectinfo.service_externalpath.push("");
  }

  removeContainerInfo(index: number) {
    this.servicePortsArray.splice(index, 1);
    this.outputData.projectinfo.service_externalpath.splice(index, 1);
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

  getContainerPorts(containerName: string): Array<ContainerPort> {
    let result: Array<ContainerPort>;
    this.containerList.forEach(container => {
      if (container.name == containerName) {
        result = container.ports;
      }
    });
    return result;
  }

  onRadNodeChange(index: number) {
    this.outputData.projectinfo.service_externalpath[index] = "";
  }

  onRadExternalpathChange(index: number) {
    this.servicePortsArray[index].nodePort = 0;
  }

  setServiceNodeport(port: string, index: number) {
    this.servicePortsArray[index].nodePort = Number(port).valueOf();
  }

  getExternalPath(index: number): string {
    return this.outputData.projectinfo.service_externalpath[index];
  }

  setExternalPath(index: number, value: string): void {
    this.outputData.projectinfo.service_externalpath[index] = value;
  }

  forward(): void {
    this.k8sService.setServiceConfig(this.outputData).then(() => {
      this.k8sService.stepSource.next({index: 6, isBack: false});
    });
  }

  backUpStep(): void {
    this.k8sService.setServiceConfig(this.outputData).then(() => {
      this.k8sService.stepSource.next({index: 3, isBack: true});
    });
  }
}