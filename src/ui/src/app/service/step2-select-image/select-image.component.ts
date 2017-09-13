import { Component, Input, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { ServiceStep1Output, ServiceStep2Output, ServiceStepComponent } from '../service-step.component';
import { K8sService } from '../service.k8s';
import { MessageService } from "../../shared/message-service/message.service";
import { Image, ImageDetail } from "../../image/image";
import { AppInitService } from "../../app.init.service";
import { Message } from "../../shared/message-service/message";
import { FormControl, FormGroup, ValidatorFn, Validators } from "@angular/forms";

enum ImageSource{
  fromBoardRegistry,
  fromDockerHub
}
const AUTO_REFRESH_IMAGE_LIST: number = 2000;
@Component({
  templateUrl: './select-image.component.html',
  styleUrls: ["./select-image.component.css"]
})
export class SelectImageComponent implements ServiceStepComponent, OnInit, OnDestroy {
  @Input() data: any;
  intervalAutoRefreshImageList: any;
  isNeedAutoRefreshImageList: boolean = false;
  imageSource: ImageSource = ImageSource.fromBoardRegistry;
  imageSourceList: Array<Image>;
  imageSelectList: Array<Image>;
  imageDetailSourceList: Map<string, Array<ImageDetail>>;
  imageDetailSelectList: Map<string, ImageDetail>;
  imageTemplateList: Array<Object> = [{name: "Docker File Template"}];
  customerCreateImage: Image;
  customerSelectImage: Image;
  outputData: ServiceStep2Output = new ServiceStep2Output();
  filesList: Map<string, Array<{path: string, file_name: string, size: number}>>;
  consoleText: string = "";
  isOpenNewImage: boolean = false;
  newImageIndex: number;
  testValidatorFn: Array<ValidatorFn> = [Validators.required, Validators.maxLength(10)];
  constructor(private k8sService: K8sService,
              private messageService: MessageService,
              private appInitService: AppInitService) {
    this.imageSelectList = Array<Image>();
    this.imageDetailSelectList = new Map<string, ImageDetail>();
    this.imageDetailSourceList = new Map<string, Array<ImageDetail>>();
    this.filesList = new Map<string, Array<{path: string, file_name: string, size: number}>>();
    this.customerCreateImage = new Image();
    this.customerSelectImage = new Image();
  }

  ngOnInit() {
    let step1Out: ServiceStep1Output = this.k8sService.getStepData(1) as ServiceStep1Output;
    this.outputData.image_dockerfile.image_author = this.appInitService.currentUser["user_name"];
    this.outputData.project_name = step1Out.project_name;
    this.outputData.image_template = "dockerfile-template";
    this.customerCreateImage.image_name = "SERVICE.STEP_2_CREATE_IMAGE";
    this.customerCreateImage["isSpecial"] = true;
    this.customerCreateImage["OnlyClick"] = true;
    this.customerSelectImage.image_name = "SERVICE.STEP_2_SELECT_IMAGE";
    this.k8sService.getImages("", 0, 0)
      .then(res => {
        this.imageSourceList = res;
        this.imageSourceList.unshift(this.customerCreateImage);
        this.imageSelectList.push(this.customerSelectImage);
      })
      .catch(err => this.messageService.dispatchError(err));
    this.intervalAutoRefreshImageList = setInterval(() => {
      if (this.isNeedAutoRefreshImageList) {
        this.k8sService.getImages("", 0, 0).then(res => {
          res.forEach(value => {
            let newImageName = `${this.outputData.project_name}/${this.outputData.image_name}`;
            if (value.image_name == newImageName) {
              this.isNeedAutoRefreshImageList = false;
              this.imageSourceList = res;
              this.imageSelectList[this.newImageIndex] = value;
              this.isOpenNewImage = false;
            }
          });
        }).catch(err => this.messageService.dispatchError(err));
      }
    }, AUTO_REFRESH_IMAGE_LIST);
  }

  ngOnDestroy() {
    this.k8sService.setStepData(2, this.outputData);
    clearInterval(this.intervalAutoRefreshImageList);
  }

  get imageRun(): Array<string> {
    return this.outputData.image_dockerfile.image_run;
  }

  get imageVolume(): Array<string> {
    return this.outputData.image_dockerfile.image_volume;
  }

  get isCanNextStep(): boolean {
    return this.imageSelectList.filter(value => {
        return value.image_name != "SERVICE.STEP_2_SELECT_IMAGE" &&
          value.image_name != "SERVICE.STEP_2_CREATE_IMAGE"
      }).length > 0;
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

  changeSelectImage(index: number, image: Image) {
    this.imageSelectList[index] = image;
    this.setImageDetailList(image.image_name);
  }

  clickSelectImage(index: number, image: Image) {
    this.isOpenNewImage = true;
    this.newImageIndex = index;
  }

  changeSelectImageDetail(imageName: string, imageDetail: ImageDetail) {
    this.imageDetailSelectList.set(imageName, imageDetail);
  }

  modifySelectImage(index: number) {
    if (index == this.imageSelectList.length - 1) {
      this.imageSelectList.push(this.customerSelectImage);
    } else {
      this.imageSelectList.splice(index, 1);
    }
  }

  buildImage() {
    this.isNeedAutoRefreshImageList = true;
    this.k8sService.buildImage(this.outputData)
      .then(res => res)
      .catch((err) => {
        this.messageService.dispatchError(err);
        this.isNeedAutoRefreshImageList = false;
      })
  }

  updateFileList(): Promise<boolean> {
    let formFileList: FormData = new FormData();
    formFileList.append('project_name', this.outputData.project_name);
    formFileList.append('image_name', this.outputData.image_name);
    formFileList.append('tag_name', this.outputData.image_tag);
    return this.k8sService.getFileList(formFileList).then(res => {
      this.filesList.set(this.outputData.image_name, res);
      let imageCopyArr = this.outputData.image_dockerfile.image_copy;
      imageCopyArr.splice(0, imageCopyArr.length);
      this.filesList.get(this.outputData.image_name).forEach(value => {
        imageCopyArr.push({
          dockerfile_copyfrom: value.path + "/" + value.file_name,
          dockerfile_copyto: "tmp"
        });
      });
      return true;
    }).catch(err => this.messageService.dispatchError(err));
  }

  async asyncGetDockerFilePreviewInfo() {
    await this.updateFileList();
    this.getDockerFilePreviewInfo();
  }

  async uploadFile(event) {
    let fileList: FileList = event.target.files;
    if (fileList.length > 0) {
      let file: File = fileList[0];
      let formData: FormData = new FormData();
      formData.append('upload_file', file, file.name);
      formData.append('project_name', this.outputData.project_name);
      formData.append('image_name', this.outputData.image_name);
      formData.append('tag_name', this.outputData.image_tag);
      this.k8sService.uploadFile(formData).then(res => {
        event.target.value = "";
        let m: Message = new Message();
        m.message = "SERVICE.STEP_2_UPLOAD_SUCCESS";
        this.messageService.inlineAlertMessage(m);
        this.asyncGetDockerFilePreviewInfo();
      }).catch(err => this.messageService.dispatchError(err));
    }
  }

  getDockerFilePreviewInfo() {
    this.k8sService.getDockerFilePreview(this.outputData)
      .then(res => {
        this.consoleText = res;
      }).catch(err => this.messageService.dispatchError(err));
  }

  forward(): void {
    this.k8sService.stepSource.next(3);
  }
}