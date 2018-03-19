/**
 * Created by liyanq on 21/11/2017.
 */

import { AfterContentChecked, Component, EventEmitter, Input, OnDestroy, OnInit, Output, QueryList, ViewChildren } from "@angular/core"
import { CsInputArrayComponent } from "../../shared/cs-components-library/cs-input-array/cs-input-array.component";
import { CsInputComponent } from "../../shared/cs-components-library/cs-input/cs-input.component";
import { BuildImageData, Image, ImageDetail } from "../image";
import { ImageService } from "../image-service/image-service";
import { MessageService } from "../../shared/message-service/message.service";
import { HttpErrorResponse, HttpEvent, HttpEventType, HttpProgressEvent, HttpResponse } from "@angular/common/http"
import { AppInitService } from "../../app.init.service";
import { Subscription } from "rxjs/Subscription";
import { WebsocketService } from "../../shared/websocket-service/websocket.service";
import { EnvType } from "../../shared/environment-value/environment-value.component";
import { ValidationErrors } from "@angular/forms";

enum ImageSource {fromBoardRegistry, fromDockerHub}

enum ImageBuildMethod {fromTemplate, fromImportFile}

const AUTO_REFRESH_IMAGE_LIST: number = 2000;

type alertType = "alert-info" | "alert-danger";

/*declared in shared-module*/
@Component({
  selector: "create-image",
  templateUrl: "./image-create.component.html",
  styleUrls: ["./image-create.component.css"]
})
export class CreateImageComponent implements OnInit, AfterContentChecked, OnDestroy {
  boardHost: string;
  _isOpen: boolean = false;
  @ViewChildren(CsInputArrayComponent) inputArrayComponents: QueryList<CsInputArrayComponent>;
  @ViewChildren(CsInputComponent) inputComponents: QueryList<CsInputComponent>;
  @Input() projectId: number = 0;
  @Input() projectName: string = "";
  @Input() imageBuildMethod: ImageBuildMethod = ImageBuildMethod.fromTemplate;
  @Output() onBuildCompleted: EventEmitter<string>;
  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();
  _isOpenEnvironment = false;
  patternNewImageName: RegExp = /^[a-z\d.-]+$/;
  patternNewImageTag: RegExp = /^[a-z\d.-]+$/;
  patternBaseImage: RegExp = /^[a-z\d/.:-]+$/;
  patternExpose: RegExp = /^[\d-\s\w/\\]+$/;
  patternVolume: RegExp = /.+/;
  patternRun: RegExp = /.+/;
  patternEntryPoint: RegExp = /.+/;
  patternCopyPath: RegExp = /.+/;
  imageSource: ImageSource = ImageSource.fromBoardRegistry;
  newImageAlertType: alertType = "alert-danger";
  imageTemplateList: Array<Object> = [{name: "Docker File Template"}];
  filesList: Map<string, Array<{path: string, file_name: string, size: number}>>;
  selectFromImportFile: File;
  intervalAutoRefreshImageList: any;
  isNeedAutoRefreshImageList: boolean = false;
  isBuildImageWIP: boolean = false;
  isInputComponentsValid: boolean = false;
  isServerHaveDockerFile: boolean = false;
  isUploadFileWIP = false;
  customerNewImage: BuildImageData;
  autoRefreshTimesCount: number = 0;
  isNewImageAlertOpen: boolean = false;
  newImageErrMessage: string = "";
  newImageErrReason: string = "";
  consoleText: string = "";
  toggleCancelBuilding: boolean = false;
  processImageSubscription: Subscription;
  uploadCopyToPath: string = "./tmp";
  uploadProgressValue: HttpProgressEvent;
  imageList: Array<Image>;
  imageDetailList: Array<ImageDetail>;
  selectedImage: Image;
  baseImageSource: number = 1;
  boardRegistry: string = "";

  constructor(private imageService: ImageService,
              private messageService: MessageService,
              private webSocketService: WebsocketService,
              private appInitService: AppInitService) {
    this.onBuildCompleted = new EventEmitter<string>();
    this.filesList = new Map<string, Array<{path: string, file_name: string, size: number}>>();
    this.boardHost = this.appInitService.systemInfo['board_host'];
    this.imageList = Array<Image>();
    this.imageDetailList = Array<ImageDetail>();
  }

