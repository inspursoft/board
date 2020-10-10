import { AfterViewInit, Component, ComponentFactoryResolver, OnInit, ViewContainerRef } from '@angular/core';
import { Observable } from 'rxjs';
import { of } from 'rxjs/internal/observable/of';
import { AbstractControl, ValidationErrors } from '@angular/forms';
import { map } from 'rxjs/operators';
import { InputArrayExType } from 'board-components-library';
import { CsModalChildMessage } from '../../../shared/cs-modal-base/cs-modal-child-base';
import { MessageService } from '../../../shared.service/message.service';
import { Container, ContainerType, EnvStruct, ServiceStep2Data, Volume } from '../../service-step.component';
import { VolumeMountsComponent } from '../volume-mounts/volume-mounts.component';
import { SharedEnvType } from '../../../shared/shared.types';
import { K8sService } from '../../service.k8s';
import { NodeAvailableResources } from '../../service.types';

@Component({
  selector: 'app-config-params',
  templateUrl: './config-params.component.html',
  styleUrls: ['./config-params.component.css']
})
export class ConfigParamsComponent extends CsModalChildMessage implements OnInit, AfterViewInit {
  patternContainerName: RegExp = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/;
  patternWorkDir: RegExp = /^~?[\w\d-\/.{}$\/:]+[\s]*$/;
  patternCpuRequest: RegExp = /^[0-9]*m$/;
  patternCpuLimit: RegExp = /^[0-9]*m$/;
  patternMemRequest: RegExp = /^[0-9]*Mi$/;
  patternMemLimit: RegExp = /^[0-9]*Mi$/;
  container: Container;
  showEnvironmentValue = false;
  fixedContainerEnv: Map<Container, Array<EnvStruct>>;
  fixedContainerPort: Map<Container, Array<number>>;
  step2Data: ServiceStep2Data;
  curContainerType: ContainerType = ContainerType.runContainer;
  isAfterViewInit = false;

  constructor(protected messageService: MessageService,
              private view: ViewContainerRef,
              private resolver: ComponentFactoryResolver,
              private k8sService: K8sService) {
    super(messageService);
  }

  ngOnInit() {
  }

