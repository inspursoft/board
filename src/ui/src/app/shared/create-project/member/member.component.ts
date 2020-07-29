import { Component } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { AppInitService } from '../../../shared.service/app-init.service';
import { Member, Project, Role } from '../../../project/project';
import { SharedService } from '../../../shared.service/shared.service';
import { AlertType } from '../../shared.types';
import { CsModalChildBase } from '../../cs-modal-base/cs-modal-child-base';
import { Observable } from 'rxjs';
import { MessageService } from '../../../shared.service/message.service';

@Component({
  selector: 'project-member',
  styleUrls: ['./member.component.css'],
  templateUrl: './member.component.html'
})
export class MemberComponent extends CsModalChildBase {
  role: Role = new Role();
  availableMembers: Array<Member>;
  assignedMembers: Array<Member>;
  selectedMember: Member = new Member();
  project: Project = new Project();
  doSet = false;
  doUnset = false;
  isActionWip = false;

  constructor(private sharedService: SharedService,
              private messageService: MessageService,
              private appInitService: AppInitService,
              private translateService: TranslateService) {
    super();
    this.availableMembers = Array<Member>();
    this.assignedMembers = Array<Member>();
  }

  get isProjectOwner(): boolean {
    return this.project.project_owner_id === this.appInitService.currentUser.user_id;
  }

  get isSelf(): boolean {
    return this.appInitService.currentUser.user_id === this.selectedMember.project_member_user_id;
  }

  get isSystemAdmin(): boolean {
    return this.appInitService.currentUser.user_system_admin === 1;
  }

  get isOnesProject(): boolean {
    return this.project.project_owner_id === this.selectedMember.project_member_user_id;
  }

  retrieveMembers() {
    this.sharedService.getAssignedMembers(this.project.project_id)
      .subscribe((res: Array<Member>) => this.assignedMembers = res);
    this.sharedService.getAvailableMembers(this.project.project_id)
      .subscribe((res: Array<Member>) => this.availableMembers = res);
  }

  openMemberModal(project: Project): Observable<string> {
    this.project = project;
    this.role.role_id = 1;
    this.retrieveMembers();
    return super.openModal();
  }

  pickUpAvailableMember(member: Member) {
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
    this.role.role_id = this.selectedMember.project_member_role_id;
  }

  pickUpAssignedMember(member: Member) {
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
    this.role.role_id = this.selectedMember.project_member_role_id;
  }

  pickUpRole(role: Role) {
    this.selectedMember.project_member_role_id = role.role_id;
    this.isActionWip = true;
    this.sharedService.addOrUpdateProjectMember(this.project.project_id,
      this.selectedMember.project_member_user_id,
      this.selectedMember.project_member_role_id).subscribe(
      () => this.displayInlineMessage('PROJECT.SUCCESSFUL_CHANGED_MEMBER_ROLE', 'info', [this.selectedMember.project_member_username]),
      () => this.displayInlineMessage('PROJECT.FAILED_TO_CHANGE_MEMBER_ROLE', 'danger')
    );
  }

  setMember(): void {
    this.selectedMember.project_member_role_id = this.role.role_id;
    this.isActionWip = true;
    this.doSet = false;
    this.sharedService.addOrUpdateProjectMember(this.project.project_id,
      this.selectedMember.project_member_user_id,
      this.selectedMember.project_member_role_id).subscribe(
      () => this.displayInlineMessage('PROJECT.SUCCESSFUL_ADDED_MEMBER', 'info', [this.selectedMember.project_member_username]),
      () => this.displayInlineMessage('PROJECT.FAILED_TO_ADD_MEMBER', 'danger'),
      () => this.retrieveMembers()
    );
  }

  unsetMember(): void {
    this.selectedMember.project_member_id = 1;
    this.isActionWip = true;
    this.doUnset = false;
    this.sharedService.deleteProjectMember(this.project.project_id, this.selectedMember.project_member_user_id).subscribe(
      () => this.displayInlineMessage('PROJECT.SUCCESSFUL_REMOVED_MEMBER', 'info', [this.selectedMember.project_member_username]),
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
