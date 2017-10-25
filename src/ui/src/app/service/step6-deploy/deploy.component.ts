/**
 * Created by liyanq on 9/17/17.
 */

import { Component, OnDestroy, OnInit } from "@angular/core"
import { K8sService } from "../service.k8s";
import { ServiceStep4Output } from "../service-step.component";
import { MessageService } from "../../shared/message-service/message.service";

const AUTO_REFRESH_DEPLOY_STATUS: number = 2000;
@Component({
  templateUrl: "./deploy.component.html",
  styleUrls: ["./deploy.component.css"]
})
export class DeployComponent implements OnInit, OnDestroy {
  isInDeployed: boolean = false;
  isInDeployIng: boolean = false;
  consoleText: string = "";
  output4: ServiceStep4Output;
  intervalAutoRefreshDeployStatus: any;
  autoRefreshTimesCount: number = 0;
  isNeedAutoRefreshDeployStatus: boolean;

  constructor(private k8sService: K8sService,
              private messageService: MessageService) {
  }

  ngOnInit() {
    this.output4 = this.k8sService.getStepData(4) as ServiceStep4Output;
    this.output4.config_phase = "deploy";
    this.intervalAutoRefreshDeployStatus = setInterval(() => {
      if (this.isNeedAutoRefreshDeployStatus) {
        this.autoRefreshTimesCount++;
        this.k8sService.getDeployStatus(this.output4.service_yaml.service_name)
          .then(res => {
            this.consoleText = JSON.stringify(res);
            if (!res["code"]){
              this.isNeedAutoRefreshDeployStatus = false;
              this.isInDeployed = true;
              this.isInDeployIng = false;
            }
          }).catch(err => {
            this.isNeedAutoRefreshDeployStatus = false;
            this.messageService.dispatchError(err);
        });
      }
    }, AUTO_REFRESH_DEPLOY_STATUS);
  }

  ngOnDestroy() {
    clearInterval(this.intervalAutoRefreshDeployStatus);
  }

  serviceDeploy() {
    if (!this.isInDeployed) {
      this.autoRefreshTimesCount = 0;
      this.isInDeployIng = true;
      this.isNeedAutoRefreshDeployStatus = true;
      this.k8sService.serviceDeployment(this.output4)
        .then(res => res)
        .catch(err => {
          this.messageService.dispatchError(err);
          this.isNeedAutoRefreshDeployStatus = false;
          this.isInDeployed = true;
          this.isInDeployIng = false;
        })
    }
  }

  deployComplete(): void {
    this.k8sService.stepSource.next(0);
  }

  backStep(): void {
    this.k8sService.stepSource.next(4);
  }
}