  ngAfterViewInit(): void {
    super.ngAfterViewInit();
    this.isAfterViewInit = true;
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

  get validContainerGpuLimitFun() {
    return this.validContainerGpuLimit.bind(this);
  }

  getVolumesDescription(index: number, container: Container): string {
    const volume = container.volumeMounts;
    if (volume.length > index) {
      const storageServer = volume[index].targetStorageService === '' ? '' :
        volume[index].targetStorageService.concat(':');
      const result = `${volume[index].containerPath}:${storageServer}${volume[index].targetPath}`;
      return result === ':' ? '' : result;
    } else {
      return '';
    }
  }

  editVolumeMount() {
    const factory = this.resolver.resolveComponentFactory(VolumeMountsComponent);
    const componentRef = this.view.createComponent(factory);
    componentRef.instance.volumeDataList = this.container.volumeMounts;
    componentRef.instance.projectName = this.step2Data.projectName;
    componentRef.instance.onConfirmEvent.subscribe((res: Array<Volume>) => this.container.volumeMounts = res);
    componentRef.instance.openModal().subscribe(() => this.view.remove(this.view.indexOf(componentRef.hostView)));
  }

  getEnvsDescription(): string {
    let result = '';
    this.container.env.forEach((value: EnvStruct) => {
      result += `${value.dockerFileEnvName}=${value.dockerFileEnvValue};`;
    });
    return result;
  }

  editEnvironment() {
    this.showEnvironmentValue = true;
  }

  get defaultContainerPorts(): Array<number> {
    if (this.fixedContainerPort.has(this.container)) {
      const fixedPorts = this.fixedContainerPort.get(this.container);
      return this.container.containerPort.filter(value => fixedPorts.indexOf(value) === -1);
    } else {
      return this.container.containerPort;
    }
  }

  getDefaultEnvsData() {
    const result = Array<SharedEnvType>();
    this.container.env.forEach((value: EnvStruct) => {
      const env = new SharedEnvType();
      env.envName = value.dockerFileEnvName;
      env.envValue = value.dockerFileEnvValue;
      env.envConfigMapKey = value.configMapKey;
      env.envConfigMapName = value.configMapName;
      result.push(env);
    });
    return result;
  }

  getDefaultEnvsFixedData(): Array<string> {
    const result = Array<string>();
    if (this.fixedContainerEnv.has(this.container)) {
      const fixedEnvs: Array<EnvStruct> = this.fixedContainerEnv.get(this.container);
      fixedEnvs.forEach(value => result.push(value.dockerFileEnvName));
    }
    return result;
  }

  setEnvironment(envsData: Array<SharedEnvType>) {
    this.container.env.splice(0, this.container.env.length);
    envsData.forEach((value: SharedEnvType) => {
      const env = new EnvStruct();
      env.dockerFileEnvName = value.envName;
      env.dockerFileEnvValue = value.envValue;
      env.configMapName = value.envConfigMapName;
      env.configMapKey = value.envConfigMapKey;
      this.container.env.push(env);
    });
  }

  setContainerPorts(event: Array<InputArrayExType>) {
    this.container.containerPort.splice(0, this.container.containerPort.length);
    event.forEach(value => {
      if (typeof value === 'number') {
        this.container.containerPort.push(value);
      }
    });
    if (this.fixedContainerPort.has(this.container)) {
      const fixedPorts = this.fixedContainerPort.get(this.container);
      fixedPorts.forEach(value => this.container.containerPort.push(value));
    }
  }

  checkSetCpuRequest(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.isAfterViewInit ? this.k8sService.getNodesAvailableSources().pipe(
      map((res: Array<NodeAvailableResources>) => {
        const isInValid = res.every(value =>
          Number.parseInt(control.value, 0) > Number.parseInt(value.cpuAvailable, 0) * 1000);
        if (isInValid) {
          return {beyondMaxLimit: 'SERVICE.STEP_2_BEYOND_MAX_VALUE'};
        } else {
          return null;
        }
      })
    ) : of(null);
  }

  checkSetMemRequest(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.isAfterViewInit ? this.k8sService.getNodesAvailableSources().pipe(
      map((res: Array<NodeAvailableResources>) => {
        const isInValid = res.every(value =>
          Number.parseInt(control.value, 0) > Number.parseInt(value.memAvailable, 0) / (1024 * 1024));
        if (isInValid) {
          return {beyondMaxLimit: 'SERVICE.STEP_2_BEYOND_MAX_VALUE'};
        } else {
          return null;
        }
      })
    ) : of(null);
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
    this.step2Data.containerList.forEach((container) => {
      container.containerPort.forEach(port => {
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
    if (control.value !== '' && this.container.cpuLimit !== '') {
      isValid = Number.parseFloat(control.value) <= Number.parseFloat(this.container.cpuLimit);
    }
    return isValid ? null : {resourceRequestInvalid: 'resourceRequestInvalid'};
  }

  validContainerMem(control: AbstractControl): ValidationErrors | null {
    let isValid = true;
    if (control.value !== '' && this.container.memLimit !== '') {
      isValid = Number.parseFloat(control.value) <= Number.parseFloat(this.container.memLimit);
    }
    return isValid ? null : {resourceRequestInvalid: 'resourceRequestInvalid'};
  }

  validContainerCpuLimit(control: AbstractControl): ValidationErrors | null {
    let isValid = true;
    if (control.value !== '' && this.container.cpuRequest !== '') {
      isValid = Number.parseFloat(control.value) >= Number.parseFloat(this.container.cpuRequest);
    }
    return isValid ? null : {resourceRequestInvalid: 'resourceRequestInvalid'};
  }

  validContainerMemLimit(control: AbstractControl): ValidationErrors | null {
    let isValid = true;
    if (control.value !== '' && this.container.memRequest !== '') {
      isValid = Number.parseFloat(control.value) >= Number.parseFloat(this.container.memRequest);
    }
    return isValid ? null : {resourceRequestInvalid: 'resourceRequestInvalid'};
  }

  validContainerGpuLimit(control: AbstractControl): ValidationErrors | null {
    return Number(control.value) >= 0 ? null : {resourceRequestInvalid: 'resourceRequestInvalid'};
  }
}
