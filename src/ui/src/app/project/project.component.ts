import { Component, OnDestroy, OnInit, ViewContainerRef } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { AppInitService } from '../app.init.service';
import { MessageService } from '../shared/message-service/message.service';
import { Message } from '../shared/message-service/message';
import { BUTTON_STYLE, GUIDE_STEP, MESSAGE_TARGET } from '../shared/shared.const';
import { Project } from './project';
import { ProjectService } from './project.service';
import { ClrDatagridSortOrder, ClrDatagridStateInterface } from "@clr/angular";
import { SharedActionService } from "../shared/shared-action.service";

@Component({
  selector: 'project',
  styleUrls: ["./project.component.css"],
  templateUrl: 'project.component.html'
})
export class ProjectComponent implements OnInit, OnDestroy {
  _subscription: Subscription;
  totalRecordCount: number;
  pageIndex: number = 1;  
  pageSize: number = 15;
  projects: Project[];
  currentUser: {[key: string]: any};
  isInLoading:boolean = false;
  descSort = ClrDatagridSortOrder.DESC;
  oldStateInfo: ClrDatagridStateInterface;
  constructor(
    private appInitService: AppInitService,
    private projectService: ProjectService,
    private messageService: MessageService,
    private sharedActionService: SharedActionService,
    private selfView: ViewContainerRef) {
    this._subscription = this.messageService.messageConfirmed$.subscribe((msg: Message) => {
      if (msg.target == MESSAGE_TARGET.DELETE_PROJECT) {
        let project = <Project>msg.data;
        this.projectService
          .deleteProject(project)
          .then(() => {
            let inlineMessage = new Message();
            inlineMessage.message = 'PROJECT.SUCCESSFUL_DELETE_PROJECT';
            this.messageService.inlineAlertMessage(inlineMessage);
            this.retrieve(this.oldStateInfo);
          })
          .catch(err => this.messageService.dispatchError(err, ''));
      }
    });
  }

  ngOnInit(): void {
    this.currentUser = this.appInitService.currentUser;
  }

  ngOnDestroy(): void {
    if(this._subscription) {
      this._subscription.unsubscribe();
    }
  }

  retrieve(state: ClrDatagridStateInterface): void {
    setTimeout(() => {
      if (state) {
        this.isInLoading = true;
        this.oldStateInfo = state;
        this.projectService
          .getProjects('', this.pageIndex, this.pageSize, state.sort.by as string, state.sort.reverse)
          .then(paginatedProjects => {
            this.totalRecordCount = paginatedProjects.pagination.total_count;
            this.projects = paginatedProjects.project_list;
            this.isInLoading = false;
          })
          .catch(err => {
            this.messageService.dispatchError(err, 'PROJECT.FAILED_TO_RETRIEVE_PROJECTS');
            this.isInLoading = false;
          });
      }
    });
  }

  createProject(): void {
    this.sharedActionService.createProjectComponent(this.selfView).subscribe((projectName: string) => {
      this.createProjectClose();
      if (projectName) {
        this.retrieve(this.oldStateInfo);
      }
    });
  }

  editProjectMember(project: Project): void {
    if (this.isSystemAdminOrOwner(project)) {
      this.sharedActionService.createProjectMemberComponent(project, this.selfView);
    }
  }

  confirmToDeleteProject(p: Project): void {
    if (this.isSystemAdminOrOwner(p)) {
      let announceMessage = new Message();
      announceMessage.title = 'PROJECT.DELETE_PROJECT';
      announceMessage.message = 'PROJECT.CONFIRM_TO_DELETE_PROJECT';
      announceMessage.params = [p.project_name];
      announceMessage.target = MESSAGE_TARGET.DELETE_PROJECT;
      announceMessage.buttons = BUTTON_STYLE.DELETION;
      announceMessage.data = p;
      this.messageService.announceMessage(announceMessage);
    }
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

  isSystemAdminOrOwner(project: Project): boolean {
    if (this.appInitService.currentUser) {
      return this.appInitService.currentUser["user_system_admin"] == 1 ||
        project.project_owner_id == this.appInitService.currentUser["user_id"];
    }
    return false;
  }

  guideNextStep(step:GUIDE_STEP){
    this.createProject();
  }

  setGuideNoneStep(){
     this.appInitService.guideStep = GUIDE_STEP.NONE_STEP;
  }

  createProjectClose() {
    if (this.isFirstLogin){
      this.appInitService.guideStep = GUIDE_STEP.SERVICE_LIST;
    }
  }
}