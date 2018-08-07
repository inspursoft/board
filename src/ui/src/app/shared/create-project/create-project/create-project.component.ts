import { Component } from '@angular/core';
import { CreateProject, Project } from "../../../project/project";
import { SharedService } from "../../shared.service";
import { Message } from "../../message-service/message";
import { MessageService } from "../../message-service/message.service";
import { HttpErrorResponse } from "@angular/common/http";
import { Subject } from "rxjs/Subject";
import { Observable } from "rxjs/Observable";
import { CsComponentBase } from "../../cs-components-library/cs-component-base";

@Component({
  selector: 'create-project',
  styleUrls: [ './create-project.component.css' ],
  templateUrl: './create-project.component.html'
})
export class CreateProjectComponent extends CsComponentBase {
  _createProjectOpened: boolean = false;
  alertClosed: boolean;
  errorMessage: string = "";
  createProject: CreateProject;
  closeNotification: Subject<string>;
  isCreateProjectWIP: boolean = false;
  projectNamePattern: string = '^[a-z0-9]+(?:[-][a-z0-9]+)*$';
  constructor(private sharedService: SharedService,
              private messageService: MessageService) {
    super();
    this.createProject = new CreateProject();
    this.closeNotification = new Subject<string>();
  }

  get createProjectOpened(): boolean{
    return this._createProjectOpened;
  }
  set createProjectOpened(value:boolean){
    this._createProjectOpened = value;
    if (!value){
      this.closeNotification.next();
    }
  }

  openModal(): Observable<string> {
    this.createProjectOpened = true;
    this.alertClosed = true;
    return this.closeNotification.asObservable();
  }

  confirm(): void {
    if (this.verifyInputValid()) {
      this.isCreateProjectWIP = true;
      let project = new Project();
      project.project_name = this.createProject.projectName;
      project.project_public = this.createProject.publicity ? 1 : 0;
      project.project_comment = this.createProject.comment;
      this.sharedService.createProject(project)
        .then(() => {
          this.isCreateProjectWIP = false;
          this.createProjectOpened = false;
          let inlineMessage = new Message();
          inlineMessage.message = 'PROJECT.SUCCESSFUL_CREATED_PROJECT';
          this.messageService.inlineAlertMessage(inlineMessage);
          this.closeNotification.next(project.project_name);
        })
        .catch((err: HttpErrorResponse) => {
          this.isCreateProjectWIP = false;
          if (err) {
            this.alertClosed = false;
            switch (err.status) {
              case 409:
                this.errorMessage = 'PROJECT.PROJECT_NAME_ALREADY_EXISTS';
                break;
              case 400:
                this.errorMessage = 'PROJECT.PROJECT_NAME_IS_ILLEGAL';
                break;
              default:
                this.errorMessage = 'ERROR.INTERNAL_ERROR';
            }
          }
        });
    }
  }
}