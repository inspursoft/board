import { Component, Output, EventEmitter } from '@angular/core';

import { CreateProject } from './create-project';

import { Project } from '../project';
import { ProjectService } from '../project.service';

import { MessageService } from '../../shared/message-service/message.service';
import { Message } from '../../shared/message-service/message';

@Component({
  selector: 'create-project',
  styleUrls: [ './create-project.component.css' ],
  templateUrl: './create-project.component.html'
})
export class CreateProjectComponent {

  createProjectOpened: boolean;
  alertClosed: boolean;
  errorMessage: string;

  @Output() reload: EventEmitter<boolean> = new EventEmitter<boolean>();

  createProject: CreateProject = new CreateProject();

  constructor(
    private projectService: ProjectService,
    private messageService: MessageService
  ){}

  openModal(): void {
    this.createProjectOpened = true;
    this.alertClosed = true;
  }

  confirm(): void {
    let project = new Project();
    project.project_name = this.createProject.projectName;
    project.project_public = this.createProject.publicity ? 1 : 0;
    project.project_comment = this.createProject.comment;

    this.projectService
      .createProject(project)
      .then(resp=>{
        this.createProjectOpened = false;
        let inlineMessage = new Message();
        inlineMessage.message = 'Successful created project.';
        this.messageService.inlineAlertMessage(inlineMessage);
        this.reload.emit(true);
      })
      .catch(err=>{
        if (err) {
          switch(err.status) {
          case 409:
            this.alertClosed = false;
            this.errorMessage = 'Project name already exists.';
            break;
          default:
            this.alertClosed = false;
            this.errorMessage = 'Unknown error.';
            break;
          }
        }
      });
  }

  closeAlert(): void {
    this.alertClosed = true;
  }
}