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
import { CsModalChildBase } from "../../../shared/cs-modal-base/cs-modal-child-base";
import { TranslateService } from "@ngx-translate/core";
import { LoadBalanceComponent } from "./loadBalance/loadBalance.component";
import { Observable, Subject } from "rxjs";

export interface IScaleInfo {
  desired_instance: number;
  available_instance: number;
}

enum ActionMethod {scale, update, locate, loadBalance}

@Component({
  selector: "service-control",
  styleUrls: ["./service-control.component.css"],
  templateUrl: "./service-control.component.html"
})
export class ServiceControlComponent extends CsModalChildBase implements OnInit {
  @ViewChild(UpdateComponent) updateComponent: UpdateComponent;
  @ViewChild(ScaleComponent) scaleComponent: ScaleComponent;
  @ViewChild(LocateComponent) locateComponent: LocateComponent;
  @ViewChild(LoadBalanceComponent) loadBalanceComponent: LoadBalanceComponent;
  service: Service;
  _isOpen: boolean = false;
  actionMethod: ActionMethod = ActionMethod.scale;
  actionEnable: boolean = false;
  isActionInWIP: boolean = false;
  closeNotification: Subject<any>;

  constructor(private k8sService: K8sService,
              private translateService: TranslateService,
              private messageService: MessageService) {
    super();
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
  }

  defaultHandleMessage(msg: string) {
    this.isOpen = false;
    this.translateService.get(msg, [this.service.service_name])
      .subscribe((res: string) => this.messageService.showAlert(res));
  }

  defaultHandleAlertMessage(msg: string) {
    this.messageService.showAlert(msg, {alertType: 'alert-warning', view: this.alertView})
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
    } else if (this.actionMethod == ActionMethod.locate){
      this.locateComponent.actionExecute();
    } else {
      this.loadBalanceComponent.actionExecute();
    }
  }
}
