/**
 * Created by liyanq on 9/17/17.
 */
import { Component, Injector } from "@angular/core"
import { ServiceStepBase } from "../service-step";
import { PHASE_ENTIRE_SERVICE, ServiceStepPhase, UIServiceStepBase } from "../service-step.component";
import { HttpErrorResponse } from "@angular/common/http";
import { GlobalAlertType, Message, RETURN_STATUS } from "../../shared/shared.types";
import { ClrLoadingButton } from "@clr/angular";

@Component({
  templateUrl: "./deploy.component.html",
  styleUrls: ["./deploy.component.css"]
})
export class DeployComponent extends ServiceStepBase {
  boardHost: string;
  isDeployed: boolean = false;
  isDeploySuccess: boolean = false;
  isInDeployWIP: boolean = false;
  serviceID: number = 0;
  deployConsole:Object;

  constructor(protected injector: Injector) {
    super(injector);
    this.boardHost = this.appInitService.systemInfo['board_host'];
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
          this.messageService.showAlert('SERVICE.STEP_6_DEPLOY_SUCCESS');
          this.isDeploySuccess = true;
          this.isInDeployWIP = false;
        })
        .catch((err: HttpErrorResponse) => {
          this.messageService.showGlobalMessage('SERVICE.STEP_6_DEPLOY_FAILED', {
            globalAlertType: GlobalAlertType.gatShowDetail,
            errorObject: err
          });
          this.isDeploySuccess = false;
          this.isInDeployWIP = false;
        })
    }
  }

  deleteDeploy(): void {
    this.messageService.showDeleteDialog('SERVICE.STEP_6_DELETE_MSG', 'SERVICE.STEP_6_DELETE_TITLE').subscribe((message: Message) => {
      if (message.returnStatus == RETURN_STATUS.rsConfirm) {
        this.isInDeployWIP = true;
        this.k8sService.deleteDeployment(this.serviceID)
          .then(() => this.k8sService.stepSource.next({index: 0, isBack: false}))
          .catch(() => this.k8sService.stepSource.next({index: 0, isBack: false}))
      }
    })
  }

  deployComplete(): void {
    this.k8sService.stepSource.next({isBack: false, index: 0});
  }

  backStep(): void {
    this.k8sService.stepSource.next({index: 3, isBack: true});
  }
}