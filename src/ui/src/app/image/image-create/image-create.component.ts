/**
 * Created by liyanq on 21/11/2017.
 */

import { Component, ElementRef, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { BuildImageData, Image, ImageDetail } from '../image';
import { ImageService } from '../image-service/image-service';
import { MessageService } from '../../shared.service/message.service';
import { HttpErrorResponse, HttpEvent, HttpEventType, HttpProgressEvent } from '@angular/common/http';
import { AppInitService } from '../../shared.service/app-init.service';
import { WebsocketService } from '../../shared.service/websocket.service';
import { EnvType } from '../../shared/environment-value/environment-value.component';
import { ValidationErrors } from '@angular/forms';
import { TranslateService } from '@ngx-translate/core';
import { CsModalChildBase } from '../../shared/cs-modal-base/cs-modal-child-base';
import { BUTTON_STYLE, CreateImageMethod, GlobalAlertType, RETURN_STATUS, Tools } from '../../shared/shared.types';
import { interval, Observable, of, Subscription } from 'rxjs';
import { catchError, map } from 'rxjs/operators';

const AUTO_REFRESH_IMAGE_LIST = 2000;

/*declared in shared-module*/
@Component({
  selector: 'create-image',
  templateUrl: './image-create.component.html',
  styleUrls: ['./image-create.component.css']
})
export class CreateImageComponent extends CsModalChildBase implements OnInit, OnDestroy {
  boardHost: string;
  @ViewChild('areaStatus') areaStatus: ElementRef;
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
  consoleText = '';
  uploadCopyToPath = '/tmp';
  uploadProgressValue: HttpProgressEvent;
  imageList: Array<Image>;
  imageDetailList: Array<ImageDetail>;
  selectedImage: Image;
  baseImageSource = 1;
  boardRegistry = '';
  processImageSubscription: Subscription;
  announceUserSubscription: Subscription;
  cancelButtonDisable = true;
  cancelInfo: { isShow: boolean, isForce: boolean, title: string, message: string };
  uploadTarPackageName = '';
  waitingMessage = '';
  waitingPoints = '';

  constructor(private imageService: ImageService,
              private messageService: MessageService,
              private webSocketService: WebsocketService,
              private translateService: TranslateService,
              private appInitService: AppInitService) {
    super();
    this.filesList = new Map<string, Array<{ path: string, file_name: string, size: number }>>();
    this.boardHost = this.appInitService.systemInfo.board_host;
    this.imageList = Array<Image>();
    this.imageDetailList = Array<ImageDetail>();
    this.cancelInfo = {isShow: false, isForce: false, title: '', message: ''};
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
        this.imageService.getImages(this.customerNewImage.image_name, 0, 0).subscribe((res: Array<Image>) => {
          res.forEach(value => {
            const newImageName = `${this.customerNewImage.project_name}/${this.customerNewImage.image_name}`;
            if (value.image_name === newImageName) {
              this.isNeedAutoRefreshImageList = false;
              this.closeNotification.next(newImageName);
              this.modalOpened = false;
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
    if (this.processImageSubscription) {
      this.processImageSubscription.unsubscribe();
    }
    if (this.announceUserSubscription) {
      this.announceUserSubscription.unsubscribe();
    }
    clearInterval(this.intervalAutoRefreshImageList);
    clearInterval(this.intervalWaitingPoints);
  }

  public initCustomerNewImage(projectId: number, projectName: string): void {
    this.customerNewImage = new BuildImageData();
    this.customerNewImage.image_dockerfile.image_author = this.appInitService.currentUser.user_name;
    this.customerNewImage.project_id = projectId;
    this.customerNewImage.project_name = projectName;
    this.customerNewImage.image_template = 'dockerfile-template';
    this.imageService.deleteImageConfig(projectName).subscribe();
  }

  public initBuildMethod(method: CreateImageMethod): void {
    this.imageBuildMethod = method;
  }

  get imageRun(): Array<string> {
    return this.customerNewImage.image_dockerfile.image_run;
  }

  set imageRun(value: Array<string>) {
    this.customerNewImage.image_dockerfile.image_run = value;
  }

  get imageVolume(): Array<string> {
    return this.customerNewImage.image_dockerfile.image_volume;
  }

  set imageVolume(value: Array<string>) {
    this.customerNewImage.image_dockerfile.image_volume = value;
  }

  get imageExpose(): Array<string> {
    return this.customerNewImage.image_dockerfile.image_expose;
  }

  set imageExpose(value: Array<string>) {
    this.customerNewImage.image_dockerfile.image_expose = value;
  }

  get envsDescription() {
    let result = '';
    this.customerNewImage.image_dockerfile.image_env.forEach(value => {
      result += value.dockerfile_envname + '=' + value.dockerfile_envvalue + ';';
    });
    return result;
  }

  get defaultEnvsData() {
    const result = Array<EnvType>();
    this.customerNewImage.image_dockerfile.image_env.forEach(value => {
      result.push(new EnvType(value.dockerfile_envname, value.dockerfile_envvalue));
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
    return Tools.isInvalidString(this.customerNewImage.image_name) || Tools.isInvalidString(this.customerNewImage.image_tag)
      || this.isBuildImageWIP || this.isUploadFileWIP;
  }

  get checkImageTagFun() {
    return this.checkImageTag.bind(this);
  }

  get checkImageNameFun() {
    return this.checkImageName.bind(this);
  }

  get cancelCaption() {
    return this.consoleText === 'IMAGE.CREATE_IMAGE_JENKINS_PREPARE' ?
      'IMAGE.CREATE_IMAGE_CANCEL_WAIT' :
      'IMAGE.CREATE_IMAGE_BUILD_CANCEL';
  }

  checkImageTag(control: HTMLInputElement): Observable<ValidationErrors | null> {
    if (this.customerNewImage.image_name === '') {
      return of(null);
    }
    return this.imageService.checkImageExist(this.customerNewImage.project_name, this.customerNewImage.image_name, control.value)
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
    if (this.customerNewImage.image_tag === '') {
      return of(null);
    }
    return this.imageService.checkImageExist(this.customerNewImage.project_name, control.value, this.customerNewImage.image_tag)
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
    if (this.consoleText === 'IMAGE.CREATE_IMAGE_JENKINS_PREPARE') {
      this.cancelInfo.isForce = true;
      this.cancelInfo.title = 'IMAGE.CREATE_IMAGE_FORCE_QUIT';
      this.cancelInfo.message = 'IMAGE.CREATE_IMAGE_FORCE_QUIT_MSG';
    } else {
      this.cancelInfo.isForce = false;
      this.cancelInfo.title = 'IMAGE.CREATE_IMAGE_BUILD_CANCEL';
      this.cancelInfo.message = 'IMAGE.CREATE_IMAGE_BUILD_CANCEL_MSG';
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
        () => this.modalOpened = false,
        () => this.modalOpened = false);
    }
  }

  uploadDockerFile(): Observable<string> {
    const formData: FormData = new FormData();
    formData.append('upload_file', this.selectedDockerFile, this.selectedDockerFile.name);
    formData.append('project_name', this.customerNewImage.project_name);
    formData.append('image_name', this.customerNewImage.image_name);
    formData.append('image_tag', this.customerNewImage.image_tag);
    return this.imageService.uploadDockerFile(formData);
  }

  buildImageByDockerFile(): Observable<any> {
    const fileInfo = {
      imageName: this.customerNewImage.image_name,
      tagName: this.customerNewImage.image_tag,
      projectName: this.customerNewImage.project_name
    };
    return this.imageService.buildImageFromDockerFile(fileInfo);
  }

  buildImageByImagePackage(): Observable<any> {
    const params = {
      imageName: this.customerNewImage.image_name,
      tagName: this.customerNewImage.image_tag,
      projectName: this.customerNewImage.project_name,
      imagePackageName: this.uploadTarPackageName
    };
    return this.imageService.buildImageFromImagePackage(params);
  }

  cleanImageConfig(err?: HttpErrorResponse) {
    this.isBuildImageWIP = false;
    this.isUploadFileWIP = false;
    this.isNeedAutoRefreshImageList = false;
    if (err) {
      const reason = err ? ((err as HttpErrorResponse).error as Error).message : '';
      this.translateService.get(`IMAGE.CREATE_IMAGE_BUILD_IMAGE_FAILED`).subscribe((msg: string) => {
        this.messageService.showAlert(`${msg}:${reason}`, {alertType: 'danger', view: this.alertView});
      });
    }
    this.imageService.deleteImageConfig(this.customerNewImage.project_name).subscribe();
  }

  buildImageResole() {
    const wsHost = `ws://${this.boardHost}/api/v1/jenkins-job/console`;
    const wsParams = `job_name=${this.customerNewImage.project_name}&token=${this.appInitService.token}`;
    this.processImageSubscription = this.webSocketService.connect(`${wsHost}?${wsParams}`)
      .subscribe((obs: MessageEvent) => {
        this.consoleText = obs.data as string;
        this.waitingMessage = 'IMAGE.CREATE_IMAGE_WAITING_BUILD';
        this.cancelButtonDisable = false;
        this.areaStatus.nativeElement.scrollTop = this.areaStatus.nativeElement.scrollHeight;
        const consoleTextArr: Array<string> = this.consoleText.split(/[\n]/g);
        if (consoleTextArr.find(value => value.indexOf('Finished: SUCCESS') > -1)) {
          this.isNeedAutoRefreshImageList = true;
          this.announceUserSubscription = interval(5 * 60 * 1000).subscribe(() => {
            this.messageService.showDialog('IMAGE.CREATE_IMAGE_UPLOAD_IMAGE_TIMEOUT',
              {title: 'IMAGE.CREATE_IMAGE_TIMEOUT', view: this.alertView, buttonStyle: BUTTON_STYLE.YES_NO})
              .subscribe(message => {
                if (message.returnStatus === RETURN_STATUS.rsCancel) {
                  this.closeNotification.next({});
                  this.modalOpened = false;
                }
              });
          });
          this.processImageSubscription.unsubscribe();
        }
        if (consoleTextArr.find(value => value.indexOf('Finished: FAILURE') > -1)) {
          this.isBuildImageWIP = false;
          this.isUploadFileWIP = false;
          this.cancelButtonDisable = true;
          this.isNeedAutoRefreshImageList = false;
          this.appInitService.setAuditLog({
            operation_user_id: this.appInitService.currentUser.user_id,
            operation_user_name: this.appInitService.currentUser.user_name,
            operation_project_id: this.customerNewImage.project_id,
            operation_project_name: this.customerNewImage.project_name,
            operation_object_type: 'images',
            operation_object_name: '',
            operation_action: 'create',
            operation_status: 'Failed'
          }).subscribe();
          this.processImageSubscription.unsubscribe();
        }
      }, err => {
        this.messageService.showGlobalMessage('websocket connect error',
          {globalAlertType: GlobalAlertType.gatShowDetail, errorObject: err}
        );
        this.modalOpened = false;
      }, () => this.modalOpened = false);
  }

  buildImage() {
    const buildImageInit = () => {
      this.cancelButtonDisable = true;
      this.isBuildImageWIP = true;
      this.waitingMessage = 'IMAGE.CREATE_IMAGE_JENKINS_PREPARE';
      this.consoleText = 'IMAGE.CREATE_IMAGE_JENKINS_PREPARE';
      setTimeout(() => this.cancelButtonDisable = false, 10000);
    };
    if (this.imageBuildMethod === CreateImageMethod.Template) {
      if (this.verifyInputExValid() &&
        this.verifyInputArrayExValid() &&
        this.verifyDropdownExValid() &&
        this.customerNewImage.image_dockerfile.image_base !== '') {
        buildImageInit();
        this.imageService.buildImageFromTemp(this.customerNewImage).subscribe(
          () => this.buildImageResole(),
          (error: HttpErrorResponse) => this.cleanImageConfig(error)
        );
      }
    } else if (this.imageBuildMethod === CreateImageMethod.DockerFile) {
      if (this.verifyInputExValid()) {
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
    formFileList.append('project_name', this.customerNewImage.project_name);
    formFileList.append('image_name', this.customerNewImage.image_name);
    formFileList.append('image_tag', this.customerNewImage.image_tag);
    return this.imageService.getFileList(formFileList)
      .pipe(
        map(res => {
          this.filesList.set(this.customerNewImage.image_name, res);
          const imageCopyArr = this.customerNewImage.image_dockerfile.image_copy;
          imageCopyArr.splice(0, imageCopyArr.length);
          this.filesList.get(this.customerNewImage.image_name).forEach(value => {
            imageCopyArr.push({
              dockerfile_copyfrom: value.file_name,
              dockerfile_copyto: this.uploadCopyToPath + '/' + value.file_name,
            });
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
        const reader = new FileReader();
        reader.onload = (ev: ProgressEvent) => {
          this.consoleText = (ev.target as FileReader).result as string;
        };
        reader.readAsText(this.selectedDockerFile);
        this.uploadDockerFile().subscribe((res: string) => {
          this.isSelectedDockerFile = true;
          this.consoleText = res;
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
        formData.append('project_name', this.customerNewImage.project_name);
        formData.append('image_name', this.customerNewImage.image_name);
        formData.append('image_tag', this.customerNewImage.image_tag);
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
    if (this.customerNewImage.image_dockerfile.image_base !== '') {
      this.imageService.getDockerFilePreview(this.customerNewImage).subscribe(
        res => this.consoleText = res,
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

  setEnvironment(envsData: Array<EnvType>) {
    const envsArray = this.customerNewImage.image_dockerfile.image_env;
    envsArray.splice(0, envsArray.length);
    envsData.forEach((value: EnvType) => {
      envsArray.push({
        dockerfile_envname: value.envName,
        dockerfile_envvalue: value.envValue,
      });
    });
    this.getDockerFilePreviewInfo();
  }

  removeFile(file: { path: string, file_name: string, size: number }) {
    const fromRemoveData: FormData = new FormData();
    fromRemoveData.append('project_name', this.customerNewImage.project_name);
    fromRemoveData.append('image_name', this.customerNewImage.image_name);
    fromRemoveData.append('image_tag', this.customerNewImage.image_tag);
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
      this.consoleText = '';
      this.imageDetailList.splice(0, this.imageDetailList.length);
      this.customerNewImage.image_dockerfile.image_base = '';
    }
  }

  setBaseImage(selectImage: Image): void {
    this.selectedImage = null;
    this.imageDetailList = null;
    this.isGetImageDetailListWip = true;
    this.imageService.getBoardRegistry().subscribe((res: string) => {
      this.boardRegistry = res.replace(/"/g, '');
      this.imageService.getImageDetailList(selectImage.image_name).subscribe((imageDetails: ImageDetail[]) => {
          this.selectedImage = selectImage;
          this.imageDetailList = imageDetails;
          this.customerNewImage.image_dockerfile.image_base =
            `${this.boardRegistry}/${this.selectedImage.image_name}:${imageDetails[0].image_tag}`;
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
      this.customerNewImage.image_dockerfile.image_base = `${this.boardRegistry}/${this.selectedImage.image_name}:${detail.image_tag}`;
      this.getDockerFilePreviewInfo();
    });
  }
}
