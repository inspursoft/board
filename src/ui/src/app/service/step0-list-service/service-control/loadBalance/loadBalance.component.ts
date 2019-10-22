import { Component, EventEmitter, Input, OnInit, Output } from "@angular/core";
import { Service } from "../../../service";
import { K8sService } from "../../../service.k8s";

@Component({
  selector: 'load-balance',
  templateUrl: './loadBalance.component.html',
  styleUrls: ['./loadBalance.component.css']
})
export class LoadBalanceComponent implements OnInit {
  @Input() isActionInWIP: boolean;
  @Input() service: Service;
  @Output() onMessage: EventEmitter<string>;
  @Output() onError: EventEmitter<any>;
  @Output() onActionIsEnabled: EventEmitter<boolean>;
  sessionAffinityFlag = false;
  oldSessionAffinityFlag = false;

  constructor(private k8sService: K8sService) {
    this.onMessage = new EventEmitter<string>();
    this.onError = new EventEmitter<any>();
    this.onActionIsEnabled = new EventEmitter<boolean>();
  }

  ngOnInit() {
    this.k8sService.getSessionAffinityFlag(this.service.service_name, this.service.service_project_name).subscribe(flag => {
      this.sessionAffinityFlag = flag;
      this.oldSessionAffinityFlag = flag
    }, (err) => this.onError.emit(err));
  }

  actionExecute() {
    this.k8sService.setSessionAffinityFlag(this.service.service_name, this.service.service_project_name, this.sessionAffinityFlag).subscribe(
      () => this.onMessage.emit('SERVICE.SERVICE_CONTROL_SESSION_AFFINITY_SUCCESSFUL'),
      (err) => this.onError.emit(err));
  }

  setSessionAffinityFlag(flag: boolean) {
    this.sessionAffinityFlag = flag;
    this.onActionIsEnabled.emit(this.sessionAffinityFlag != this.oldSessionAffinityFlag);
  }
}