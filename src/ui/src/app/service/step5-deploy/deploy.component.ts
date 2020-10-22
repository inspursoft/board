/**
 * Created by liyanq on 9/17/17.
 */
import { Component, Injector, OnInit } from '@angular/core';
import { Observable } from 'rxjs';
import { PHASE_ENTIRE_SERVICE, PHASE_EXTERNAL_SERVICE, ServiceStep3Data, ServiceStepPhase } from '../service-step.component';
import { HttpErrorResponse } from '@angular/common/http';
import { GlobalAlertType, Message, RETURN_STATUS } from '../../shared/shared.types';
import { ServiceStepComponentBase } from '../service-step';
import { ServiceType } from '../service.types';

@Component({
  templateUrl: './deploy.component.html',
  styleUrls: ['./deploy.component.css']
})
export class DeployComponent extends ServiceStepComponentBase implements OnInit {
  isDeployed = false;
  isDeploySuccess = false;
  isInDeployWIP = false;
  isDeleteInWIP = false;
  serviceID = 0;
  deployConsole: object;
  serviceType: ServiceType = ServiceType.ServiceTypeUnknown;

  constructor(protected injector: Injector) {
    super(injector);
  }

  ngOnInit(): void {
    this.k8sService.getServiceConfig(PHASE_EXTERNAL_SERVICE, ServiceStep3Data).subscribe(
      (serviceStep3Data: ServiceStep3Data) => this.serviceType = serviceStep3Data.serviceType
    );
  }

  get stepPhase(): ServiceStepPhase {
    return PHASE_ENTIRE_SERVICE;
  }

  serviceDeploy() {
    if (!this.isDeployed) {
      this.isDeployed = true;
      this.isInDeployWIP = true;
      let obsDeploy: Observable<object>;
      if (this.serviceType === ServiceType.ServiceTypeStatefulSet) {
        obsDeploy = this.k8sService.serviceStatefulDeployment();
      } else {
        obsDeploy = this.k8sService.serviceDeployment();
      }
      obsDeploy.subscribe(res => {
        this.serviceID = Reflect.get(res, 'service_id');
        this.deployConsole = res;
        this.messageService.showAlert('SERVICE.STEP_5_DEPLOY_SUCCESS');
        this.isDeploySuccess = true;
        this.isInDeployWIP = false;
      }, (err: HttpErrorResponse) => {
        this.messageService.showGlobalMessage('SERVICE.STEP_5_DEPLOY_FAILED', {
          globalAlertType: GlobalAlertType.gatShowDetail,
          errorObject: err
        });
        this.isDeploySuccess = false;
        this.isInDeployWIP = false;
      });
    }
  }

  deleteDeploy(): void {
    this.messageService.showDeleteDialog('SERVICE.STEP_5_DELETE_MSG', 'SERVICE.STEP_5_DELETE_TITLE').subscribe(
      (message: Message) => {
        if (message.returnStatus === RETURN_STATUS.rsConfirm) {
          this.isDeleteInWIP = true;
          this.k8sService.deleteDeployment(this.serviceID).subscribe(
            () => this.k8sService.stepSource.next({index: 0, isBack: false}),
            () => this.k8sService.stepSource.next({index: 0, isBack: false})
          );
        }
      });
  }

  deployComplete(): void {
    this.k8sService.stepSource.next({isBack: false, index: 0});
  }

  backStep(): void {
    this.k8sService.stepSource.next({index: 3, isBack: true});
  }
}
