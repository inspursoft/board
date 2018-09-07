import { K8sService } from "./service.k8s";
import { Injector } from "@angular/core";
import { AppInitService } from "../app.init.service";
import { MessageService } from "../shared/message-service/message.service";
import { ServiceStepPhase, UiServiceFactory, UIServiceStepBase } from "./service-step.component";
import { Router } from "@angular/router";
import { CsComponentBase } from "../shared/cs-components-library/cs-component-base";
import { Message, RETURN_STATUS } from "../shared/shared.types";

export abstract class ServiceStepBase extends CsComponentBase {
  protected k8sService: K8sService;
  protected appInitService: AppInitService;
  protected messageService: MessageService;
  protected uiBaseData: UIServiceStepBase;
  protected router: Router;
  public isBack: boolean = false;

  protected constructor(protected injector: Injector) {
    super();
    this.k8sService = injector.get(K8sService);
    this.appInitService = injector.get(AppInitService);
    this.messageService = injector.get(MessageService);
    this.router = injector.get(Router);
    this.uiBaseData = UiServiceFactory.getInstance(this.stepPhase);//init empty object for template
  }

  public cancelBuildService(): void {
    this.messageService.showYesNoDialog('SERVICE.ASK_TEXT','SERVICE.ASK_TITLE').subscribe((message: Message) => {
      if (message.returnStatus == RETURN_STATUS.rsConfirm) {
        this.k8sService.cancelBuildService();
      }
    });
  }

  get stepPhase(): ServiceStepPhase {
    return null;
  }

  get uiData(): UIServiceStepBase {
    return this.uiBaseData;
  }
}