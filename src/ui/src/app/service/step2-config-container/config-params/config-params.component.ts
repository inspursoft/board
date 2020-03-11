import { Component, ComponentFactoryResolver, OnInit, ViewContainerRef } from '@angular/core';
import { Observable } from 'rxjs';
import { AbstractControl, ValidationErrors } from '@angular/forms';
import { map } from 'rxjs/operators';
import { InputArrayExType } from 'board-components-library';
import { CsModalChildMessage } from '../../../shared/cs-modal-base/cs-modal-child-base';
import { MessageService } from '../../../shared.service/message.service';
import { Container, ContainerType, EnvStruct, UIServiceStep2, Volume } from '../../service-step.component';
import { VolumeMountsComponent } from '../volume-mounts/volume-mounts.component';
import { EnvType } from '../../../shared/environment-value/environment-value.component';
import { NodeAvailableResources } from '../../../shared/shared.types';
import { K8sService } from '../../service.k8s';

@Component({
  selector: 'app-config-params',
  templateUrl: './config-params.component.html',
  styleUrls: ['./config-params.component.css']
})
export class ConfigParamsComponent extends CsModalChildMessage implements OnInit {
  patternContainerName: RegExp = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/;
  patternWorkDir: RegExp = /^~?[\w\d-\/.{}$\/:]+[\s]*$/;
  patternCpuRequest: RegExp = /^[0-9]*m$/;
  patternCpuLimit: RegExp = /^[0-9]*m$/;
  patternMemRequest: RegExp = /^[0-9]*Mi$/;
  patternMemLimit: RegExp = /^[0-9]*Mi$/;
  container: Container;
  showEnvironmentValue = false;
  fixedContainerEnv: Map<string, Array<EnvStruct>>;
  fixedContainerPort: Map<Container, Array<number>>;
  step2Data: UIServiceStep2;
  curContainerType: ContainerType = ContainerType.runContainer;

  constructor(protected messageService: MessageService,
              private view: ViewContainerRef,
              private resolver: ComponentFactoryResolver,
              private k8sService: K8sService) {
    super(messageService);
  }

  ngOnInit() {
  }

  get checkSetCpuRequestFun() {
    return this.checkSetCpuRequest.bind(this);
  }

  get checkSetMemRequestFun() {
    return this.checkSetMemRequest.bind(this);
  }

  get validContainerNameFun() {
    return this.validContainerName.bind(this);
  }

  get validContainerPortsFun() {
    return this.validContainerPorts.bind(this);
  }

  get validContainerCpuFun() {
    return this.validContainerCpu.bind(this);
  }

  get validContainerMemFun() {
    return this.validContainerMem.bind(this);
  }

  get validContainerCpuLimitFun() {
    return this.validContainerCpuLimit.bind(this);
  }

  get validContainerMemLimitFun() {
    return this.validContainerMemLimit.bind(this);
  }

  getVolumesDescription(index: number, container: Container): string {
    const volume = container.volume_mounts;
    if (volume.length > index) {
      const storageServer = volume[index].targetStorageService === '' ? '' :
        volume[index].targetStorageService.concat(':');
      const result = `${volume[index].containerPath}:${storageServer}${volume[index].containerPath}`;
      return result === ':' ? '' : result;
    } else {
      return '';
    }
  }

  editVolumeMount() {
    const factory = this.resolver.resolveComponentFactory(VolumeMountsComponent);
    const componentRef = this.view.createComponent(factory);
    componentRef.instance.volumeDataList = this.container.volume_mounts;
    componentRef.instance.projectName = this.step2Data.projectName;
    componentRef.instance.onConfirmEvent.subscribe((res: Array<Volume>) => this.container.volume_mounts = res);
    componentRef.instance.openModal().subscribe(() => this.view.remove(this.view.indexOf(componentRef.hostView)));
  }

  getEnvsDescription(): string {
    let result = '';
    this.container.env.forEach((value: EnvStruct) => {
      result += `${value.dockerfile_envname}=${value.dockerfile_envvalue};`;
    });
    return result;
  }

  editEnvironment() {
    this.showEnvironmentValue = true;
  }

  getDefaultEnvsData() {
    const result = Array<EnvType>();
    this.container.env.forEach((value: EnvStruct) => {
      const env = new EnvType(value.dockerfile_envname, value.dockerfile_envvalue);
      env.envConfigMapKey = value.configmap_key;
      env.envConfigMapName = value.configmap_name;
      result.push(env);
    });
    return result;
  }

