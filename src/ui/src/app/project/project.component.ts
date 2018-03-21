import { Component, OnInit, OnDestroy, ViewChild, ChangeDetectorRef } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { AppInitService } from '../app.init.service';
import { MessageService } from '../shared/message-service/message.service';
import { Message } from '../shared/message-service/message';
import { MESSAGE_TARGET, BUTTON_STYLE, GUIDE_STEP } from '../shared/shared.const';
import { Project } from './project';
import { ProjectService } from './project.service';
import { CreateProjectComponent } from './create-project/create-project.component';
import { MemberComponent } from './member/member.component';
import { ActivatedRoute } from "@angular/router";

@Component({
  selector: 'project',
  templateUrl: 'project.component.html'
})
export class ProjectComponent implements OnInit, OnDestroy {
  totalRecordCount: number;
  pageIndex: number = 1;  
  pageSize: number = 15;
  projects: Project[];

  @ViewChild(CreateProjectComponent) createProjectModal;
  @ViewChild(MemberComponent) memberModal;

  _subscription: Subscription;
  currentUser: {[key: string]: any};
  isInLoading:boolean = false;
  constructor(
    private activatedRoute:ActivatedRoute,
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
    this.activatedRoute.fragment.subscribe((res)=>{
      if (res && res =="create"){
        this.createProject();
      }
    });
  }

  ngOnDestroy(): void {
    if(this._subscription) {
      this._subscription.unsubscribe();
    }
  }
  retrieve(): void {
    setTimeout(()=>{
      this.isInLoading = true;
      this.projectService
        .getProjects('', this.pageIndex, this.pageSize)
        .then(paginatedProjects=>{
          this.totalRecordCount = paginatedProjects.pagination.total_count;
          this.projects = paginatedProjects.project_list;
          this.isInLoading = false;
        })
        .catch(err=>{
          this.messageService.dispatchError(err, 'PROJECT.FAILED_TO_RETRIEVE_PROJECTS');
          this.isInLoading = false;
        });
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
    announceMessage.title = 'PROJECT.DELETE_PROJECT';
    announceMessage.message = 'PROJECT.CONFIRM_TO_DELETE_PROJECT';
    announceMessage.params = [p.project_name];
    announceMessage.target = MESSAGE_TARGET.DELETE_PROJECT;
    announceMessage.buttons = BUTTON_STYLE.DELETION;
    announceMessage.data = p;
    this.messageService.announceMessage(announceMessage);
  }

  toggleProjectPublic(project: Project, $event:MouseEvent): void {
    let oldPublic = project.project_public;
    this.projectService
      .togglePublicity(project.project_id, project.project_public === 1? 0 : 1)
      .then(()=>{
        let toggleMessage = new Message();
        toggleMessage.message = 'PROJECT.SUCCESSFUL_TOGGLE_PROJECT'; 
        this.messageService.inlineAlertMessage(toggleMessage);
        project.project_public = oldPublic === 1 ? 0 : 1;
      })
      .catch(err=>{
        ($event.srcElement as HTMLInputElement).checked = oldPublic === 1;
        this.messageService.dispatchError(err);
      });
  }

  get isFirstLogin(): boolean{
    return this.appInitService.isFirstLogin;
  }

  get guideStep(): GUIDE_STEP{
    return this.appInitService.guideStep;
  }

  guideNextStep(step:GUIDE_STEP){
    this.createProject();
  }

  setGuideNoneStep(){
     this.appInitService.guideStep = GUIDE_STEP.NONE_STEP;
  }

  createProjectClose(step: GUIDE_STEP){
    if (this.isFirstLogin){
      this.appInitService.guideStep = GUIDE_STEP.SERVICE_LIST;
    }
  }
}