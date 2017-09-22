/**
 * Created by liyanq on 9/17/17.
 */

import { Component, OnInit } from "@angular/core"
import { K8sService } from "../service.k8s";
import { ServiceStep4Output } from "../service-step.component";
import { MessageService } from "../../shared/message-service/message.service";

@Component({
  templateUrl: "./deploy.component.html",
  styleUrls: ["./deploy.component.css"]
})
export class DeployComponent implements OnInit {
  isInDeployed: boolean = false;
  isInDeployIng: boolean = false;
  consoleText: string = "";
  output4: ServiceStep4Output;

  constructor(private k8sService: K8sService,
              private messageService: MessageService) {
  }

  ngOnInit() {
    this.output4 = this.k8sService.getStepData(4) as ServiceStep4Output;
    this.output4.config_phase = "deploy";
  }

  serviceDeploy() {
    if (!this.isInDeployed) {
      this.isInDeployIng = true;
      this.k8sService.serviceDeployment(this.output4)
        .then(res => {
          this.consoleText = JSON.stringify(res);
          this.isInDeployed = true;
          this.isInDeployIng = false;
        })
        .catch(err => {
          this.messageService.dispatchError(err);
          this.isInDeployed = true;
          this.isInDeployIng = false;
        })
    }
  }

  deployComplete(): void {
    this.k8sService.stepSource.next(0);
  }
}