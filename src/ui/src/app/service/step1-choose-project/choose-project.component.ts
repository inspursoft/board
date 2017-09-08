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
  dropdownText: string = '';
  projectsList: Array<Project>;
  getProjectIng: boolean = false;
  outputData: ServiceStep1Output = new ServiceStep1Output();

  constructor(private k8sService: K8sService,
              private router: Router,
              private messageService: MessageService) {
  }

  ngOnInit() {
    this.getProjectIng = true;
    this.k8sService.clearStepData();
    this.k8sService.getProjects()
      .then(res => {
        this.getProjectIng = false;
        if (res.length > 0) {
          this.projectsList = res;
          this.dropdownText = res[0].project_name;
          this.outputData.project_name = res[0].project_name;
          this.outputData.project_id = res[0].project_id;
        } else {
          this.dropdownText = "SERVICE.STEP_1_NONE_PROJECT";
        }
      })
      .catch(err => {
        this.getProjectIng = false;
        this.messageService.dispatchError(err);
      });
  }

  ngOnDestroy() {
    this.k8sService.setStepData(1, this.outputData);
  }

  forward() {
    this.k8sService.stepSource.next(2);
  }

  selectProject(project: Project) {
    this.outputData.project_name = project.project_name;
    this.outputData.project_id = project.project_id;
    this.dropdownText = project.project_name;
  }

  redirectToCreateProject() {
    this.router.navigate(["/projects"]);

  }
}