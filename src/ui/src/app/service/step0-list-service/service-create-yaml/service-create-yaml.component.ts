import { Component, OnInit, Output, EventEmitter } from '@angular/core';
import { Router } from '@angular/router';
import { HttpErrorResponse } from "@angular/common/http";
import { K8sService } from "../../service.k8s";
import { MessageService } from "../../../shared/message-service/message.service";
import { Project } from "../../../project/project";
import { Observable } from "rxjs/Observable";
import "rxjs/add/observable/concat";

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
  filesDataMap: Map<string, File>;
  isAlreadyCheck: boolean = false;
  isCheckWIP: boolean = false;
  isRightYamlFile: boolean = false;
  isUploadFileSuccess: boolean = false;
  errorMessage: string = "";
  successMessage: string = "";
  @Output() onCancelEvent: EventEmitter<any>;

  constructor(private k8sService: K8sService,
              private router: Router,
              private messageService: MessageService) {
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
        if (this.filesDataMap.has("deployment") && this.filesDataMap.has("service")) {
          this.checkYamlFile();
        }
      } else {
        (event.target as HTMLInputElement).value = '';
        this.errorMessage = "SERVICE.SERVICE_YAML_INVALID_FILE";
      }
    }
  }

  checkYamlFile() {
    this.isCheckWIP = true;
    let formData = new FormData();
    let deploymentFile = this.filesDataMap.get("deployment");
    let serviceFile = this.filesDataMap.get("service");
    formData.append("deployment_file", deploymentFile, deploymentFile.name);
    formData.append("service_file", serviceFile, serviceFile.name);
    this.k8sService.checkCreateServiceYaml(this.selectedProjectName, formData)
      .subscribe((res: Object) => {//Todo:get service name from res

        },
        (error: HttpErrorResponse) => {
          this.isCheckWIP = false;
          this.isAlreadyCheck = true;
          this.isRightYamlFile = false;
          this.messageService.dispatchError(error)
        },
        () => {
          this.isCheckWIP = false;
          this.isAlreadyCheck = true;
          this.isRightYamlFile = true;
          this.successMessage = "SERVICE.SERVICE_YAML_VALID_FILE";
          this.newServiceName = "This is new service name";
        });
  }

  clickSelectProject(project: Project) {
    this.router.navigate(["/projects"]);
  }

  changeSelectProject(project: Project) {
    this.selectedProjectName = project.project_name;
    this.selectedProjectId = project.project_id;
  }

  btnCancelClick(event: MouseEvent) {
    this.onCancelEvent.emit(event);
  }

  btnUploadClick(event: MouseEvent) {
    this.errorMessage = "";
    this.successMessage = "";
    let obs1 = this.uploadOneYamlFile("deployment");
    let obs2 = this.uploadOneYamlFile("service");
    let obsConcat = Observable.concat(obs1, obs2);
    obsConcat.subscribe((res: Object) => {//Todo:get service id from res

    }, (error: any) => {
      this.isUploadFileSuccess = false;
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

  uploadOneYamlFile(yamlType: "deployment" | "service"): Observable<Object> {
    let formData = new FormData();
    let file: File = this.filesDataMap.get(yamlType);
    formData.append("upload_file", file, file.name);
    return this.k8sService.uploadServiceYamlFile(this.newServiceName, this.selectedProjectName, formData, yamlType);
  }

  get isBtnUploadDisabled(): boolean {
    return this.selectedProjectId == 0
      || !this.isAlreadyCheck
      || this.isUploadFileSuccess
      || !this.isRightYamlFile;
  }
}
