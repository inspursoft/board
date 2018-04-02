import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { Router } from '@angular/router';
import { HttpErrorResponse } from "@angular/common/http";
import { K8sService } from "../../service.k8s";
import { MessageService } from "../../../shared/message-service/message.service";
import { Project } from "../../../project/project";
import { Service } from "../../service";
import { AppInitService } from "../../../app.init.service";

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
  errorMessage: string = "";
  successMessage: string = "";
  @Output() onCancelEvent: EventEmitter<any>;

  constructor(private k8sService: K8sService,
              private router: Router,
              private messageService: MessageService,
              private appInitService: AppInitService) {
    this.projectsList = Array<Project>();
    this.onCancelEvent = new EventEmitter<any>();
    this.filesDataMap = new Map<string, File>();
  }

  ngOnInit() {
    this.k8sService.getProjects()
      .then((res: Array<Project>) => {
        let createNewProject: Project = new Project();
        createNewProject.project_name = "IMAGE.CREATE_IMAGE_CREATE_PROJECT";
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
      this.errorMessage = "";
      this.successMessage = "";
      let file: File = fileList[0];
      if (file.name.endsWith(".yaml")) {//Todo:unchecked with ie11
        if (isDeploymentYaml) {
          this.filesDataMap.delete("deployment");
          this.filesDataMap.set("deployment", file);
        } else {
          this.filesDataMap.delete("service");
          this.filesDataMap.set("service", file);
        }
      } else {
        (event.target as HTMLInputElement).value = '';
        this.errorMessage = "SERVICE.SERVICE_YAML_INVALID_FILE";
      }
    }
  }

  clickSelectProject(project: Project) {
    this.router.navigate(["/projects"],{queryParams: {token: this.appInitService.token}, fragment: "create"});
  }

  changeSelectProject(project: Project) {
    this.selectedProjectName = project.project_name;
    this.selectedProjectId = project.project_id;
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
      .catch(err => {
        this.isToggleServiceWIP = false;
        this.messageService.dispatchError(err);
      });
  }

  btnUploadClick(event: MouseEvent) {
    this.errorMessage = "";
    this.successMessage = "";
    let formData = new FormData();
    let deploymentFile = this.filesDataMap.get("deployment");
    let serviceFile = this.filesDataMap.get("service");
    formData.append("deployment_file", deploymentFile, deploymentFile.name);
    formData.append("service_file", serviceFile, serviceFile.name);
    this.isUploadFileWIP = true;
    this.k8sService.uploadServiceYamlFile(this.selectedProjectName, formData)
      .subscribe((res: Service) => {
        this.newServiceName = res.service_name;
        this.newServiceId = res.service_id;
        this.isUploadFileWIP = false;
      }, (error: any) => {
        this.isUploadFileSuccess = false;
        this.isUploadFileWIP = false;
        if (error && error instanceof HttpErrorResponse && (error as HttpErrorResponse).status == 400) {
          this.errorMessage = (error as HttpErrorResponse).error;
        } else {
          this.messageService.dispatchError(error);
        }
      }, () => {
        this.successMessage = "SERVICE.SERVICE_YAML_UPLOAD_SUCCESS";
        this.isUploadFileSuccess = true;
      })
  }

  get isBtnUploadDisabled(): boolean {
    return this.selectedProjectId == 0
      || !this.filesDataMap.has('deployment')
      || this.isUploadFileSuccess
      || !this.filesDataMap.has('service');
  }
}
