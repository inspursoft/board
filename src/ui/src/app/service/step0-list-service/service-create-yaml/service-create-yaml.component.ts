import { Component, EventEmitter, OnInit, Output, ViewContainerRef } from '@angular/core';
import { K8sService } from "../../service.k8s";
import { MessageService } from "../../../shared/message-service/message.service";
import { Project } from "../../../project/project";
import { Service } from "../../service";
import { AppInitService } from "../../../app.init.service";
import { MESSAGE_TYPE } from "../../../shared/shared.const";
import { Message } from "../../../shared/message-service/message";
import { HttpErrorResponse } from "@angular/common/http";
import { SharedService } from "../../../shared/shared.service";
import { SharedActionService } from "../../../shared/shared-action.service";

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
  filesDataMap: Map<string, File>;
  isUploadFileWIP: boolean = false;
  isToggleServiceWIP: boolean = false;
  isUploadFileSuccess: boolean = false;
  isFileInEdit: boolean = false;
  curFileContent: string = "";
  curFileName: FileType;
  dropdownDefaultText: string;
  @Output() onCancelEvent: EventEmitter<any>;

  constructor(private k8sService: K8sService,
              private selfView: ViewContainerRef,
              private sharedService: SharedService,
              private sharedActionService: SharedActionService,
              private messageService: MessageService,
              private appInitService: AppInitService) {
    this.projectsList = Array<Project>();
    this.onCancelEvent = new EventEmitter<any>();
    this.filesDataMap = new Map<string, File>();
  }

  ngOnInit() {
    this.dropdownDefaultText = "IMAGE.CREATE_IMAGE_SELECT_PROJECT";
    this.k8sService.getProjects()
      .then((res: Array<Project>) => {
        let createNewProject: Project = new Project();
        createNewProject.project_name = "IMAGE.CREATE_IMAGE_CREATE_PROJECT";
        createNewProject.project_id = -1;
        createNewProject["isSpecial"] = true;
        createNewProject["OnlyClick"] = true;
        this.projectsList.push(createNewProject);
        if (res && res.length > 0) {
          this.projectsList = this.projectsList.concat(res);
        }
      }).catch(err => this.messageService.dispatchError(err));
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
        let msg = new Message();
        msg.type = MESSAGE_TYPE.COMMON_ERROR;
        msg.message = "SERVICE.SERVICE_YAML_INVALID_FILE";
        this.messageService.inlineAlertMessage(msg);
      }
    }
  }

  setDropdownDefaultText(): void {
    let selected = this.projectsList.find((project: Project) => project.project_id === this.selectedProjectId);
    this.dropdownDefaultText = selected ? selected.project_name : "IMAGE.CREATE_IMAGE_CREATE_PROJECT";
  }

  clickSelectProject() {
    this.sharedActionService.createProjectComponent(this.selfView).subscribe((projectName: string) => {
      if (projectName) {
        this.sharedService.getOneProject(projectName).then((res: Array<Project>) => {
          this.selectedProjectId = res[0].project_id;
          this.selectedProjectName = res[0].project_name;
          let project = this.projectsList.shift();
          this.projectsList.unshift(res[0]);
          this.projectsList.unshift(project);
          this.setDropdownDefaultText();
        })
      }
    });
  }

  changeSelectProject(project: Project) {
    this.selectedProjectName = project.project_name;
    this.selectedProjectId = project.project_id;
    this.setDropdownDefaultText();
  }

  btnCancelClick(event: MouseEvent) {
    this.onCancelEvent.emit(event);
  }

  btnCreateClick(event: MouseEvent) {
    this.isToggleServiceWIP = true;
    this.k8sService.toggleServiceStatus(this.newServiceId, 1)
      .then(() => {
        this.isToggleServiceWIP = false;
        this.onCancelEvent.emit(event);
      })
      .catch((err:HttpErrorResponse) => {
        this.isToggleServiceWIP = false;
        let msg = new Message();
        msg.type = MESSAGE_TYPE.SHOW_DETAIL;
        msg.message = err.error;
        msg.errorObject = err;
        this.messageService.globalMessage(msg);
      });
  }

  btnUploadClick() {
    let formData = new FormData();
    let deploymentFile = this.filesDataMap.get(DEPLOYMENT);
    let serviceFile = this.filesDataMap.get(SERVICE);
    formData.append("deployment_file", deploymentFile, deploymentFile.name);
    formData.append("service_file", serviceFile, serviceFile.name);
    this.isUploadFileWIP = true;
    this.k8sService.uploadServiceYamlFile(this.selectedProjectName, formData)
      .subscribe((res: Service) => {
        this.newServiceName = res.service_name;
        this.newServiceId = res.service_id;
        this.isUploadFileWIP = false;
      }, (error: HttpErrorResponse) => {
        this.isUploadFileSuccess = false;
        this.isUploadFileWIP = false;
        let msg = new Message();
        msg.type = MESSAGE_TYPE.SHOW_DETAIL;
        msg.message = error.error;
        msg.errorObject = error;
        this.messageService.globalMessage(msg);
      }, () => {
        this.isUploadFileSuccess = true;
        let msg = new Message();
        msg.type = MESSAGE_TYPE.COMMON_ERROR;
        msg.message = "SERVICE.SERVICE_YAML_UPLOAD_SUCCESS";
        this.messageService.inlineAlertMessage(msg);
      });
  }

  get isBtnUploadDisabled(): boolean {
    return this.selectedProjectId == 0
      || !this.filesDataMap.has('deployment')
      || this.isUploadFileSuccess
      || this.isFileInEdit
      || !this.filesDataMap.has('service');
  }

  get isEditDeploymentEnable(): boolean{
    return !this.isUploadFileSuccess
      && !this.isUploadFileWIP
      && !this.isFileInEdit
      && this.filesDataMap.get('deployment') != undefined
  }

  get isEditServiceEnable(): boolean{
    return !this.isUploadFileSuccess
      && !this.isUploadFileWIP
      && !this.isFileInEdit
      && this.filesDataMap.get('service') != undefined
  }

  editFile(fileName: FileType): void {
    this.isFileInEdit = true;
    this.curFileName = fileName;
    let file = this.filesDataMap.get(fileName);
    let reader = new FileReader();
    reader.onload = (ev: ProgressEvent) => {
      this.curFileContent = (ev.target as FileReader).result;
    };
    reader.readAsText(file);
  }

  saveFile():void{
    this.isFileInEdit = false;
    this.filesDataMap.delete(this.curFileName);
    let writer = new File(Array.from(this.curFileContent), this.curFileName);
    this.filesDataMap.set(this.curFileName, writer);
  }

}
