import { Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { AppInitService } from '../app.init.service';
import { MessageService } from '../shared/message-service/message.service';
import { Message } from '../shared/message-service/message';
import { BUTTON_STYLE, GUIDE_STEP, MESSAGE_TARGET } from '../shared/shared.const';
import { Project } from './project';
import { ProjectService } from './project.service';
import { CreateProjectComponent } from './create-project/create-project.component';
import { MemberComponent } from './member/member.component';
import { ActivatedRoute } from "@angular/router";
import { ClrDatagridSortOrder, ClrDatagridStateInterface } from "@clr/angular";

@Component({
  selector: 'project',
  styleUrls: ["./project.component.css"],
  templateUrl: 'project.component.html'
})
export class ProjectComponent implements OnInit, OnDestroy {
  @ViewChild(CreateProjectComponent) createProjectModal;
  @ViewChild(MemberComponent) memberModal;
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
    private activatedRoute: ActivatedRoute,
    private appInitService: AppInitService,
    private projectService: ProjectService,
    private messageService: MessageService) {
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
    this.createProjectModal.openModal();
  }

  editProjectMember(p: Project): void {
    if (this.isSystemAdminOrOwner(p)) {
      this.memberModal.openModal(p);
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

  createProjectClose(step: GUIDE_STEP){
    if (this.isFirstLogin){
      this.appInitService.guideStep = GUIDE_STEP.SERVICE_LIST;
    }
  }
}