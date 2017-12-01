import { K8sService } from "./service.k8s";
import { Injector } from "@angular/core";
import { AppInitService } from "../app.init.service";
import { MessageService } from "../shared/message-service/message.service";
import { UIServiceStepBase, ServiceStepPhase, UiServiceFactory } from "./service-step.component";
import { Router } from "@angular/router";
import { Message } from "../shared/message-service/message";
import { BUTTON_STYLE } from "../shared/shared.const";

export abstract class ServiceStepBase {
  protected k8sService: K8sService;
  protected appInitService: AppInitService;
  protected messageService: MessageService;
  protected uiBaseData: UIServiceStepBase;
  protected router: Router;
  public isBack: boolean = false;

  constructor(protected injector: Injector) {
    this.k8sService = injector.get(K8sService);
    this.appInitService = injector.get(AppInitService);
    this.messageService = injector.get(MessageService);
    this.router = injector.get(Router);
    this.uiBaseData = UiServiceFactory.getInstance(this.stepPhase);//init empty object for template
  }

  public cancelBuildService(): void {
    let confirmSubscription = this.messageService.messageConfirmed$.subscribe(next => {
      this.k8sService.cancelBuildService();
      confirmSubscription.unsubscribe();
    });
    let m: Message = new Message();
    m.title = "SERVICE.ASK_TITLE";
    m.buttons = BUTTON_STYLE.YES_NO;
    m.message = "SERVICE.ASK_TEXT";
    this.messageService.announceMessage(m);
  }

  get stepPhase(): ServiceStepPhase {
    return null;
  }

  get uiData(): UIServiceStepBase {
    return this.uiBaseData;
  }
}