import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { ServiceStep1Output, ServiceStep2Output, ServiceStepComponent } from '../service-step.component';
import { K8sService } from '../service.k8s';
import { MessageService } from "../../shared/message-service/message.service";
import { Image, ImageDetail } from "../../image/image";
import { AppInitService } from "../../app.init.service";

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
  outputData: ServiceStep2Output = new ServiceStep2Output();

  constructor(private k8sService: K8sService,
              private messageService: MessageService,
              private appInitService: AppInitService) {
    this.imageSelectList = Array<Image>();
    this.imageDetailSelectList = new Map<string, ImageDetail>();
    this.imageDetailSourceList = new Map<string, Array<ImageDetail>>();
  }

  ngOnInit() {
    let step1Out: ServiceStep1Output = this.k8sService.getStepData(1) as ServiceStep1Output;
    this.outputData.image_dockerfile.image_author = this.appInitService.currentUser["user_name"];
    this.outputData.project_name = step1Out.project_name;
    this.outputData.image_template = "dockerfile-template";
    this.k8sService.getImages("", 0, 0).then(res => {
      if (res.length > 0) {
        this.imageSourceList = res;
        this.imageSelectList.push(res[0]);
        this.setImageDetailList(res[0].image_name);
      }
    }).catch(err => this.messageService.dispatchError(err));
    this.intervalAutoRefreshImageList = setInterval(() => {
      if (this.isNeedAutoRefreshImageList) {
        this.k8sService.getImages("", 0, 0).then(res => {
          res.forEach(value => {
            let newImageName = `${this.outputData.project_name}/${this.outputData.image_name}`;
            if (value.image_name == newImageName) {
              this.isNeedAutoRefreshImageList = false;
              this.imageSourceList = res;
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

  changeSelectImageDetail(imageName: string, imageDetail: ImageDetail) {
    this.imageDetailSelectList.set(imageName, imageDetail);
  }

  modifySelectImage(index: number) {
    if (index == this.imageSelectList.length - 1) {
      this.imageSelectList.push(this.imageSourceList[0]);
    } else {
      this.imageSelectList.splice(index, 1);
    }
  }

  buildImage() {
    this.isNeedAutoRefreshImageList = true;
    console.log(this.outputData);
    this.outputData.image_dockerfile.image_volume = this.outputData.image_dockerfile.image_volume
      .filter(value => {
        return value != ""
      });
    this.outputData.image_dockerfile.image_run = this.outputData.image_dockerfile.image_run
      .filter(value => {
        return value != ""
      });
    console.log(this.outputData);
    this.k8sService.buildImage(this.outputData).then(res => {
      //show log...
    }).catch((err) => {
      this.messageService.dispatchError(err);
      this.isNeedAutoRefreshImageList = false;
    })
  }

  get imageRun(): Array<string> {
    return this.outputData.image_dockerfile.image_run;
  }

  get imageVolume(): Array<string> {
    return this.outputData.image_dockerfile.image_volume;
  }

  isImageDetailExist(image: Image): boolean {
    return this.imageDetailSourceList.get(image.image_name) &&
      this.imageDetailSourceList.get(image.image_name).length > 0;
  }


  forward(): void {
    this.k8sService.stepSource.next(3);
  }
}