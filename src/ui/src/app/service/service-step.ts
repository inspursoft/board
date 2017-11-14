import { K8sService } from "./service.k8s";
import { Injector, OnDestroy } from "@angular/core";
import { AppInitService } from "../app.init.service";
import { MessageService } from "../shared/message-service/message.service";
import { Container, DeploymentServiceData, Volume } from "./service-step.component";
import { Router } from "@angular/router";
import { Message } from "../shared/message-service/message";
import { BUTTON_STYLE } from "../shared/shared.const";
import { Subscription } from "rxjs/Subscription";

export abstract class ServiceStepBase implements OnDestroy {
  _confirmSubscription: Subscription;
  protected k8sService: K8sService;
  protected appInitService: AppInitService;
  protected messageService: MessageService;
  protected outputData: DeploymentServiceData;
  protected router: Router;
  public isBack: boolean = false;

  constructor(protected injector: Injector) {
    this.k8sService = injector.get(K8sService);
    this.appInitService = injector.get(AppInitService);
    this.messageService = injector.get(MessageService);
    this.router = injector.get(Router);
    this.outputData = new DeploymentServiceData();
    this._confirmSubscription = this.messageService.messageConfirmed$.subscribe(next => {
      this.k8sService.cancelBuildService();
    });
  }

  ngOnDestroy() {
    this._confirmSubscription.unsubscribe();
  }

  public cancelBuildService(): void {
    let m: Message = new Message();
    m.title = "SERVICE.ASK_TITLE";
    m.buttons = BUTTON_STYLE.YES_NO;
    m.message = "SERVICE.ASK_TEXT";
    this.messageService.announceMessage(m);
  }

  get containerList(): Array<Container> {
    //safe get at index > 1
    return this.outputData.deployment_yaml.spec.template.spec.containers;
  }

  get deployVolumes(): Array<Volume> {
    //safe get at index > 1
    return this.outputData.deployment_yaml.spec.template.spec.volumes;
  }

  get newServiceId(): number {
    return this.k8sService.newServiceId;
  }

  set newServiceId(value: number) {
    this.k8sService.newServiceId = value;
  }
}