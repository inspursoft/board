/**
 * Created by liyanq on 21/11/2017.
 */

import { AfterContentChecked, Component, ElementRef, OnDestroy, OnInit, QueryList, ViewChild, ViewChildren } from "@angular/core"
import { CsInputArrayComponent } from "../../shared/cs-components-library/cs-input-array/cs-input-array.component";
import { BuildImageData, Image, ImageDetail } from "../image";
import { ImageService } from "../image-service/image-service";
import { MessageService } from "../../shared/message-service/message.service";
import { HttpErrorResponse, HttpEvent, HttpEventType, HttpProgressEvent, HttpResponse } from "@angular/common/http"
import { AppInitService } from "../../app.init.service";
import { Subscription } from "rxjs/Subscription";
import { WebsocketService } from "../../shared/websocket-service/websocket.service";
import { EnvType } from "../../shared/environment-value/environment-value.component";
import { ValidationErrors } from "@angular/forms";
import { TranslateService } from "@ngx-translate/core";
import { CsModalChildBase } from "../../shared/cs-modal-base/cs-modal-child-base";
import { CreateImageMethod } from "../../shared/shared.types";

const AUTO_REFRESH_IMAGE_LIST: number = 2000;

/*declared in shared-module*/
@Component({
  selector: "create-image",
  templateUrl: "./image-create.component.html",
  styleUrls: ["./image-create.component.css"]
})
export class CreateImageComponent extends CsModalChildBase implements OnInit, OnDestroy {
  boardHost: string;
  @ViewChildren(CsInputArrayComponent) inputArrayComponents: QueryList<CsInputArrayComponent>;
  @ViewChild("areaStatus") areaStatus: ElementRef;
  imageBuildMethod: CreateImageMethod = CreateImageMethod.Template;
  isOpenEnvironment = false;
  patternNewImageName: RegExp = /^[a-z\d.-]+$/;
  patternNewImageTag: RegExp = /^[a-z\d.-]+$/;
  patternBaseImage: RegExp = /^[a-z\d/.:-]+$/;
  patternExpose: RegExp = /^[\d-\s\w/\\]+$/;
  patternVolume: RegExp = /.+/;
  patternRun: RegExp = /.+/;
  patternEntryPoint: RegExp = /.+/;
  patternCopyPath: RegExp = /.+/;
  imageTemplateList: Array<Object> = [{name: "Docker File Template"}];
  filesList: Map<string, Array<{path: string, file_name: string, size: number}>>;
  selectFromImportFile: File;
  intervalAutoRefreshImageList: any;
  isNeedAutoRefreshImageList: boolean = false;
  isBuildImageWIP: boolean = false;
  isServerHaveDockerFile: boolean = false;
  isUploadFileWIP = false;
  customerNewImage: BuildImageData;
  consoleText: string = "";
  uploadCopyToPath: string = "/tmp";
  uploadProgressValue: HttpProgressEvent;
  imageList: Array<Image>;
  imageDetailList: Array<ImageDetail>;
  selectedImage: Image;
  baseImageSource: number = 1;
  boardRegistry: string = "";
  processImageSubscription: Subscription;
  cancelButtonDisable = true;
  cancelInfo: {isShow: boolean, isForce: boolean, title: string, message: string};

  constructor(private imageService: ImageService,
              private messageService: MessageService,
              private webSocketService: WebsocketService,
              private translateService: TranslateService,
              private appInitService: AppInitService) {
    super();
    this.filesList = new Map<string, Array<{path: string, file_name: string, size: number}>>();
    this.boardHost = this.appInitService.systemInfo['board_host'];
    this.imageList = Array<Image>();
    this.imageDetailList = Array<ImageDetail>();
    this.cancelInfo = {isShow: false, isForce: false, title: "", message: ""};
  }

