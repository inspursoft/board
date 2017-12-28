import { Component, Output, EventEmitter, ViewChild } from '@angular/core';
import { NgForm } from '@angular/forms';
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

  @ViewChild('createProjectForm') projectForm: NgForm;

  @Output() reload: EventEmitter<boolean> = new EventEmitter<boolean>();

  createProject: CreateProject = new CreateProject();

  constructor(
    private projectService: ProjectService,
    private messageService: MessageService
  ){}

  openModal(): void {
    this.createProjectOpened = true;
    this.alertClosed = true;
    this.projectForm.resetForm();
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
        inlineMessage.message = 'PROJECT.SUCCESSFUL_CREATED_PROJECT';
        this.messageService.inlineAlertMessage(inlineMessage);
        this.reload.emit(true);
      })
      .catch(err=>{
        if (err) {
          this.alertClosed = false;
          switch(err.status) {
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

  closeAlert(): void {
    this.alertClosed = true;
  }
}