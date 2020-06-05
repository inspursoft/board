import { Component, ComponentFactoryResolver, Input, OnInit, ViewContainerRef } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { ValidationErrors } from '@angular/forms';
import { map } from 'rxjs/operators';
import { CsModalChildBase } from '../../shared/cs-modal-base/cs-modal-child-base';
import { JobContainer, JobEnv, JobVolumeMounts } from '../job.type';
import { JobVolumeMountsComponent } from '../job-volume-mounts/job-volume-mounts.component';
import { EnvType } from '../../shared/environment-value/environment-value.component';
import { NodeAvailableResources } from '../../shared/shared.types';
import { JobService } from '../job.service';
import { MessageService } from '../../shared.service/message.service';

@Component({
  selector: 'app-job-container-config',
  templateUrl: './job-container-config.component.html',
  styleUrls: ['./job-container-config.component.css']
})
export class JobContainerConfigComponent extends CsModalChildBase implements OnInit {
  @Input() container: JobContainer;
  @Input() containerList: Array<JobContainer>;
  @Input() projectName: string;
  @Input() projectId: number;
  @Input() isEditModel = false;
  patternContainerName: RegExp = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?$/;
  patternWorkdir: RegExp = /^~?[\w\d-\/.{}$\/:]+[\s]*$/;
  patternCpuRequest: RegExp = /^[0-9]*m$/;
  patternCpuLimit: RegExp = /^[0-9]*m$/;
  patternMemRequest: RegExp = /^[0-9]*Mi$/;
  patternMemLimit: RegExp = /^[0-9]*Mi$/;
  isLoading = false;
  volumesDescriptions: Array<string>;
  showEnvironmentValue = false;
  createSuccess: Subject<JobContainer>;

  constructor(private factoryResolver: ComponentFactoryResolver,
              private view: ViewContainerRef,
              private jobService: JobService,
              private messageService: MessageService) {
    super();
    this.volumesDescriptions = Array<string>();
    this.createSuccess = new Subject();
  }

  ngOnInit() {
    this.generateDescriptions();
  }

  createNewContainer() {
    if (this.isExistsContainerNames()) {
      this.messageService.showAlert('JOB.JOB_CREATE_CONTAINER_NAME_REPEAT', {alertType: 'warning', view: this.alertView});
    } else if (this.isExistsContainerPorts()) {
      this.messageService.showAlert('JOB.JOB_CREATE_CONTAINER_PORT_REPEAT', {alertType: 'warning', view: this.alertView});
    } else if (this.isInvalidContainerCpuAndMem()) {
      this.messageService.showAlert('JOB.JOB_CREATE_CONTAINER_REQUEST_ERROR', {alertType: 'warning', view: this.alertView});
    } else if (this.verifyInputExValid()) {
      this.createSuccess.next(this.container);
      this.modalOpened = false;
    }
  }

  isInvalidContainerCpuAndMem(): boolean {
    let cpuValid = true;
    let memValid = true;
    if (this.container.cpuRequest && this.container.cpuLimit) {
      cpuValid = Number.parseFloat(this.container.cpuRequest) < Number.parseFloat(this.container.cpuLimit);
    }
    if (this.container.memRequest && this.container.memLimit) {
      memValid = Number.parseFloat(this.container.memRequest) < Number.parseFloat(this.container.memLimit);
    }
    return !(cpuValid && memValid);
  }

  isExistsContainerNames(): boolean {
    const findRepeat = this.containerList.find((findValue: JobContainer) => findValue.name === this.container.name);
    return findRepeat !== undefined && !this.isEditModel;
  }

  isExistsContainerPorts(): boolean {
    let isExists = false;
    const portBuf = new Set<number>();
    this.containerList.forEach((container) => {
      if (container.name !== this.container.name) {
        container.containerPort.forEach(port => portBuf.add(port));
      }
    });
    this.container.containerPort.forEach((port: number) => {
      if (!isExists) {
        isExists = portBuf.has(port);
      }
    });
    return isExists;
  }

  getEnvsDescription(): string {
    let result = '';
    this.container.env.forEach((value: JobEnv) => {
      result += `${value.dockerfileEnvName}=${value.dockerfileEnvValue};`;
    });
    return result;
  }

  getDefaultEnvsData(): Array<EnvType> {
    const result = Array<EnvType>();
    this.container.env.forEach((value: JobEnv) => {
      const env = new EnvType(value.dockerfileEnvName, value.dockerfileEnvValue);
      env.envConfigMapKey = value.configMapKey;
      env.envConfigMapName = value.configMapName;
      result.push(env);
    });
    return result;
  }

  get checkSetCpuRequestFun() {
    return this.checkSetCpuRequest.bind(this);
  }

  get checkSetMemRequestFun() {
    return this.checkSetMemRequest.bind(this);
  }

  checkSetCpuRequest(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.jobService.getNodesAvailableSources()
      .pipe(map((res: Array<NodeAvailableResources>) => {
        const isInValid = res.every(value => Number.parseInt(control.value, 10) > Number.parseInt(value.cpu_available, 10) * 1000);
        if (isInValid) {
          return {beyondMaxLimit: 'JOB.JOB_CREATE_BEYOND_MAX_VALUE'};
        } else {
          return null;
        }
      }));
  }

  checkSetMemRequest(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.jobService.getNodesAvailableSources()
      .pipe(map((res: Array<NodeAvailableResources>) => {
        const isInValid = res.every(value => Number.parseInt(control.value, 10) > Number.parseInt(value.mem_available, 10) / (1024 * 1024));
        if (isInValid) {
          return {beyondMaxLimit: 'JOB.JOB_CREATE_BEYOND_MAX_VALUE'};
        } else {
          return null;
        }
      }));
  }

  setEnvironment(envsData: Array<EnvType>) {
    const envsArray = this.container.env;
    envsArray.splice(0, envsArray.length);
    envsData.forEach((value: EnvType) => {
      const env = new JobEnv();
      env.dockerfileEnvName = value.envName;
      env.dockerfileEnvValue = value.envValue;
      env.configMapName = value.envConfigMapName;
      env.configMapKey = value.envConfigMapKey;
      envsArray.push(env);
    });
  }

  generateDescriptions() {
    this.volumesDescriptions.splice(0, this.volumesDescriptions.length);
    this.container.volumeMounts.forEach((volume: JobVolumeMounts) => {
      if (volume.volumeType === 'nfs') {
        const des = `NFS[${volume.volumeName}]`;
        this.volumesDescriptions.push(des);
      } else if (volume.volumeType === 'pvc') {
        const des = `PVC[${volume.volumeName}]`;
        this.volumesDescriptions.push(des);
      }
    });
  }

  editVolumeMount() {
    const factory = this.factoryResolver.resolveComponentFactory(JobVolumeMountsComponent);
    const componentRef = this.view.createComponent(factory);
    componentRef.instance.volumeDataList = this.container.volumeMounts;
    componentRef.instance.confirmEvent.subscribe((res: Array<JobVolumeMounts>) => {
      this.container.volumeMounts = res;
      this.generateDescriptions();
    });
    componentRef.instance.openModal().subscribe(() => this.view.remove(this.view.indexOf(componentRef.hostView)));
  }

  setContainerPort(ports: any) {
    this.container.containerPort = ports;
  }
}