  ngOnInit() {
    this.intervalAutoRefreshImageList = setInterval(() => {
      if (this.isNeedAutoRefreshImageList && this.isBuildImageWIP) {
        this.imageService.getImages(this.customerNewImage.image_name, 0, 0).then(res => {
          res.forEach(value => {
            let newImageName = `${this.customerNewImage.project_name}/${this.customerNewImage.image_name}`;
            if (value.image_name == newImageName) {
              this.isNeedAutoRefreshImageList = false;
              this.closeNotification.next(newImageName);
              this.modalOpened = false;
            }
          });
        }).catch(() => this.modalOpened = false);
      }
    }, AUTO_REFRESH_IMAGE_LIST);
    this.imageService.getImages("", 0, 0)
      .then(res => this.imageList = res || [])
      .catch(() => {
        this.modalOpened = false;
        this.messageService.showAlert('IMAGE.CREATE_IMAGE_UPDATE_IMAGE_LIST_FAILED', {alertType: 'alert-danger', view: this.alertView});
      });
  }

  ngOnDestroy() {
    if (this.processImageSubscription) {
      this.processImageSubscription.unsubscribe();
    }
    clearInterval(this.intervalAutoRefreshImageList);
  }

  public initCustomerNewImage(projectId: number, projectName: string): void {
    this.customerNewImage = new BuildImageData();
    this.customerNewImage.image_dockerfile.image_author = this.appInitService.currentUser["user_name"];
    this.customerNewImage.project_id = projectId;
    this.customerNewImage.project_name = projectName;
    this.customerNewImage.image_template = "dockerfile-template";
    this.imageService.deleteImageConfig(projectName).subscribe();
  }

