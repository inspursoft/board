/**
 * Created by liyanq on 04/12/2017.
 */

import { Component, OnInit, ViewChild } from "@angular/core"
import { Service } from "../../service";
import { K8sService } from "../../service.k8s";
import { MessageService } from "../../../shared/message-service/message.service";
import { UpdateComponent } from "./update/update.component";
import { LocateComponent } from "./locate/locate.component";
import { ScaleComponent } from "./scale/scale.component";
import { Message } from "../../../shared/message-service/message";
import { Subject } from "rxjs/Subject";
import { Observable } from "rxjs/Observable";

export interface IScaleInfo {
  desired_instance: number;
  available_instance: number;
}

enum ActionMethod {scale, update, locate}

@Component({
  selector: "service-control",
  styleUrls: ["./service-control.component.css"],
  templateUrl: "./service-control.component.html"
})
export class ServiceControlComponent implements OnInit {
  @ViewChild(UpdateComponent) updateComponent: UpdateComponent;
  @ViewChild(ScaleComponent) scaleComponent: ScaleComponent;
  @ViewChild(LocateComponent) locateComponent: LocateComponent;
  service: Service;
  _isOpen: boolean = false;
  actionMethod: ActionMethod = ActionMethod.scale;
  actionEnable: boolean = false;
  isActionInWIP: boolean = false;
  alertMessage: string = "";
  closeNotification: Subject<any>;

  constructor(private k8sService: K8sService,
              private messageService: MessageService) {
    this.closeNotification = new Subject<any>();
  }

  ngOnInit() {
  }

  get isOpen():boolean{
    return this._isOpen;
  }

  set isOpen(value: boolean){
    this._isOpen = value;
    if (!this._isOpen){
      this.closeNotification.next();
    }
  }

  public openModal(): Observable<any> {
    this.isOpen = true;
    return this.closeNotification.asObservable();
  }

  defaultDispatchErr(err) {
    this.isOpen = false;
    this.messageService.dispatchError(err);
  }

  defaultHandleMessage(msg: Message) {
    this.isOpen = false;
    this.messageService.inlineAlertMessage(msg);
  }

  defaultHandleAlertMessage(msg: string) {
    this.alertMessage = msg;
  }

  defaultHandleActionEnabled(enabled: boolean){
    this.actionEnable = enabled;
  }

  actionExecute() {
    this.isActionInWIP = true;
    if (this.actionMethod == ActionMethod.update) {
      this.updateComponent.actionExecute();
    } else if (this.actionMethod == ActionMethod.scale) {
      this.scaleComponent.actionExecute();
    } else {
      this.locateComponent.actionExecute();
    }
  }
}