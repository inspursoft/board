import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { K8sService } from '../../../service.k8s';
import { Service, ServiceType } from '../../../service.types';

@Component({
  selector: 'app-locate',
  templateUrl: './locate.component.html',
  styleUrls: ['./locate.component.css']
})
export class LocateComponent implements OnInit {
  @Input() isActionInWIP: boolean;
  @Input() service: Service;
  @Output() messageEvent: EventEmitter<string>;
  @Output() errorEvent: EventEmitter<any>;
  @Output() actionIsEnabledEvent: EventEmitter<boolean>;
  dropdownDefaultText = 'SERVICE.SERVICE_CONTROL_LOCATE_SELECT';
  nodeSelectorList: Array<string>;
  nodeSelector = '';

  constructor(private k8sService: K8sService) {
    this.messageEvent = new EventEmitter<string>();
    this.errorEvent = new EventEmitter<any>();
    this.nodeSelectorList = Array<string>();
    this.actionIsEnabledEvent = new EventEmitter<boolean>();
  }

  ngOnInit() {
    this.actionIsEnabledEvent.emit(false);
    if (this.service.serviceType === ServiceType.ServiceTypeEdgeComputing) {
      this.k8sService.getEdgeNodes().subscribe(
        (res: Array<{ description: string }>) => res.forEach(node => this.nodeSelectorList.push(node.description)),
        (err) => this.errorEvent.emit(err));
    } else {
      this.k8sService.getNodeSelectors().subscribe(
        (res: Array<{ name: string, status: number }>) => res.forEach(value => this.nodeSelectorList.push(value.name)),
        (err) => this.errorEvent.emit(err)
      );
    }
    this.k8sService.getLocate(this.service.serviceProjectName, this.service.serviceName).subscribe(
      res => {
        if (res && res !== '') {
          this.dropdownDefaultText = res;
        }
      }, (err) => this.errorEvent.emit(err));
  }

  actionExecute() {
    this.k8sService.setLocate(this.nodeSelector, this.service.serviceProjectName, this.service.serviceName).subscribe(
      () => this.messageEvent.emit('SERVICE.SERVICE_CONTROL_LOCATE_SUCCESSFUL'),
      (err) => this.errorEvent.emit(err));
  }

  setNodeSelector(selector: string) {
    this.nodeSelector = selector;
    this.actionIsEnabledEvent.emit(this.nodeSelector !== '' && this.nodeSelector !== this.dropdownDefaultText);
  }
}
