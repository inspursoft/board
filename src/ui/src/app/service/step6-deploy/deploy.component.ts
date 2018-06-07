/**
 * Created by liyanq on 9/17/17.
 */
import { Component, Injector, OnDestroy, OnInit } from "@angular/core"
import { Subscription } from "rxjs/Subscription";
import { Message } from "../../shared/message-service/message";
import { BUTTON_STYLE } from "../../shared/shared.const";
import { ServiceStepBase } from "../service-step";
import { HttpErrorResponse } from "@angular/common/http"
import { PHASE_ENTIRE_SERVICE, ServiceStepPhase, UIServiceStepBase } from "../service-step.component";

@Component({
  templateUrl: "./deploy.component.html",
  styleUrls: ["./deploy.component.css"]
})
export class DeployComponent extends ServiceStepBase implements OnInit, OnDestroy {
  boardHost: string;
  isDeployed: boolean = false;
  isDeploySuccess: boolean = false;
  isInDeployWIP: boolean = false;
  serviceID: number = 0;
  consoleText: string = "";
  _confirmSubscription: Subscription;

  constructor(protected injector: Injector) {
    super(injector);
    this.boardHost = this.appInitService.systemInfo['board_host'];
  }

  ngOnInit() {
    this._confirmSubscription = this.messageService.messageConfirmed$.subscribe((next: Message) => {
      this.k8sService.deleteDeployment(this.serviceID)
        .then(() => this.k8sService.stepSource.next({index: 0, isBack: false}))
        .catch(err => {
          this.messageService.dispatchError(err);
          this.k8sService.stepSource.next({index: 0, isBack: false});
        })
    });
  }

  ngOnDestroy() {
    this._confirmSubscription.unsubscribe();
  }

  get stepPhase(): ServiceStepPhase {
    return PHASE_ENTIRE_SERVICE;
  }

  get uiData(): UIServiceStepBase {
    return this.uiBaseData;
  }

  serviceDeploy() {
    if (!this.isDeployed) {
      this.isDeployed = true;
      this.isInDeployWIP = true;
      this.consoleText = "SERVICE.STEP_6_DEPLOYING";
      this.k8sService.serviceDeployment()
        .then(res => {
          this.serviceID = res['service_id'];
          this.consoleText = JSON.stringify(res);
          this.isDeploySuccess = true;
          this.isInDeployWIP = false;
        })
        .catch(err => {
          this.messageService.dispatchError(err,"",true);
          this.isDeploySuccess = false;
          this.isInDeployWIP = false;
        })
    }
  }

  deleteDeploy(): void {
    let m: Message = new Message();
    m.title = "SERVICE.STEP_6_CANCEL_TITLE";
    m.buttons = BUTTON_STYLE.DELETION;
    m.message = "SERVICE.STEP_6_CANCEL_MSG";
    this.messageService.announceMessage(m);
  }

  deployComplete(): void {
    this.k8sService.stepSource.next({isBack: false, index: 0});
  }

  backStep(): void {
    this.k8sService.stepSource.next({index: 4, isBack: true});
  }
}