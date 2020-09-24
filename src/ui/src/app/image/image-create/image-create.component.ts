/**
 * Created by liyanq on 21/11/2017.
 */

import { Component, OnDestroy, OnInit, Output, ViewChild } from '@angular/core';
import { ValidationErrors } from '@angular/forms';
import { TranslateService } from '@ngx-translate/core';
import { HttpErrorResponse, HttpEvent, HttpEventType, HttpProgressEvent } from '@angular/common/http';
import { interval, Observable, of, Subject, Subscription } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { ImageService } from '../image.service';
import { MessageService } from '../../shared.service/message.service';
import { AppInitService } from '../../shared.service/app-init.service';
import { WebsocketService } from '../../shared.service/websocket.service';
import { CsModalChildBase } from '../../shared/cs-modal-base/cs-modal-child-base';
import { BUTTON_STYLE, RETURN_STATUS, SharedEnvType, Tools } from '../../shared/shared.types';
import { BuildImageData, CreateImageMethod, Image, ImageCopy, ImageDetail, ImageEnv } from '../image.types';
import { JobLogComponent } from '../job-log/job-log.component';

const AUTO_REFRESH_IMAGE_LIST = 2000;

/*declared in shared-module*/
@Component({
  templateUrl: './image-create.component.html',
  styleUrls: ['./image-create.component.css']
})
export class CreateImageComponent extends CsModalChildBase implements OnInit, OnDestroy {
  @Output() refreshNotification: Subject<any>;
  @ViewChild(JobLogComponent) jobLogComponent: JobLogComponent;
  imageBuildMethod: CreateImageMethod = CreateImageMethod.Template;
  createImageMethod = CreateImageMethod;
  isOpenEnvironment = false;
  patternNewImageName: RegExp = /^[a-z\d.-]+$/;
  patternNewImageTag: RegExp = /^[a-z\d.-]+$/;
  patternBaseImage: RegExp = /^[a-z\d/.:-]+$/;
  patternExpose: RegExp = /^[\d-\s\w/\\]+$/;
  patternVolume: RegExp = /.+/;
  patternRun: RegExp = /.+/;
  patternEntryPoint: RegExp = /.+/;
  patternCopyPath: RegExp = /.+/;
  imageTemplateList: Array<object> = [{name: 'Docker File Template'}];
  filesList: Map<string, Array<{ path: string, file_name: string, size: number }>>;
  selectedDockerFile: File;
  intervalAutoRefreshImageList: any;
  intervalWaitingPoints: any;
  isNeedAutoRefreshImageList = false;
  isBuildImageWIP = false;
  isSelectedDockerFile = false;
  isUploadFileWIP = false;
  isGetImageDetailListWip = false;
  customerNewImage: BuildImageData;
  uploadCopyToPath = '/tmp';
  uploadProgressValue: HttpProgressEvent;
  imageList: Array<Image>;
  imageDetailList: Array<ImageDetail>;
  selectedImage: Image;
  baseImageSource = 1;
  boardRegistry = '';
  announceUserSubscription: Subscription;
  cancelButtonDisable = true;
  cancelInfo: { isShow: boolean, isForce: boolean, title: string, message: string };
  uploadTarPackageName = '';
  waitingMessage = '';
  waitingPoints = '';
  ws: WebSocket;

  constructor(private imageService: ImageService,
              private messageService: MessageService,
              private webSocketService: WebsocketService,
              private translateService: TranslateService,
              private appInitService: AppInitService) {
    super();
    this.filesList = new Map<string, Array<{ path: string, file_name: string, size: number }>>();
    this.imageList = Array<Image>();
    this.imageDetailList = Array<ImageDetail>();
    this.cancelInfo = {isShow: false, isForce: false, title: '', message: ''};
    this.refreshNotification = new Subject<any>();
  }

