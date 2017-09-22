import { Component, Input, OnDestroy, OnInit, QueryList, ViewChildren } from '@angular/core';
import {
  ServiceStep1Output,
  ServiceStep2Output,
  ServiceStep2Type,
  ServiceStep2NewImageType,
  ServiceStepComponent
} from '../service-step.component';
import { K8sService } from '../service.k8s';
import { MessageService } from "../../shared/message-service/message.service";
import { Image, ImageDetail } from "../../image/image";
import { AppInitService } from "../../app.init.service";
import { Message } from "../../shared/message-service/message";
import { EnvType } from "../environment-value/environment-value.component";
import { CsInputArrayComponent } from "../cs-input-array/cs-input-array.component";
import { CsInputComponent } from "../cs-input/cs-input.component";

enum ImageSource{fromBoardRegistry, fromDockerHub}
const AUTO_REFRESH_IMAGE_LIST: number = 2000;
type alertType = "alert-info" | "alert-danger";
@Component({
  templateUrl: './select-image.component.html',
  styleUrls: ["./select-image.component.css"]
})
export class SelectImageComponent implements ServiceStepComponent, OnInit, OnDestroy {
  @Input() data: any;
  @ViewChildren(CsInputArrayComponent) inputArrayComponents: QueryList<CsInputArrayComponent>;
  @ViewChildren(CsInputComponent) inputComponents: QueryList<CsInputComponent>;
  patternNewImageName: RegExp = /^[a-z\d.]+$/;
  patternNewImageTag: RegExp = /^[a-z\d.]+$/;
  patternBaseImage: RegExp = /^[a-z\d.]+$/;
  patternExpose: RegExp = /^[\d-\s\w/\\]+$/;
  patternVolume: RegExp = /^[a-zA-Z_]+$/;
  patternRun: RegExp = /^[a-zA-Z_]+$/;
  patternEntryPoint:RegExp = /^[a-zA-Z_]+$/;
  _isOpenEnvironment = false;
  intervalAutoRefreshImageList: any;
  isNeedAutoRefreshImageList: boolean = false;
  imageSource: ImageSource = ImageSource.fromBoardRegistry;
  imageSourceList: Array<Image>;
  imageSelectList: Array<Image>;
  imageDetailSourceList: Map<string, Array<ImageDetail>>;
  imageDetailSelectList: Map<string, ImageDetail>;
  imageTemplateList: Array<Object> = [{name: "Docker File Template"}];
  customerNewImage: ServiceStep2NewImageType;
  customerCreateImage: Image;
  outputData: ServiceStep2Output;
  filesList: Map<string, Array<{path: string, file_name: string, size: number}>>;
  consoleText: string = "";
  isOpenNewImage: boolean = false;
  newImageErrMessage: string = "";
  newImageAlertType: alertType = "alert-danger";
  isNewImageAlertOpen: boolean = false;
  newImageIndex: number;

  constructor(private k8sService: K8sService,
              private messageService: MessageService,
              private appInitService: AppInitService) {
    this.outputData = Array<ServiceStep2Type>();
    this.imageSelectList = Array<Image>();
    this.imageDetailSelectList = new Map<string, ImageDetail>();
    this.imageDetailSourceList = new Map<string, Array<ImageDetail>>();
    this.filesList = new Map<string, Array<{path: string, file_name: string, size: number}>>();
    this.customerCreateImage = new Image();
  }

  ngOnInit() {
    this.customerCreateImage.image_name = "SERVICE.STEP_2_CREATE_IMAGE";
    this.customerCreateImage["isSpecial"] = true;
    this.customerCreateImage["OnlyClick"] = true;
    this.k8sService.getImages("", 0, 0)
      .then(res => {
        this.imageSourceList = res;
        this.imageSourceList.unshift(this.customerCreateImage);
      })
      .catch(err => this.messageService.dispatchError(err));
    this.intervalAutoRefreshImageList = setInterval(() => {
      if (this.isNeedAutoRefreshImageList) {
        this.isNewImageAlertOpen = false;
        this.k8sService.getImages("", 0, 0).then(res => {
          res.forEach(value => {
            let newImageName = `${this.customerNewImage.project_name}/${this.customerNewImage.image_name}`;
            if (value.image_name == newImageName) {
              this.isNeedAutoRefreshImageList = false;
              this.imageSourceList = res;
              this.imageSelectList[this.newImageIndex] = value;
              this.setImageDetailList(value.image_name);
              this.isOpenNewImage = false;
            }
          });
        }).catch(err => {
          if (err && err.status == 401) {
            this.isOpenNewImage = false;
            this.messageService.dispatchError(err);
          } else {
            this.newImageAlertType = "alert-danger";
            this.newImageErrMessage = "SERVICE.STEP_2_UPDATE_IMAGE_LIST_FAILED";
            this.isNewImageAlertOpen = true;
          }
        });
      }
    }, AUTO_REFRESH_IMAGE_LIST);
  }

