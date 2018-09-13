import { Component, Injector, OnInit } from '@angular/core';
import {
  Container,
  EnvStruct,
  ImageIndex,
  PHASE_CONFIG_CONTAINERS,
  PHASE_SELECT_IMAGES,
  UiServiceFactory,
  UIServiceStep2,
  UIServiceStep3
} from '../service-step.component';
import { BuildImageDockerfileData, Image, ImageDetail } from "../../image/image";
import { ServiceStepBase } from "../service-step";
import { CreateImageComponent } from "../../image/image-create/image-create.component";

@Component({
  templateUrl: './select-image.component.html',
  styleUrls: ["./select-image.component.css"]
})
export class SelectImageComponent extends ServiceStepBase implements OnInit {
  patternContainerName: RegExp = /^[a-zA-Z\d_-]+$/;
  patternWorkdir: RegExp = /^~?[\w\d-\/.{}$\/:]+[\s]*$/;
  imageSourceList: Array<Image>;
  imageDetailSourceList: Map<string, Array<ImageDetail>>;
  imageTagNotReadyList: Map<string, boolean>;
  workBufferList: Array<{imageIndex: ImageIndex, container: Container}>;
  containerIsInEdit: Map<Container, boolean>;
  fixedContainerPort: Map<string, Array<number>>;
  fixedContainerEnv: Map<string, Array<EnvStruct>>;
  stepSelectImageData: UIServiceStep2;
  stepConfigContainerData: UIServiceStep3;

  constructor(protected injector: Injector) {
    super(injector);
    this.workBufferList = Array<{imageIndex: ImageIndex, container: Container}>();
    this.imageDetailSourceList = new Map<string, Array<ImageDetail>>();
    this.imageTagNotReadyList = new Map<string, boolean>();
    this.containerIsInEdit = new Map<Container, boolean>();
    this.fixedContainerPort = new Map<string, Array<number>>();
    this.fixedContainerEnv = new Map<string, Array<EnvStruct>>();
    this.stepSelectImageData = UiServiceFactory.getInstance(PHASE_SELECT_IMAGES) as UIServiceStep2;
    this.stepConfigContainerData = UiServiceFactory.getInstance(PHASE_CONFIG_CONTAINERS) as UIServiceStep3;
  }

  ngOnInit() {
    let promiseSelectImage = this.k8sService.getServiceConfig(PHASE_SELECT_IMAGES);
    let promiseConfigContainer = this.k8sService.getServiceConfig(PHASE_CONFIG_CONTAINERS);
    Promise.all([promiseSelectImage, promiseConfigContainer]).then(([resSelectImage, resConfigContainer]) => {
      this.stepSelectImageData = resSelectImage as UIServiceStep2;
      this.stepConfigContainerData = resConfigContainer as UIServiceStep3;
      if (this.stepSelectImageData.imageList.length == 0) {
        this.addEmptyWorkItem();
      }
      this.stepSelectImageData.imageList.forEach((imageIndex: ImageIndex) => {
        this.getImageDetailList(imageIndex.image_name).then();
        let imageIndexBuf = new ImageIndex();
        imageIndexBuf.image_name = imageIndex.image_name;
        imageIndexBuf.image_tag = imageIndex.image_tag;
        imageIndexBuf.project_name = imageIndex.project_name;
        let container = this.stepConfigContainerData.containerList.find(value =>
          value.image.image_name == imageIndex.image_name && value.image.image_tag == imageIndex.image_tag);
        let containerBuf = new Container();
        containerBuf.image.image_name = imageIndex.image_name;
        containerBuf.image.image_tag = imageIndex.image_tag;
        containerBuf.image.project_name = imageIndex.project_name;
        containerBuf.volume_mount.container_path = container.volume_mount.container_path;
        containerBuf.volume_mount.target_path = container.volume_mount.target_path;
        containerBuf.volume_mount.volume_name = container.volume_mount.volume_name;
        containerBuf.volume_mount.target_storage_service = container.volume_mount.target_storage_service;
        containerBuf.name = container.name;
        containerBuf.command = container.command;
        containerBuf.working_dir = container.working_dir;
        container.container_port.forEach(port => containerBuf.container_port.push(port));
        container.env.forEach(env => {
          let envBuf = new EnvStruct();
          envBuf.dockerfile_envname = env.dockerfile_envname;
          envBuf.dockerfile_envvalue = env.dockerfile_envvalue;
          containerBuf.env.push(envBuf);
        });
        this.workBufferList.push({imageIndex: imageIndexBuf, container: containerBuf});
        this.containerIsInEdit.set(containerBuf, false);
        this.setContainerFixedInfo(containerBuf);
      });
    });
    this.k8sService.getImages("", 0, 0).then(res => {
      this.imageSourceList = res;
      this.unshiftCustomerCreateImage();
    })
  }

