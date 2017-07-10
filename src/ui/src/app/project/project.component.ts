import { Component, OnInit, ViewChild } from '@angular/core';
import { Project } from './project';
import { ProjectService } from './project.service';
import { CreateProjectComponent } from './create-project/create-project.component';
import { MemberComponent } from './member/member.component';

@Component({
  selector: 'project',
  templateUrl: 'project.component.html'
})
export class ProjectComponent implements OnInit {
  
  projects: Project[];

  @ViewChild(CreateProjectComponent) createProjectModal;
  @ViewChild(MemberComponent) memberModal;

  constructor(private projectService: ProjectService){}

  ngOnInit(): void {
    this.projectService
      .getProjects()
      .then(projects=>this.projects = projects);
  }

  createProject(): void {
    this.createProjectModal.openModal();
  }

  editProjectMember(p: Project): void {
    this.memberModal.openModal(p);
  }
}