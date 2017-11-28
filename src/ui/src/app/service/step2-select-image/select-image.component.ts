import { Component, OnInit, Injector } from '@angular/core';
import { Container } from '../service-step.component';
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
    this.k8sService.getServiceConfig(this.newServiceId, this.outputData).then(res => {
      this.outputData = res;
      this.containerList.forEach((container: Container) => {
        let index = container.image.indexOf(":");
        let imageName = container.image.slice(0, index);
        let imageTag = container.image.slice(index + 1);
        this.imageSelectList.push({image_name: imageName, image_comment: "", image_deleted: 0});
        this.setImageDetailList(imageName, imageTag);
      })
    });
    this.k8sService.getImages("", 0, 0)
      .then(res => {
        this.imageSourceList = res;
        this.unshiftCustomerCreateImage();
      })
      .catch(err => this.messageService.dispatchError(err));
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
    return this.outputData.projectinfo.project_name;
  }

  get projectId(): number {
    return this.outputData.projectinfo.project_id;
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
    this.imageSelectList.forEach((image: Image) => {
      let outValue = this.containerList.find((container: Container) => {
        return container.image.startsWith(image.image_name);
      });
      if (!outValue && image.image_name != "SERVICE.STEP_2_SELECT_IMAGE") {
        let newContainer = new Container();
        let firstIndex = image.image_name.indexOf("/");
        let imageTag = this.imageDetailSelectList.get(image.image_name).image_tag;
        newContainer.name = image.image_name.slice(firstIndex + 1, image.image_name.length);
        newContainer.image = image.image_name + ":" + imageTag;
        this.containerList.push(newContainer);
      }
    });
    this.k8sService.setServiceConfig(this.outputData).then(res => {
      this.k8sService.stepSource.next({index: 3, isBack: false});
    });
  }
}