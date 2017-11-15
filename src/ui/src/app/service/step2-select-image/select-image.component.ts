import { Component, OnDestroy, OnInit, QueryList, ViewChildren, AfterContentChecked, Injector } from '@angular/core';
import { ServiceStep2NewImageType, Container } from '../service-step.component';
import { Image, ImageDetail } from "../../image/image";
import { Message } from "../../shared/message-service/message";
import { EnvType } from "../environment-value/environment-value.component";
import { CsInputArrayComponent } from "../cs-input-array/cs-input-array.component";
import { CsInputComponent } from "../cs-input/cs-input.component";
import { WebsocketService } from "../../shared/websocket-service/websocket.service";
import { Subscription } from "rxjs/Subscription";
import { ServiceStepBase } from "../service-step";

enum ImageSource{fromBoardRegistry, fromDockerHub}
const AUTO_REFRESH_IMAGE_LIST: number = 2000;
// const PROCESS_IMAGE_CONSOLE_URL = `ws://10.165.22.61:8088/api/v1/jenkins-job/console?job_name=process_image`;
const PROCESS_IMAGE_CONSOLE_URL = `ws://localhost/api/v1/jenkins-job/console?job_name=process_image`;
type alertType = "alert-info" | "alert-danger";
@Component({
  templateUrl: './select-image.component.html',
  styleUrls: ["./select-image.component.css"]
})
export class SelectImageComponent extends ServiceStepBase implements OnInit, OnDestroy, AfterContentChecked {
  @ViewChildren(CsInputArrayComponent) inputArrayComponents: QueryList<CsInputArrayComponent>;
  @ViewChildren(CsInputComponent) inputComponents: QueryList<CsInputComponent>;
  patternNewImageName: RegExp = /^[a-z\d.-]+$/;
  patternNewImageTag: RegExp = /^[a-z\d.-]+$/;
  patternBaseImage: RegExp = /^[a-z\d.:-]+$/;
  patternExpose: RegExp = /^[\d-\s\w/\\]+$/;
  patternVolume: RegExp = /.+/;
  patternRun: RegExp = /.+/;
  patternEntryPoint: RegExp = /.+/;
  _isOpenEnvironment = false;
  _isOpenNewImage = false;
  intervalAutoRefreshImageList: any;
  isNeedAutoRefreshImageList: boolean = false;
  imageInBuilding: boolean = false;
  isInputComponentsValid: boolean = false;
  autoRefreshTimesCount: number = 0;
  imageSource: ImageSource = ImageSource.fromBoardRegistry;
  imageSourceList: Array<Image>;
  imageSelectList: Array<Image>;
  imageDetailSourceList: Map<string, Array<ImageDetail>>;
  imageDetailSelectList: Map<string, ImageDetail>;
  imageTagNotReadyList: Map<string, boolean>;
  imageTemplateList: Array<Object> = [{name: "Docker File Template"}];
  customerNewImage: ServiceStep2NewImageType;
  filesList: Map<string, Array<{path: string, file_name: string, size: number}>>;
  consoleText: string = "";
  newImageErrMessage: string = "";
  newImageAlertType: alertType = "alert-danger";
  isNewImageAlertOpen: boolean = false;
  isUploadFileIng = false;
  newImageIndex: number;
  lastJobNumber: number = 0;
  processImageSubscription: Subscription;

