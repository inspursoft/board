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
  isOpen: boolean = false;
  actionMethod: ActionMethod = ActionMethod.scale;
  actionEnable: boolean = false;
  isActionInWIP: boolean = false;
  alertMessage: string = "";

  constructor(private k8sService: K8sService,
              private messageService: MessageService) {

  }

  ngOnInit() {
  }

  public openModal() {
    this.alertMessage = '';
    this.isActionInWIP = false;
    this.isOpen = true;
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