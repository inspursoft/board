import { Component, EventEmitter, Input, OnInit, Output, ViewChild, ViewContainerRef } from '@angular/core';
import { K8sService } from '../../../service.k8s';
import { IScaleInfo } from '../service-control.component';
import { BUTTON_STYLE, Message, RETURN_STATUS } from '../../../../shared/shared.types';
import { MessageService } from '../../../../shared.service/message.service';
import { CsComponentBase } from '../../../../shared/cs-components-library/cs-component-base';
import { Service, ServiceHPA } from '../../../service.types';

enum ScaleMethod {smManually, smAuto}

@Component({
  selector: 'app-scale',
  templateUrl: './scale.component.html',
  styleUrls: ['./scale.component.css']
})
export class ScaleComponent extends CsComponentBase implements OnInit {
  @ViewChild('alertView', {read: ViewContainerRef}) alertView: ViewContainerRef;
  @Input() isActionInWIP: boolean;
  @Input() service: Service;
  @Output() isActionInWIPChange: EventEmitter<boolean>;
  @Output() messageEvent: EventEmitter<string>;
  @Output() errorEvent: EventEmitter<any>;
  @Output() actionIsEnabledEvent: EventEmitter<boolean>;
  scaleModule: ScaleMethod = ScaleMethod.smManually;
  dropDownListNum: Array<number>;
  scaleNum = 0;
  scaleInfo: IScaleInfo = {desired_instance: 0, available_instance: 0};
  autoScaleConfig: Array<ServiceHPA>;
  patternHpaName: RegExp = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$/;

  constructor(private k8sService: K8sService,
              private messageService: MessageService) {
    super();
    this.dropDownListNum = Array<number>();
    this.autoScaleConfig = Array<ServiceHPA>();
    this.messageEvent = new EventEmitter<string>();
    this.errorEvent = new EventEmitter<any>();
    this.actionIsEnabledEvent = new EventEmitter<boolean>();
    this.isActionInWIPChange = new EventEmitter<boolean>();
  }

  ngOnInit() {
    this.actionIsEnabledEvent.emit(false);
    for (let i = 1; i <= 10; i++) {
      this.dropDownListNum.push(i);
    }
    this.k8sService.getServiceScaleInfo(this.service.serviceId).subscribe(
      (scaleInfo: IScaleInfo) => {
        this.scaleInfo = scaleInfo;
        this.scaleNum = this.scaleInfo.available_instance;
        this.actionEnabled();
      }, (err) => this.errorEvent.emit(err)
    );
    this.k8sService.getAutoScaleConfig(this.service.serviceId).subscribe(
      (res: Array<ServiceHPA>) => {
        this.autoScaleConfig = res;
        if (this.autoScaleConfig.length > 0) {
          this.scaleModule = ScaleMethod.smAuto;
        }
        this.actionEnabled();
      }, (err) => this.errorEvent.emit(err)
    );
  }

  setScaleMethod(scaleMethod: ScaleMethod): void {
    this.scaleModule = scaleMethod;
    this.actionEnabled();
  }

  actionExecute() {
    if (this.verifyInputExValid()) {
      if (this.scaleModule === ScaleMethod.smManually) {
        this.isActionInWIPChange.emit(true);
        this.k8sService.setServiceScale(this.service.serviceId, this.scaleNum).subscribe(
          () => this.messageEvent.emit('SERVICE.SERVICE_CONTROL_SCALE_SUCCESSFUL'),
          (err) => this.errorEvent.emit(err)
        );
      } else {
        this.autoScaleConfig.forEach((config: ServiceHPA) => {
          if (config.minPod > config.maxPod) {
            this.messageService.showAlert('SERVICE.SERVICE_CONTROL_HPA_WARNING', {view: this.alertView, alertType: 'warning'});
            this.isActionInWIPChange.emit(false);
          } else {
            this.isActionInWIPChange.emit(true);
            if (config.isEdit) {
              Reflect.deleteProperty(config, 'isEdit');
              this.k8sService.modifyAutoScaleConfig(this.service.serviceId, config)
                .subscribe(() => this.messageEvent.emit('SERVICE.SERVICE_CONTROL_SCALE_SUCCESSFUL'),
                  err => this.errorEvent.emit(err));
            } else {
              Reflect.deleteProperty(config, 'isEdit');
              this.k8sService.setAutoScaleConfig(this.service.serviceId, config)
                .subscribe(() => this.messageEvent.emit('SERVICE.SERVICE_CONTROL_SCALE_SUCCESSFUL'),
                  err => this.errorEvent.emit(err));
            }
          }
        });
      }
    } else {
      this.isActionInWIPChange.emit(false);
    }
  }

  actionEnabled() {
    const enabled = this.scaleModule === ScaleMethod.smManually
      ? this.scaleNum > 0 && this.scaleNum !== this.scaleInfo.available_instance
      : this.autoScaleConfig.length > 0;
    this.actionIsEnabledEvent.emit(enabled);
  }

  addOneHpa() {
    if (this.autoScaleConfig.length === 0) {
      this.autoScaleConfig.push(new ServiceHPA());
      this.actionEnabled();
    }
  }

  deleteOneHpa(hpa: ServiceHPA) {
    if (!this.isActionInWIP) {
      this.messageService.showDialog('SERVICE.SERVICE_CONTROL_HPA_CONFIRM_DELETE',
        {view: this.alertView, buttonStyle: BUTTON_STYLE.DELETION, title: 'GLOBAL_ALERT.DELETE'}).subscribe((res: Message) => {
        if (res.returnStatus === RETURN_STATUS.rsConfirm) {
          this.isActionInWIPChange.emit(true);
          this.k8sService.deleteAutoScaleConfig(this.service.serviceId, hpa).subscribe(
            () => this.messageEvent.emit('SERVICE.SERVICE_CONTROL_HPA_DELETE_SUCCESS'),
            err => this.errorEvent.emit(err),
            () => this.actionEnabled());
        }
      });
    }
  }
}
