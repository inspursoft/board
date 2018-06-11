import { K8sService } from "./service.k8s";
import { Injector, OnDestroy } from "@angular/core";
import { AppInitService } from "../app.init.service";
import { MessageService } from "../shared/message-service/message.service";
import { ServiceStepPhase, UiServiceFactory, UIServiceStepBase } from "./service-step.component";
import { Router } from "@angular/router";
import { Message } from "../shared/message-service/message";
import { BUTTON_STYLE, MESSAGE_TARGET } from "../shared/shared.const";
import { Subscription } from "rxjs/Subscription";

export abstract class ServiceStepBase implements OnDestroy{
  protected k8sService: K8sService;
  protected appInitService: AppInitService;
  protected messageService: MessageService;
  protected uiBaseData: UIServiceStepBase;
  protected router: Router;
  protected confirmSubscription: Subscription;
  public isBack: boolean = false;

  protected constructor(protected injector: Injector) {
    this.k8sService = injector.get(K8sService);
    this.appInitService = injector.get(AppInitService);
    this.messageService = injector.get(MessageService);
    this.router = injector.get(Router);
    this.uiBaseData = UiServiceFactory.getInstance(this.stepPhase);//init empty object for template
  }

  ngOnDestroy(){
    if (this.confirmSubscription) {
      this.confirmSubscription.unsubscribe();
    }
  }

  public cancelBuildService(): void {
    if (this.confirmSubscription) {
      this.confirmSubscription.unsubscribe();
    }
    this.confirmSubscription = this.messageService.messageConfirmed$
      .subscribe((next: Message) => {
        if (next.target == MESSAGE_TARGET.CANCEL_BUILD_SERVICE) {
          this.k8sService.cancelBuildService();
        }
      });
    let msg: Message = new Message();
    msg.title = "SERVICE.ASK_TITLE";
    msg.buttons = BUTTON_STYLE.YES_NO;
    msg.message = "SERVICE.ASK_TEXT";
    msg.target = MESSAGE_TARGET.CANCEL_BUILD_SERVICE;
    this.messageService.announceMessage(msg);
  }

  get stepPhase(): ServiceStepPhase {
    return null;
  }

  get uiData(): UIServiceStepBase {
    return this.uiBaseData;
  }
}