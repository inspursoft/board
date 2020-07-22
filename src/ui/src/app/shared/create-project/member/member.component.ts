import { Component } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { Observable } from 'rxjs';
import { AppInitService } from '../../../shared.service/app-init.service';
import { SharedService } from '../../../shared.service/shared.service';
import { AlertType, SharedMember, SharedProject, SharedRole } from '../../shared.types';
import { CsModalChildBase } from '../../cs-modal-base/cs-modal-child-base';
import { MessageService } from '../../../shared.service/message.service';

@Component({
  styleUrls: ['./member.component.css'],
  templateUrl: './member.component.html'
})
export class MemberComponent extends CsModalChildBase {
  role: SharedRole;
  availableMembers: Array<SharedMember>;
  assignedMembers: Array<SharedMember>;
  selectedMember: SharedMember;
  project: SharedProject;
  doSet = false;
  doUnset = false;
  isActionWip = false;

  constructor(private sharedService: SharedService,
              private messageService: MessageService,
              private appInitService: AppInitService,
              private translateService: TranslateService) {
    super();
    this.project = new SharedProject();
    this.role = new SharedRole();
    this.selectedMember = new SharedMember();
    this.availableMembers = new Array<SharedMember>();
    this.assignedMembers = new Array<SharedMember>();
  }

  get isProjectOwner(): boolean {
    return this.project.projectOwnerId === this.appInitService.currentUser.user_id;
  }

  get isSelf(): boolean {
    return this.appInitService.currentUser.user_id === this.selectedMember.userId;
  }

  get isSystemAdmin(): boolean {
    return this.appInitService.currentUser.user_system_admin === 1;
  }

  get isOnesProject(): boolean {
    return this.project.projectOwnerId === this.selectedMember.userId;
  }

  retrieveMembers() {
    this.sharedService.getAssignedMembers(this.project.projectId)
      .subscribe((res: Array<SharedMember>) => this.assignedMembers = res);
    this.sharedService.getAvailableMembers(this.project.projectId)
      .subscribe((res: Array<SharedMember>) => this.availableMembers = res);
  }

  openMemberModal(project: SharedProject): Observable<string> {
    this.project = project;
    this.role.roleId = 1;
    this.retrieveMembers();
    return super.openModal();
  }

  pickUpAvailableMember(member: SharedMember) {
    this.selectedMember = member;
    this.doSet = false;
    this.doUnset = false;
    if ((this.isSelf && this.isProjectOwner) || (this.isSystemAdmin && this.isOnesProject)) {
      this.doSet = false;
    } else {
      if (this.isProjectOwner || this.isSystemAdmin) {
        this.doSet = true;
      }
    }
    this.role.roleId = this.selectedMember.roleId;
  }

  pickUpAssignedMember(member: SharedMember) {
    this.selectedMember = member;
    this.doUnset = false;
    this.doSet = false;
    if ((this.isSelf && this.isProjectOwner) || (this.isSystemAdmin && this.isOnesProject)) {
      this.doUnset = false;
    } else {
      if (this.isProjectOwner || this.isSystemAdmin) {
        this.doUnset = true;
      }
    }
    this.role.roleId = this.selectedMember.roleId;
  }

  pickUpRole(role: SharedRole) {
    this.selectedMember.roleId = role.roleId;
    this.isActionWip = true;
    this.sharedService.addOrUpdateProjectMember(this.project.projectId,
      this.selectedMember.userId,
      this.selectedMember.roleId).subscribe(
      () => this.displayInlineMessage('PROJECT.SUCCESSFUL_CHANGED_MEMBER_ROLE', 'info', [this.selectedMember.userName]),
      () => this.displayInlineMessage('PROJECT.FAILED_TO_CHANGE_MEMBER_ROLE', 'danger')
    );
  }

  setMember(): void {
    this.selectedMember.roleId = this.role.roleId;
    this.isActionWip = true;
    this.doSet = false;
    this.sharedService.addOrUpdateProjectMember(this.project.projectId,
      this.selectedMember.userId,
      this.selectedMember.roleId).subscribe(
      () => this.displayInlineMessage('PROJECT.SUCCESSFUL_ADDED_MEMBER', 'info', [this.selectedMember.userName]),
      () => this.displayInlineMessage('PROJECT.FAILED_TO_ADD_MEMBER', 'danger'),
      () => this.retrieveMembers()
    );
  }

  unsetMember(): void {
    this.selectedMember.id = 1;
    this.isActionWip = true;
    this.doUnset = false;
    this.sharedService.deleteProjectMember(this.project.projectId, this.selectedMember.userId).subscribe(
      () => this.displayInlineMessage('PROJECT.SUCCESSFUL_REMOVED_MEMBER', 'info', [this.selectedMember.userName]),
      () => this.displayInlineMessage('PROJECT.FAILED_TO_REMOVE_MEMBER', 'danger'),
      () => this.retrieveMembers()
    );
  }

  displayInlineMessage(message: string, alertType: AlertType, params?: object): void {
    this.translateService.get(message, params).subscribe((res: string) => {
      this.isActionWip = false;
      this.messageService.showAlert(res, {alertType, view: this.alertView});
    });
  }
}
