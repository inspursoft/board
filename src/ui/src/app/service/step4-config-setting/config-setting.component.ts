import { Component, Input, OnInit, OnDestroy } from '@angular/core';
import {
  ServiceStep1Output,
  ServiceStep2Output, ServiceStep3Output, ServiceStep3Type, ServiceStep4Output,
  ServiceStepComponent
} from '../service-step.component';
import { K8sService } from '../service.k8s';


@Component({
  styleUrls: ["./config-setting.component.css"],
  templateUrl: './config-setting.component.html'
})
export class ConfigSettingComponent implements ServiceStepComponent, OnInit, OnDestroy {
  @Input() data: any;
  dropDownListNum: Array<number>;
  selectContainerPorts: Map<string, Array<number>>;
  step2Output: ServiceStep2Output;
  step3Output: ServiceStep3Output;
  step4Output: ServiceStep4Output;
  showAdvanced: boolean = false;
  showExternal: boolean = false;
  showCollaborative: boolean = false;

  constructor(private k8sService: K8sService) {
    this.dropDownListNum = Array<number>();
    this.selectContainerPorts = new Map<string, Array<number>>();
    this.step4Output = new ServiceStep4Output();
  }

  ngOnDestroy() {
    this.k8sService.setStepData(4, this.step4Output);
  }

  ngOnInit() {
    let step1Out = this.k8sService.getStepData(1) as ServiceStep1Output;
    this.step2Output = this.k8sService.getStepData(2) as ServiceStep2Output;
    this.step3Output = this.k8sService.getStepData(3) as ServiceStep3Output;
    this.step4Output.service_id = step1Out.service_id;
    this.step4Output.project_name = step1Out.project_name;
    this.step4Output.project_id = step1Out.project_id;
    this.step4Output.deployment_yaml.container_list = this.step3Output;
    let volumeList = this.step4Output.deployment_yaml.volume_list;
    this.step3Output.forEach((value: ServiceStep3Type) => {
      value.container_volumes.forEach(volume => {
        volumeList.push({
          volume_name: volume.target_storagename,
          volume_path: volume.target_dir,
          server_name: ""
        })
      });
    });
    for (let i = 1; i <= 100; i++) {
      this.dropDownListNum.push(i)
    }
  }

  get serviceExternalArray() {
    return this.step4Output.service_yaml.service_external;
  }

  get serviceSelectorsArray() {
    //get selectors list api ...
    return this.step4Output.service_yaml.service_selectors;
  }

  setServiceName(serviceName: string) {
    this.step4Output.deployment_yaml.deployment_name = serviceName;
    this.step4Output.service_yaml.service_name = serviceName;
    let selectors = this.step4Output.service_yaml.service_selectors;
    if (selectors.length == 0) {
      selectors.push(serviceName);
    } else {
      selectors[0] = serviceName;
    }
  }

  addContainerInfo() {
    let serviceExternal = this.serviceExternalArray;
    serviceExternal.push({
      service_containername: "",
      service_externalpath: "",
      service_containerport: 0,
      service_nodeport: 0
    })
  }

  removeContainerInfo(index: number) {
    let serviceExternal = this.serviceExternalArray;
    serviceExternal.splice(index, 1);
  }

  getContainerDropdownText(index: number): string {
    let serviceExternal = this.serviceExternalArray;
    let result = serviceExternal[index].service_containername;
    return result == "" ? "SERVICE.STEP_4_SELECT_CONTAINER" : result;
  }

  getContainerPortDropdownText(index: number): string {
    let serviceExternal = this.step4Output.service_yaml.service_external;
    let result = serviceExternal[index].service_containerport;
    return result == 0 ? "SERVICE.STEP_4_SELECT_PORT" : result.toString();
  }

  setExternalInfo(item: ServiceStep3Type, index: number) {
    let serviceExternal = this.serviceExternalArray;
    serviceExternal[index].service_containername = item.container_name;
    this.selectContainerPorts.set(item.container_name, item.container_ports);
  }

  getContainerPorts(containerName: string): Array<number> {
    return this.selectContainerPorts.get(containerName);
  }

  onRadNodeChange(index: number) {
    let serviceExternal = this.serviceExternalArray;
    serviceExternal[index].service_externalpath = "";
  }

  onRadExternalpathChange(index: number) {
    let serviceExternal = this.serviceExternalArray;
    serviceExternal[index].service_nodeport = 0;
  }

  setServiceNodeport(port: string, index: number) {
    this.serviceExternalArray[index].service_nodeport = Number(port).valueOf();
  }


  forward(): void {
    this.k8sService.stepSource.next(5);
  }
}