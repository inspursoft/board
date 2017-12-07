/**
 * Created by liyanq on 04/12/2017.
 */

import { Component, EventEmitter, Input, Output, OnInit } from "@angular/core"
import { Service } from "../../service";
import { K8sService } from "../../service.k8s";
import { MessageService } from "../../../shared/message-service/message.service";

enum ScaleMethod{smNone, smManually, smAuto}
@Component({
  selector: "service-control",
  styleUrls: ["./service-control.component.css"],
  templateUrl: "./service-control.component.html"
})
export class ServiceControlComponent implements OnInit {
  _isOpen: boolean = false;
  dropDownListNum: Array<number>;
  scaleModule: ScaleMethod = ScaleMethod.smNone;
  scaleNum: number = 0;
  @Input() service: Service;

  constructor(private k8sService: K8sService,
              private messageService: MessageService) {
    this.dropDownListNum = Array<number>();
  }

  ngOnInit() {
    for (let i = 1; i <= 100; i++) {
      this.dropDownListNum.push(i)
    }
  }

  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();

  @Input() get isOpen() {
    return this._isOpen;
  }

  set isOpen(open: boolean) {
    this._isOpen = open;
    this.isOpenChange.emit(this._isOpen);
  }

  get reason(): string {
    let result: string = this.service["service_comment"];
    if (result.toLowerCase().startsWith("reason:")) {
      result = result.slice(7)
    }
    return result;
  }

  getServiceStatus(status: number): string {
    switch (status) {
      case 0:
        return 'SERVICE.STATUS_PREPARING';
      case 1:
        return 'SERVICE.STATUS_RUNNING';
      case 2:
        return 'SERVICE.STATUS_STOPPED';
      case 3:
        return 'SERVICE.STATUS_WARNING';
    }
  }

  getStatusClass(status: number) {
    return {
      'running': status == 1,
      'stopped': status == 2,
      'warning': status == 3
    }
  }

  setScaleMethod(scaleMethod: ScaleMethod): void {
    console.log(scaleMethod);
    this.scaleModule = scaleMethod;
  }

  setServiceScale() {
    this.k8sService.setServiceScale(this.service.service_id, this.scaleNum)
      .then(res => this.isOpen = false)
      .catch(err => {
        this.isOpen = false;
        this.messageService.dispatchError(err);
      })
  }

}