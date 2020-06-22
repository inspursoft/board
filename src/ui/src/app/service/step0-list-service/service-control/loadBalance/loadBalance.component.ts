import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { K8sService } from '../../../service.k8s';
import { Service } from '../../../service.types';

@Component({
  selector: 'app-load-balance',
  templateUrl: './loadBalance.component.html',
  styleUrls: ['./loadBalance.component.css']
})
export class LoadBalanceComponent implements OnInit {
  @Input() isActionInWIP: boolean;
  @Input() service: Service;
  @Output() messageEvent: EventEmitter<string>;
  @Output() errorEvent: EventEmitter<any>;
  @Output() actionIsEnabledEvent: EventEmitter<boolean>;
  sessionAffinityFlag = false;
  oldSessionAffinityFlag = false;

  constructor(private k8sService: K8sService) {
    this.messageEvent = new EventEmitter<string>();
    this.errorEvent = new EventEmitter<any>();
    this.actionIsEnabledEvent = new EventEmitter<boolean>();
  }

  ngOnInit() {
    this.k8sService.getSessionAffinityFlag(this.service.serviceName, this.service.serviceProjectName).subscribe(
      flag => {
        this.sessionAffinityFlag = flag;
        this.oldSessionAffinityFlag = flag;
      },
      (err) => this.errorEvent.emit(err)
    );
  }

  actionExecute() {
    const serviceName = this.service.serviceName;
    const projectName = this.service.serviceProjectName;
    this.k8sService.setSessionAffinityFlag(serviceName, projectName, this.sessionAffinityFlag).subscribe(
      () => this.messageEvent.emit('SERVICE.SERVICE_CONTROL_SESSION_AFFINITY_SUCCESSFUL'),
      (err) => this.errorEvent.emit(err)
    );
  }

  setSessionAffinityFlag(flag: boolean) {
    this.sessionAffinityFlag = flag;
    this.actionIsEnabledEvent.emit(this.sessionAffinityFlag !== this.oldSessionAffinityFlag);
  }
}
