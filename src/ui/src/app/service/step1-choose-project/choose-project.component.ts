import { Component, Injector, OnInit } from '@angular/core';
import { DeploymentServiceData } from '../service-step.component';
import { Project } from "../../project/project";
import { Router } from "@angular/router";
import { ServiceStepBase } from "../service-step";

@Component({
  styleUrls: ["./choose-project.component.css"],
  templateUrl: './choose-project.component.html'
})
export class ChooseProjectComponent extends ServiceStepBase implements OnInit {
  projectsList: Array<Project>;

  constructor(protected injector: Injector) {
    super(injector);
    this.projectsList = Array<Project>();
  }

  ngOnInit() {
    this.outputData = new DeploymentServiceData();
    this.k8sService.getProjects()
      .then(res => {
        let createNewProject: Project = new Project();
        createNewProject.project_name = "SERVICE.STEP_1_CREATE_PROJECT";
        createNewProject["isSpecial"] = true;
        createNewProject["OnlyClick"] = true;
        this.projectsList.push(createNewProject);
        if (res && res.length > 0) {
          this.projectsList = this.projectsList.concat(res);
        }
      })
      .catch(err => this.messageService.dispatchError(err));
  }

  forward() {
    this.k8sService.getServiceID({
      project_name: this.outputData.projectinfo.project_name,
      project_id: this.outputData.projectinfo.project_id
    }).then(res => {
      let serviceId = Number(res).valueOf();
      this.newServiceId = serviceId;
      this.outputData.projectinfo.service_id = serviceId;
      this.k8sService.setServiceConfig(this.outputData).then((isCompleted) => {
        this.k8sService.stepSource.next({index: 2, isBack: false});
      });
    }).catch(err => this.messageService.dispatchError(err));
  }

  clickSelectProject(project: Project) {
    this.router.navigate(["/projects"]);
  }

  changeSelectProject(project: Project) {
    this.outputData.projectinfo.project_name = project.project_name;
    this.outputData.projectinfo.project_id = project.project_id;
  }
}