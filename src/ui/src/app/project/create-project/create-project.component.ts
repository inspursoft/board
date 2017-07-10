import { Component } from '@angular/core';
import { Project } from '../project';


@Component({
  selector: 'create-project',
  styleUrls: [ './create-project.component.css' ],
  templateUrl: './create-project.component.html'
})
export class CreateProjectComponent {
  createProjectOpened: boolean;
  project: Project = new Project();

  openModal(): void {
    this.createProjectOpened = true;
  }

  confirm(): void {
    this.createProjectOpened = false;
  }
}