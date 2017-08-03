import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';

import { AppInitService } from '../app.init.service';
import { MessageService } from '../shared/message-service/message.service';
import { ConfirmationDialogComponent } from '../shared/confirmation-dialog/confirmation-dialog.component';
import { Message } from '../shared/message-service/message';
import { MESSAGE_TARGET, BUTTON_STYLE } from '../shared/shared.const';

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
  currentUser: {[key: string]: any};

  constructor(
    private appInitService: AppInitService,
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
            inlineMessage.message = 'PROJECT.SUCCESSFUL_DELETE_PROJECT';
            this.messageService.inlineAlertMessage(inlineMessage);
            this.retrieve();
          })
          .catch(err=>this.messageService.dispatchError(err, ''));
      }
    });
  }

  ngOnInit(): void {
    this.currentUser = this.appInitService.currentUser;
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
      })
      .catch(err=>this.messageService.dispatchError(err, 'PROJECT.FAILED_TO_RETRIEVE_PROJECTS'));
  }

  createProject(): void {
    this.createProjectModal.openModal();
  }

  editProjectMember(p: Project): void {
    this.memberModal.openModal(p);
  }

  confirmToDeleteProject(p: Project): void {
    let announceMessage = new Message();
    announceMessage.title = 'PROJECT.DELETE_PROJECT';
    announceMessage.message = 'PROJECT.CONFIRM_TO_DELETE_PROJECT';
    announceMessage.params = [p.project_name];
    announceMessage.target = MESSAGE_TARGET.DELETE_PROJECT;
    announceMessage.buttons = BUTTON_STYLE.DELETION;
    announceMessage.data = p;
    this.messageService.announceMessage(announceMessage);
  }

  toggleProjectPublic(p: Project): void {
    p.project_public = (p.project_public === 1 ? 0 : 1);
    let toggleMessage = new Message();
    this.projectService
      .togglePublicity(p)
      .then(()=>{
        toggleMessage.message = 'PROJECT.SUCCESSFUL_TOGGLE_PROJECT'; 
        this.messageService.inlineAlertMessage(toggleMessage);
      })
      .catch(err=>this.messageService.dispatchError(err, ''));
  }
}