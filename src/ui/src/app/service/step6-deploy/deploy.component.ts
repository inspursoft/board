/**
 * Created by liyanq on 9/17/17.
 */
import { Component, ComponentFactoryResolver, Injector, OnDestroy, OnInit, ViewChild, ViewContainerRef } from "@angular/core"
import { Subscription } from "rxjs/Subscription";
import { Message } from "../../shared/message-service/message";
import { BUTTON_STYLE, MESSAGE_TARGET } from "../../shared/shared.const";
import { ServiceStepBase } from "../service-step";
import { PHASE_ENTIRE_SERVICE, ServiceStepPhase, UIServiceStepBase } from "../service-step.component";
import { CsSyntaxHighlighterComponent } from "../../shared/cs-components-library/cs-syntax-highlighter/cs-syntax-highlighter.component";

@Component({
  templateUrl: "./deploy.component.html",
  styleUrls: ["./deploy.component.css"]
})
export class DeployComponent extends ServiceStepBase implements OnInit, OnDestroy {
  @ViewChild("consoleView", {read: ViewContainerRef}) consoleView: ViewContainerRef;
  boardHost: string;
  isDeployed: boolean = false;
  isDeploySuccess: boolean = false;
  isInDeployWIP: boolean = false;
  serviceID: number = 0;
  _confirmSubscription: Subscription;

  constructor(protected injector: Injector, private resolver: ComponentFactoryResolver) {
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
          let factory = this.resolver.resolveComponentFactory(CsSyntaxHighlighterComponent);
          let componentRef = this.consoleView.createComponent(factory);
          componentRef.instance.language = 'json';
          componentRef.instance.jsonContent = res;
          let msg: Message = new Message();
          msg.message = "SERVICE.STEP_6_DEPLOY_SUCCESS";
          this.messageService.inlineAlertMessage(msg);
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