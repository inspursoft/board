/**
 * Created by liyanq on 21/11/2017.
 */

import {
  Component, QueryList, ViewChildren, OnInit, OnDestroy,
  AfterContentChecked, Output, EventEmitter, Input
} from "@angular/core"
import { CsInputArrayComponent } from "../../shared/cs-components-library/cs-input-array/cs-input-array.component";
import { CsInputComponent } from "../../shared/cs-components-library/cs-input/cs-input.component";
import { BuildImageData } from "../image";
import { ImageService } from "../image-service/image-service";
import { MessageService } from "../../shared/message-service/message.service";
import { Response } from "@angular/http"
import { AppInitService } from "../../app.init.service";
import { Subscription } from "rxjs/Subscription";
import { WebsocketService } from "../../shared/websocket-service/websocket.service";
import { EnvType } from "../../shared/environment-value/environment-value.component";

enum ImageSource{fromBoardRegistry, fromDockerHub}
const AUTO_REFRESH_IMAGE_LIST: number = 2000;
const PROCESS_IMAGE_CONSOLE_URL = `ws://10.165.22.61:8088/api/v1/jenkins-job/console?job_name=process_image`;
// const PROCESS_IMAGE_CONSOLE_URL = `ws://localhost/api/v1/jenkins-job/console?job_name=process_image`;
type alertType = "alert-info" | "alert-danger";

/*declared in shared-module*/
@Component({
  selector: "create-image",
  templateUrl: "./image-create.component.html",
  styleUrls: ["./image-create.component.css"]
})
export class CreateImageComponent implements OnInit, AfterContentChecked, OnDestroy {
  _isOpen: boolean = false;
  @ViewChildren(CsInputArrayComponent) inputArrayComponents: QueryList<CsInputArrayComponent>;
  @ViewChildren(CsInputComponent) inputComponents: QueryList<CsInputComponent>;
  @Input() projectId: number = 0;
  @Input() projectName: string = "";
  @Output() onBuildCompleted: EventEmitter<string>;
  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();


  _isOpenEnvironment = false;
  patternNewImageName: RegExp = /^[a-z\d.-]+$/;
  patternNewImageTag: RegExp = /^[a-z\d.-]+$/;
  patternBaseImage: RegExp = /^[a-z\d.:-]+$/;
  patternExpose: RegExp = /^[\d-\s\w/\\]+$/;
  patternVolume: RegExp = /.+/;
  patternRun: RegExp = /.+/;
  patternEntryPoint: RegExp = /.+/;
  imageSource: ImageSource = ImageSource.fromBoardRegistry;
  newImageAlertType: alertType = "alert-danger";
  imageTemplateList: Array<Object> = [{name: "Docker File Template"}];
  filesList: Map<string, Array<{path: string, file_name: string, size: number}>>;
  intervalAutoRefreshImageList: any;
  isNeedAutoRefreshImageList: boolean = false;
  imageInBuilding: boolean = false;
  isInputComponentsValid: boolean = false;
  isUploadFileIng = false;
  customerNewImage: BuildImageData;
  autoRefreshTimesCount: number = 0;
  isNewImageAlertOpen: boolean = false;
  newImageErrMessage: string = "";
  consoleText: string = "";
  lastJobNumber: number = 0;
  processImageSubscription: Subscription;

  constructor(private imageService: ImageService,
              private messageService: MessageService,
              private webSocketService: WebsocketService,
              private appInitService: AppInitService) {
    this.onBuildCompleted = new EventEmitter<string>();
    this.filesList = new Map<string, Array<{path: string, file_name: string, size: number}>>();
  }

