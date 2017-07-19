import { Component, Input, OnInit } from '@angular/core';
import { Project } from '../project';
import { ProjectService } from '../project.service';
import { Member } from './member';
import { Role } from './role';
import { Subject } from 'rxjs/Subject';

import { ROLES } from '../../shared/shared.const';

@Component({
  selector: 'project-member',
  templateUrl: 'member.component.html'
})
export class MemberComponent implements OnInit {
  projectMemberOpened: boolean;

  role: Role = new Role();
  members: Member[];
  availableMembers: Member[];

  selectedMember: Member;
  
  project: Project = new Project();

  doSet: boolean;
  doUnset: boolean;

  hasChanged: boolean;
  changedMessage: string;

  memberSubject: Subject<Member[]> = new Subject<Member[]>();

  constructor(private projectService: ProjectService){}

  ngOnInit(): void {}

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
            if(this.availableMembers && this.availableMembers.length > 0) {
              this.selectedMember = this.availableMembers[0];
              this.doSet = true;
              this.doUnset = false;
            }
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
    if(this.doUnset){
      this.role.role_id = this.selectedMember.project_member_role_id;    
    }
  }

  pickUpRole(r: Role) {
    this.selectedMember.project_member_role_id = r.role_id;
    this.projectService
      .addOrUpdateProjectMember(this.project.project_id, 
        this.selectedMember.project_member_user_id, 
        this.selectedMember.project_member_role_id)
      .then(()=>this.displayInlineMessage('Successful changed member ' + this.selectedMember.project_member_username + ' with role ' + ROLES[this.selectedMember.project_member_role_id]))
      .catch(err=>console.error('Failed to delete member user_id:' + this.selectedMember.project_member_user_id));
  }

  setMember(): void {
    this.availableMembers.forEach(m=>{
      if(m.project_member_user_id === this.selectedMember.project_member_user_id) {
        m.project_member_role_id = this.role.role_id;
        this.projectService
          .addOrUpdateProjectMember(this.project.project_id, 
            this.selectedMember.project_member_user_id, 
            this.selectedMember.project_member_role_id)
          .then(()=>this.displayInlineMessage('Successful added member ' + this.selectedMember.project_member_username + ' with role ' + ROLES[this.selectedMember.project_member_role_id]));
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
            this.displayInlineMessage('Successful deleted member ' + this.selectedMember.project_member_username);
            this.doSet = true;
          })
          .catch(err=>console.error('Failed to delete member user_id:' + this.selectedMember.project_member_user_id));
      }
    });
    this.memberSubject.next(this.availableMembers);
  }

  displayInlineMessage(message: string): void {
    this.hasChanged = true;
    this.changedMessage = message;
    setTimeout(()=>this.hasChanged = false, 2*1000);
  }
}