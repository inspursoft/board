import { Component } from '@angular/core';
import { Observable } from 'rxjs';
import { HttpErrorResponse } from '@angular/common/http';
import { SharedService } from '../../../shared.service/shared.service';
import { CsModalChildBase } from '../../cs-modal-base/cs-modal-child-base';
import { MessageService } from '../../../shared.service/message.service';
import { SharedCreateProject } from '../../shared.types';

@Component({
  styleUrls: [ './create-project.component.css' ],
  templateUrl: './create-project.component.html'
})
export class CreateProjectComponent extends CsModalChildBase {
  createProject: SharedCreateProject;
  isCreateProjectWIP = false;
  projectNamePattern = /^[a-z0-9]+(?:[-][a-z0-9]+)*$/;
  constructor(private sharedService: SharedService,
              private messageService: MessageService) {
    super();
    this.createProject = new SharedCreateProject();
  }

  openCreateProjectModal(): Observable<string> {
    this.modalOpened = true;
    return this.closeNotification.asObservable();
  }

  confirm(): void {
    if (this.verifyInputExValid()) {
      this.isCreateProjectWIP = true;
      const project = new SharedCreateProject();
      project.projectName = this.createProject.projectName;
      project.publicity = this.createProject.publicity;
      project.comment = this.createProject.comment;
      this.sharedService.createProject(project).subscribe(
        () => this.messageService.showAlert('PROJECT.SUCCESSFUL_CREATED_PROJECT'),
        (err: HttpErrorResponse) => {
          this.isCreateProjectWIP = false;
          if (err.status === 409) {
            this.messageService.showAlert('PROJECT.PROJECT_NAME_ALREADY_EXISTS', {alertType: 'danger'});
          } else if (err.status === 400) {
            this.messageService.showAlert('PROJECT.PROJECT_NAME_IS_ILLEGAL', {alertType: 'danger'});
          }
          this.modalOpened = false;
        },
        () => {
          this.closeNotification.next(project.projectName);
          this.modalOpened = false;
        }
      );
    }
  }
}
