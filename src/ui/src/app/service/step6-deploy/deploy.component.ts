/**
 * Created by liyanq on 9/17/17.
 */

import { Component, OnDestroy, OnInit } from "@angular/core"
import { K8sService } from "../service.k8s";
import { ServiceStep4Output } from "../service-step.component";
import { MessageService } from "../../shared/message-service/message.service";
import { Subscription } from "rxjs/Subscription";
import { Message } from "../../shared/message-service/message";
import { BUTTON_STYLE } from "../../shared/shared.const";
import { WebsocketService } from "../../shared/websocket-service/websocket.service";
import { AppInitService } from "../../app.init.service";

const PROCESS_SERVICE_CONSOLE_URL = `ws://10.165.22.61:8088/api/v1/jenkins-job/console?job_name=process_service`;
// const PROCESS_SERVICE_CONSOLE_URL = `ws://localhost/api/v1/jenkins-job/console?job_name=process_service`;
@Component({
  templateUrl: "./deploy.component.html",
  styleUrls: ["./deploy.component.css"]
})
export class DeployComponent implements OnInit, OnDestroy {
  isDeployed: boolean = false;
  isInDeploySuccess: boolean = false;
  isInDeployIng: boolean = false;
  consoleText: string = "";
  output4: ServiceStep4Output;
  processImageSubscription: Subscription;
  _confirmSubscription: Subscription;

  constructor(private k8sService: K8sService,
              private appInitService: AppInitService,
              private webSocketService: WebsocketService,
              private messageService: MessageService) {
  }

  ngOnInit() {
    this.output4 = this.k8sService.getStepData(4) as ServiceStep4Output;
    this.output4.projectinfo.config_phase = "deploy";
    this._confirmSubscription = this.messageService.messageConfirmed$.subscribe((next: Message) => {
      this.k8sService.deleteDeployment(this.output4.projectinfo.service_id)
        .then(isDelete => {
          if (this.processImageSubscription){
            this.processImageSubscription.unsubscribe();
          }
          this.k8sService.stepSource.next(0);
        })
        .catch(err => {
          this.messageService.dispatchError(err);
          if (this.processImageSubscription){
            this.processImageSubscription.unsubscribe();
          }
          this.k8sService.stepSource.next(0);
        })
    });
  }

  ngOnDestroy() {
    this._confirmSubscription.unsubscribe();
  }

  serviceDeploy() {
    if (!this.isDeployed) {
      this.isDeployed = true;
      this.isInDeployIng = true;
      this.consoleText = "Deploying...";
      this.k8sService.serviceDeployment(this.output4)
        .then(res => {
          setTimeout(() => {
            this.processImageSubscription = this.webSocketService
              .connect(PROCESS_SERVICE_CONSOLE_URL + `&token=${this.appInitService.token}`)
              .subscribe(obs => {
                this.consoleText = <string>obs.data;
                let consoleTextArr: Array<string> = this.consoleText.split(/[\n]/g);
                if (consoleTextArr.find(value => value.indexOf("Finished: SUCCESS") > -1)) {
                  this.isInDeploySuccess = true;
                  this.isInDeployIng = false;
                  this.processImageSubscription.unsubscribe();
                }
                if (consoleTextArr.find(value => value.indexOf("Finished: FAILURE") > -1)) {
                  this.isInDeploySuccess = false;
                  this.isInDeployIng = false;
                  this.processImageSubscription.unsubscribe();
                }
              }, err => err, () => {
                this.isInDeploySuccess = false;
                this.isInDeployIng = false;
              });
          }, 10000);
        })
        .catch(err => {
          this.messageService.dispatchError(err);
          this.isInDeploySuccess = false;
          this.isInDeployIng = false;
        })
    }
  }

  deployComplete(): void {
    this.k8sService.stepSource.next(0);
  }

  deleteDeploy(): void {
    let m: Message = new Message();
    m.title = "SERVICE.STEP_6_CANCEL_TITLE";
    m.buttons = BUTTON_STYLE.DELETION;
    m.message = "SERVICE.STEP_6_CANCEL_MSG";
    this.messageService.announceMessage(m);
  }

  backStep(): void {
    this.k8sService.stepSource.next(4);
  }
}