import { Component, Input, OnInit, OnDestroy } from '@angular/core';
import {
  ImageDockerfile,
  ServiceStep2Output, ServiceStep2Type, ServiceStep3Output, ServiceStep3Type,
  ServiceStepComponent
} from '../service-step.component';
import { K8sService } from '../service.k8s';
import { MessageService } from "../../shared/message-service/message.service";
import { EnvType } from "../environment-value/environment-value.component";
import { Message } from "../../shared/message-service/message";

enum ContainerStatus{csNew, csSelectedImage}
@Component({
  templateUrl: './edit-container.component.html',
  styleUrls: ["./edit-container.component.css"]
})
export class EditContainerComponent implements ServiceStepComponent, OnInit, OnDestroy {
  @Input() data: any;
  step2Output: ServiceStep2Output;
  step3Output: ServiceStep3Output;
  containerStatusList: Array<ContainerStatus>;
  showVolumeMounts = false;
  showEnvironmentValue = false;

  constructor(private k8sService: K8sService,
              private messageService: MessageService) {
    this.step3Output = Array<ServiceStep3Type>();
    this.containerStatusList = Array<ContainerStatus>();
  }

  ngOnInit() {
    this.step2Output = this.k8sService.getStepData(2) as ServiceStep2Output;
    this.containerStatusList.push(ContainerStatus.csNew)
  }

  ngOnDestroy() {
    this.k8sService.setStepData(3, this.step3Output);
  }

  get isCanNextStep(): boolean {
    return this.step3Output.length > 0;
  }

  get selfObject() {
    return this;
  }

  modifySelectContainer(index: number) {
    if (index == this.containerStatusList.length - 1) {
      this.containerStatusList.push(ContainerStatus.csNew);
    } else {
      this.containerStatusList.splice(index, 1);
      this.step3Output.splice(index, 1);
    }
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
    config.container_name = image.image_name;
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
              config.container_envs.push({env_name: value.dockerfile_envvalue, env_value: value.dockerfile_envvalue})
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
      result += `${value.container_dir} : ${storageServer}${value.target_dir}`
    });
    return result;
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
      result.push(new EnvType(value.env_value, value.env_name))
    });
    return result;
  }

  setEnvironment(index: number, envsData: Array<EnvType>) {
    let envsArray = this.step3Output[index].container_envs;
    envsArray.splice(0, envsArray.length);
    envsData.forEach((value: EnvType) => {
      envsArray.push({env_name: value.envName, env_value: value.envValue})
    })
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