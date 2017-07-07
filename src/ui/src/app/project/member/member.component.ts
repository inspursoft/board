import { Component, OnInit } from '@angular/core';
import { Member } from './member';
import { Subject } from 'rxjs/Subject';

const MEMBERS = [
  {
    userId: 1,
    username: 'Zhao',
    isMember: false
  },
  {
    userId: 3,
    username: 'Qian',
    isMember: false
  },
  {
    userId: 5,
    username: 'Sun',
    isMember: false
  }
]

@Component({
  selector: 'project-member',
  templateUrl: 'member.component.html'
})
export class MemberComponent implements OnInit {
  projectMemberOpened: boolean;
  members: Member[];
  
  selectedMember: Member;
  
  doSet: boolean;
  doUnset: boolean;

  memberSubject: Subject<Member[]> = new Subject<Member[]>();

  ngOnInit(): void {
    this.members = MEMBERS;
    
    if(this.members && this.members.length > 0) {
      this.selectedMember = this.members[0];
      this.doSet = true;
      this.doUnset = false;
    }
    this.memberSubject.subscribe(changedMembers=>{
      this.members = changedMembers;
    });
  }

  openModal(): void {
    this.projectMemberOpened = true;
  }

  confirm(): void {
    this.projectMemberOpened = false;
  }

  pickUpMember(m: Member) {
    this.selectedMember = m;
  }

  setMember(): void {
    this.members.forEach(m=>{
      if(m.userId === this.selectedMember.userId) {
        m.isMember = true;
      }
    });
    this.memberSubject.next(this.members);
  }

  unsetMember(): void {
    this.members.forEach(m=>{
      if(m.userId === this.selectedMember.userId) {
        m.isMember = false;
      }
    });
    this.memberSubject.next(this.members);
  }
}