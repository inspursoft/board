import { Component, EventEmitter, Input, OnInit, Output, ViewChild, ViewContainerRef } from '@angular/core';
import { Service } from "../../../service";
import { K8sService } from "../../../service.k8s";
import { IScaleInfo } from "../service-control.component";
import { BUTTON_STYLE, Message, RETURN_STATUS, ServiceHPA } from "../../../../shared/shared.types";
import { MessageService } from "../../../../shared.service/message.service";
import { CsComponentBase } from "../../../../shared/cs-components-library/cs-component-base";

enum ScaleMethod {smManually, smAuto}

@Component({
  selector: 'scale',
  templateUrl: './scale.component.html',
  styleUrls: ['./scale.component.css']
})
export class ScaleComponent extends CsComponentBase implements OnInit {
  @ViewChild('alertView', {read: ViewContainerRef}) alertView: ViewContainerRef;
  @Input('isActionInWIP') isActionInWIP: boolean;
  @Input('service') service: Service;
  @Output('isActionInWIPChange') onActionInWIPChange: EventEmitter<boolean>;
  @Output("onMessage") onMessage: EventEmitter<string>;
  @Output("onError") onError: EventEmitter<any>;
  @Output("onActionIsEnabled") onActionIsEnabled: EventEmitter<boolean>;
  scaleModule: ScaleMethod = ScaleMethod.smManually;
  dropDownListNum: Array<number>;
  scaleNum: number = 0;
  scaleInfo: IScaleInfo = {desired_instance: 0, available_instance: 0};
  autoScaleConfig: Array<ServiceHPA>;
  patternHpaName: RegExp = /^[a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*$/;

  constructor(private k8sService: K8sService,
              private messageService: MessageService) {
    super();
    this.dropDownListNum = Array<number>();
    this.autoScaleConfig = Array<ServiceHPA>();
    this.onMessage = new EventEmitter<string>();
    this.onError = new EventEmitter<any>();
    this.onActionIsEnabled = new EventEmitter<boolean>();
    this.onActionInWIPChange = new EventEmitter<boolean>();
  }

  ngOnInit() {
    this.onActionIsEnabled.emit(false);
    for (let i = 1; i <= 10; i++) {
      this.dropDownListNum.push(i)
    }
    this.k8sService.getServiceScaleInfo(this.service.service_id).subscribe((scaleInfo: IScaleInfo) => {
      this.scaleInfo = scaleInfo;
      this.scaleNum = this.scaleInfo.available_instance;
      this.actionEnabled();
    });
    this.k8sService.getAutoScaleConfig(this.service.service_id).subscribe((res: Array<ServiceHPA>) => {
      this.autoScaleConfig = res;
      if (this.autoScaleConfig.length > 0) {
        this.scaleModule = ScaleMethod.smAuto;
      }
      this.actionEnabled();
    });
  }

  setScaleMethod(scaleMethod: ScaleMethod): void {
    this.scaleModule = scaleMethod;
    this.actionEnabled();
  }

  actionExecute() {
    if (this.verifyInputExValid()) {
      if (this.scaleModule == ScaleMethod.smManually) {
        this.onActionInWIPChange.emit(true);
        this.k8sService.setServiceScale(this.service.service_id, this.scaleNum).subscribe(
          () => this.onMessage.emit('SERVICE.SERVICE_CONTROL_SCALE_SUCCESSFUL'),
          (err) => this.onError.emit(err)
        );
      } else {
        this.autoScaleConfig.forEach((config: ServiceHPA) => {
          if (config.min_pod > config.max_pod) {
            this.messageService.showAlert('SERVICE.SERVICE_CONTROL_HPA_WARNING', {view: this.alertView, alertType: 'warning'});
            this.onActionInWIPChange.emit(false);
          } else {
            this.onActionInWIPChange.emit(true);
            if (config.isEdit) {
              Reflect.deleteProperty(config, 'isEdit');
              this.k8sService.modifyAutoScaleConfig(this.service.service_id, config)
                .subscribe(() => this.onMessage.emit('SERVICE.SERVICE_CONTROL_SCALE_SUCCESSFUL'),
                  err => this.onError.emit(err))
            } else {
              Reflect.deleteProperty(config, 'isEdit');
              this.k8sService.setAutoScaleConfig(this.service.service_id, config)
                .subscribe(() => this.onMessage.emit('SERVICE.SERVICE_CONTROL_SCALE_SUCCESSFUL'),
                  err => this.onError.emit(err))
            }
          }
        });
      }
    } else {
      this.onActionInWIPChange.emit(false);
    }
  }

  actionEnabled() {
    let enabled = this.scaleModule == ScaleMethod.smManually
      ? this.scaleNum > 0 && this.scaleNum != this.scaleInfo.available_instance
      : this.autoScaleConfig.length > 0;
    this.onActionIsEnabled.emit(enabled);
  }

  addOneHpa() {
    if (this.autoScaleConfig.length == 0) {
      this.autoScaleConfig.push(new ServiceHPA());
      this.actionEnabled();
    }
  }

  deleteOneHpa(hpa: ServiceHPA) {
    if (!this.isActionInWIP) {
      this.messageService.showDialog('SERVICE.SERVICE_CONTROL_HPA_CONFIRM_DELETE',
        {view: this.alertView, buttonStyle: BUTTON_STYLE.DELETION, title: 'GLOBAL_ALERT.DELETE'}).subscribe((res: Message) => {
        if (res.returnStatus == RETURN_STATUS.rsConfirm) {
          this.onActionInWIPChange.emit(true);
          this.k8sService.deleteAutoScaleConfig(this.service.service_id, hpa).subscribe(
            () => this.onMessage.emit('SERVICE.SERVICE_CONTROL_HPA_DELETE_SUCCESS'),
            err => this.onError.emit(err),
            () => this.actionEnabled());
        }
      })
    }
  }
}
