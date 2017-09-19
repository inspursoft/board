import { Component, Input, OnDestroy, OnInit } from '@angular/core';
import { ServiceStep1Output, ServiceStepComponent } from '../service-step.component';
import { K8sService } from '../service.k8s';
import { Project } from "../../project/project";
import { MessageService } from "../../shared/message-service/message.service";
import { Router } from "@angular/router";

@Component({
  styleUrls: ["./choose-project.component.css"],
  templateUrl: './choose-project.component.html'
})
export class ChooseProjectComponent implements ServiceStepComponent, OnInit, OnDestroy {
  @Input() data: any;
  projectsList: Array<Project>;
  outputData: ServiceStep1Output = new ServiceStep1Output();

  constructor(private k8sService: K8sService,
              private router: Router,
              private messageService: MessageService) {
    this.projectsList = Array<Project>();
  }

  ngOnInit() {
    this.k8sService.clearStepData();
    this.k8sService.getProjects()
      .then(res => {
        let createNewProject: Project = new Project();
        createNewProject.project_name = "SERVICE.STEP_1_CREATE_PROJECT";
        createNewProject["isSpecial"] = true;
        createNewProject["OnlyClick"] = true;
        this.projectsList.push(createNewProject);
        if (res && res.length > 0){
          this.projectsList = this.projectsList.concat(res);
        }
      })
      .catch(err => this.messageService.dispatchError(err));
  }

  ngOnDestroy() {
    this.k8sService.setStepData(1, this.outputData);
  }

  forward() {
    this.k8sService.stepSource.next(2);
  }

  clickSelectProject(project: Project) {
    this.router.navigate(["/projects"]);
  }

  changeSelectProject(project: Project) {
    this.outputData.project_name = project.project_name;
    this.outputData.project_id = project.project_id;
  }
}