  ngOnDestroy() {
    this.k8sService.setStepData(2, this.outputData);
    clearInterval(this.intervalAutoRefreshImageList);
  }

  set isOpenEnvironment(value) {
    this.isOpenNewImage = !value;
    this._isOpenEnvironment = value;
  }

  get isOpenEnvironment() {
    return this._isOpenEnvironment;
  }

  get imageRun(): Array<string> {
    return this.customerNewImage.image_dockerfile.image_run;
  }

  get imageVolume(): Array<string> {

    return this.customerNewImage.image_dockerfile.image_volume;
  }

  get imageExpose(): Array<string> {
    return this.customerNewImage.image_dockerfile.image_expose;
  }

  get isCanNextStep(): boolean {
    let hasSelectImage = this.imageSelectList.filter(value => {
        return value.image_name != "SERVICE.STEP_2_SELECT_IMAGE";
      }).length > 0;
    return hasSelectImage && this.imageDetailSelectList.size > 0;
  }

  get selfObject() {
    return this;
  }

  get envsDescription() {
    let result: string = "";
    this.customerNewImage.image_dockerfile.image_env.forEach(value => {
      result += value.dockerfile_envname + "=" + value.dockerfile_envvalue + ";"
    });
    return result;
  }

  get defaultEnvsData() {
    let result = Array<EnvType>();
    this.customerNewImage.image_dockerfile.image_env.forEach(value => {
      result.push(new EnvType(value.dockerfile_envname, value.dockerfile_envvalue))
    });
    return result;
  }

  get newImageFormValid(): boolean {
    let result = true;
    this.inputArrayComponents.forEach(item => {
      if (!item.arrayIsValid()) {
        result = false;
      }
    });
    this.inputComponents.forEach(item => {
      if (!item.valid) {
        result = false;
      }
    });
    return result;
  }

  shieldEnter($event: KeyboardEvent) {
    if ($event.charCode == 13) {
      (<any>$event.target).blur();
      this.getDockerFilePreviewInfo();
    }
  }

  setImageDetailList(imageName: string): void {
    this.k8sService.getImageDetailList(imageName).then((res: ImageDetail[]) => {
      for (let item of res) {
        item['image_detail'] = JSON.parse(item['image_detail']);
        item['image_size_number'] = Number.parseFloat((item['image_size_number'] / (1024 * 1024)).toFixed(2));
        item['image_size_unit'] = 'MB';
      }
      this.imageDetailSourceList.set(res[0].image_name, res);
      this.imageDetailSelectList.set(res[0].image_name, res[0]);
    }).catch(err => this.messageService.dispatchError(err));
  }

  canChangeSelectImage(image: Image) {
    if (this.imageSelectList.indexOf(image) > -1) {
      let m: Message = new Message();
      m.message = "SERVICE.STEP_2_IMAGE_SELECTED";
      this.messageService.inlineAlertMessage(m);
      return false;
    }
    return true;
  }

  changeSelectImage(index: number, image: Image) {
    this.imageSelectList[index] = image;
    this.setImageDetailList(image.image_name);
  }

  clickSelectImage(index: number, image: Image) {
    let step1Out: ServiceStep1Output = this.k8sService.getStepData(1) as ServiceStep1Output;
    this.customerNewImage = new ServiceStep2NewImageType();
    this.customerNewImage.image_dockerfile.image_author = this.appInitService.currentUser["user_name"];
    this.customerNewImage.project_id = step1Out.project_id;
    this.customerNewImage.project_name = step1Out.project_name;
    this.customerNewImage.image_template = "dockerfile-template";
    this.isOpenNewImage = true;
    this.newImageIndex = index;
  }

  changeSelectImageDetail(imageName: string, imageDetail: ImageDetail) {
    this.imageDetailSelectList.set(imageName, imageDetail);
  }

  minusSelectImage(index: number) {
    let image = this.imageSelectList[index];
    this.imageDetailSelectList.delete(image.image_name);
    this.imageSelectList.splice(index, 1);
  }

  addSelectImage() {
    let customerSelectImage = new Image();
    customerSelectImage.image_name = "SERVICE.STEP_2_SELECT_IMAGE";
    this.imageSelectList.push(customerSelectImage);
  }

