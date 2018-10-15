import { Component, Injector, OnInit, ViewContainerRef } from '@angular/core';
import { PHASE_SELECT_PROJECT, ServiceStepPhase, UIServiceStep1 } from '../service-step.component';
import { Project } from "../../project/project";
import { ServiceStepBase } from "../service-step";
import { SharedActionService } from "../../shared/shared-action.service";
import { SharedService } from "../../shared/shared.service";

@Component({
  styleUrls: ["./choose-project.component.css"],
  templateUrl: './choose-project.component.html'
})
export class ChooseProjectComponent extends ServiceStepBase implements OnInit {
  projectsList: Array<Project>;
  dropdownDefaultText: string = "";

  constructor(protected injector: Injector,
              private sharedService: SharedService,
              private sharedActionService: SharedActionService) {
    super(injector);
    this.projectsList = Array<Project>();
    this.dropdownDefaultText = "SERVICE.STEP_TITLE_1";
  }

  ngOnInit() {
    if (this.isBack) {
      this.k8sService.getServiceConfig(this.stepPhase).subscribe((res: UIServiceStep1) => {
        this.uiBaseData = res;
        this.uiData.projectId = res.projectId;
      })
    } else {
      this.k8sService.deleteServiceConfig().subscribe(res => res);
    }
    this.k8sService.getProjects()
      .subscribe((res: Array<Project>) => {
        let createNewProject: Project = new Project();
        createNewProject.project_name = "SERVICE.STEP_1_CREATE_PROJECT";
        createNewProject.project_id = -1;
        createNewProject["isSpecial"] = true;
        createNewProject["OnlyClick"] = true;
        this.projectsList.push(createNewProject);
        if (res && res.length > 0) {
          this.projectsList = this.projectsList.concat(res);
        }
        this.setDropdownDefaultText();
      })
  }

  setDropdownDefaultText(): void {
    let selected = this.projectsList.find((project: Project) => project.project_id === this.uiData.projectId);
    this.dropdownDefaultText = selected ? selected.project_name : "SERVICE.STEP_TITLE_1";
  }

  get stepPhase(): ServiceStepPhase {
    return PHASE_SELECT_PROJECT;
  }

  get uiData(): UIServiceStep1 {
    return this.uiBaseData as UIServiceStep1;
  }

  forward() {
    this.k8sService.setServiceConfig(this.uiData.uiToServer()).subscribe(
      () => this.k8sService.stepSource.next({index: 2, isBack: false})
    );
  }

  clickSelectProject() {
    this.sharedActionService.createProjectComponent(this.selfView).subscribe((projectName: string) => {
      if (projectName) {
        this.sharedService.getOneProject(projectName).subscribe((res: Array<Project>) => {
          this.uiData.projectId = res[0].project_id;
          let project = this.projectsList.shift();
          this.projectsList.unshift(res[0]);
          this.projectsList.unshift(project);
          this.setDropdownDefaultText();
        })
      }
    });
  }

  changeSelectProject(project: Project) {
    this.uiData.projectId = project.project_id;
    this.setDropdownDefaultText();
  }
}