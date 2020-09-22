/**
 * Created by liyanq on 04/12/2017.
 */
import { AfterViewInit, Component, OnInit, ViewChild, ViewContainerRef } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { K8sService } from '../../service.k8s';
import { MessageService } from '../../../shared.service/message.service';
import { UpdateComponent } from './update/update.component';
import { LocateComponent } from './locate/locate.component';
import { ScaleComponent } from './scale/scale.component';
import { CsModalChildBase } from '../../../shared/cs-modal-base/cs-modal-child-base';
import { LoadBalanceComponent } from './loadBalance/loadBalance.component';
import { ConsoleComponent } from './console/console.component';
import { GlobalAlertType } from '../../../shared/shared.types';
import { HttpErrorResponse, HttpProgressEvent } from '@angular/common/http';
import { Service } from '../../service.types';

export interface IScaleInfo {
  desired_instance: number;
  available_instance: number;
}

enum ActionMethod {scale, update, locate, loadBalance, console}

@Component({
  selector: 'app-service-control',
  styleUrls: ['./service-control.component.css'],
  templateUrl: './service-control.component.html'
})
export class ServiceControlComponent extends CsModalChildBase implements OnInit, AfterViewInit {
  @ViewChild(UpdateComponent) updateComponent: UpdateComponent;
  @ViewChild(ScaleComponent) scaleComponent: ScaleComponent;
  @ViewChild(LocateComponent) locateComponent: LocateComponent;
  @ViewChild(LoadBalanceComponent) loadBalanceComponent: LoadBalanceComponent;
  @ViewChild(ConsoleComponent) consoleComponent: ConsoleComponent;
  service: Service;
  actionMethod: ActionMethod = ActionMethod.scale;
  actionEnable = false;
  isActionInWIP = false;
  uploadProgressValue: HttpProgressEvent;

  constructor(private k8sService: K8sService,
              private translateService: TranslateService,
              private messageService: MessageService) {
    super();
  }

  ngOnInit() {
  }

  ngAfterViewInit(): void {
    const header = document.getElementsByClassName('modal-header');
    if (header.length > 0) {
      (header.item(0) as HTMLElement).style.padding = `0px`;
    }
  }

  defaultDispatchErr(err: HttpErrorResponse) {
    this.modalOpened = false;
    const errMsg = typeof err.error === 'string' ? err.error : err.message;
    this.messageService.showGlobalMessage(errMsg, {
        alertType: 'danger',
        globalAlertType: GlobalAlertType.gatShowDetail,
        errorObject: err
      }
    );
  }

  defaultHandleMessage(msg: string) {
    this.modalOpened = false;
    this.translateService.get(msg, [this.service.serviceName]).subscribe(
      (res: string) => this.messageService.showAlert(res)
    );
  }

  defaultHandleAlertMessage(msg: string) {
    this.messageService.showAlert(msg, {
      alertType: 'warning',
      view: this.alertView
    });
  }

  defaultHandleActionEnabled(enabled: boolean) {
    this.actionEnable = enabled;
  }

  defaultHandleProgress(progress: HttpProgressEvent) {
    this.uploadProgressValue = progress;
  }

  actionExecute() {
    this.isActionInWIP = true;
    if (this.actionMethod === ActionMethod.update) {
      this.updateComponent.actionExecute();
    } else if (this.actionMethod === ActionMethod.scale) {
      this.scaleComponent.actionExecute();
    } else if (this.actionMethod === ActionMethod.locate) {
      this.locateComponent.actionExecute();
    } else if (this.actionMethod === ActionMethod.console) {
      this.modalOpened = false;
    } else {
      this.loadBalanceComponent.actionExecute();
    }
  }
}
