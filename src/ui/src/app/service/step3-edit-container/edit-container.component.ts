import { Component, Input, OnInit, OnDestroy, AfterContentChecked, QueryList, ViewChildren } from '@angular/core';
import {
  ImageDockerfile,
  ServiceStep2Output, ServiceStep2Type, ServiceStep3Output, ServiceStep3Type,
  ServiceStepComponent
} from '../service-step.component';
import { K8sService } from '../service.k8s';
import { MessageService } from "../../shared/message-service/message.service";
import { EnvType } from "../environment-value/environment-value.component";
import { Message } from "../../shared/message-service/message";
import { CsInputComponent } from "../cs-input/cs-input.component";

enum ContainerStatus{csNew, csSelectedImage}
@Component({
  templateUrl: './edit-container.component.html',
  styleUrls: ["./edit-container.component.css"]
})
export class EditContainerComponent implements ServiceStepComponent, OnInit, OnDestroy, AfterContentChecked {
  @Input() data: any;
  @ViewChildren(CsInputComponent) inputComponents: QueryList<CsInputComponent>;
  patternContainerName: RegExp = /^[a-zA-Z_-]+$/;
  patternWorkdir: RegExp = /^~?[\w\d-\/.{}$\/:]+[\s]*$/;
  step2Output: ServiceStep2Output;
  step3Output: ServiceStep3Output;
  containerStatusList: Array<ContainerStatus>;
  showVolumeMounts = false;
  showEnvironmentValue = false;
  isInputComponentsValid = false;

  constructor(private k8sService: K8sService,
              private messageService: MessageService) {
    this.step3Output = Array<ServiceStep3Type>();
    this.containerStatusList = Array<ContainerStatus>();
  }

  ngOnInit() {
    this.step2Output = this.k8sService.getStepData(2) as ServiceStep2Output;
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
  }

  get isCanNextStep(): boolean {
    return this.step3Output.length > 0 && this.isInputComponentsValid;
  }

  get selfObject() {
    return this;
  }

  minusSelectContainer(index: number) {
    this.containerStatusList.splice(index, 1);
    this.step3Output.splice(index, 1);
  }

  addSelectContainer() {
    this.containerStatusList.push(ContainerStatus.csNew);
  }

  changeSelectImage(index: number, status: ContainerStatus, image: ServiceStep2Type) {
    let config: ServiceStep3Type;
    if (status == ContainerStatus.csNew) {
      config = new ServiceStep3Type();
      this.containerStatusList[index] = ContainerStatus.csSelectedImage;
      this.step3Output.push(config);
    } else {
      config = this.step3Output[index];
    }
    let firstIndex = image.image_name.indexOf("/");
    config.container_name = image.image_name.slice(firstIndex + 1, image.image_name.length);
    config.container_baseimage = image.image_name + ":" + image.image_tag;
    let firstPos = image.image_name.indexOf("/");
    let projectName = image.image_name.slice(0, firstPos);
    if (projectName == image.project_name) {
      let imageName = image.image_name.slice(firstPos + 1);
      this.k8sService.getContainerDefaultInfo(imageName, image.image_tag, image.project_name)
        .then((res: ImageDockerfile) => {
          if (res.image_cmd) {
            config.container_command.push(res.image_cmd);//copy cmd
          }
          if (res.image_env) {
            res.image_env.forEach(value => {//copy env
              config.container_envs.push({env_name: value.dockerfile_envname, env_value: value.dockerfile_envvalue})
            });
          }
          if (res.image_expose) {
            res.image_expose.forEach(value => {//copy port
              config.container_ports.push(Number(value).valueOf());
            });
          }
        })
        .catch(err => this.messageService.dispatchError(err));
    }
  }

  getVolumesDescription(index: number): string {
    let volumesArr = this.step3Output[index].container_volumes;
    let result: string = "";
    volumesArr.forEach(value => {
      let storageServer = value.target_storageServer == "" ? "" : value.target_storageServer.concat(":");
      result += `${value.container_dir}:${storageServer}${value.target_dir}`
    });
    return result == ":" ? "" : result;
  }

  getSelectImageDefaultText(index: number) {
    if (this.containerStatusList[index] == ContainerStatus.csNew) {
      return "SERVICE.STEP_3_SELECT_IMAGE"
    }
    return this.step3Output[index].container_name;
  }

  canChangeSelectImage(item: ServiceStep2Type) {
    let baseImage = item.image_name + ":" + item.image_tag;
    let hasItem = this.step3Output.find(value => {
      return value.container_baseimage == baseImage;
    });
    if (hasItem) {
      let m: Message = new Message();
      m.message = "SERVICE.STEP_3_IMAGE_SELECTED";
      this.messageService.inlineAlertMessage(m);
      return false;
    }
    return true;
  }

  getEnvsDescription(index: number): string {
    let envsArr = this.step3Output[index].container_envs;
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

  shieldEnter($event: KeyboardEvent) {
    if ($event.charCode == 13) {
      (<any>$event.target).blur();
    }
  }

  forward(): void {
    this.k8sService.stepSource.next(4);
  }
}