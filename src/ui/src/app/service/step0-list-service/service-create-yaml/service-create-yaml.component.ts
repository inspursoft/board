import { Component, EventEmitter, OnInit, Output, ViewContainerRef } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { K8sService } from '../../service.k8s';
import { MessageService } from '../../../shared.service/message.service';
import { SharedService } from '../../../shared.service/shared.service';
import { SharedActionService } from '../../../shared.service/shared-action.service';
import { EXECUTE_STATUS, GlobalAlertType } from '../../../shared/shared.types';
import { Service, ServiceProject } from '../../service.types';

export const DEPLOYMENT = 'deployment';
export const SERVICE = 'service';
type FileType = 'deployment' | 'service';

@Component({
  selector: 'app-service-create-yaml',
  templateUrl: './service-create-yaml.component.html',
  styleUrls: ['./service-create-yaml.component.css']
})
export class ServiceCreateYamlComponent implements OnInit {
  selectedProjectName = '';
  selectedProjectId = 0;
  projectsList: Array<ServiceProject>;
  newServiceName = '';
  newServiceId = 0;
  filesDataMap: Map<FileType, Blob>;
  uploadFileStatus = EXECUTE_STATUS.esNotExe;
  createServiceStatus = EXECUTE_STATUS.esNotExe;
  isFileInEdit = false;
  curFileContent = '';
  curFileName: FileType;
  @Output() cancelEvent: EventEmitter<any>;

  constructor(private k8sService: K8sService,
              private selfView: ViewContainerRef,
              private sharedService: SharedService,
              private sharedActionService: SharedActionService,
              private messageService: MessageService) {
    this.projectsList = Array<ServiceProject>();
    this.cancelEvent = new EventEmitter<any>();
    this.filesDataMap = new Map<FileType, Blob>();
  }

  ngOnInit() {
    this.k8sService.getProjects().subscribe((res: Array<ServiceProject>) => this.projectsList = res);
  }

  get cancelBtnCaption(): string {
    return (this.createServiceStatus === EXECUTE_STATUS.esFailed ||
      this.uploadFileStatus === EXECUTE_STATUS.esFailed ||
      this.uploadFileStatus === EXECUTE_STATUS.esSuccess) ? 'BUTTON.DELETE' : 'BUTTON.CANCEL';
  }

  get createBtnDisabled(): boolean {
    return this.createServiceStatus === EXECUTE_STATUS.esExecuting ||
      this.uploadFileStatus === EXECUTE_STATUS.esExecuting ||
      this.uploadFileStatus !== EXECUTE_STATUS.esSuccess;
  }

  get uploadFileBtnDisabled(): boolean {
    return this.selectedProjectId === 0 ||
      this.uploadFileStatus === EXECUTE_STATUS.esExecuting ||
      this.uploadFileStatus === EXECUTE_STATUS.esSuccess;
  }

  get isBtnUploadDisabled(): boolean {
    return this.selectedProjectId === 0
      || this.uploadFileStatus === EXECUTE_STATUS.esExecuting
      || this.uploadFileStatus === EXECUTE_STATUS.esSuccess
      || !this.filesDataMap.has(DEPLOYMENT)
      || this.isFileInEdit
      || !this.filesDataMap.has(SERVICE);
  }

  get isEditDeploymentEnable(): boolean {
    return this.uploadFileStatus === EXECUTE_STATUS.esNotExe
      && !this.isFileInEdit
      && this.filesDataMap.get(DEPLOYMENT) !== undefined;
  }

  get isEditServiceEnable(): boolean {
    return this.uploadFileStatus === EXECUTE_STATUS.esNotExe
      && !this.isFileInEdit
      && this.filesDataMap.get(SERVICE) !== undefined;
  }

  uploadFile(event: Event, isDeploymentYaml: boolean) {
    const fileList: FileList = (event.target as HTMLInputElement).files;
    if (fileList.length > 0) {
      const file: File = fileList[0];
      if (file.name.endsWith('.yaml')) {// Todo:unchecked with ie11
        if (isDeploymentYaml) {
          this.filesDataMap.delete(DEPLOYMENT);
          this.filesDataMap.set(DEPLOYMENT, file);
        } else {
          this.filesDataMap.delete(SERVICE);
          this.filesDataMap.set(SERVICE, file);
        }
      } else {
        (event.target as HTMLInputElement).value = '';
        this.messageService.showAlert('SERVICE.SERVICE_YAML_INVALID_FILE', {alertType: 'warning'});
      }
    } else {
      isDeploymentYaml ? this.filesDataMap.delete(DEPLOYMENT) : this.filesDataMap.delete(SERVICE);
    }
  }

