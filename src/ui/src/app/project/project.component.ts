import { Component, ViewContainerRef } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { TranslateService } from '@ngx-translate/core';
import { ClrDatagridSortOrder, ClrDatagridStateInterface } from '@clr/angular';
import { AppInitService } from '../shared.service/app-init.service';
import { MessageService } from '../shared.service/message.service';
import { GUIDE_STEP } from '../shared/shared.const';
import { ProjectService } from './project.service';
import { SharedActionService } from '../shared.service/shared-action.service';
import { Message, RETURN_STATUS, SharedProject } from '../shared/shared.types';
import { PaginationProject } from './project.types';

@Component({
  styleUrls: ['./project.component.css'],
  templateUrl: 'project.component.html'
})
export class ProjectComponent {
  totalRecordCount: number;
  pageIndex = 1;
  pageSize = 15;
  projects: PaginationProject;
  isInLoading = false;
  descSort = ClrDatagridSortOrder.DESC;
  oldStateInfo: ClrDatagridStateInterface;

  constructor(private appInitService: AppInitService,
              private projectService: ProjectService,
              private messageService: MessageService,
              private sharedActionService: SharedActionService,
              private translateService: TranslateService,
              private selfView: ViewContainerRef) {
    this.projects = new PaginationProject();
  }

  retrieve(state: ClrDatagridStateInterface): void {
    setTimeout(() => {
      if (state) {
        this.isInLoading = true;
        this.oldStateInfo = state;
        this.projectService.getProjects('', this.pageIndex, this.pageSize, state.sort.by as string, state.sort.reverse).subscribe(
          (res: PaginationProject) => {
            this.totalRecordCount = res.pagination.TotalCount;
            this.projects = res;
          },
          () => this.isInLoading = false,
          () => this.isInLoading = false
        );
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

  editProjectMember(project: SharedProject): void {
    if (this.isSystemAdminOrOwner(project)) {
      this.sharedActionService.createProjectMemberComponent(project, this.selfView);
    }
  }

  confirmToDeleteProject(project: SharedProject): void {
    if (this.isSystemAdminOrOwner(project)) {
      this.translateService.get('PROJECT.CONFIRM_TO_DELETE_PROJECT', [project.projectName]).subscribe((msg: string) => {
        this.messageService.showDeleteDialog(msg, 'PROJECT.DELETE_PROJECT').subscribe((message: Message) => {
          if (message.returnStatus === RETURN_STATUS.rsConfirm) {
            this.projectService.deleteProject(project).subscribe(() => {
              this.messageService.showAlert('PROJECT.SUCCESSFUL_DELETE_PROJECT');
              this.retrieve(this.oldStateInfo);
            }, (error: HttpErrorResponse) => {
              if (error.status === 422) {
                this.messageService.showAlert('PROJECT.FAILED_TO_DELETE_PROJECT_ERROR', {alertType: 'warning'});
              }
            });
          }
        });
      });
    }
  }

  toggleProjectPublic(project: SharedProject, $event: MouseEvent): void {
    const oldPublic = project.projectPublic;
    this.projectService.togglePublicity(project.projectId, project.projectPublic === 1 ? 0 : 1).subscribe(() => {
        this.messageService.showAlert('PROJECT.SUCCESSFUL_TOGGLE_PROJECT');
        project.projectPublic = oldPublic === 1 ? 0 : 1;
      }, () => ($event.srcElement as HTMLInputElement).checked = oldPublic === 1
    );
  }

  get isFirstLogin(): boolean {
    return this.appInitService.isFirstLogin;
  }

  get guideStep(): GUIDE_STEP {
    return this.appInitService.guideStep;
  }

  isSystemAdminOrOwner(project: SharedProject): boolean {
    return this.appInitService.currentUser.user_system_admin === 1 ||
      project.projectOwnerId === this.appInitService.currentUser.user_id;
  }

  guideNextStep(step: GUIDE_STEP) {
    this.createProject();
  }

  setGuideNoneStep() {
    this.appInitService.guideStep = GUIDE_STEP.NONE_STEP;
  }

  createProjectClose() {
    if (this.isFirstLogin) {
      this.appInitService.guideStep = GUIDE_STEP.SERVICE_LIST;
    }
  }
}