  getDefaultEnvsFixedData(): Array<string> {
    const result = Array<string>();
    if (this.fixedContainerEnv.has(this.container.image.image_name)) {
      const fixedEnvs: Array<EnvStruct> = this.fixedContainerEnv.get(this.container.image.image_name);
      fixedEnvs.forEach(value => result.push(value.dockerfile_envname));
    }
    return result;
  }

  setEnvironment(envsData: Array<EnvType>) {
    this.container.env.splice(0, this.container.env.length);
    envsData.forEach((value: EnvType) => {
      const env = new EnvStruct();
      env.dockerfile_envname = value.envName;
      env.dockerfile_envvalue = value.envValue;
      env.configmap_name = value.envConfigMapName;
      env.configmap_key = value.envConfigMapKey;
      this.container.env.push(env);
    });
  }

  setContainerPorts(event: Array<InputArrayExType>) {
    this.container.container_port.splice(0, this.container.container_port.length);
    event.forEach(value => {
      if (typeof value === 'number') {
        this.container.container_port.push(value);
      }
    });
  }

  checkSetCpuRequest(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.k8sService.getNodesAvailableSources().pipe(
      map((res: Array<NodeAvailableResources>) => {
        const isInValid = res.every(value =>
          Number.parseInt(control.value, 0) > Number.parseInt(value.cpu_available, 0) * 1000);
        if (isInValid) {
          return {beyondMaxLimit: 'SERVICE.STEP_2_BEYOND_MAX_VALUE'};
        } else {
          return null;
        }
      })
    );
  }

  checkSetMemRequest(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.k8sService.getNodesAvailableSources().pipe(
      map((res: Array<NodeAvailableResources>) => {
        const isInValid = res.every(value =>
          Number.parseInt(control.value, 0) > Number.parseInt(value.mem_available, 0) / (1024 * 1024));
        if (isInValid) {
          return {beyondMaxLimit: 'SERVICE.STEP_2_BEYOND_MAX_VALUE'};
        } else {
          return null;
        }
      })
    );
  }

  validContainerName(control: AbstractControl): ValidationErrors | null {
    const isValid = this.step2Data.containerList.find(container =>
      this.container !== container && container.name === control.value) === undefined;
    return isValid ? null : {containerNameRepeat: 'containerNameRepeat'};
  }

  validContainerPorts(control: AbstractControl): ValidationErrors | null {
    let isValid = true;
    const portBuf = new Set<number>();
    portBuf.add(Number.parseInt(control.value, 0));
    this.step2Data.containerList.forEach((container, index) => {
      container.container_port.forEach(port => {
        if (portBuf.has(port)) {
          isValid = false;
        } else {
          portBuf.add(port);
        }
      });
    });
    return isValid ? null : {containerPortRepeat: 'containerPortRepeat'};
  }

  validContainerCpu(control: AbstractControl): ValidationErrors | null {
    let isValid = true;
    if (control.value !== '' && this.container.cpu_limit !== '') {
      isValid = Number.parseFloat(control.value) <= Number.parseFloat(this.container.cpu_limit);
    }
    return isValid ? null : {resourceRequestInvalid: 'resourceRequestInvalid'};
  }

  validContainerMem(control: AbstractControl): ValidationErrors | null {
    let isValid = true;
    if (control.value !== '' && this.container.mem_limit !== '') {
      isValid = Number.parseFloat(control.value) <= Number.parseFloat(this.container.mem_limit);
    }
    return isValid ? null : {resourceRequestInvalid: 'resourceRequestInvalid'};
  }

  validContainerCpuLimit(control: AbstractControl): ValidationErrors | null {
    let isValid = true;
    if (control.value !== '' && this.container.cpu_request !== '') {
      isValid = Number.parseFloat(control.value) >= Number.parseFloat(this.container.cpu_request);
    }
    return isValid ? null : {resourceRequestInvalid: 'resourceRequestInvalid'};
  }

  validContainerMemLimit(control: AbstractControl): ValidationErrors | null {
    let isValid = true;
    if (control.value !== '' && this.container.mem_request !== '') {
      isValid = Number.parseFloat(control.value) >= Number.parseFloat(this.container.mem_request);
    }
    return isValid ? null : {resourceRequestInvalid: 'resourceRequestInvalid'};
  }

  setParams() {
    this.modalOpened = false;
  }
}
