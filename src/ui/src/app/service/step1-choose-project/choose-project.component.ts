import { Component, Injector, OnInit } from '@angular/core';
import { PHASE_SELECT_PROJECT, ServiceStepPhase, UIServiceStep1 } from '../service-step.component';
import { Project } from "../../project/project";
import { ServiceStepBase } from "../service-step";
import { SharedActionService } from "../../shared.service/shared-action.service";
import { SharedService } from "../../shared.service/shared.service";

@Component({
  styleUrls: ["./choose-project.component.css"],
  templateUrl: './choose-project.component.html'
})
export class ChooseProjectComponent extends ServiceStepBase implements OnInit {
  projectsList: Array<Project>;
  curActiveProject: Project;

  constructor(protected injector: Injector,
              private sharedService: SharedService,
              private sharedActionService: SharedActionService) {
    super(injector);
    this.projectsList = Array<Project>();
  }

  ngOnInit() {
    if (this.isBack) {
      this.k8sService.getServiceConfig(this.stepPhase).subscribe((res: UIServiceStep1) => {
        this.uiBaseData = res;
        this.k8sService.getProjects().subscribe((res: Array<Project>) => {
          this.projectsList = res;
          this.curActiveProject = this.projectsList.find(value => value.project_id === this.uiData.projectId);
        })
      })
    } else {
      this.k8sService.deleteServiceConfig().subscribe(res => res);
      this.k8sService.getProjects().subscribe((res: Array<Project>) => this.projectsList = res);
    }
  }

  ngAfterViewInit(): void {

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
          this.uiData.projectName = res[0].project_name;
          this.projectsList.push(res[0]);
        })
      }
    });
  }

  changeSelectProject(project: Project) {
    this.uiData.projectId = project.project_id;
    this.uiData.projectName = project.project_name;
  }
}
