import { Component, Input, OnInit } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';

import { AppInitService } from '../../app.init.service';

import { Project } from '../project';
import { ProjectService } from '../project.service';
import { Member } from './member';
import { Role } from './role';
import { Subject } from 'rxjs/Subject';

import { ROLES } from '../../shared/shared.const';
import { MessageService } from '../../shared/message-service/message.service';

@Component({
  selector: 'project-member',
  templateUrl: 'member.component.html'
})
export class MemberComponent implements OnInit {

  currentUser: {[key: string]: any};

  projectMemberOpened: boolean;

  role: Role = new Role();
  members: Member[];
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

  memberSubject: Subject<Member[]> = new Subject<Member[]>();

  constructor(
    private appInitService: AppInitService,
    private projectService: ProjectService,
    private translateService: TranslateService,
    private messageService: MessageService
  ){}
  
  ngOnInit(): void {
    this.currentUser = this.appInitService.currentUser;
  }

  retrieveAvailableMembers() {
    this.projectService
      .getProjectMembers(this.project.project_id)
      .then(members=>{
        this.projectService
          .getAvailableMembers()
          .then(availableMembers=>{
            this.availableMembers = availableMembers;
            this.availableMembers.forEach(am=>{
              members.forEach(m=>{
                if (m.project_member_user_id === am.project_member_user_id) {
                 am.project_member_id = m.project_member_id;
                 am.project_member_role_id = m.project_member_role_id;
                 am.isMember = true;
                }
              });
            });
            this.memberSubject.subscribe(changedMembers=>{
              this.availableMembers = changedMembers;
            });
          })
          .catch(err=>console.error('Failed to load available members.'));
      }).catch(err=>console.error('Failed to load current members.'));
  }

  openModal(p: Project): void {
    this.projectMemberOpened = true;
    this.project = p;
    this.role.role_id = 1;
    this.retrieveAvailableMembers();
  }

  confirm(): void {
    this.projectMemberOpened = false;
  }

  pickUpMember(m: Member) {
    this.selectedMember = m;
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

  pickUpRole(r: Role) {
    this.selectedMember.project_member_role_id = r.role_id;
    this.projectService
      .addOrUpdateProjectMember(this.project.project_id, 
        this.selectedMember.project_member_user_id, 
        this.selectedMember.project_member_role_id)
      .then(()=>{
        this.alertType = 'alert-info';
        this.displayInlineMessage('PROJECT.SUCCESSFUL_CHANGED_MEMBER_ROLE', [this.selectedMember.project_member_username]);
      })
      .catch(err=>{
        this.alertType = 'alert-danger';
        this.displayInlineMessage('PROJECT.FAILED_TO_CHANGE_MEMBER_ROLE');
      });
  }

  setMember(): void {
    this.availableMembers.forEach(m=>{
      if(m.project_member_user_id === this.selectedMember.project_member_user_id) {
        m.project_member_role_id = this.role.role_id;
        this.projectService
          .addOrUpdateProjectMember(this.project.project_id, 
            this.selectedMember.project_member_user_id, 
            this.selectedMember.project_member_role_id)
          .then(()=>{
            this.alertType = 'alert-info';
            this.displayInlineMessage('PROJECT.SUCCESSFUL_ADDED_MEMBER',[this.selectedMember.project_member_username])
          })
          .catch(err=>{
            this.alertType = 'alert-danger';
            this.displayInlineMessage('PROJECT.FAILED_TO_ADD_MEMBER');
          });
        m.isMember = true;
      }
    });
    this.memberSubject.next(this.availableMembers);
  }

  unsetMember(): void {
    this.availableMembers.forEach(m=>{
      if(m.project_member_user_id === this.selectedMember.project_member_user_id) {
        this.selectedMember.project_member_id = 1;
        m.isMember = false;
        this.projectService
          .deleteProjectMember(this.project.project_id, this.selectedMember.project_member_user_id)
          .then(()=>{
            this.alertType = 'alert-info';
            this.displayInlineMessage('PROJECT.SUCCESSFUL_REMOVED_MEMBER', [this.selectedMember.project_member_username]);
          })
          .catch(err=>{
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
      .subscribe(res=>{
        this.changedMessage = res;
      });
    setTimeout(()=>this.hasChanged = false, 2*1000);
  }
}