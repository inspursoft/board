import { Component, ComponentFactoryResolver, Input, OnInit, Type, ViewContainerRef } from '@angular/core';
import { CsModalChildBase } from "../../shared/cs-modal-base/cs-modal-child-base";
import { JobContainer, JobEnv, JobVolumeMounts } from "../job.type";
import { JobVolumeMountsComponent } from "../job-volume-mounts/job-volume-mounts.component";
import { EnvType } from "../../shared/environment-value/environment-value.component";
import { Observable, Subject } from "rxjs";
import { ValidationErrors } from "@angular/forms";
import { map } from "rxjs/operators";
import { NodeAvailableResources } from "../../shared/shared.types";
import { JobService } from "../job.service";
import { MessageService } from "../../shared.service/message.service";

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
  }

  createNewContainer() {
    if (this.isExistsContainerNames()) {
      this.messageService.showAlert('JOB.JOB_CREATE_CONTAINER_NAME_REPEAT', {alertType: "warning", view: this.alertView});
    } else if (this.isExistsContainerPorts()) {
      this.messageService.showAlert('JOB.JOB_CREATE_CONTAINER_PORT_REPEAT', {alertType: "warning", view: this.alertView});
    } else if (this.isInvalidContainerCpuAndMem()) {
      this.messageService.showAlert('JOB.JOB_CREATE_CONTAINER_REQUEST_ERROR', {alertType: "warning", view: this.alertView});
    } else if (this.verifyInputValid() && this.verifyInputArrayValid()) {
      this.createSuccess.next(this.container);
      this.modalOpened = false;
    }
  }

  isInvalidContainerCpuAndMem(): boolean {
    let cpuValid = true;
    let memValid = true;
    if (this.container.cpu_request && this.container.cpu_limit) {
      cpuValid = Number.parseFloat(this.container.cpu_request) < Number.parseFloat(this.container.cpu_limit);
    }
    if (this.container.mem_request && this.container.mem_limit) {
      memValid = Number.parseFloat(this.container.mem_request) < Number.parseFloat(this.container.mem_limit)
    }
    return !(cpuValid && memValid);
  }

  isExistsContainerNames(): boolean {
    let findRepeat = this.containerList.find((findValue: JobContainer) => findValue.name === this.container.name);
    return findRepeat !== undefined && !this.isEditModel;
  }

  isExistsContainerPorts(): boolean {
    let isExists = false;
    let portBuf = new Set<number>();
    this.containerList.forEach((container) => {
      if (container.name !== this.container.name){
        container.container_port.forEach(port => portBuf.add(port))
      }
    });
    this.container.container_port.forEach((port: number) => {
      if (!isExists) {
        isExists = portBuf.has(port);
      }
    });
    return isExists;
  }

  getEnvsDescription(): string {
    let result: string = "";
    this.container.env.forEach((value: JobEnv) => {
      result += `${value.dockerfile_envname}=${value.dockerfile_envvalue};`
    });
    return result;
  }

  getDefaultEnvsData(): Array<EnvType> {
    let result = Array<EnvType>();
    this.container.env.forEach((value: JobEnv) => {
      let env = new EnvType(value.dockerfile_envname, value.dockerfile_envvalue);
      env.envConfigMapKey = value.configmap_key;
      env.envConfigMapName = value.configmap_name;
      result.push(env)
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
        let isInValid = res.every(value => Number.parseInt(control.value) > Number.parseInt(value.cpu_available) * 1000);
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
        let isInValid = res.every(value => Number.parseInt(control.value) > Number.parseInt(value.mem_available) / (1024 * 1024));
        if (isInValid) {
          return {beyondMaxLimit: 'JOB.JOB_CREATE_BEYOND_MAX_VALUE'};
        } else {
          return null;
        }
      }));
  }

  setEnvironment(envsData: Array<EnvType>) {
    let envsArray = this.container.env;
    envsArray.splice(0, envsArray.length);
    envsData.forEach((value: EnvType) => {
      let env = new JobEnv();
      env.dockerfile_envname = value.envName;
      env.dockerfile_envvalue = value.envValue;
      env.configmap_name = value.envConfigMapName;
      env.configmap_key = value.envConfigMapKey;
      envsArray.push(env);
    });
  }

  editVolumeMount() {
    let factory = this.factoryResolver.resolveComponentFactory(JobVolumeMountsComponent);
    let componentRef = this.view.createComponent(factory);
    componentRef.instance.volumeDataList = this.container.volume_mounts;
    componentRef.instance.onConfirmEvent.subscribe((res: Array<JobVolumeMounts>) => {
      this.container.volume_mounts = res;
      this.volumesDescriptions.splice(0, this.volumesDescriptions.length);
      res.forEach((volume: JobVolumeMounts) => {
        if (volume.volume_type === 'nfs') {
          const des = `NFS[${volume.volume_name}:${volume.target_storage_service}${volume.container_path}]`;
          this.volumesDescriptions.push(des);
        } else if (volume.volume_type === 'pvc') {
          const des = `PVC[${volume.volume_name}:${volume.target_pvc}:${volume.container_path}]`;
          this.volumesDescriptions.push(des);
        }
      });
    });
    componentRef.instance.openModal().subscribe(() => this.view.remove(this.view.indexOf(componentRef.hostView)));
  }

}
