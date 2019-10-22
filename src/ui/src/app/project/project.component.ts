import { Component, ViewContainerRef } from '@angular/core';
import { HttpErrorResponse } from "@angular/common/http";
import { TranslateService } from "@ngx-translate/core";
import { ClrDatagridSortOrder, ClrDatagridStateInterface } from "@clr/angular";
import { AppInitService } from '../shared.service/app-init.service';
import { MessageService } from '../shared.service/message.service';
import { GUIDE_STEP } from '../shared/shared.const';
import { Project } from './project';
import { ProjectService } from './project.service';
import { SharedActionService } from "../shared.service/shared-action.service";
import { Message, RETURN_STATUS } from "../shared/shared.types";

@Component({
  selector: 'project',
  styleUrls: ["./project.component.css"],
  templateUrl: 'project.component.html'
})
export class ProjectComponent {
  totalRecordCount: number;
  pageIndex: number = 1;
  pageSize: number = 15;
  projects: Project[];
  isInLoading: boolean = false;
  descSort = ClrDatagridSortOrder.DESC;
  oldStateInfo: ClrDatagridStateInterface;

  constructor(private appInitService: AppInitService,
              private projectService: ProjectService,
              private messageService: MessageService,
              private sharedActionService: SharedActionService,
              private translateService: TranslateService,
              private selfView: ViewContainerRef) {
  }

  retrieve(state: ClrDatagridStateInterface): void {
    setTimeout(() => {
      if (state) {
        this.isInLoading = true;
        this.oldStateInfo = state;
        this.projectService.getProjects('', this.pageIndex, this.pageSize, state.sort.by as string, state.sort.reverse).subscribe(
          paginatedProjects => {
            this.totalRecordCount = paginatedProjects.pagination.total_count;
            this.projects = paginatedProjects.project_list;
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

  editProjectMember(project: Project): void {
    if (this.isSystemAdminOrOwner(project)) {
      this.sharedActionService.createProjectMemberComponent(project, this.selfView);
    }
  }

  confirmToDeleteProject(project: Project): void {
    if (this.isSystemAdminOrOwner(project)) {
      this.translateService.get('PROJECT.CONFIRM_TO_DELETE_PROJECT', [project.project_name]).subscribe((msg: string) => {
        this.messageService.showDeleteDialog(msg, 'PROJECT.DELETE_PROJECT').subscribe((message: Message) => {
          if (message.returnStatus == RETURN_STATUS.rsConfirm) {
            this.projectService.deleteProject(project).subscribe(() => {
              this.messageService.showAlert('PROJECT.SUCCESSFUL_DELETE_PROJECT');
              this.retrieve(this.oldStateInfo);
            }, (error: HttpErrorResponse) => {
              if (error.status == 422) {
                this.messageService.showAlert('PROJECT.FAILED_TO_DELETE_PROJECT_ERROR', {alertType: "warning"})
              }
            })
          }
        });
      });
    }
  }

  toggleProjectPublic(project: Project, $event: MouseEvent): void {
    let oldPublic = project.project_public;
    this.projectService.togglePublicity(project.project_id, project.project_public === 1 ? 0 : 1).subscribe(() => {
        this.messageService.showAlert('PROJECT.SUCCESSFUL_TOGGLE_PROJECT');
        project.project_public = oldPublic === 1 ? 0 : 1;
      }, () => ($event.srcElement as HTMLInputElement).checked = oldPublic === 1
    );
  }

  get isFirstLogin(): boolean {
    return this.appInitService.isFirstLogin;
  }

  get guideStep(): GUIDE_STEP {
    return this.appInitService.guideStep;
  }

  isSystemAdminOrOwner(project: Project): boolean {
    return this.appInitService.currentUser.user_system_admin == 1 ||
      project.project_owner_id == this.appInitService.currentUser.user_id;
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