  ngOnInit() {
    this.customerNewImage = new BuildImageData();
    this.toggleCancelBuilding = false;
    this.customerNewImage.image_dockerfile.image_author = this.appInitService.currentUser["user_name"];
    this.customerNewImage.project_id = this.projectId;
    this.customerNewImage.project_name = this.projectName;
    this.customerNewImage.image_template = "dockerfile-template";
    this.intervalAutoRefreshImageList = setInterval(() => {
      if (this.isNeedAutoRefreshImageList && this.isBuildImageWIP) {
        this.autoRefreshTimesCount++;
        this.isNewImageAlertOpen = false;
        this.imageService.getImages("", 0, 0).then(res => {
          res.forEach(value => {
            let newImageName = `${this.customerNewImage.project_name}/${this.customerNewImage.image_name}`;
            if (value.image_name == newImageName) {
              this.isNeedAutoRefreshImageList = false;
              this.onBuildCompleted.emit(newImageName);
              this.isOpen = false;
            }
          });
        }).catch(err => {
          if (err && err instanceof HttpErrorResponse && (err as HttpErrorResponse).status == 401) {
            this.isOpen = false;
            this.messageService.dispatchError(err);
          } else {
            this.isBuildImageWIP = false;
            this.newImageAlertType = "alert-danger";
            this.newImageErrMessage = "IMAGE.CREATE_IMAGE_UPDATE_IMAGE_LIST_FAILED";
            this.isNewImageAlertOpen = true;
          }
        });
      }
    }, AUTO_REFRESH_IMAGE_LIST);
    this.imageService.getImages("", 0, 0)
      .then(res => {
        this.imageList = res || [];
        if (this.imageList.length > 0) {
          this.selectedImage = this.imageList[0];
        }
      })
      .catch(err => {
        this.isOpen = false;
        this.messageService.dispatchError(err)
      });
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

  ngOnDestroy() {
    if (this.processImageSubscription) {
      this.processImageSubscription.unsubscribe();
    }
    clearInterval(this.intervalAutoRefreshImageList);
  }

  @Input() get isOpen() {
    return this._isOpen;
  }

  set isOpen(open: boolean) {
    this._isOpen = open;
    this.isOpenChange.emit(this._isOpen);
  }

  set isOpenEnvironment(value) {
    // this.isOpen = !value;
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
    let baseDisabled = this.isBuildImageWIP ||
      !this.isInputComponentsValid ||
      this.isUploadFileWIP;
    let fromTemplate = baseDisabled || this.customerNewImage.image_dockerfile.image_base == "";
    let fromDockerFile = baseDisabled || (!this.selectFromImportFile && !this.isServerHaveDockerFile);
    return this.imageBuildMethod == ImageBuildMethod.fromTemplate ? fromTemplate : fromDockerFile;
  }

  get checkImageTagFun() {
    return this.checkImageTag.bind(this);
  }

  get checkImageNameFun() {
    return this.checkImageName.bind(this);
  }

  checkImageTag(control: HTMLInputElement): Promise<ValidationErrors> {
    if (this.customerNewImage.image_name == "") {
      return Promise.resolve(null);
    }
    return this.imageService.checkImageExist(this.projectName, this.customerNewImage.image_name, control.value)
      .then(() => null)
      .catch(err => {
        if (err && err instanceof HttpErrorResponse && (err as HttpErrorResponse).status == 409) {
          return {imageTagExist: "IMAGE.CREATE_IMAGE_TAG_EXIST"}
        }
        this.isOpen = false;
        this.messageService.dispatchError(err);
      });
  }

  checkImageName(control: HTMLInputElement): Promise<ValidationErrors> {
    if (this.customerNewImage.image_tag == "") {
      return Promise.resolve(null);
    }
    return this.imageService.checkImageExist(this.projectName, control.value, this.customerNewImage.image_tag)
      .then(() => null)
      .catch(err => {
        if (err && err instanceof HttpErrorResponse && (err as HttpErrorResponse).status == 409) {
          return {imageNameExist: "IMAGE.CREATE_IMAGE_NAME_EXIST"}
        }
        this.isOpen = false;
        this.messageService.dispatchError(err);
      });
  }

  cancelBuildImage() {
    this.imageService.cancelConsole("process_image").then(() => {
      this.isOpen = false;
    });
  }

  uploadDockerFile(): Promise<boolean> {
    if (this.selectFromImportFile) {
      this.isUploadFileWIP = true;
      let formData: FormData = new FormData();
      formData.append("upload_file", this.selectFromImportFile, this.selectFromImportFile.name);
      formData.append("project_name", this.customerNewImage.project_name);
      formData.append("image_name", this.customerNewImage.image_name);
      formData.append("tag_name", this.customerNewImage.image_tag);
      return this.imageService.uploadDockerFile(formData)
        .then(() => {
          this.isUploadFileWIP = false;
          return true;
        })
    }
  }

  buildImageByTemplate(): Promise<any> {
    this.toggleCancelBuilding = true;
    return this.imageService.buildImageFromTemp(this.customerNewImage);
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

  buildImageResole() {
    this.processImageSubscription = this.webSocketService
      .connect(`ws://${this.boardHost}/api/v1/jenkins-job/console?job_name=process_image&token=${this.appInitService.token}`)
      .subscribe((obs: MessageEvent) => {
        this.consoleText = <string>obs.data;
        let consoleTextArr: Array<string> = this.consoleText.split(/[\n]/g);
        if (consoleTextArr.find(value => value.indexOf("Finished: SUCCESS") > -1)) {
          this.isNeedAutoRefreshImageList = true;
          this.autoRefreshTimesCount = 0;
          this.processImageSubscription.unsubscribe();
        }
        if (consoleTextArr.find(value => value.indexOf("Finished: FAILURE") > -1)) {
          this.isBuildImageWIP = false;
          this.isNeedAutoRefreshImageList = false;
          this.newImageAlertType = "alert-danger";
          this.newImageErrMessage = "IMAGE.CREATE_IMAGE_BUILD_IMAGE_FAILED";
          this.isNewImageAlertOpen = true;
          this.processImageSubscription.unsubscribe();
        }
      }, err => err, () => {
        this.isOpen = false;
      });
  }

  buildImageReject(err: any) {
    this.isBuildImageWIP = false;
    this.isUploadFileWIP = false;
    this.isNeedAutoRefreshImageList = false;
    if (err && err instanceof HttpErrorResponse && (err as HttpErrorResponse).status == 401) {
      this.isOpen = false;
      this.messageService.dispatchError(err);
    } else {
      this.newImageAlertType = "alert-danger";
      this.newImageErrMessage = "IMAGE.CREATE_IMAGE_BUILD_IMAGE_FAILED";
      this.newImageErrReason = err instanceof HttpErrorResponse ? (err as HttpErrorResponse).error: "";
      this.isNewImageAlertOpen = true;
    }
  }

  buildImage() {
    this.isNewImageAlertOpen = false;
    this.isBuildImageWIP = true;
    this.consoleText = "IMAGE.CREATE_IMAGE_JENKINS_PREPARE";
    this.newImageErrReason = "";
    let buildImageFun: () => Promise<any> = this.imageBuildMethod == ImageBuildMethod.fromTemplate ?
      this.buildImageByTemplate.bind(this) :
      this.buildImageByDockerFile.bind(this);
    buildImageFun()
      .then(this.buildImageResole.bind(this))
      .catch(this.buildImageReject.bind(this));
  }

  updateFileList(): Promise<any> {
    this.isNewImageAlertOpen = false;
    let formFileList: FormData = new FormData();
    formFileList.append('project_name', this.customerNewImage.project_name);
    formFileList.append('image_name', this.customerNewImage.image_name);
    formFileList.append('tag_name', this.customerNewImage.image_tag);
    return this.imageService.getFileList(formFileList).then(res => {
      this.filesList.set(this.customerNewImage.image_name, res);
      let imageCopyArr = this.customerNewImage.image_dockerfile.image_copy;
      imageCopyArr.splice(0, imageCopyArr.length);
      this.filesList.get(this.customerNewImage.image_name).forEach(value => {
        imageCopyArr.push({
          dockerfile_copyfrom: value.path + "/" + value.file_name,
          dockerfile_copyto: this.uploadCopyToPath
        });
      });
    }).catch(err => {
      if (err && err instanceof HttpErrorResponse && (err as HttpErrorResponse).status == 401) {
        this.isOpen = false;
        this.messageService.dispatchError(err);
      } else {
        this.newImageAlertType = "alert-danger";
        this.newImageErrMessage = "IMAGE.CREATE_IMAGE_UPDATE_IMAGE_LIST_FAILED";
        this.isNewImageAlertOpen = true;
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
      this.isNewImageAlertOpen = false;
      let file:File = fileList[0];
      if (file.name !== "Dockerfile"){
        (event.target as HTMLInputElement).value = "";
        this.selectFromImportFile = null;
        this.newImageAlertType = "alert-danger";
        this.newImageErrMessage = "IMAGE.CREATE_IMAGE_FILE_NAME_ERROR";
        this.isNewImageAlertOpen = true;
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
        .catch(() => {//need't handle this error
        });
    }
  }

  uploadFile(event: Event) {
    let fileList: FileList = (event.target as HTMLInputElement).files;
    if (fileList.length > 0) {
      this.isNewImageAlertOpen = false;
      this.isUploadFileWIP = true;
      let file: File = fileList[0];
      let formData: FormData = new FormData();
      formData.append('upload_file', file, file.name);
      formData.append('project_name', this.customerNewImage.project_name);
      formData.append('image_name', this.customerNewImage.image_name);
      formData.append('tag_name', this.customerNewImage.image_tag);
      this.imageService.uploadFile(formData).subscribe((res: HttpEvent<Object>) => {
        if (res.type == HttpEventType.UploadProgress) {
          this.uploadProgressValue = res;
        } else if (res.type == HttpEventType.Response) {
          (event.target as HTMLInputElement).value = "";
          this.newImageAlertType = "alert-info";
          this.newImageErrMessage = "IMAGE.CREATE_IMAGE_UPLOAD_SUCCESS";
          this.isNewImageAlertOpen = true;
          this.isUploadFileWIP = false;
          this.updateFileListAndPreviewInfo();
        }
      }, (error: any) => {
        this.isUploadFileWIP = false;
        if (error && (error instanceof HttpErrorResponse) && (error as HttpErrorResponse).status == 401) {
          this.isOpen = false;
          this.messageService.dispatchError(error);
        } else {
          if (error && (error instanceof HttpErrorResponse)) {
            this.newImageErrReason = `:${(error as HttpErrorResponse).message}`;
          }
          (event.target as HTMLInputElement).value = "";
          this.newImageAlertType = "alert-danger";
          this.newImageErrMessage = "IMAGE.CREATE_IMAGE_UPLOAD_FAILED";
          this.isNewImageAlertOpen = true;
        }
      });
    }
  }

  getDockerFilePreviewInfo() {
    if (this.customerNewImage.image_dockerfile.image_base != "") {
      this.isNewImageAlertOpen = false;
      this.imageService.getDockerFilePreview(this.customerNewImage)
        .then(res => {
          this.consoleText = res;
        }).catch(err => {
        if (err && err instanceof HttpErrorResponse && (err as HttpErrorResponse).status == 401) {
          this.isOpen = false;
          this.messageService.dispatchError(err);
        } else {
          this.newImageAlertType = "alert-danger";
          this.newImageErrMessage = "IMAGE.CREATE_IMAGE_UPDATE_DOCKER_FILE_FAILED";
          this.isNewImageAlertOpen = true;
        }
      });
    }
  }

  shieldEnter($event: KeyboardEvent) {
    if ($event.charCode == 13) {
      (<any>$event.target).blur();
      this.getDockerFilePreviewInfo();
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
    this.isNewImageAlertOpen = false;
    let fromRemoveData: FormData = new FormData();
    fromRemoveData.append("project_name", this.customerNewImage.project_name);
    fromRemoveData.append("image_name", this.customerNewImage.image_name);
    fromRemoveData.append("tag_name", this.customerNewImage.image_tag);
    fromRemoveData.append("file_name", file.file_name);
    this.imageService.removeFile(fromRemoveData)
      .then(() => this.updateFileListAndPreviewInfo())
      .catch(err => {
        if (err && (err instanceof HttpErrorResponse) && (err as HttpErrorResponse).status == 401) {
          this.isOpen = false;
          this.messageService.dispatchError(err);
        } else {
          this.newImageAlertType = "alert-danger";
          this.newImageErrMessage = "IMAGE.CREATE_IMAGE_REMOVE_FILE_FAILED";
          this.isNewImageAlertOpen = true;
        }
      });
  }

  resetBuildMethod(method: ImageBuildMethod) {
    this.imageBuildMethod = method;
    this.consoleText = "";
    if (method == ImageBuildMethod.fromTemplate) {
      this.selectFromImportFile = null;
    }
  }

  cleanBaseImageInfo(isGetBoardRegistry: boolean = false): void {
    this.selectedImage = null;
    this.consoleText = "";
    this.imageDetailList.splice(0,this.imageDetailList.length);
    this.customerNewImage.image_dockerfile.image_base = "";
  }

  setBaseImage($event: Image): void {
    this.selectedImage = $event;
    this.imageService.getBoardRegistry().subscribe((res: string) => {
      this.boardRegistry = res.replace(/"/g,"");
      this.imageService.getImageDetailList(this.selectedImage.image_name)
        .then((res: ImageDetail[]) => {
          this.imageDetailList = res;
          this.customerNewImage.image_dockerfile.image_base = `${this.boardRegistry}/${this.selectedImage.image_name}:${res[0].image_tag}`;
          this.getDockerFilePreviewInfo();
        })
        .catch(err => {
          this.isOpen = false;
          this.messageService.dispatchError(err)
        });
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