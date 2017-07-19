import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';

import { Subscription } from 'rxjs/Subscription';

import { MessageService } from '../shared/message-service/message.service';
import { ConfirmationDialogComponent } from '../shared/confirmation-dialog/confirmation-dialog.component';
import { Message } from '../shared/message-service/message';
import { MESSAGE_TARGET } from '../shared/shared.const';

import { Project } from './project';
import { ProjectService } from './project.service';
import { CreateProjectComponent } from './create-project/create-project.component';
import { MemberComponent } from './member/member.component';


@Component({
  selector: 'project',
  templateUrl: 'project.component.html'
})
export class ProjectComponent implements OnInit, OnDestroy {
  
  projects: Project[];

  @ViewChild(CreateProjectComponent) createProjectModal;
  @ViewChild(MemberComponent) memberModal;

  _subscription: Subscription;

  constructor(
    private projectService: ProjectService,
    private messageService: MessageService
  ){
    this._subscription = this.messageService.messageConfirmed$.subscribe(m=>{
      let confirmationMessage = <Message>m;
      if(confirmationMessage) {
        let project = <Project>confirmationMessage.data;
        this.projectService
          .deleteProject(project)
          .then(()=>{
            let inlineMessage = new Message();
            inlineMessage.message = 'Successful deleted project.';
            this.messageService.inlineAlertMessage(inlineMessage);
            this.retrieve();
          })
          .catch(err=>{
            let globalMessage = new Message();
            globalMessage.message = 'Unexpected error:' + err;
            this.messageService.globalMessage(err);
          });
      }
    });
  }

  ngOnInit(): void {
    this.retrieve();
  }

  ngOnDestroy(): void {
    if(this._subscription) {
      this._subscription.unsubscribe();
    }
  }

  retrieve(): void {
    this.projectService
      .getProjects()
      .then(projects=>{
        this.projects = projects;
      });
  }

  createProject(): void {
    this.createProjectModal.openModal();
  }

  editProjectMember(p: Project): void {
    this.memberModal.openModal(p);
  }

  confirmToDeleteProject(p: Project): void {
    let announceMessage = new Message();
    announceMessage.title = 'Delete Project';
    announceMessage.message = 'Are you sure to delete project?';
    announceMessage.target = MESSAGE_TARGET.DELETE_PROJECT;
    announceMessage.data = p;
    this.messageService.announceMessage(announceMessage);
  }

  toggleProjectPublic(p: Project): void {
    p.project_public = (p.project_public === 1 ? 0 : 1);
    let toggleMessage = new Message();
    toggleMessage.title = 'Toggle Project Public';
    this.projectService
      .togglePublicity(p)
      .then(()=>{
        toggleMessage.message = 'Successful toggle project to ' + ((p.project_public === 1) ? 'public' : 'private'); 
        this.messageService.inlineAlertMessage(toggleMessage);
      })
      .catch(err=>{
        toggleMessage.message = 'Failed to toggle project, due to ' + err.responseText;
        this.messageService.inlineAlertMessage(toggleMessage);
      });
  }
}