  ngOnInit() {
    this.waitingMessage = 'IMAGE.CREATE_IMAGE_WAITING_UPLOAD';
    this.intervalWaitingPoints = setInterval(() => {
      if (this.isBuildImageWIP) {
        if (this.waitingPoints === '') {
          this.waitingPoints = '.';
        } else if (this.waitingPoints === '.') {
          this.waitingPoints = '..';
        } else if (this.waitingPoints === '..') {
          this.waitingPoints = '...';
        } else {
          this.waitingPoints = '';
        }
      }
    }, 1000);
    this.intervalAutoRefreshImageList = setInterval(() => {
      if (this.isNeedAutoRefreshImageList && this.isBuildImageWIP) {
        this.waitingMessage = 'IMAGE.CREATE_IMAGE_WAITING_UPLOAD';
        this.imageService.getImages(this.customerNewImage.imageName, 0, 0).subscribe((res: Array<Image>) => {
          res.forEach(value => {
            const newImageName = `${this.customerNewImage.projectName}/${this.customerNewImage.imageName}`;
            if (value.imageName === newImageName) {
              this.isNeedAutoRefreshImageList = false;
              this.refreshNotification.next(newImageName);
              this.messageService.showGlobalMessage('IMAGE.CREATE_IMAGE_SUCCESS', {
                alertType: 'success',
                view: this.alertView
              });
              this.waitingMessage = '';
              this.isBuildImageWIP = false;
            }
          });
        }, () => this.modalOpened = false);
      }
    }, AUTO_REFRESH_IMAGE_LIST);
    this.imageService.getImages('', 0, 0).subscribe(
      (res: Array<Image>) => this.imageList = res || [],
      () => this.modalOpened = false);
  }

  ngOnDestroy() {
    if (this.ws && (this.ws.readyState === WebSocket.OPEN ||
      this.ws.readyState === WebSocket.CONNECTING)) {
      this.ws.close();
    }
    if (this.announceUserSubscription) {
      this.announceUserSubscription.unsubscribe();
    }
    clearInterval(this.intervalAutoRefreshImageList);
    clearInterval(this.intervalWaitingPoints);
  }

  public initCustomerNewImage(projectId: number, projectName: string): void {
    this.customerNewImage = new BuildImageData();
    this.customerNewImage.imageDockerFile.imageAuthor = this.appInitService.currentUser.userName;
    this.customerNewImage.projectId = projectId;
    this.customerNewImage.projectName = projectName;
    this.customerNewImage.imageTemplate = 'dockerfile-template';
    this.imageService.deleteImageConfig(projectName).subscribe();
  }

  public initBuildMethod(method: CreateImageMethod): void {
    this.imageBuildMethod = method;
  }

  get imageRun(): Array<any> {
    return this.customerNewImage.imageDockerFile.imageRun;
  }

  set imageRun(value: Array<any>) {
    this.customerNewImage.imageDockerFile.imageRun = value;
  }

  get imageVolume(): Array<any> {
    return this.customerNewImage.imageDockerFile.imageVolume;
  }

  set imageVolume(value: Array<any>) {
    this.customerNewImage.imageDockerFile.imageVolume = value;
  }

  get imageExpose(): Array<any> {
    return this.customerNewImage.imageDockerFile.imageExpose;
  }

  set imageExpose(value: Array<any>) {
    this.customerNewImage.imageDockerFile.imageExpose = value;
  }

  get envsDescription() {
    let result = '';
    this.customerNewImage.imageDockerFile.imageEnv.forEach(value => {
      result += value.envName + '=' + value.envValue + ';';
    });
    return result;
  }

  get defaultEnvsData() {
    const result = Array<SharedEnvType>();
    this.customerNewImage.imageDockerFile.imageEnv.forEach(value => {
      const env = new SharedEnvType();
      env.envName = value.envName;
      env.envValue = value.envValue;
      result.push(env);
    });
    return result;
  }

  get isBuildDisabled() {
    return this.isBuildImageWIP || this.isUploadFileWIP;
  }

  get nameAndTagDisabledDockerFile() {
    return this.isBuildDisabled || this.isSelectedDockerFile;
  }

  get isUploadDisabled(): boolean {
    return Tools.isInvalidString(this.customerNewImage.imageName) || Tools.isInvalidString(this.customerNewImage.imageTag)
      || this.isBuildImageWIP || this.isUploadFileWIP;
  }

  get checkImageTagFun() {
    return this.checkImageTag.bind(this);
  }

  get checkImageNameFun() {
    return this.checkImageName.bind(this);
  }

  get cancelCaption() {
    return this.waitingMessage === 'IMAGE.CREATE_IMAGE_JENKINS_PREPARE' ?
      'IMAGE.CREATE_IMAGE_CANCEL_WAIT' :
      'IMAGE.CREATE_IMAGE_BUILD_CANCEL';
  }

