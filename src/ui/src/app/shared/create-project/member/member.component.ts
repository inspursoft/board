import { Component, OnInit } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { AppInitService } from "../../../app.init.service";
import { Member, Project, Role } from "../../../project/project";
import { Subject } from 'rxjs/Subject';
import { SharedService } from "../../shared.service";
import { MessageService } from "../../message-service/message.service";
import { HttpErrorResponse } from "@angular/common/http";
import { Observable } from "rxjs/Observable";

@Component({
  selector: 'project-member',
  templateUrl: './member.component.html'
})
export class MemberComponent implements OnInit {
  currentUser: {[key: string]: any};
  _projectMemberOpened: boolean = false;
  role: Role = new Role();
  availableMembers: Member[];
  selectedMember: Member = new Member();
  project: Project = new Project();
  isLeftPane: boolean;
  isRightPane: boolean;
  doSet: boolean;
  doUnset: boolean;
  alertType: string;
  hasChanged: boolean;
  changedMessage: string;
  closeNotification: Subject<any>;
  memberSubject: Subject<Member[]> = new Subject<Member[]>();

  constructor(private sharedService: SharedService,
              private messageService: MessageService,
              private appInitService: AppInitService,
              private translateService: TranslateService) {
    this.closeNotification = new Subject<any>();
  }
  
  ngOnInit(): void {
    this.currentUser = this.appInitService.currentUser;
  }

  retrieveAvailableMembers() {
    this.sharedService.getProjectMembers(this.project.project_id)
      .then((members: Array<Member>) => {
        this.sharedService.getAvailableMembers()
          .then((availableMembers: Array<Member>) => {
            this.availableMembers = availableMembers;
            this.availableMembers.forEach((am: Member) => {
              members.forEach((member: Member) => {
                if (member.project_member_user_id === am.project_member_user_id) {
                  am.project_member_id = member.project_member_id;
                  am.project_member_role_id = member.project_member_role_id;
                 am.isMember = true;
                }
              });
            });
            this.memberSubject.subscribe((changedMembers: Array<Member>) => {
              this.availableMembers = changedMembers;
            });
          })
          .catch((err: HttpErrorResponse) => {
            this.projectMemberOpened = false;
            this.messageService.dispatchError(err);
          });
      })
      .catch((err: HttpErrorResponse) => {
        this.projectMemberOpened = false;
        this.messageService.dispatchError(err);
      });
  }

  openModal(project: Project): Observable<boolean> {
    this.projectMemberOpened = true;
    this.project = project;
    this.role.role_id = 1;
    this.retrieveAvailableMembers();
    return this.closeNotification.asObservable();
  }

  get projectMemberOpened(): boolean {
    return this._projectMemberOpened;
  }

  set projectMemberOpened(value: boolean) {
    this._projectMemberOpened = value;
    if (!value) {
      this.closeNotification.next();
    }
  }

  confirm(): void {
    this.projectMemberOpened = false;
  }

  pickUpMember(member: Member) {
    this.selectedMember = member;
    this.doSet = false;
    this.doUnset = false;
    let isProjectOwner = (this.project.project_owner_id === this.currentUser.user_id);
    let isSelf = (this.currentUser.user_id === this.selectedMember.project_member_user_id);
    let isSystemAdmin = (this.currentUser.user_system_admin === 1);
    let isOnesProject = (this.project.project_owner_id === this.selectedMember.project_member_user_id);
    if((isSelf && isProjectOwner) || (isSystemAdmin && isOnesProject)) {
      this.doSet = false;
      this.doUnset = false;
    } else { 
      if(isProjectOwner || isSystemAdmin) {
        this.doSet = this.isLeftPane;
        this.doUnset = this.isRightPane;
      }
    }
    this.role.role_id = this.selectedMember.project_member_role_id;
  }

  pickUpRole(role: Role) {
    this.selectedMember.project_member_role_id = role.role_id;
    this.sharedService.addOrUpdateProjectMember(this.project.project_id,
        this.selectedMember.project_member_user_id, 
        this.selectedMember.project_member_role_id)
      .then(()=>{
        this.alertType = 'alert-info';
        this.displayInlineMessage('PROJECT.SUCCESSFUL_CHANGED_MEMBER_ROLE', [this.selectedMember.project_member_username]);
      })
      .catch(() => {
        this.alertType = 'alert-danger';
        this.displayInlineMessage('PROJECT.FAILED_TO_CHANGE_MEMBER_ROLE');
      });
  }

  setMember(): void {
    this.availableMembers.forEach((member: Member) => {
      if (member.project_member_user_id === this.selectedMember.project_member_user_id) {
        member.project_member_role_id = this.role.role_id;
        this.sharedService.addOrUpdateProjectMember(this.project.project_id,
            this.selectedMember.project_member_user_id, 
            this.selectedMember.project_member_role_id)
          .then(()=>{
            this.alertType = 'alert-info';
            this.displayInlineMessage('PROJECT.SUCCESSFUL_ADDED_MEMBER',[this.selectedMember.project_member_username])
          })
          .catch(() => {
            this.alertType = 'alert-danger';
            this.displayInlineMessage('PROJECT.FAILED_TO_ADD_MEMBER');
          });
        member.isMember = true;
      }
    });
    this.memberSubject.next(this.availableMembers);
  }

  unsetMember(): void {
    this.availableMembers.forEach((member: Member) => {
      if (member.project_member_user_id === this.selectedMember.project_member_user_id) {
        this.selectedMember.project_member_id = 1;
        member.isMember = false;
        this.sharedService
          .deleteProjectMember(this.project.project_id, this.selectedMember.project_member_user_id)
          .then(()=>{
            this.alertType = 'alert-info';
            this.displayInlineMessage('PROJECT.SUCCESSFUL_REMOVED_MEMBER', [this.selectedMember.project_member_username]);
          })
          .catch(() => {
            this.alertType = 'alert-danger';
            this.displayInlineMessage('PROJECT.FAILED_TO_REMOVE_MEMBER');
          });
      }
    });
    this.memberSubject.next(this.availableMembers);
  }

  displayInlineMessage(message: string, params?: object): void {
    this.hasChanged = true;
    this.translateService.get(message, params || [])
      .subscribe(res => this.changedMessage = res);
    setTimeout(()=>this.hasChanged = false, 2*1000);
  }
}