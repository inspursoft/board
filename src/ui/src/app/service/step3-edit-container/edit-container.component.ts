import { Component, Input, OnInit, OnDestroy, AfterContentChecked, QueryList, ViewChildren } from '@angular/core';
import {
  ImageDockerfile,
  ServiceStep2Output, ServiceStep2Type, ServiceStep3Output, Container, VolumeMount,
  ServiceStepComponent
} from '../service-step.component';
import { K8sService } from '../service.k8s';
import { EnvType } from "../environment-value/environment-value.component";
import { CsInputComponent } from "../cs-input/cs-input.component";
import { CsInputArrayComponent } from "../cs-input-array/cs-input-array.component";
import { VolumeOutPut } from "./volume-mounts/volume-mounts.component";

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
  step3TypeStatus: Map<Container, boolean>;
  showVolumeMounts = false;
  showEnvironmentValue = false;
  isInputComponentsValid = false;
  fixedEnvKeys: Array<string>;
  containerPort: Map<string, Array<number>>;
  fixedContainerPort: Map<string, Array<number>>;
  curContainerIndex: number;

  constructor(private k8sService: K8sService) {
    this.step3Output = Array<Container>();
    this.step3TypeStatus = new Map<Container, boolean>();
    this.containerPort = new Map<string, Array<number>>();
    this.fixedEnvKeys = Array<string>();
    this.fixedContainerPort = new Map<string, Array<number>>();
  }

  ngOnInit() {
    this.step2Output = this.k8sService.getStepData(2) as ServiceStep2Output;
    if (this.k8sService.getStepData(3)) {
      this.step3Output = this.k8sService.getStepData(3) as ServiceStep3Output;
      this.step3Output.forEach((value: Container) => {
        this.step3TypeStatus.set(value, false);
        this.setDefaultContainerInfo(value, false);
      });
    } else {
      this.step2Output.forEach((value: ServiceStep2Type) => {
        let container = new Container();
        let firstIndex = value.image_name.indexOf("/");
        container.name = value.image_name.slice(firstIndex + 1, value.image_name.length);
        container.image = value.image_name + ":" + value.image_tag;
        this.step3TypeStatus.set(container, false);
        this.setDefaultContainerInfo(container, true);
        this.step3Output.push(container);
      });
    }
  }

  ngOnDestroy() {
    this.step3Output.forEach((container: Container) => {
      let ports = this.containerPort.get(container.name);
      ports.forEach(port => {
        if (!container.ports.find(value => value.containerPort == port)) {
          container.ports.push({containerPort: port});
        }
      });
    });
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

  setDefaultContainerInfo(container: Container, isNew: boolean): void {
    let imageName: string = container.image.split(":")[0];
    let imageTag: string = container.image.split(":")[1];
    let firstIndex = imageName.indexOf("/");
    let projectName = imageName.slice(0, firstIndex);
    if (isNew) {
      this.containerPort.set(container.name, Array<number>());
    } else {
      let ports: Array<number> = Array();
      container.ports.forEach(value => {
        ports.push(value.containerPort);
      });
      this.containerPort.set(container.name, ports);
    }
    if (this.step2Output[0].project_name == projectName) {
      this.k8sService.getContainerDefaultInfo(imageName, imageTag, projectName)
        .then((res: ImageDockerfile) => {
          this.step3TypeStatus.set(container, true);
          if (res.image_cmd && isNew) {
            container.command.push(res.image_cmd);//copy cmd
          }
          if (res.image_env) {
            res.image_env.forEach(value => {//copy env
              if (isNew) {
                container.env.push({
                  name: value.dockerfile_envname,
                  value: value.dockerfile_envvalue
                });
              }
              this.fixedEnvKeys.push(value.dockerfile_envname);
            });
          }
          if (res.image_expose) {
            let fixedPorts: Array<number> = Array();
            res.image_expose.forEach(value => {//copy port
              let port: number = Number(value).valueOf();
              fixedPorts.push(port);
              this.containerPort.get(container.name).push(port);
            });
            this.fixedContainerPort.set(container.name, fixedPorts);
          }
        }).catch(() => {
      });
    }
  }

  get isCanNextStep(): boolean {
    return this.step3Output.length > 0 && this.isInputComponentsValid;
  }

  getVolumesDescription(container: Container): string {
    let volumesArr = container.volumeMounts;
    let result: string = "";
    volumesArr.forEach(value => {
      let storageServer = value.name == "" ? "" : value.name.concat(":");
      result += `${value.mountPath}:${storageServer}${value.mountPath}`
    });
    return result == ":" ? "" : result;
  }

  getEnvsDescription(container: Container): string {
    let envsArr = container.env;
    let result: string = "";
    envsArr.forEach(value => {
      result += `${value.name}=${value.value};`
    });
    return result;
  }

  getDefaultEnvsData(index: number) {
    let result = Array<EnvType>();
    this.step3Output[index].env.forEach(value => {
      result.push(new EnvType(value.name, value.value))
    });
    return result;
  }

  setEnvironment(index: number, envsData: Array<EnvType>) {
    let envsArray = this.step3Output[index].env;
    envsArray.splice(0, envsArray.length);
    envsData.forEach((value: EnvType) => {
      envsArray.push({name: value.envName, value: value.envValue})
    });
  }

  setVolumeMount(data: VolumeOutPut, index: number) {
    let volumeArr = this.step3Output[index].volumeMounts;
    if (volumeArr.length == 0) {
      volumeArr.push({
        name: data.out_name,
        mountPath: data.out_mountPath,
        ui_nfs_server: data.out_medium,
        ui_nfs_path: data.out_path
      })
    } else {
      volumeArr[0].name = data.out_name;
      volumeArr[0].mountPath = data.out_mountPath;
      volumeArr[0].ui_nfs_server = data.out_medium;
      volumeArr[0].ui_nfs_path = data.out_path;
    }
  }

  getVolumeMountData(index: number): VolumeOutPut {
    let volumeArr = this.step3Output[index].volumeMounts;
    if (volumeArr.length == 0) {
      return {
        out_name: "",
        out_mountPath: "",
        out_path: "",
        out_medium: ""
      };
    } else {
      return {
        out_name: volumeArr[0].name,
        out_mountPath: volumeArr[0].mountPath,
        out_path: volumeArr[0].ui_nfs_path,
        out_medium: volumeArr[0].ui_nfs_server
      };
    }
  }

  toggleShowStatus(container: Container): void {
    let status = this.step3TypeStatus.get(container);
    this.step3TypeStatus.set(container, !status);
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