  buildImage() {
    this.isNewImageAlertOpen = false;
    this.isNeedAutoRefreshImageList = true;
    this.k8sService.buildImage(this.customerNewImage)
      .then(res => res)
      .catch((err) => {
        this.isNeedAutoRefreshImageList = false;
        if (err && err.status == 401) {
          this.isOpenNewImage = false;
          this.messageService.dispatchError(err);
        } else {
          this.newImageAlertType = "alert-danger";
          this.newImageErrMessage = "SERVICE.STEP_2_BUILD_IMAGE_FAILED";
          this.isNewImageAlertOpen = true;
        }
      })
  }

  updateFileList(): Promise<boolean> {
    this.isNewImageAlertOpen = false;
    let formFileList: FormData = new FormData();
    formFileList.append('project_name', this.customerNewImage.project_name);
    formFileList.append('image_name', this.customerNewImage.image_name);
    formFileList.append('tag_name', this.customerNewImage.image_tag);
    return this.k8sService.getFileList(formFileList).then(res => {
      this.filesList.set(this.customerNewImage.image_name, res);
      let imageCopyArr = this.customerNewImage.image_dockerfile.image_copy;
      imageCopyArr.splice(0, imageCopyArr.length);
      this.filesList.get(this.customerNewImage.image_name).forEach(value => {
        imageCopyArr.push({
          dockerfile_copyfrom: value.path + "/" + value.file_name,
          dockerfile_copyto: "/tmp"
        });
      });
      return true;
    }).catch(err => {
      if (err && err.status == 401) {
        this.isOpenNewImage = false;
        this.messageService.dispatchError(err);
      } else {
        this.newImageAlertType = "alert-danger";
        this.newImageErrMessage = "SERVICE.STEP_2_UPDATE_FILE_LIST_FAILED";
        this.isNewImageAlertOpen = true;
      }
    });
  }

  async asyncGetDockerFilePreviewInfo() {
    await this.updateFileList();
    this.getDockerFilePreviewInfo();
  }

  async uploadFile(event) {
    let fileList: FileList = event.target.files;
    if (fileList.length > 0) {
      this.isNewImageAlertOpen = false;
      let file: File = fileList[0];
      let formData: FormData = new FormData();
      formData.append('upload_file', file, file.name);
      formData.append('project_name', this.customerNewImage.project_name);
      formData.append('image_name', this.customerNewImage.image_name);
      formData.append('tag_name', this.customerNewImage.image_tag);
      this.k8sService.uploadFile(formData).then(res => {
        event.target.value = "";
        this.newImageAlertType = "alert-info";
        this.newImageErrMessage = "SERVICE.STEP_2_UPLOAD_SUCCESS";
        this.isNewImageAlertOpen = true;
        this.asyncGetDockerFilePreviewInfo();
      }).catch(err => {
        if (err && err.status == 401) {
          this.isOpenNewImage = false;
          this.messageService.dispatchError(err);
        } else {
          this.newImageAlertType = "alert-danger";
          this.newImageErrMessage = "SERVICE.STEP_2_UPLOAD_FAILED";
          this.isNewImageAlertOpen = true;
        }
      });
    }
  }

  getDockerFilePreviewInfo() {
    if (this.customerNewImage.image_dockerfile.image_base != "") {
      this.isNewImageAlertOpen = false;
      this.k8sService.getDockerFilePreview(this.customerNewImage)
        .then(res => {
          this.consoleText = res;
        }).catch(err => {
        if (err && err.status == 401) {
          this.isOpenNewImage = false;
          this.messageService.dispatchError(err);
        } else {
          this.newImageAlertType = "alert-danger";
          this.newImageErrMessage = "SERVICE.STEP_2_UPDATE_DOCKER_FILE_FAILED";
          this.isNewImageAlertOpen = true;
        }
      });
    }
  }

  setEnvironment(envsData: Array<EnvType>) {
    let envsArray = this.customerNewImage.image_dockerfile.image_env;
    envsArray.splice(0, envsArray.length);
    envsData.forEach((value: EnvType) => {
      envsArray.push({
        dockerfile_envname: value.envName,
        dockerfile_envvalue: value.envValue
      })
    });
    this.getDockerFilePreviewInfo();
  }

  forward(): void {
    let step1Out: ServiceStep1Output = this.k8sService.getStepData(1) as ServiceStep1Output;
    this.imageSelectList.forEach(value => {
      if (value.image_name != "SERVICE.STEP_2_SELECT_IMAGE") {
        let serviceStep2 = new ServiceStep2Type();
        serviceStep2.image_name = value.image_name;
        serviceStep2.image_tag = this.imageDetailSelectList.get(value.image_name).image_tag;
        serviceStep2.project_name = step1Out.project_name;
        serviceStep2.image_template = "dockerfile-template";
        this.outputData.push(serviceStep2);
      }
    });
    this.k8sService.stepSource.next(3);
  }
}