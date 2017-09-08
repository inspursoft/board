import { Component, Directive, Input, OnInit } from '@angular/core';
import { ServiceStepComponent } from '../service-step.component';
import { K8sService } from '../service.k8s';
import { Image } from "../../image/image";
import { MessageService } from "../../shared/message-service/message.service";
import { ContainerField, ContainerFieldStatus } from "./container-field/container-field.component";

enum ContainerStatus{csNew, csSelectedImage}

class ContainerConfigure {
  containerName: ContainerField<string>;
  workingDir: ContainerField<string>;
  volumeMounts: ContainerField<string>;
  env: ContainerField<string>;
  containerPort: ContainerField<number>;
  command: ContainerField<string>;
  cpu: ContainerField<string>;
  memory: ContainerField<string>;

  constructor() {
    this.containerName = new ContainerField(ContainerFieldStatus.cfsView, "", "");
    this.workingDir = new ContainerField(ContainerFieldStatus.cfsView, "", "");
    this.volumeMounts = new ContainerField(ContainerFieldStatus.cfsView, "", "");
    this.env = new ContainerField(ContainerFieldStatus.cfsView, "", "");
    this.containerPort = new ContainerField(ContainerFieldStatus.cfsView, 0, 0);
    this.command = new ContainerField(ContainerFieldStatus.cfsView, "", "");
    this.cpu = new ContainerField(ContainerFieldStatus.cfsView, "", "");
    this.memory = new ContainerField(ContainerFieldStatus.cfsView, "", "");
  }
}

@Component({
  templateUrl: './edit-container.component.html',
  styleUrls: ["./edit-container.component.css"]
})
export class EditContainerComponent implements ServiceStepComponent, OnInit {
  @Input() data: any;
  imageSourceList: Array<Image>;
  containerList: Array<{status: ContainerStatus, image: Image}>;
  containerDetailList: Map<string, ContainerConfigure>;
  showVolumeMounts = false;
  showEnvironmentValue = false;

  constructor(private k8sService: K8sService,
              private messageService: MessageService) {
    this.containerList = Array<{status: ContainerStatus, image: Image}>();
    this.containerDetailList = new Map<string, ContainerConfigure>();
  }

  ngOnInit() {
    this.k8sService.getImages("", 0, 0).then(res => {
      this.imageSourceList = res;
      this.containerList.push({status: ContainerStatus.csNew, image: null});
    }).catch(err => this.messageService.dispatchError(err))
  }


  modifySelectContainer(index: number) {
    if (index == this.containerList.length - 1) {
      this.containerList.push({status: ContainerStatus.csNew, image: null});
    } else {
      this.containerList.splice(index, 1);
    }
  }

  changeSelectImage(container: {status: ContainerStatus, image: Image}, image: Image) {
    container.image = image;
    container.status = 1;
    //get detail by image name
    let obj = new ContainerConfigure();
    obj.containerName.value = image.image_name;
    obj.containerName.defaultValue = image.image_name;
    this.containerDetailList.set(image.image_name, obj);
  }

  getDetail(container: {status: ContainerStatus, image: Image}): ContainerConfigure {
    return this.containerDetailList.get(container.image.image_name);
  }

  forward(): void {
    this.k8sService.stepSource.next(4);
  }
}