  public initBuildMethod(method: CreateImageMethod): void {
    this.imageBuildMethod = method;
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

  get isBuildDisabled() {
    return this.isBuildImageWIP || this.isUploadFileWIP;
  }

  get checkImageTagFun() {
    return this.checkImageTag.bind(this);
  }

  get checkImageNameFun() {
    return this.checkImageName.bind(this);
  }

  get cancelCaption(){
    return this.consoleText == "IMAGE.CREATE_IMAGE_JENKINS_PREPARE" ?
      "IMAGE.CREATE_IMAGE_CANCEL_WAIT":
      "IMAGE.CREATE_IMAGE_BUILD_CANCEL";
  }

  checkImageTag(control: HTMLInputElement): Promise<ValidationErrors> {
    if (this.customerNewImage.image_name == "") {
      return Promise.resolve(null);
    }
    return this.imageService.checkImageExist(this.customerNewImage.image_name, this.customerNewImage.image_name, control.value)
      .then(() => null)
      .catch((err: HttpErrorResponse) => {
        if (err.status == 409) {
          this.messageService.cleanNotification();
          return {imageTagExist: "IMAGE.CREATE_IMAGE_TAG_EXIST"}
        } else if (err.status == 404) {
          this.messageService.cleanNotification();
        } else {
          this.modalOpened = false;
        }
      });
  }

  checkImageName(control: HTMLInputElement): Promise<ValidationErrors> {
    if (this.customerNewImage.image_tag == "") {
      return Promise.resolve(null);
    }
    return this.imageService.checkImageExist(this.customerNewImage.image_name, control.value, this.customerNewImage.image_tag)
      .then(() => null)
      .catch((err: HttpErrorResponse) => {
        if (err.status == 409) {
          this.messageService.cleanNotification();
          return {imageNameExist: "IMAGE.CREATE_IMAGE_NAME_EXIST"}
        } else if (err.status == 404) {
          this.messageService.cleanNotification();
        } else {
          this.modalOpened = false;
        }
      });
  }

  cancelBuildImage() {
    if (this.consoleText == "IMAGE.CREATE_IMAGE_JENKINS_PREPARE") {
      this.cancelInfo.isForce = true;
      this.cancelInfo.title = "IMAGE.CREATE_IMAGE_FORCE_QUIT";
      this.cancelInfo.message = "IMAGE.CREATE_IMAGE_FORCE_QUIT_MSG";
    } else {
      this.cancelInfo.isForce = false;
      this.cancelInfo.title = "IMAGE.CREATE_IMAGE_BUILD_CANCEL";
      this.cancelInfo.message = "IMAGE.CREATE_IMAGE_BUILD_CANCEL_MSG";
    }
    this.cancelInfo.isShow = true;
  }

  cancelBuildImageBehavior() {
    this.cancelInfo.isShow = false;
    if (this.cancelInfo.isForce) {
      this.modalOpened = false;
    } else {
      this.imageService.cancelConsole(this.customerNewImage.image_name).subscribe(
        () => this.cleanImageConfig(),
        () => this.modalOpened = false);
    }
  }

  uploadDockerFile(): Promise<boolean> {
    if (this.selectFromImportFile) {
      this.isUploadFileWIP = true;
      let formData: FormData = new FormData();
      formData.append("upload_file", this.selectFromImportFile, this.selectFromImportFile.name);
      formData.append("project_name", this.customerNewImage.project_name);
      formData.append("image_name", this.customerNewImage.image_name);
      formData.append("image_tag", this.customerNewImage.image_tag);
      return this.imageService.uploadDockerFile(formData)
        .then(() => {
          this.isUploadFileWIP = false;
          return true;
        })
    }
  }

  async buildImageByDockerFile(): Promise<any> {
    let fileInfo = {
      imageName: this.customerNewImage.image_name,
      tagName: this.customerNewImage.image_tag,
      projectName: this.customerNewImage.project_name
    };
    if (this.isServerHaveDockerFile) {
      return this.imageService.buildImageFromDockerFile(fileInfo);
    } else {
      if (await this.uploadDockerFile()) {
        return this.imageService.buildImageFromDockerFile(fileInfo);
      } else {
        return Promise.reject(false);
      }
    }
  }

  cleanImageConfig(err?: HttpErrorResponse) {
    this.isBuildImageWIP = false;
    this.isUploadFileWIP = false;
    this.isNeedAutoRefreshImageList = false;
    if (err) {
      let reason = err ? ((err as HttpErrorResponse).error as Error).message : "";
      this.translateService.get(`IMAGE.CREATE_IMAGE_BUILD_IMAGE_FAILED`).subscribe((msg: string) => {
        this.messageService.showAlert(`${msg}:${reason}`, {alertType: 'alert-danger', view: this.alertView});
      });
    }
    this.imageService.deleteImageConfig(this.customerNewImage.project_name).subscribe();
    this.updateFileList().then();
  }

  buildImageResole() {
    this.processImageSubscription = this.webSocketService
      .connect(`ws://${this.boardHost}/api/v1/jenkins-job/console?job_name=${this.customerNewImage.project_name}&token=${this.appInitService.token}`)
      .subscribe((obs: MessageEvent) => {
        this.consoleText = <string>obs.data;
        this.cancelButtonDisable = false;
        this.areaStatus.nativeElement.scrollTop = this.areaStatus.nativeElement.scrollHeight;
        let consoleTextArr: Array<string> = this.consoleText.split(/[\n]/g);
        if (consoleTextArr.find(value => value.indexOf("Finished: SUCCESS") > -1)) {
          this.isNeedAutoRefreshImageList = true;
          this.processImageSubscription.unsubscribe();
        }
        if (consoleTextArr.find(value => value.indexOf("Finished: FAILURE") > -1)) {
          this.isBuildImageWIP = false;
          this.isUploadFileWIP = false;
          this.cancelButtonDisable = true;
          this.isNeedAutoRefreshImageList = false;
          this.appInitService.setAuditLog({
            operation_user_id: this.appInitService.currentUser["user_id"],
            operation_user_name: this.appInitService.currentUser["user_name"],
            operation_project_id: this.customerNewImage.project_id,
            operation_project_name: this.customerNewImage.project_name,
            operation_object_type: "images",
            operation_object_name: "",
            operation_action: "create",
            operation_status: "Failed"
          }).subscribe();
          this.processImageSubscription.unsubscribe();
        }
      }, err => err, () => this.modalOpened = false);
  }

  buildImage() {
    let buildImageInit = () => {
      this.cancelButtonDisable = true;
      this.isBuildImageWIP = true;
      this.consoleText = "IMAGE.CREATE_IMAGE_JENKINS_PREPARE";
      setTimeout(() => this.cancelButtonDisable = false, 10000);
    };
    if (this.imageBuildMethod == CreateImageMethod.Template) {
      if (this.verifyInputValid() &&
        this.verifyInputArrayValid() &&
        this.verifyDropdownValid() &&
        this.customerNewImage.image_dockerfile.image_base != "") {
        buildImageInit();
        this.imageService.buildImageFromTemp(this.customerNewImage)
          .then(this.buildImageResole.bind(this))
          .catch(this.cleanImageConfig.bind(this));
      }
    } else if (this.verifyInputValid()) {
      if (this.selectFromImportFile) {
        buildImageInit();
        this.buildImageByDockerFile()
          .then(this.buildImageResole.bind(this))
          .catch(this.cleanImageConfig.bind(this));
      } else {
        this.messageService.showAlert('IMAGE.CREATE_IMAGE_SELECT_DOCKER_FILE', {alertType: 'alert-warning', view: this.alertView});
      }
    }
  }

  updateFileList(): Promise<any> {
    this.filesList.clear();
    let formFileList: FormData = new FormData();
    formFileList.append('project_name', this.customerNewImage.project_name);
    formFileList.append('image_name', this.customerNewImage.image_name);
    formFileList.append('image_tag', this.customerNewImage.image_tag);
    return this.imageService.getFileList(formFileList).then(res => {
      this.filesList.set(this.customerNewImage.image_name, res);
      let imageCopyArr = this.customerNewImage.image_dockerfile.image_copy;
      imageCopyArr.splice(0, imageCopyArr.length);
      this.filesList.get(this.customerNewImage.image_name).forEach(value => {
        imageCopyArr.push({
          dockerfile_copyfrom: value.file_name,
          dockerfile_copyto: this.uploadCopyToPath + "/" + value.file_name,
        });
      });
    }).catch((err: HttpErrorResponse) => {
      if (err.status == 401) {
        this.modalOpened = false;
      } else {
        this.messageService.showAlert('IMAGE.CREATE_IMAGE_UPDATE_IMAGE_LIST_FAILED', {alertType: 'alert-danger', view: this.alertView});
      }
    });
  }

  updateFileListAndPreviewInfo() {
    this.updateFileList().then(() => {
      this.getDockerFilePreviewInfo();
    });
  }

  selectDockerFile(event: Event) {
    let fileList: FileList = (event.target as HTMLInputElement).files;
    if (fileList.length > 0) {
      let file:File = fileList[0];
      if (file.name !== "Dockerfile"){
        (event.target as HTMLInputElement).value = "";
        this.selectFromImportFile = null;
        this.messageService.showAlert('IMAGE.CREATE_IMAGE_FILE_NAME_ERROR', {alertType: 'alert-danger', view: this.alertView});
      } else {
        this.selectFromImportFile = file;
        let reader = new FileReader();
        reader.onload = (ev: ProgressEvent) => {
          this.consoleText = (ev.target as FileReader).result;
        };
        reader.readAsText(this.selectFromImportFile);
      }

    }
  }

  downloadDockerFile(): void {
    this.selectFromImportFile = null;
    this.consoleText = "";
    this.isServerHaveDockerFile = false;
    if (this.customerNewImage.image_name && this.customerNewImage.image_tag) {
      let downloadInfo = {
        imageName: this.customerNewImage.image_name,
        tagName: this.customerNewImage.image_tag,
        projectName: this.customerNewImage.project_name
      };
      this.imageService.downloadDockerFile(downloadInfo)
        .then((res: HttpResponse<string>) => {
          this.consoleText = res.body;
          this.isServerHaveDockerFile = true;
        })
        .catch(()=>this.messageService.cleanNotification());
    }
  }

  uploadFile(event: Event) {
    let fileList: FileList = (event.target as HTMLInputElement).files;
    if (fileList.length > 0) {
      let file: File = fileList[0];
      if (file.size > 1024 * 1024 * 500) {
        (event.target as HTMLInputElement).value = "";
        this.messageService.showAlert('IMAGE.CREATE_IMAGE_UPDATE_FILE_SIZE', {alertType: 'alert-danger', view: this.alertView});
      } else {
        let formData: FormData = new FormData();
        this.isUploadFileWIP = true;
        formData.append('upload_file', file, file.name);
        formData.append('project_name', this.customerNewImage.project_name);
        formData.append('image_name', this.customerNewImage.image_name);
        formData.append('image_tag', this.customerNewImage.image_tag);
        this.imageService.uploadFile(formData).subscribe((res: HttpEvent<Object>) => {
          if (res.type == HttpEventType.UploadProgress) {
            this.uploadProgressValue = res;
          } else if (res.type == HttpEventType.Response) {
            (event.target as HTMLInputElement).value = "";
            this.messageService.showAlert('IMAGE.CREATE_IMAGE_UPLOAD_SUCCESS', {view: this.alertView});
            this.isUploadFileWIP = false;
            this.updateFileListAndPreviewInfo();
          }
        }, (error: HttpErrorResponse) => {
          this.isUploadFileWIP = false;
          if (error.status == 401) {
            this.modalOpened = false;
          } else {
            (event.target as HTMLInputElement).value = "";
            let newImageErrReason = (error.error as Error).message;
            this.translateService.get('IMAGE.CREATE_IMAGE_UPLOAD_FAILED').subscribe((msg: string) => {
              this.messageService.showAlert(`${msg}:${newImageErrReason}`, {alertType: 'alert-danger', view: this.alertView});
            });
          }
        });
      }
    }
  }

  getDockerFilePreviewInfo() {
    if (this.customerNewImage.image_dockerfile.image_base != "") {
      this.imageService.getDockerFilePreview(this.customerNewImage)
        .then(res => this.consoleText = res)
        .catch((err: HttpErrorResponse) => {
          if (err.status == 401) {
            this.modalOpened = false;
        } else {
            this.messageService.showAlert('IMAGE.CREATE_IMAGE_UPDATE_DOCKER_FILE_FAILED', {
              alertType: 'alert-danger',
              view: this.alertView
            });
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

  removeFile(file: {path: string, file_name: string, size: number}) {
    let fromRemoveData: FormData = new FormData();
    fromRemoveData.append("project_name", this.customerNewImage.project_name);
    fromRemoveData.append("image_name", this.customerNewImage.image_name);
    fromRemoveData.append("image_tag", this.customerNewImage.image_tag);
    fromRemoveData.append("file_name", file.file_name);
    this.imageService.removeFile(fromRemoveData).subscribe(
      () => this.messageService.showAlert('IMAGE.CREATE_IMAGE_REMOVE_FILE_SUCCESS', {view: this.alertView}),
      (err: HttpErrorResponse) => {
        if (err.status == 401) {
          this.modalOpened = false;
        } else {
          this.messageService.showAlert('IMAGE.CREATE_IMAGE_REMOVE_FILE_FAILED', {alertType: 'alert-danger', view: this.alertView});
        }
      },
      () => this.updateFileListAndPreviewInfo());
  }

  resetBuildMethod(method: CreateImageMethod) {
    this.imageBuildMethod = method;
    this.consoleText = "";
    if (method == CreateImageMethod.Template) {
      this.selectFromImportFile = null;
    }
  }

  cleanBaseImageInfo(isGetBoardRegistry: boolean = false): void {
    if ((this.baseImageSource == 1 && isGetBoardRegistry) ||
      (this.baseImageSource == 2 && !isGetBoardRegistry)) {
      this.selectedImage = null;
      this.consoleText = "";
      this.imageDetailList.splice(0, this.imageDetailList.length);
      this.customerNewImage.image_dockerfile.image_base = "";
    }
  }

  setBaseImage(selectImage: Image): void {
    this.selectedImage = null;
    this.imageDetailList = null;
    this.imageService.getBoardRegistry().subscribe((res: string) => {
      this.boardRegistry = res.replace(/"/g,"");
      this.imageService.getImageDetailList(selectImage.image_name)
        .then((res: ImageDetail[]) => {
          this.selectedImage = selectImage;
          this.imageDetailList = res;
          this.customerNewImage.image_dockerfile.image_base = `${this.boardRegistry}/${this.selectedImage.image_name}:${res[0].image_tag}`;
          this.getDockerFilePreviewInfo();
        })
        .catch(() => this.modalOpened = false);
    });
  }

  setBaseImageDetail(detail: ImageDetail): void {
    this.imageService.getBoardRegistry().subscribe((res: string) => {
      this.boardRegistry = res.replace(/"/g,"");
      this.customerNewImage.image_dockerfile.image_base = `${this.boardRegistry}/${this.selectedImage.image_name}:${detail.image_tag}`;
      this.getDockerFilePreviewInfo();
    });
  }
}