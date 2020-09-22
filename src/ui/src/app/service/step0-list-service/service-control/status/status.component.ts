import { Component, EventEmitter, Input, OnDestroy, OnInit, Output } from '@angular/core';
import { Subject } from 'rxjs';
import { SERVICE_STATUS } from '../../../../shared/shared.const';
import { K8sService } from '../../../service.k8s';
import { IScaleInfo } from '../service-control.component';
import { Service } from '../../../service.types';

@Component({
  selector: 'app-status',
  templateUrl: './status.component.html',
  styleUrls: ['./status.component.css']
})
export class StatusComponent implements OnInit, OnDestroy {
  @Output() errorEvent: EventEmitter<any>;
  @Input() service: Service;
  @Input() scaleInfo: IScaleInfo = {
    desired_instance: 0,
    available_instance: 0
  };
  onDestroy: Subject<any>;

  constructor(private k8sService: K8sService) {
    this.onDestroy = new Subject<any>();
    this.errorEvent = new EventEmitter<any>();
  }

  ngOnDestroy() {
    this.onDestroy.next();
  }

  ngOnInit() {
    this.onDestroy.subscribe(() => this.refreshScaleInfo());
    this.refreshScaleInfo();
  }

  refreshScaleInfo() {
    this.k8sService.getServiceScaleInfo(this.service.serviceId).subscribe(
      (scaleInfo: IScaleInfo) => {
        this.scaleInfo = scaleInfo;
      },
      (err) => this.errorEvent.emit(err)
    );
  }

  getStatusClass(status: SERVICE_STATUS) {
    return {
      running: status === SERVICE_STATUS.RUNNING,
      stopped: status === SERVICE_STATUS.STOPPED,
      warning: status === SERVICE_STATUS.WARNING
    };
  }

  getServiceStatus(status: SERVICE_STATUS): string {
    switch (status) {
      case SERVICE_STATUS.PREPARING:
        return 'SERVICE.STATUS_PREPARING';
      case SERVICE_STATUS.RUNNING:
        return 'SERVICE.STATUS_RUNNING';
      case SERVICE_STATUS.STOPPED:
        return 'SERVICE.STATUS_STOPPED';
      case SERVICE_STATUS.WARNING:
        return 'SERVICE.STATUS_UNCOMPLETED';
    }
  }

  get reason(): string {
    let result: string = this.service.serviceComment;
    if (result.toLowerCase().startsWith('reason:')) {
      result = result.slice(7);
    }
    return result;
  }

}
