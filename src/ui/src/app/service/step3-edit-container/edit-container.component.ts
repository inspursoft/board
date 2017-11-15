import { Component, OnInit, AfterContentChecked, QueryList, ViewChildren, Injector } from '@angular/core';
import { ImageDockerfile, Container, ServiceContainerList, DeploymentServiceData } from '../service-step.component';
import { EnvType } from "../environment-value/environment-value.component";
import { CsInputComponent } from "../cs-input/cs-input.component";
import { CsInputArrayComponent } from "../cs-input-array/cs-input-array.component";
import { VolumeOutPut } from "./volume-mounts/volume-mounts.component";
import { ServiceStepBase } from "../service-step";

@Component({
  templateUrl: './edit-container.component.html',
  styleUrls: ["./edit-container.component.css"]
})
export class EditContainerComponent extends ServiceStepBase implements OnInit, AfterContentChecked {
  @ViewChildren(CsInputComponent) inputComponents: QueryList<CsInputComponent>;
  @ViewChildren(CsInputArrayComponent) inputArrayComponents: QueryList<CsInputArrayComponent>;
  patternContainerName: RegExp = /^[a-zA-Z\d_-]+$/;
  patternWorkdir: RegExp = /^~?[\w\d-\/.{}$\/:]+[\s]*$/;
  step3TypeStatus: Map<Container, boolean>;
  showVolumeMounts = false;
  showEnvironmentValue = false;
  isInputComponentsValid = false;
  fixedEnvKeys: Array<string>;
  containerPort: Map<string, Array<number>>;
  fixedContainerPort: Map<string, Array<number>>;
  curContainerIndex: number;

  constructor(protected injector: Injector) {
    super(injector);
    this.step3TypeStatus = new Map<Container, boolean>();
    this.containerPort = new Map<string, Array<number>>();
    this.fixedEnvKeys = Array<string>();
    this.fixedContainerPort = new Map<string, Array<number>>();
  }

  ngOnInit() {
    this.k8sService.getServiceConfig(this.newServiceId, this.outputData).then(res => {
      this.outputData = res;
      this.containerList.forEach((container: Container) => {
        this.step3TypeStatus.set(container, false);
        this.setDefaultContainerInfo(container);
      });
    });
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

  setDefaultContainerInfo(container: Container): void {
    let isNew = !this.isBack;
    let ports: Array<number> = Array();
    container.ports.forEach(value => {
      ports.push(value.containerPort);
    });
    this.containerPort.set(container.name, ports);
    let imageName: string = container.image.split(":")[0];
    let imageTag: string = container.image.split(":")[1];
    let firstIndex = imageName.indexOf("/");
    let projectName = firstIndex > -1 ? imageName.slice(0, firstIndex) : "";
    this.k8sService.getContainerDefaultInfo(imageName, imageTag, projectName)
      .then((res: ImageDockerfile) => {
        this.step3TypeStatus.set(container, true);
        if (res.image_cmd && isNew) {
          container.command.push(res.image_cmd);//copy cmd
        }
        if (res.image_env) {
          res.image_env.forEach(value => {//copy env
            this.fixedEnvKeys.push(value.dockerfile_envname);
            if (isNew) {
              container.env.push({name: value.dockerfile_envname, value: value.dockerfile_envvalue});
            }
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

  get isCanNextStep(): boolean {
    if (this.outputData) {
      let containerList: ServiceContainerList = this.outputData.deployment_yaml.spec.template.spec.containers;
      return containerList.length > 0 && this.isInputComponentsValid;
    } else {
      return false;
    }
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
    this.containerList[index].env.forEach(value => {
      result.push(new EnvType(value.name, value.value))
    });
    return result;
  }

  setEnvironment(index: number, envsData: Array<EnvType>) {
    let envsArray = this.containerList[index].env;
    envsArray.splice(0, envsArray.length);
    envsData.forEach((value: EnvType) => {
      envsArray.push({name: value.envName, value: value.envValue})
    });
  }

  setVolumeMount(data: VolumeOutPut, index: number) {
    let volumeArr = this.containerList[index].volumeMounts;
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
    let volumeArr = this.containerList[index].volumeMounts;
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

  setContainerPorts(): void {
    this.containerList.forEach((container: Container) => {
      container.ports.splice(0,container.ports.length);
      let ports = this.containerPort.get(container.name);
      ports.forEach(port => {
        container.ports.push({containerPort: port});
      });
    });
  }

  backStep(): void {
    this.setContainerPorts();
    this.k8sService.setServiceConfig(this.outputData).then(res => {
      this.k8sService.stepSource.next({index: 2, isBack: true});
    });
  }

  forward(): void {
    this.setContainerPorts();
    this.k8sService.setServiceConfig(this.outputData).then(res => {
      this.k8sService.stepSource.next({index: 4, isBack: false});
    });
  }
}