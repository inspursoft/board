import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Service } from "../../../service";
import { K8sService } from "../../../service.k8s";
import { IScaleInfo } from "../service-control.component";
import "rxjs/add/observable/of"
enum ScaleMethod {smManually, smAuto}

@Component({
  selector: 'scale',
  templateUrl: './scale.component.html',
  styleUrls: ['./scale.component.css']
})
export class ScaleComponent implements OnInit {
  @Input('isActionInWIP') isActionInWIP: boolean;
  @Input('service') service: Service;
  @Output("onMessage") onMessage: EventEmitter<string>;
  @Output("onError") onError: EventEmitter<any>;
  @Output("onActionIsEnabled") onActionIsEnabled: EventEmitter<boolean>;
  scaleModule: ScaleMethod = ScaleMethod.smManually;
  dropDownListNum: Array<number>;
  scaleNum: number = 0;
  scaleInfo: IScaleInfo = {desired_instance: 0, available_instance: 0};

  constructor(private k8sService: K8sService) {
    this.dropDownListNum = Array<number>();
    this.onMessage = new EventEmitter<string>();
    this.onError = new EventEmitter<any>();
    this.onActionIsEnabled = new EventEmitter<boolean>();
  }

  ngOnInit() {
    this.onActionIsEnabled.emit(false);
    for (let i = 1; i <= 10; i++) {
      this.dropDownListNum.push(i)
    }
    this.k8sService.getServiceScaleInfo(this.service.service_id)
      .subscribe((scaleInfo: IScaleInfo) => {
        this.scaleInfo = scaleInfo;
        this.scaleNum = this.scaleInfo.available_instance;
        this.actionEnabled();
      })
  }

  setScaleMethod(scaleMethod: ScaleMethod): void {
    this.scaleModule = scaleMethod;
  }

  actionExecute() {
    this.isActionInWIP = true;
    this.k8sService.setServiceScale(this.service.service_id, this.scaleNum)
      .then(() => this.onMessage.emit('SERVICE.SERVICE_CONTROL_SCALE_SUCCESSFUL'))
      .catch((err) => this.onError.emit(err));
  }

  actionEnabled(){
    this.onActionIsEnabled.emit(this.scaleNum > 0 && this.scaleNum != this.scaleInfo.available_instance);
  }
}
