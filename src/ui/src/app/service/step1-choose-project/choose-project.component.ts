import { Component, Injector, OnInit } from '@angular/core';
import { PHASE_SELECT_PROJECT, ServiceStepPhase, UIServiceStep1 } from '../service-step.component';
import { Project } from "../../project/project";
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

  get stepPhase(): ServiceStepPhase {
    return PHASE_SELECT_PROJECT;
  }

  get uiData():UIServiceStep1{
    return this.uiBaseData as UIServiceStep1;
  }

  forward() {
    this.k8sService.setServiceConfig(this.uiData.uiToServer()).then((isCompleted) => {
      this.k8sService.stepSource.next({index: 2, isBack: false});
    }).catch(err => this.messageService.dispatchError(err));
  }

  clickSelectProject(project: Project) {
    this.router.navigate(["/projects"]);
  }

  changeSelectProject(project: Project) {
    this.uiData.projectId = project.project_id;
  }
}