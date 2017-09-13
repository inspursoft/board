import { Component, Directive, Input, OnInit } from '@angular/core';
import { ServiceStepComponent } from '../service-step.component';
import { K8sService } from '../service.k8s';
import { Image } from "../../image/image";
import { MessageService } from "../../shared/message-service/message.service";
import { CsInputFiled, CsInputStatus } from "../cs-input/cs-input.component";

enum ContainerStatus{csNew, csSelectedImage}
class ContainerConfigure {
  containerName: CsInputFiled<string>;
  workingDir: CsInputFiled<string>;
  volumeMounts: CsInputFiled<string>;
  env: CsInputFiled<string>;
  containerPort: CsInputFiled<number>;
  command: CsInputFiled<string>;
  cpu: CsInputFiled<string>;
  memory: CsInputFiled<string>;

  constructor() {
    this.containerName = new CsInputFiled(CsInputStatus.isView, "", "");
    this.workingDir = new CsInputFiled(CsInputStatus.isView, "", "");
    this.volumeMounts = new CsInputFiled(CsInputStatus.isView, "", "");
    this.env = new CsInputFiled(CsInputStatus.isView, "", "");
    this.containerPort = new CsInputFiled(CsInputStatus.isView, 0, 0);
    this.command = new CsInputFiled(CsInputStatus.isView, "", "");
    this.cpu = new CsInputFiled(CsInputStatus.isView, "", "");
    this.memory = new CsInputFiled(CsInputStatus.isView, "", "");
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
    container.status = ContainerStatus.csSelectedImage;
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