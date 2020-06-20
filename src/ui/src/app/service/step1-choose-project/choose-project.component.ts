import { Component, Injector, OnInit } from '@angular/core';
import { PHASE_SELECT_PROJECT, ServiceStep1Data, ServiceStepPhase } from '../service-step.component';
import { ServiceStepComponentBase } from '../service-step';
import { SharedActionService } from '../../shared.service/shared-action.service';
import { SharedService } from '../../shared.service/shared.service';
import { ServiceProject } from '../service.types';

@Component({
  styleUrls: ['./choose-project.component.css'],
  templateUrl: './choose-project.component.html'
})
export class ChooseProjectComponent extends ServiceStepComponentBase implements OnInit {
  projectsList: Array<ServiceProject>;
  curActiveProject: ServiceProject;
  stepData: ServiceStep1Data;

  constructor(protected injector: Injector,
              private sharedService: SharedService,
              private sharedActionService: SharedActionService) {
    super(injector);
    this.projectsList = Array<ServiceProject>();
    this.stepData = new ServiceStep1Data();
  }

  ngOnInit() {
    if (this.isBack) {
      this.k8sService.getServiceConfig(this.stepPhase, ServiceStep1Data).subscribe((res: ServiceStep1Data) => {
        this.stepData = res;
        this.k8sService.getProjects().subscribe((projects: Array<ServiceProject>) => {
          this.projectsList = projects;
          this.curActiveProject = this.projectsList.find(value => value.projectId === this.stepData.projectId);
        });
      });
    } else {
      this.k8sService.deleteServiceConfig().subscribe(res => res);
      this.k8sService.getProjects().subscribe((res: Array<ServiceProject>) => this.projectsList = res);
    }
  }

  get stepPhase(): ServiceStepPhase {
    return PHASE_SELECT_PROJECT;
  }

  forward() {
    this.k8sService.setServiceStepConfig(this.stepData).subscribe(
      () => this.k8sService.stepSource.next({index: 2, isBack: false})
    );
  }

  clickSelectProject() {
    this.sharedActionService.createProjectComponent(this.selfView).subscribe((projectName: string) => {
      if (projectName) {
        this.k8sService.getOneProject(projectName).subscribe((res: Array<ServiceProject>) => {
          this.stepData.projectId = res[0].projectId;
          this.stepData.projectName = res[0].projectName;
          this.projectsList.push(res[0]);
        });
      }
    });
  }

  changeSelectProject(project: ServiceProject) {
    this.stepData.projectId = project.projectId;
    this.stepData.projectName = project.projectName;
  }
}