  checkImageTag(control: HTMLInputElement): Observable<ValidationErrors | null> {
    if (this.customerNewImage.imageName === '' ||
      control.value === this.customerNewImage.imageTag) {
      return of(null);
    }
    return this.imageService.checkImageExist(this.customerNewImage.projectName, this.customerNewImage.imageName, control.value)
      .pipe(
        map(() => null),
        catchError((err: HttpErrorResponse) => {
          if (err.status === 409) {
            this.messageService.cleanNotification();
            return of({imageTagExists: 'IMAGE.CREATE_IMAGE_TAG_EXIST'});
          } else if (err.status === 404) {
            this.messageService.cleanNotification();
          } else {
            this.modalOpened = false;
          }
          return of(null);
        }));
  }

  checkImageName(control: HTMLInputElement): Observable<ValidationErrors | null> {
    if (this.customerNewImage.imageTag === '' ||
      control.value === this.customerNewImage.imageName) {
      return of(null);
    }
    return this.imageService.checkImageExist(this.customerNewImage.projectName, control.value, this.customerNewImage.imageTag)
      .pipe(
        map(() => null),
        catchError((err: HttpErrorResponse) => {
          if (err.status === 409) {
            this.messageService.cleanNotification();
            return of({imageNameExists: 'IMAGE.CREATE_IMAGE_NAME_EXIST'});
          } else if (err.status === 404) {
            this.messageService.cleanNotification();
          } else {
            this.modalOpened = false;
          }
          return of(null);
        }));
  }

  cancelBuildImage() {
    if (this.waitingMessage === 'IMAGE.CREATE_IMAGE_WAITING_BUILD') {
      this.cancelInfo.isForce = false;
      this.cancelInfo.title = 'IMAGE.CREATE_IMAGE_BUILD_CANCEL';
      this.cancelInfo.message = 'IMAGE.CREATE_IMAGE_BUILD_CANCEL_MSG';
    } else {
      this.cancelInfo.isForce = true;
      this.cancelInfo.title = 'IMAGE.CREATE_IMAGE_FORCE_QUIT';
      this.cancelInfo.message = 'IMAGE.CREATE_IMAGE_FORCE_QUIT_MSG';
    }
    this.cancelInfo.isShow = true;
  }

  cancelBuildImageBehavior() {
    this.cancelInfo.isShow = false;
    if (this.cancelInfo.isForce) {
      this.modalOpened = false;
    } else {
      this.imageService.cancelConsole(this.customerNewImage.projectName).subscribe(
        () => this.cleanImageConfig(),
        () => this.modalOpened = false,
        () => this.modalOpened = false);
    }
  }

  uploadDockerFile(): Observable<string> {
    const formData: FormData = new FormData();
    formData.append('upload_file', this.selectedDockerFile, this.selectedDockerFile.name);
    formData.append('project_name', this.customerNewImage.projectName);
    formData.append('image_name', this.customerNewImage.imageName);
    formData.append('image_tag', this.customerNewImage.imageTag);
    return this.imageService.uploadDockerFile(formData);
  }

  buildImageByDockerFile(): Observable<any> {
    const fileInfo = {
      imageName: this.customerNewImage.imageName,
      tagName: this.customerNewImage.imageTag,
      projectName: this.customerNewImage.projectName
    };
    return this.imageService.buildImageFromDockerFile(fileInfo);
  }

  buildImageByImagePackage(): Observable<any> {
    const params = {
      imageName: this.customerNewImage.imageName,
      tagName: this.customerNewImage.imageTag,
      projectName: this.customerNewImage.projectName,
      imagePackageName: this.uploadTarPackageName
    };
    return this.imageService.buildImageFromImagePackage(params);
  }

  cleanImageConfig(err?: HttpErrorResponse) {
    this.isBuildImageWIP = false;
    this.isUploadFileWIP = false;
    this.isNeedAutoRefreshImageList = false;
    if (err) {
      const reason = err ? err.error as string : '';
      this.translateService.get(`IMAGE.CREATE_IMAGE_BUILD_IMAGE_FAILED`).subscribe(
        (msg: string) => this.messageService.showGlobalMessage(`${msg}:${reason}`, {
          alertType: 'danger',
          view: this.alertView
        }));
    }
    this.imageService.deleteImageConfig(this.customerNewImage.projectName).subscribe();
  }

