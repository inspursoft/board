import { Component, EventEmitter, OnInit, Output, ViewContainerRef } from '@angular/core';
import { K8sService } from "../../service.k8s";
import { MessageService } from "../../../shared.service/message.service";
import { Project } from "../../../project/project";
import { Service } from "../../service";
import { HttpErrorResponse } from "@angular/common/http";
import { SharedService } from "../../../shared.service/shared.service";
import { SharedActionService } from "../../../shared.service/shared-action.service";
import { EXECUTE_STATUS, GlobalAlertType } from "../../../shared/shared.types";

export const DEPLOYMENT = "deployment";
export const SERVICE = "service";
type FileType = "deployment" | "service";

@Component({
  selector: 'service-create-yaml',
  templateUrl: './service-create-yaml.component.html',
  styleUrls: ['./service-create-yaml.component.css']
})
export class ServiceCreateYamlComponent implements OnInit {
  selectedProjectName: string = "";
  selectedProjectId: number = 0;
  projectsList: Array<Project>;
  newServiceName: string = "";
  newServiceId: number = 0;
  filesDataMap: Map<FileType, Blob>;
  uploadFileStatus = EXECUTE_STATUS.esNotExe;
  createServiceStatus = EXECUTE_STATUS.esNotExe;
  isFileInEdit: boolean = false;
  curFileContent: string = "";
  curFileName: FileType;
  @Output() onCancelEvent: EventEmitter<any>;

  constructor(private k8sService: K8sService,
              private selfView: ViewContainerRef,
              private sharedService: SharedService,
              private sharedActionService: SharedActionService,
              private messageService: MessageService) {
    this.projectsList = Array<Project>();
    this.onCancelEvent = new EventEmitter<any>();
    this.filesDataMap = new Map<FileType, Blob>();
  }

  ngOnInit() {
    this.k8sService.getProjects().subscribe((res: Array<Project>) => this.projectsList = res)
  }

  uploadFile(event: Event, isDeploymentYaml: boolean) {
    let fileList: FileList = (event.target as HTMLInputElement).files;
    if (fileList.length > 0) {
      let file: File = fileList[0];
      if (file.name.endsWith(".yaml")) {//Todo:unchecked with ie11
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

  setDropdownDefaultText(): void {
    let selected = this.projectsList.find((project: Project) => project.project_id === this.selectedProjectId);
  }

  clickSelectProject() {
    this.sharedActionService.createProjectComponent(this.selfView).subscribe((projectName: string) => {
      if (projectName) {
        this.sharedService.getOneProject(projectName).subscribe((res: Array<Project>) => {
          this.selectedProjectId = res[0].project_id;
          this.selectedProjectName = res[0].project_name;
          this.projectsList.unshift(res[0]);
        })
      }
    });
  }

  changeSelectProject(project: Project) {
    this.selectedProjectName = project.project_name;
    this.selectedProjectId = project.project_id;
  }

  btnCancelClick(event: MouseEvent) {
    if (this.createServiceStatus == EXECUTE_STATUS.esFailed){
      this.k8sService.deleteService(this.newServiceId).subscribe(
        ()=>this.onCancelEvent.emit(event),
        ()=>this.onCancelEvent.emit(event)
      );
    } else {
      this.onCancelEvent.emit(event);
    }
  }

  btnCreateClick(event: MouseEvent) {
    this.createServiceStatus = EXECUTE_STATUS.esExecuting;
    this.k8sService.toggleServiceStatus(this.newServiceId, 1).subscribe(
      () => {
        this.createServiceStatus = EXECUTE_STATUS.esSuccess;
        this.onCancelEvent.emit(event);
      },
      (err: HttpErrorResponse) => {
        this.createServiceStatus = EXECUTE_STATUS.esFailed;
        this.messageService.showGlobalMessage('SERVICE.SERVICE_YAML_CREATE_FAILED', {
          errorObject: err,
          globalAlertType: GlobalAlertType.gatShowDetail
        })
      });
  }

  btnUploadClick() {
    let formData = new FormData();
    let deploymentFile = this.filesDataMap.get(DEPLOYMENT);
    let serviceFile = this.filesDataMap.get(SERVICE);
    formData.append("deployment_file", deploymentFile, `${DEPLOYMENT}.yaml`);
    formData.append("service_file", serviceFile, `${SERVICE}.yaml`);
    this.uploadFileStatus = EXECUTE_STATUS.esExecuting;
    this.k8sService.uploadServiceYamlFile(this.selectedProjectName, formData)
      .subscribe((res: Service) => {
        this.newServiceName = res.service_name;
        this.newServiceId = res.service_id;
      }, (error: HttpErrorResponse) => {
        this.uploadFileStatus = EXECUTE_STATUS.esFailed;
        this.messageService.showGlobalMessage('SERVICE.SERVICE_YAML_UPLOAD_FAILED', {
          errorObject: error,
          globalAlertType: GlobalAlertType.gatShowDetail
        })
      }, () => {
        this.uploadFileStatus = EXECUTE_STATUS.esSuccess;
        this.messageService.showAlert('SERVICE.SERVICE_YAML_UPLOAD_SUCCESS');
      });
  }

  get isBtnUploadDisabled(): boolean {
    return this.selectedProjectId == 0
      || this.uploadFileStatus == EXECUTE_STATUS.esExecuting
      || this.uploadFileStatus == EXECUTE_STATUS.esSuccess
      || !this.filesDataMap.has(DEPLOYMENT)
      || this.isFileInEdit
      || !this.filesDataMap.has(SERVICE);
  }

  get isEditDeploymentEnable(): boolean{
    return this.uploadFileStatus == EXECUTE_STATUS.esNotExe
      && !this.isFileInEdit
      && this.filesDataMap.get(DEPLOYMENT) != undefined
  }

  get isEditServiceEnable(): boolean{
    return this.uploadFileStatus == EXECUTE_STATUS.esNotExe
      && !this.isFileInEdit
      && this.filesDataMap.get(SERVICE) != undefined
  }

  editFile(fileName: FileType): void {
    this.isFileInEdit = true;
    this.curFileName = fileName;
    let file = this.filesDataMap.get(fileName);
    let reader = new FileReader();
    reader.onload = (ev: ProgressEvent) => {
      this.curFileContent = (ev.target as FileReader).result as string;
    };
    reader.readAsText(file);
  }

  saveFile():void{
    this.isFileInEdit = false;
    this.filesDataMap.delete(this.curFileName);
    try {
      let writer = new File(Array.from(this.curFileContent), this.curFileName);
      this.filesDataMap.set(this.curFileName, writer);
    } catch (e) {
      let writer = new MSBlobBuilder();
      writer.append(this.curFileContent);
      this.filesDataMap.set(this.curFileName, writer.getBlob());
    }
  }
}
