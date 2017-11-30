import { Component, OnInit, Injector } from '@angular/core';
import {
  PHASE_SELECT_IMAGES,
  PHASE_SELECT_PROJECT,
  ImageIndex,
  ServiceStepPhase,
  UIServiceStep2
} from '../service-step.component';
import { Image, ImageDetail } from "../../image/image";
import { Message } from "../../shared/message-service/message";
import { ServiceStepBase } from "../service-step";

@Component({
  templateUrl: './select-image.component.html',
  styleUrls: ["./select-image.component.css"]
})
export class SelectImageComponent extends ServiceStepBase implements OnInit {
  isOpenNewImage = false;
  imageSourceList: Array<Image>;
  imageSelectList: Array<Image>;
  imageDetailSourceList: Map<string, Array<ImageDetail>>;
  imageDetailSelectList: Map<string, ImageDetail>;
  imageTagNotReadyList: Map<string, boolean>;
  newImageIndex: number;

  constructor(protected injector: Injector) {
    super(injector);
    this.imageSelectList = Array<Image>();
    this.imageDetailSelectList = new Map<string, ImageDetail>();
    this.imageDetailSourceList = new Map<string, Array<ImageDetail>>();
    this.imageTagNotReadyList = new Map<string, boolean>();
  }

  ngOnInit() {
    this.k8sService.getServiceConfig(this.stepPhase).then(res => {
      this.uiBaseData = res;
      this.uiData.imageList.forEach((image: ImageIndex) => {
        this.imageSelectList.push({image_name: image.image_name, image_comment: "", image_deleted: 0});
        this.setImageDetailList(image.image_name, image.image_tag);
      })
    }).catch(err => this.messageService.dispatchError(err));
    this.k8sService.getImages("", 0, 0)
      .then(res => {
        this.imageSourceList = res;
        this.unshiftCustomerCreateImage();
      })
      .catch(err => this.messageService.dispatchError(err));
  }

  get stepPhase(): ServiceStepPhase {
    return PHASE_SELECT_IMAGES;
  }

  get uiData(): UIServiceStep2 {
    return this.uiBaseData as UIServiceStep2;
  }

  get isCanNextStep(): boolean {
    let hasSelectImage = this.imageSelectList.filter(value => {
        return value.image_name != "SERVICE.STEP_2_SELECT_IMAGE";
      }).length == this.imageSelectList.length;
    return hasSelectImage && this.imageDetailSelectList.size > 0;
  }

  get selfObject() {
    return this;
  }

  get projectName(): string {
    return this.uiData.projectName;
  }

  get projectId(): number {
    return this.uiData.projectId;
  }

  onBuildImageCompleted(imageName: string) {
    this.k8sService.getImages("", 0, 0).then(res => {
      res.forEach(value => {
        if (value.image_name == imageName) {
          this.imageSourceList = res;
          this.unshiftCustomerCreateImage();
          this.imageSelectList[this.newImageIndex] = value;
          this.setImageDetailList(value.image_name);
        }
      });
    }).catch(err => this.messageService.dispatchError(err));
  }

  unshiftCustomerCreateImage() {
    let customerCreateImage: Image = new Image();
    customerCreateImage.image_name = "SERVICE.STEP_2_CREATE_IMAGE";
    customerCreateImage["isSpecial"] = true;
    customerCreateImage["OnlyClick"] = true;
    this.imageSourceList.unshift(customerCreateImage);
  }

  setImageDetailList(imageName: string, imageTag?: string): void {
    this.imageTagNotReadyList.set(imageName, false);
    this.k8sService.getImageDetailList(imageName).then((res: ImageDetail[]) => {
      if (res && res.length > 0) {
        for (let item of res) {
          item['image_detail'] = JSON.parse(item['image_detail']);
          item['image_size_number'] = Number.parseFloat((item['image_size_number'] / (1024 * 1024)).toFixed(2));
          item['image_size_unit'] = 'MB';
        }
        this.imageDetailSourceList.set(res[0].image_name, res);
        if (imageTag) {
          let detail = res.find(value => value.image_tag == imageTag);
          this.imageDetailSelectList.set(res[0].image_name, detail);
        } else {
          this.imageDetailSelectList.set(res[0].image_name, res[0]);
        }
      } else {
        this.imageTagNotReadyList.set(imageName, true);
      }
    }).catch(err => this.messageService.dispatchError(err));
  }

  canChangeSelectImage(image: Image) {
    if (this.imageSelectList.find(value => value.image_name == image.image_name)) {
      let m: Message = new Message();
      m.message = "IMAGE.CREATE_IMAGE_EXIST";
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

  forward(): void {
    this.uiData.imageList.splice(0, this.uiData.imageList.length);//empty list
    this.imageSelectList.forEach((image: Image) => {
      let selectedImage = this.uiData.imageList.find((imageIndex: ImageIndex) => imageIndex.image_name == image.image_name);
      if (selectedImage) {
        selectedImage.image_tag = this.imageDetailSelectList.get(selectedImage.image_name).image_tag
      } else if (image.image_name != "SERVICE.STEP_2_SELECT_IMAGE") {
        let newImageIndex = new ImageIndex();
        newImageIndex.image_name = image.image_name;
        newImageIndex.image_tag = this.imageDetailSelectList.get(image.image_name).image_tag;
        this.uiData.imageList.push(newImageIndex);
      }
    });
    this.k8sService.setServiceConfig(this.uiData.uiToServer()).then(res => {
      this.k8sService.stepSource.next({index: 3, isBack: false});
    });
  }
}