  mountWebSocket() {
    this.ws.onopen = (ev: Event): any => {
      this.jobLogComponent.clear();
    };

    this.ws.onmessage = (ev: MessageEvent): any => {
      this.waitingMessage = 'IMAGE.CREATE_IMAGE_WAITING_BUILD';
      this.cancelButtonDisable = false;
      const receivedMessage = ev.data as string;
      let consoleTextArr = Array<string>();
      if (receivedMessage && receivedMessage.length > 0) {
        consoleTextArr = receivedMessage.split(/\r\n|\r|\n/);
        this.jobLogComponent.appendContentArray(consoleTextArr);
      }
      if (consoleTextArr.find(value =>
        value.indexOf('Job succeeded') > -1 ||
        value.indexOf('Finished: SUCCESS') > -1)) {
        this.isNeedAutoRefreshImageList = true;
        this.announceUserSubscription = interval(30 * 60 * 1000).subscribe(() => {
          if (this.isBuildImageWIP) {
            this.messageService.showDialog('IMAGE.CREATE_IMAGE_UPLOAD_IMAGE_TIMEOUT', {
              title: 'IMAGE.CREATE_IMAGE_TIMEOUT',
              view: this.alertView,
              buttonStyle: BUTTON_STYLE.YES_NO
            }).subscribe(message => {
              if (message.returnStatus === RETURN_STATUS.rsCancel) {
                this.refreshNotification.next({});
                this.modalOpened = false;
              }
            });
          }
        });
      } else if (consoleTextArr.find(value =>
        value.indexOf('ERROR: Job failed') > -1 ||
        value.indexOf('Finished: FAILURE') > -1)) {
        this.isBuildImageWIP = false;
        this.isUploadFileWIP = false;
        this.cancelButtonDisable = true;
        this.isNeedAutoRefreshImageList = false;
        this.appInitService.setAuditLog({
          operation_user_id: this.appInitService.currentUser.userId,
          operation_user_name: this.appInitService.currentUser.userName,
          operation_project_id: this.customerNewImage.projectId,
          operation_project_name: this.customerNewImage.projectName,
          operation_object_type: 'images',
          operation_object_name: '',
          operation_action: 'create',
          operation_status: 'Failed'
        }).subscribe();
        this.messageService.showGlobalMessage('IMAGE.CREATE_IMAGE_FAILED', {
          alertType: 'danger',
          view: this.alertView
        });
      }
    };

    this.ws.onerror = (ev: Event): any => {
      this.messageService.showGlobalMessage('websocket connect error');
      this.modalOpened = false;
    };
  }

  buildImageResole() {
    const boardHost = this.appInitService.systemInfo.boardHost;
    const wsHost = `${this.appInitService.getWebsocketPrefix}://${boardHost}:${window.location.port}/api/v1/jenkins-job/console`;
    const wsParams = `job_name=${this.customerNewImage.projectName}&token=${this.appInitService.token}`;
    try {
      this.ws = new WebSocket(`${wsHost}?${wsParams}`);
      this.mountWebSocket();
    } catch (e) {
      this.isBuildImageWIP = false;
      this.waitingMessage = '';
      this.messageService.showGlobalMessage(e.toString(), {view: this.alertView});
    }
  }

  buildImage() {
    const buildImageInit = () => {
      this.cancelButtonDisable = true;
      this.isBuildImageWIP = true;
      this.waitingMessage = 'IMAGE.CREATE_IMAGE_JENKINS_PREPARE';
      setTimeout(() => this.cancelButtonDisable = false, 10000);
    };
    if (this.imageBuildMethod === CreateImageMethod.Template) {
      if (this.verifyInputExValid() &&
        this.verifyInputArrayExValid() &&
        this.verifyDropdownExValid() &&
        this.customerNewImage.imageDockerFile.imageBase !== '') {
        buildImageInit();
        this.imageService.buildImageFromTemp(this.customerNewImage).subscribe(
          () => this.buildImageResole(),
          (error: HttpErrorResponse) => this.cleanImageConfig(error)
        );
      }
    } else if (this.imageBuildMethod === CreateImageMethod.DockerFile) {
      if (this.verifyInputExValid()) {
        if (this.isSelectedDockerFile) {
          buildImageInit();
          this.buildImageByDockerFile().subscribe(
            () => this.buildImageResole(),
            (error: HttpErrorResponse) => this.cleanImageConfig(error)
          );
        } else {
          this.messageService.showAlert('IMAGE.CREATE_IMAGE_SELECT_DOCKER_FILE', {
            alertType: 'warning',
            view: this.alertView
          });
        }
      }
    } else if (this.imageBuildMethod === CreateImageMethod.ImagePackage) {
      if (this.verifyInputExValid()) {
        if (this.uploadTarPackageName !== '') {
          buildImageInit();
          this.buildImageByImagePackage().subscribe(
            () => this.buildImageResole(),
            (error: HttpErrorResponse) => this.cleanImageConfig(error)
          );
        } else {
          this.messageService.showAlert('IMAGE.CREATE_IMAGE_SELECT_IMAGE_PACKAGE', {
            alertType: 'warning',
            view: this.alertView
          });
        }
      }
    }
  }