  constructor(private webSocketService: WebsocketService, protected injector: Injector) {
    super(injector);
    this.imageSelectList = Array<Image>();
    this.imageDetailSelectList = new Map<string, ImageDetail>();
    this.imageDetailSourceList = new Map<string, Array<ImageDetail>>();
    this.imageTagNotReadyList = new Map<string, boolean>();
    this.filesList = new Map<string, Array<{path: string, file_name: string, size: number}>>();
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
    this.intervalAutoRefreshImageList = setInterval(() => {
      if (this.isNeedAutoRefreshImageList && this.imageInBuilding) {
        this.autoRefreshTimesCount++;
        this.isNewImageAlertOpen = false;
        this.k8sService.getImages("", 0, 0).then(res => {
          res.forEach(value => {
            let newImageName = `${this.customerNewImage.project_name}/${this.customerNewImage.image_name}`;
            if (value.image_name == newImageName) {
              this.isNeedAutoRefreshImageList = false;
              this.imageSourceList = Object.create(res);
              this.unshiftCustomerCreateImage();
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
            this.imageInBuilding = false;
            this.newImageAlertType = "alert-danger";
            this.newImageErrMessage = "SERVICE.STEP_2_UPDATE_IMAGE_LIST_FAILED";
            this.isNewImageAlertOpen = true;
          }
        });
      }
    }, AUTO_REFRESH_IMAGE_LIST);
  }

  ngOnDestroy() {
    if (this.processImageSubscription) {
      this.processImageSubscription.unsubscribe();
    }
    clearInterval(this.intervalAutoRefreshImageList);
  }

  ngAfterContentChecked() {
    this.isInputComponentsValid = true;
    if (this.inputArrayComponents) {
      this.inputArrayComponents.forEach(item => {
        if (!item.valid) {
          this.isInputComponentsValid = false;
        }
      });
    }
    if (this.isInputComponentsValid && this.inputComponents) {
      this.inputComponents.forEach(item => {
        if (!item.valid) {
          this.isInputComponentsValid = false;
        }
      });
    }
  }

  get isOpenNewImage(): boolean {
    return this._isOpenNewImage;
  }

  set isOpenNewImage(value: boolean) {
    this._isOpenNewImage = value;
    if (value) {
      this.resetNewImage();
    }
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
      }).length == this.imageSelectList.length;
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

  resetNewImage() {
    this.consoleText = "";
    this.isNewImageAlertOpen = false;
    this.isNeedAutoRefreshImageList = false;
    this.imageInBuilding = false;
  }

  unshiftCustomerCreateImage() {
    let customerCreateImage: Image = new Image();
    customerCreateImage.image_name = "SERVICE.STEP_2_CREATE_IMAGE";
    customerCreateImage["isSpecial"] = true;
    customerCreateImage["OnlyClick"] = true;
    this.imageSourceList.unshift(customerCreateImage);
  }

  shieldEnter($event: KeyboardEvent) {
    if ($event.charCode == 13) {
      (<any>$event.target).blur();
      this.getDockerFilePreviewInfo();
    }
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
    this.customerNewImage = new ServiceStep2NewImageType();
    this.customerNewImage.image_dockerfile.image_author = this.appInitService.currentUser["user_name"];
    this.customerNewImage.project_id = this.outputData.projectinfo.project_id;
    this.customerNewImage.project_name = this.outputData.projectinfo.project_name;
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

  async buildImage() {
    this.isNewImageAlertOpen = false;
    this.imageInBuilding = true;
    this.lastJobNumber = 0;
    this.consoleText = "Jenkins preparing...";
    this.k8sService.buildImage(this.customerNewImage)
      .then(() => {
        setTimeout(() => {
          this.processImageSubscription = this.webSocketService
            .connect(PROCESS_IMAGE_CONSOLE_URL + `&token=${this.appInitService.token}`)
            .subscribe(obs => {
              this.consoleText = <string>obs.data;
              if (this.lastJobNumber == 0) {
                this.k8sService.getLastJobId("process_image").then(res => {
                  this.lastJobNumber = res;
                });
              }
              let consoleTextArr: Array<string> = this.consoleText.split(/[\n]/g);
              if (consoleTextArr.find(value => value.indexOf("Finished: SUCCESS") > -1)) {
                this.isNeedAutoRefreshImageList = true;
                this.autoRefreshTimesCount = 0;
                this.processImageSubscription.unsubscribe();
              }
              if (consoleTextArr.find(value => value.indexOf("Finished: FAILURE") > -1)) {
                this.imageInBuilding = false;
                this.isNeedAutoRefreshImageList = false;
                this.newImageAlertType = "alert-danger";
                this.newImageErrMessage = "SERVICE.STEP_2_BUILD_IMAGE_FAILED";
                this.isNewImageAlertOpen = true;
                this.processImageSubscription.unsubscribe();
              }
            }, err => err, () => {
              this.isOpenNewImage = false;
            });
        }, 10000);
      })
      .catch((err) => {
        this.imageInBuilding = false;
        this.isNeedAutoRefreshImageList = false;
        if (err && err.status == 401) {
          this.isOpenNewImage = false;
          this.messageService.dispatchError(err);
        } else {
          this.newImageAlertType = "alert-danger";
          this.newImageErrMessage = "SERVICE.STEP_2_BUILD_IMAGE_FAILED";
          this.isNewImageAlertOpen = true;
        }
      });
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
      this.isUploadFileIng = true;
      let file: File = fileList[0];
      let formData: FormData = new FormData();
      formData.append('upload_file', file, file.name);
      formData.append('project_name', this.customerNewImage.project_name);
      formData.append('image_name', this.customerNewImage.image_name);
      formData.append('tag_name', this.customerNewImage.image_tag);
      this.k8sService.uploadFile(formData).then(() => {
        event.target.value = "";
        this.newImageAlertType = "alert-info";
        this.newImageErrMessage = "SERVICE.STEP_2_UPLOAD_SUCCESS";
        this.isNewImageAlertOpen = true;
        this.isUploadFileIng = false;
        this.asyncGetDockerFilePreviewInfo();
      }).catch(err => {
        if (err && err.status == 401) {
          this.isOpenNewImage = false;
          this.messageService.dispatchError(err);
        } else {
          this.newImageAlertType = "alert-danger";
          this.newImageErrMessage = "SERVICE.STEP_2_UPLOAD_FAILED";
          this.isNewImageAlertOpen = true;
          this.isUploadFileIng = false;
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
        dockerfile_envvalue: value.envValue,
      })
    });
    this.getDockerFilePreviewInfo();
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

  cancelBuildImage() {
    if (this.lastJobNumber > 0) {
      this.k8sService.cancelConsole("process_image", this.lastJobNumber).then(() => {
        this.isOpenNewImage = false;
      });
      this.lastJobNumber = -1;
    }
  }

  removeFile(file: {path: string, file_name: string, size: number}) {
    this.isNewImageAlertOpen = false;
    let fromRemoveData: FormData = new FormData();
    fromRemoveData.append("project_name", this.customerNewImage.project_name);
    fromRemoveData.append("image_name", this.customerNewImage.image_name);
    fromRemoveData.append("tag_name", this.customerNewImage.image_tag);
    fromRemoveData.append("file_name", file.file_name);
    this.k8sService.removeFile(fromRemoveData)
      .then(() => this.asyncGetDockerFilePreviewInfo())
      .catch(err => {
        if (err && err.status == 401) {
          this.isOpenNewImage = false;
          this.messageService.dispatchError(err);
        } else {
          this.newImageAlertType = "alert-danger";
          this.newImageErrMessage = "SERVICE.STEP_2_REMOVE_FILE_FAILED";
          this.isNewImageAlertOpen = true;
        }
      });
  }
}