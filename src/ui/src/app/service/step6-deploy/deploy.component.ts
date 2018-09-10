/**
 * Created by liyanq on 9/17/17.
 */
import { Component, Injector, OnDestroy, OnInit } from "@angular/core"
import { Subscription } from "rxjs/Subscription";
import { Message } from "../../shared/message-service/message";
import { BUTTON_STYLE, MESSAGE_TARGET, MESSAGE_TYPE } from "../../shared/shared.const";
import { ServiceStepBase } from "../service-step";
import { PHASE_ENTIRE_SERVICE, ServiceStepPhase, UIServiceStepBase } from "../service-step.component";
import { HttpErrorResponse } from "@angular/common/http";

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
  _confirmSubscription: Subscription;
  deployConsole:Object;

  constructor(protected injector: Injector) {
    super(injector);
    this.boardHost = this.appInitService.systemInfo['board_host'];
  }

  ngOnInit() {
    this._confirmSubscription = this.messageService.messageConfirmed$.subscribe((msg: Message) => {
      if (msg.target == MESSAGE_TARGET.DELETE_SERVICE_DEPLOYMENT) {
        this.k8sService.deleteDeployment(this.serviceID)
          .then(() => this.k8sService.stepSource.next({index: 0, isBack: false}))
          .catch(err => {
            this.messageService.dispatchError(err);
            this.k8sService.stepSource.next({index: 0, isBack: false});
          })
      }
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
      this.k8sService.serviceDeployment()
        .then(res => {
          this.serviceID = res['service_id'];
          this.deployConsole = res;
          let msg: Message = new Message();
          msg.message = "SERVICE.STEP_6_DEPLOY_SUCCESS";
          this.messageService.inlineAlertMessage(msg);
          this.isDeploySuccess = true;
          this.isInDeployWIP = false;
        })
        .catch((err: HttpErrorResponse) => {
          let msg = new Message();
          msg.type = MESSAGE_TYPE.SHOW_DETAIL;
          msg.message = (typeof err.error == "object") ? (err.error as Error).message : err.error;
          msg.errorObject = err;
          this.messageService.globalMessage(msg);
          this.isDeploySuccess = false;
          this.isInDeployWIP = false;
        })
    }
  }

  deleteDeploy(): void {
    let msg: Message = new Message();
    msg.title = "SERVICE.STEP_6_DELETE_TITLE";
    msg.buttons = BUTTON_STYLE.DELETION;
    msg.message = "SERVICE.STEP_6_DELETE_MSG";
    msg.target = MESSAGE_TARGET.DELETE_SERVICE_DEPLOYMENT;
    this.messageService.announceMessage(msg);
  }

  deployComplete(): void {
    this.k8sService.stepSource.next({isBack: false, index: 0});
  }

  backStep(): void {
    this.k8sService.stepSource.next({index: 4, isBack: true});
  }
}