  updateFileList(): Observable<any> {
    this.filesList.clear();
    const formFileList: FormData = new FormData();
    formFileList.append('project_name', this.customerNewImage.projectName);
    formFileList.append('image_name', this.customerNewImage.imageName);
    formFileList.append('image_tag', this.customerNewImage.imageTag);
    return this.imageService.getFileList(formFileList)
      .pipe(
        map(res => {
          this.filesList.set(this.customerNewImage.imageName, res);
          const imageCopyArr = this.customerNewImage.imageDockerFile.imageCopy;
          imageCopyArr.splice(0, imageCopyArr.length);
          this.filesList.get(this.customerNewImage.imageName).forEach(value => {
            const copy = new ImageCopy();
            copy.copyFrom = value.file_name;
            copy.copyTo = this.uploadCopyToPath + '/' + value.file_name;
            imageCopyArr.push(copy);
          });
        }),
        catchError((err: HttpErrorResponse) => {
          if (err.status === 401) {
            this.modalOpened = false;
          } else {
            this.messageService.showAlert('IMAGE.CREATE_IMAGE_UPDATE_IMAGE_LIST_FAILED', {
              alertType: 'danger',
              view: this.alertView
            });
          }
          return null;
        }));
  }

  updateFileListAndPreviewInfo() {
    this.updateFileList().subscribe(() => {
      this.getDockerFilePreviewInfo();
    });
  }

  selectDockerFile(event: Event) {
    const fileList: FileList = (event.target as HTMLInputElement).files;
    if (fileList.length > 0 && this.verifyInputExValid()) {
      const file: File = fileList[0];
      if (file.name !== 'Dockerfile') {
        (event.target as HTMLInputElement).value = '';
        this.messageService.showAlert('IMAGE.CREATE_IMAGE_FILE_NAME_ERROR', {
          alertType: 'danger',
          view: this.alertView
        });
      } else {
        this.selectedDockerFile = file;
        this.uploadDockerFile().subscribe((res: string) => {
          this.isSelectedDockerFile = true;
          this.jobLogComponent.clear();
          this.jobLogComponent.appendContent(res);
          this.messageService.showAlert('IMAGE.CREATE_IMAGE_FILE_UPLOAD_SUCCESS', {view: this.alertView});
        }, (err: HttpErrorResponse) => this.messageService.showAlert(err.error, {
          alertType: 'danger',
          view: this.alertView
        }));
      }
    } else {
      (event.target as HTMLInputElement).value = '';
    }
  }

  uploadFile(event: Event) {
    const fileList: FileList = (event.target as HTMLInputElement).files;
    if (fileList.length > 0 && this.verifyInputExValid()) {
      const file: File = fileList[0];
      if (file.size > 1024 * 1024 * 500) {
        (event.target as HTMLInputElement).value = '';
        this.messageService.showAlert('IMAGE.CREATE_IMAGE_UPDATE_FILE_SIZE', {
          alertType: 'danger',
          view: this.alertView
        });
      } else {
        const formData: FormData = new FormData();
        this.isUploadFileWIP = true;
        formData.append('upload_file', file, file.name);
        formData.append('project_name', this.customerNewImage.projectName);
        formData.append('image_name', this.customerNewImage.imageName);
        formData.append('image_tag', this.customerNewImage.imageTag);
        this.imageService.uploadFile(formData).subscribe((res: HttpEvent<object>) => {
          if (res.type === HttpEventType.UploadProgress) {
            this.uploadProgressValue = res;
          } else if (res.type === HttpEventType.Response) {
            (event.target as HTMLInputElement).value = '';
            this.uploadTarPackageName = file.name;
            this.isUploadFileWIP = false;
            this.updateFileListAndPreviewInfo();
            this.messageService.showAlert('IMAGE.CREATE_IMAGE_UPLOAD_SUCCESS', {view: this.alertView});
          }
        }, (error: HttpErrorResponse) => {
          this.isUploadFileWIP = false;
          if (error.status === 401) {
            this.modalOpened = false;
          } else {
            (event.target as HTMLInputElement).value = '';
            const newImageErrReason = (error.error as Error).message;
            this.translateService.get('IMAGE.CREATE_IMAGE_UPLOAD_FAILED').subscribe((msg: string) => {
              this.messageService.showAlert(`${msg}:${newImageErrReason}`, {
                alertType: 'danger',
                view: this.alertView
              });
            });
          }
        });
      }
    }
  }