  changeSelectImage(index: number, image: Image) {
    let buf = this.workBufferList[index];
    buf.imageIndex.image_name = image.image_name;
    buf.imageIndex.project_name = this.stepSelectImageData.projectName;
    let containerBuf = new Container();
    containerBuf.image.image_name = image.image_name;
    containerBuf.name = image.image_name;
    containerBuf.image.project_name = this.stepSelectImageData.projectName;
    buf.container = containerBuf;
    if (this.imageDetailSourceList.has(image.image_name)) {
      let detailList: Array<ImageDetail> = this.imageDetailSourceList.get(image.image_name);
      containerBuf.image.image_tag = detailList[0].image_tag;
      buf.imageIndex.image_tag = detailList[0].image_tag;
      this.setDefaultContainerInfo(containerBuf);
      this.setContainerFixedInfo(containerBuf);
    } else {
      this.getImageDetailList(image.image_name).then((res: ImageDetail[]) => {
        buf.imageIndex.image_tag = res[0].image_tag;
        containerBuf.image.image_tag = res[0].image_tag;
        this.setDefaultContainerInfo(containerBuf);
        this.setContainerFixedInfo(containerBuf);
      })
    }
  }

  changeSelectImageDetail(imageName: string, imageDetail: ImageDetail) {
    let workBuf = this.workBufferList.find(value => value.container.image.image_name == imageName);
    workBuf.container.image.image_tag = imageDetail.image_tag;
    this.setDefaultContainerInfo(workBuf.container);
    this.setContainerFixedInfo(workBuf.container);
  }

  getImageDetailList(imageName: string): Promise<ImageDetail[]> {
    this.imageTagNotReadyList.set(imageName, false);
    return this.k8sService.getImageDetailList(imageName).then((res: ImageDetail[]) => {
      if (res && res.length > 0) {
        for (let item of res) {
          item['image_detail'] = JSON.parse(item['image_detail']);
          item['image_size_number'] = Number.parseFloat((item['image_size_number'] / (1024 * 1024)).toFixed(2));
          item['image_size_unit'] = 'MB';
        }
        this.imageDetailSourceList.set(imageName, res);
      } else {
        this.imageTagNotReadyList.set(imageName, true);
      }
      return res;
    })
  }

  setContainerFixedInfo(container: Container): void {
    this.k8sService.getContainerDefaultInfo(container.image.image_name, container.image.image_tag, container.image.project_name)
      .then((res: BuildImageDockerfileData) => {
        if (res.image_env) {
          let fixedEnvs: Array<EnvStruct> = Array<EnvStruct>();
          res.image_env.forEach(value => {
            let env = new EnvStruct();
            env.dockerfile_envname = value.dockerfile_envname;
            env.dockerfile_envvalue = value.dockerfile_envvalue;
            fixedEnvs.push(env);
          });
          this.fixedContainerEnv.set(container.image.image_name,fixedEnvs);
        }
        if (res.image_expose) {
          let fixedPorts: Array<number> = Array();
          res.image_expose.forEach(value => {
            let port: number = Number(value).valueOf();
            fixedPorts.push(port);
          });
          this.fixedContainerPort.set(container.image.image_name, fixedPorts);
        }
      }).catch(() => this.messageService.cleanNotification());
  }


