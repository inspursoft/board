import { Component, Input, OnInit, OnDestroy, AfterContentChecked, QueryList, ViewChildren } from '@angular/core';
import {
  ImageDockerfile,
  ServiceStep2Output, ServiceStep2Type, ServiceStep3Output, ServiceStep3Type,
  ServiceStepComponent
} from '../service-step.component';
import { K8sService } from '../service.k8s';
import { EnvType } from "../environment-value/environment-value.component";
import { CsInputComponent } from "../cs-input/cs-input.component";
import { CsInputArrayComponent } from "../cs-input-array/cs-input-array.component";

@Component({
  templateUrl: './edit-container.component.html',
  styleUrls: ["./edit-container.component.css"]
})
export class EditContainerComponent implements ServiceStepComponent, OnInit, OnDestroy, AfterContentChecked {
  @Input() data: any;
  @ViewChildren(CsInputComponent) inputComponents: QueryList<CsInputComponent>;
  @ViewChildren(CsInputArrayComponent) inputArrayComponents: QueryList<CsInputArrayComponent>;
  patternContainerName: RegExp = /^[a-zA-Z\d_-]+$/;
  patternWorkdir: RegExp = /^~?[\w\d-\/.{}$\/:]+[\s]*$/;
  step2Output: ServiceStep2Output;
  step3Output: ServiceStep3Output;
  step3TypeStatus: Map<ServiceStep3Type, boolean>;
  showVolumeMounts = false;
  showEnvironmentValue = false;
  isInputComponentsValid = false;
  fixedEnvKeys: Array<string>;
  fixedContainerPort: Array<number>;
  curContainerIndex:number;

  constructor(private k8sService: K8sService) {
    this.step3Output = Array<ServiceStep3Type>();
    this.step3TypeStatus = new Map<ServiceStep3Type, boolean>();
    this.fixedEnvKeys = Array<string>();
    this.fixedContainerPort = Array<number>();
  }

  ngOnInit() {
    this.step2Output = this.k8sService.getStepData(2) as ServiceStep2Output;
    if (this.k8sService.getStepData(3)) {
      this.step3Output = this.k8sService.getStepData(3) as ServiceStep3Output;
      this.step3Output.forEach((value: ServiceStep3Type) => {
        this.step3TypeStatus.set(value, false);
        this.setDefaultContainerInfo(value, false);
      });
    } else {
      this.step2Output.forEach((value: ServiceStep2Type) => {
        let config = new ServiceStep3Type();
        let firstIndex = value.image_name.indexOf("/");
        config.container_name = value.image_name.slice(firstIndex + 1, value.image_name.length);
        config.container_baseimage = value.image_name + ":" + value.image_tag;
        this.step3TypeStatus.set(config, false);
        this.setDefaultContainerInfo(config, true);
        this.step3Output.push(config);
      });
    }
  }

  ngOnDestroy() {
    this.k8sService.setStepData(3, this.step3Output);
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
    if (this.inputArrayComponents) {
      this.inputArrayComponents.forEach(item => {
        if (!item.valid) {
          this.isInputComponentsValid = false;
        }
      });
    }
  }

  setDefaultContainerInfo(config: ServiceStep3Type, isNew: boolean): void {
    let imageName: string = config.container_baseimage.split(":")[0];
    let imageTag: string = config.container_baseimage.split(":")[1];
    let firstIndex = imageName.indexOf("/");
    let projectName = imageName.slice(0, firstIndex);
    if (this.step2Output[0].project_name == projectName) {
      this.k8sService.getContainerDefaultInfo(imageName, imageTag, projectName)
        .then((res: ImageDockerfile) => {
          this.step3TypeStatus.set(config, true);
          if (res.image_cmd && isNew) {
            config.container_command.push(res.image_cmd);//copy cmd
          }
          if (res.image_env) {
            res.image_env.forEach(value => {//copy env
              if (isNew) {
                config.container_envs.push({
                  env_name: value.dockerfile_envname,
                  env_value: value.dockerfile_envvalue
                });
              }
              this.fixedEnvKeys.push(value.dockerfile_envname);
            });
          }
          if (res.image_expose) {
            res.image_expose.forEach(value => {//copy port
              let port: number = Number(value).valueOf();
              this.fixedContainerPort.push(port);
              if (isNew) {
                config.container_ports.push(port);
              }
            });
          }
        }).catch(() => {
      });
    }
  }

  get isCanNextStep(): boolean {
    return this.step3Output.length > 0 && this.isInputComponentsValid;
  }

  getVolumesDescription(item: ServiceStep3Type): string {
    let volumesArr = item.container_volumes;
    let result: string = "";
    volumesArr.forEach(value => {
      let storageServer = value.target_storageServer == "" ? "" : value.target_storageServer.concat(":");
      result += `${value.container_dir}:${storageServer}${value.target_dir}`
    });
    return result == ":" ? "" : result;
  }

  getEnvsDescription(item: ServiceStep3Type): string {
    let envsArr = item.container_envs;
    let result: string = "";
    envsArr.forEach(value => {
      result += `${value.env_name}=${value.env_value};`
    });
    return result;
  }

  getDefaultEnvsData(index: number) {
    let result = Array<EnvType>();
    this.step3Output[index].container_envs.forEach(value => {
      result.push(new EnvType(value.env_name, value.env_value))
    });
    return result;
  }

  setEnvironment(index: number, envsData: Array<EnvType>) {
    let envsArray = this.step3Output[index].container_envs;
    envsArray.splice(0, envsArray.length);
    envsData.forEach((value: EnvType) => {
      envsArray.push({env_name: value.envName, env_value: value.envValue})
    });
  }

  setVolumeMount(data: Object, index: number) {
    let volumeArr = this.step3Output[index].container_volumes;
    if (volumeArr.length == 0) {
      volumeArr.push(data as any)
    } else {
      volumeArr[0] = data as any;
    }
  }

  getVolumeMountData(index: number) {
    let volumeArr = this.step3Output[index].container_volumes;
    if (volumeArr.length == 0) {
      return {
        container_dir: "",
        target_storagename: "",
        target_storageServer: "",
        target_dir: ""
      };
    } else {
      return volumeArr[0];
    }
  }

  toggleShowStatus(item: ServiceStep3Type):void {
    let status = this.step3TypeStatus.get(item);
    this.step3TypeStatus.set(item, !status);
  }

  shieldEnter($event: KeyboardEvent) {
    if ($event.charCode == 13) {
      (<any>$event.target).blur();
    }
  }

  backStep(): void {
    this.k8sService.stepSource.next(2);
  }

  forward(): void {
    this.k8sService.stepSource.next(4);
  }
}