  clickSelectProject() {
    this.sharedActionService.createProjectComponent(this.selfView).subscribe((projectName: string) => {
      if (projectName) {
        this.k8sService.getOneProject(projectName).subscribe((res: Array<ServiceProject>) => {
          this.selectedProjectId = res[0].projectId;
          this.selectedProjectName = res[0].projectName;
          this.projectsList.unshift(res[0]);
        });
      }
    });
  }

  changeSelectProject(project: ServiceProject) {
    this.selectedProjectName = project.projectName;
    this.selectedProjectId = project.projectId;
  }

  btnCancelClick(event: MouseEvent) {
    if (this.createServiceStatus === EXECUTE_STATUS.esFailed || (
      this.uploadFileStatus === EXECUTE_STATUS.esSuccess &&
      this.createServiceStatus === EXECUTE_STATUS.esNotExe)) {
      this.k8sService.deleteService(this.newServiceId).subscribe(
        () => this.cancelEvent.emit(event),
        () => this.cancelEvent.emit(event)
      );
    } else {
      this.cancelEvent.emit(event);
    }
  }

  btnCreateClick(event: MouseEvent) {
    this.createServiceStatus = EXECUTE_STATUS.esExecuting;
    this.k8sService.toggleServiceStatus(this.newServiceId, 1).subscribe(
      () => {
        this.createServiceStatus = EXECUTE_STATUS.esSuccess;
        this.cancelEvent.emit(event);
      },
      (err: HttpErrorResponse) => {
        this.createServiceStatus = EXECUTE_STATUS.esFailed;
        this.btnCancelClick(event); // issue#2270
        this.messageService.showGlobalMessage('SERVICE.SERVICE_YAML_CREATE_FAILED', {
          errorObject: err,
          globalAlertType: GlobalAlertType.gatShowDetail
        });
      });
  }

  btnUploadClick() {
    const formData = new FormData();
    const deploymentFile = this.filesDataMap.get(DEPLOYMENT);
    const serviceFile = this.filesDataMap.get(SERVICE);
    formData.append('deployment_file', deploymentFile, `${DEPLOYMENT}.yaml`);
    formData.append('service_file', serviceFile, `${SERVICE}.yaml`);
    this.uploadFileStatus = EXECUTE_STATUS.esExecuting;
    this.k8sService.uploadServiceYamlFile(this.selectedProjectName, formData)
      .subscribe((res: Service) => {
        this.newServiceName = res.serviceName;
        this.newServiceId = res.serviceId;
      }, (error: HttpErrorResponse) => {
        this.uploadFileStatus = EXECUTE_STATUS.esFailed;
        this.messageService.showGlobalMessage('SERVICE.SERVICE_YAML_UPLOAD_FAILED', {
          errorObject: error,
          globalAlertType: GlobalAlertType.gatShowDetail
        });
      }, () => {
        this.uploadFileStatus = EXECUTE_STATUS.esSuccess;
        this.messageService.showAlert('SERVICE.SERVICE_YAML_UPLOAD_SUCCESS');
      });
  }

  editFile(fileName: FileType): void {
    this.isFileInEdit = true;
    this.curFileName = fileName;
    const file = this.filesDataMap.get(fileName);
    const reader = new FileReader();
    reader.onload = (ev: ProgressEvent) => {
      this.curFileContent = (ev.target as FileReader).result as string;
    };
    reader.readAsText(file);
  }

  saveFile(): void {
    this.isFileInEdit = false;
    this.filesDataMap.delete(this.curFileName);
    try {
      const writer = new File(Array.from(this.curFileContent), this.curFileName);
      this.filesDataMap.set(this.curFileName, writer);
    } catch (e) {
      const writer = new MSBlobBuilder();
      writer.append(this.curFileContent);
      this.filesDataMap.set(this.curFileName, writer.getBlob());
    }
  }
}