  getDockerFilePreviewInfo() {
    if (this.customerNewImage.imageDockerFile.imageBase !== '') {
      this.imageService.getDockerFilePreview(this.customerNewImage).subscribe(
        res => {
          this.jobLogComponent.clear();
          this.jobLogComponent.appendContent(res);
        },
        (err: HttpErrorResponse) => {
          if (err.status === 401) {
            this.modalOpened = false;
          } else {
            this.messageService.showAlert('IMAGE.CREATE_IMAGE_UPDATE_DOCKER_FILE_FAILED', {
              alertType: 'danger',
              view: this.alertView
            });
          }
        }
      );
    }
  }

  setEnvironment(envsData: Array<SharedEnvType>) {
    const envsArray = this.customerNewImage.imageDockerFile.imageEnv;
    envsArray.splice(0, envsArray.length);
    envsData.forEach((value: SharedEnvType) => {
      const env = new ImageEnv();
      env.envName = value.envName;
      env.envValue = value.envValue;
      envsArray.push(env);
    });
    this.getDockerFilePreviewInfo();
  }

  removeFile(file: { path: string, file_name: string, size: number }) {
    const fromRemoveData: FormData = new FormData();
    fromRemoveData.append('project_name', this.customerNewImage.projectName);
    fromRemoveData.append('image_name', this.customerNewImage.imageName);
    fromRemoveData.append('image_tag', this.customerNewImage.imageTag);
    fromRemoveData.append('file_name', file.file_name);
    this.imageService.removeFile(fromRemoveData).subscribe(
      () => this.messageService.showAlert('IMAGE.CREATE_IMAGE_REMOVE_FILE_SUCCESS', {view: this.alertView}),
      (err: HttpErrorResponse) => {
        if (err.status === 401) {
          this.modalOpened = false;
        } else {
          this.messageService.showAlert('IMAGE.CREATE_IMAGE_REMOVE_FILE_FAILED', {
            alertType: 'danger',
            view: this.alertView
          });
        }
      },
      () => this.updateFileListAndPreviewInfo());
  }

  cleanBaseImageInfo(isGetBoardRegistry: boolean = false): void {
    if ((this.baseImageSource === 1 && isGetBoardRegistry) ||
      (this.baseImageSource === 2 && !isGetBoardRegistry)) {
      this.selectedImage = null;
      this.imageDetailList.splice(0, this.imageDetailList.length);
      this.customerNewImage.imageDockerFile.imageBase = '';
    }
  }

  setBaseImage(selectImage: Image): void {
    this.selectedImage = null;
    this.imageDetailList = null;
    this.isGetImageDetailListWip = true;
    this.imageService.getBoardRegistry().subscribe((res: string) => {
      this.boardRegistry = res.replace(/"/g, '');
      this.imageService.getImageDetailList(selectImage.imageName).subscribe((imageDetails: ImageDetail[]) => {
          this.selectedImage = selectImage;
          this.imageDetailList = imageDetails;
          this.customerNewImage.imageDockerFile.imageBase =
            `${this.boardRegistry}/${this.selectedImage.imageName}:${imageDetails[0].imageTag}`;
          this.getDockerFilePreviewInfo();
        },
        () => this.modalOpened = false,
        () => this.isGetImageDetailListWip = false
      );
    });
  }

  setBaseImageDetail(detail: ImageDetail): void {
    this.imageService.getBoardRegistry().subscribe((res: string) => {
      this.boardRegistry = res.replace(/"/g, '');
      this.customerNewImage.imageDockerFile.imageBase = `${this.boardRegistry}/${this.selectedImage.imageName}:${detail.imageTag}`;
      this.getDockerFilePreviewInfo();
    });
  }
}