  setDefaultContainerInfo(container: Container): void {
    this.k8sService.getContainerDefaultInfo(container.image.image_name, container.image.image_tag, container.image.project_name)
      .then((res: BuildImageDockerfileData) => {
        if (res.image_cmd) {
          container.command = res.image_cmd;
        }
        if (res.image_env) {
          res.image_env.forEach(value => {
            let env = new EnvStruct();
            env.dockerfile_envname = value.dockerfile_envname;
            env.dockerfile_envvalue = value.dockerfile_envvalue;
            container.env.push(env);
          });
        }
        if (res.image_expose) {
          res.image_expose.forEach(value => {
            let port: number = Number(value).valueOf();
            container.container_port.push(port);
          });
        }
      }).catch(() => this.messageService.cleanNotification());
  }

  forward(): void {
    this.stepSelectImageData.imageList.splice(0, this.stepSelectImageData.imageList.length);
    this.stepConfigContainerData.containerList.splice(0, this.stepConfigContainerData.containerList.length);
    this.workBufferList.forEach(workBuf => {
      this.stepSelectImageData.imageList.push(workBuf.imageIndex);
      this.stepConfigContainerData.containerList.push(workBuf.container);
    });
    let obsSelectImage = this.k8sService.setServiceConfig(this.stepSelectImageData.uiToServer());
    let obsConfigContainer = this.k8sService.setServiceConfig(this.stepConfigContainerData.uiToServer());
    Promise.all([obsSelectImage, obsConfigContainer]).then(() => this.k8sService.stepSource.next({index: 3, isBack: false}));
  }

  get isCanNextStep(): boolean {
    return this.workBufferList
      .filter(value => value.imageIndex.image_name != "SERVICE.STEP_2_SELECT_IMAGE")
      .length == this.workBufferList.length;
  }

  get selfObject() {
    return this;
  }

  unshiftCustomerCreateImage() {
    let customerCreateImage: Image = new Image();
    customerCreateImage.image_name = "SERVICE.STEP_2_CREATE_IMAGE";
    customerCreateImage["isSpecial"] = true;
    customerCreateImage["OnlyClick"] = true;
    this.imageSourceList.unshift(customerCreateImage);
  }

  canChangeSelectImage(image: Image) {
    if (this.workBufferList.find(value => value.imageIndex.image_name == image.image_name)) {
      this.messageService.showAlert('IMAGE.CREATE_IMAGE_EXIST', {alertType: "alert-warning"});
      return false;
    }
    return true;
  }

  createNewCustomImage(index: number) {
    let newImageIndex = index;
    let component = this.createNewModal(CreateImageComponent);
    component.initCustomerNewImage(this.stepSelectImageData.projectId, this.stepSelectImageData.projectName);
    component.closeNotification.subscribe((imageName: string) => {
      if (imageName) {
        this.k8sService.getImages("", 0, 0).then(res => {
          res.forEach(value => {
            if (value.image_name === imageName) {
              this.imageSourceList = Object.create(res);
              this.unshiftCustomerCreateImage();
              this.changeSelectImage(newImageIndex,value);
            }
          });
        })
      }
    })
  }

  minusSelectImage(index: number) {
    if (index > 0){
      this.workBufferList.splice(index, 1);
    }
  }

  addEmptyWorkItem() {
    let imageIndexBuf = new ImageIndex();
    imageIndexBuf.image_name = 'SERVICE.STEP_2_SELECT_IMAGE';
    let containerBuf = new Container();
    this.workBufferList.push({imageIndex: imageIndexBuf, container: containerBuf});
  }

  getVolumesDescription(container: Container): string {
    let volume = container.volume_mount;
    let storageServer = volume.target_storage_service == "" ? "" : volume.target_storage_service.concat(":");
    let result = `${volume.container_path}:${storageServer}${volume.target_path}`;
    return result == ":" ? "" : result;
  }

  getEnvsDescription(container: Container): string {
    let envsArr = container.env;
    let result: string = "";
    envsArr.forEach((value: EnvStruct) => {
      result += `${value.dockerfile_envname}=${value.dockerfile_envvalue};`
    });
    return result;
  }

  toggleContainerEditStatus(container: Container):void {
    let oldStatus = this.containerIsInEdit.get(container);
    this.containerIsInEdit.set(container,!oldStatus);
  }

  backStep(){
    this.k8sService.stepSource.next({index: 1, isBack: true});
  }
}