  ngOnInit() {
    this.customerNewImage = new BuildImageData();
    this.customerNewImage.image_dockerfile.image_author = this.appInitService.currentUser["user_name"];
    this.customerNewImage.project_id = this.projectId;
    this.customerNewImage.project_name = this.projectName;
    this.customerNewImage.image_template = "dockerfile-template";
    this.intervalAutoRefreshImageList = setInterval(() => {
      if (this.isNeedAutoRefreshImageList && this.imageInBuilding) {
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
          if (err && err instanceof Response && (err as Response).status == 401) {
            this.isOpen = false;
            this.messageService.dispatchError(err);
          } else {
            this.imageInBuilding = false;
            this.newImageAlertType = "alert-danger";
            this.newImageErrMessage = "IMAGE.CREATE_IMAGE_UPDATE_IMAGE_LIST_FAILED";
            this.isNewImageAlertOpen = true;
          }
        });
      }
    }, AUTO_REFRESH_IMAGE_LIST);
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
    this.isOpen = !value;
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

  cancelBuildImage() {
    if (this.lastJobNumber > 0) {
      this.imageService.cancelConsole("process_image", this.lastJobNumber).then(() => {
        this.isOpen = false;
      });
      this.lastJobNumber = -1;
    }
  }

  async buildImage() {
    this.isNewImageAlertOpen = false;
    this.imageInBuilding = true;
    this.lastJobNumber = 0;
    this.consoleText = "Jenkins preparing...";
    this.imageService.buildImage(this.customerNewImage)
      .then(() => {
        setTimeout(() => {
          this.processImageSubscription = this.webSocketService
            .connect(PROCESS_IMAGE_CONSOLE_URL + `&token=${this.appInitService.token}`)
            .subscribe(obs => {
              this.consoleText = <string>obs.data;
              if (this.lastJobNumber == 0) {
                this.imageService.getLastJobId("process_image").then(res => {
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
                this.newImageErrMessage = "IMAGE.CREATE_IMAGE_BUILD_IMAGE_FAILED";
                this.isNewImageAlertOpen = true;
                this.processImageSubscription.unsubscribe();
              }
            }, err => err, () => {
              this.isOpen = false;
            });
        }, 10000);
      })
      .catch((err) => {
        this.imageInBuilding = false;
        this.isNeedAutoRefreshImageList = false;
        if (err && err instanceof Response && (err as Response).status == 401) {
          this.isOpen = false;
          this.messageService.dispatchError(err);
        } else {
          this.newImageAlertType = "alert-danger";
          this.newImageErrMessage = "IMAGE.CREATE_IMAGE_BUILD_IMAGE_FAILED";
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
    return this.imageService.getFileList(formFileList).then(res => {
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
      if (err && err instanceof Response && (err as Response).status == 401) {
        this.isOpen = false;
        this.messageService.dispatchError(err);
      } else {
        this.newImageAlertType = "alert-danger";
        this.newImageErrMessage = "IMAGE.CREATE_IMAGE_UPDATE_IMAGE_LIST_FAILED";
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
      this.imageService.uploadFile(formData).then(() => {
        event.target.value = "";
        this.newImageAlertType = "alert-info";
        this.newImageErrMessage = "IMAGE.CREATE_IMAGE_UPLOAD_SUCCESS";
        this.isNewImageAlertOpen = true;
        this.isUploadFileIng = false;
        this.asyncGetDockerFilePreviewInfo();
      }).catch(err => {
        if (err && (err instanceof Response) && (err as Response).status == 401) {
          this.isOpen = false;
          this.messageService.dispatchError(err);
        } else {
          this.newImageAlertType = "alert-danger";
          this.newImageErrMessage = "IMAGE.CREATE_IMAGE_UPLOAD_FAILED";
          this.isNewImageAlertOpen = true;
          this.isUploadFileIng = false;
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
        if (err && err instanceof Response && (err as Response).status == 401) {
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
      .then(() => this.asyncGetDockerFilePreviewInfo())
      .catch(err => {
        if (err && (err instanceof Response) && (err as Response).status == 401) {
          this.isOpen = false;
          this.messageService.dispatchError(err);
        } else {
          this.newImageAlertType = "alert-danger";
          this.newImageErrMessage = "IMAGE.CREATE_IMAGE_REMOVE_FILE_FAILED";
          this.isNewImageAlertOpen = true;
        }
      });
  }

}