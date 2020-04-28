import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Service } from '../../../service';
import { K8sService } from '../../../service.k8s';

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
    this.k8sService.getNodeSelectors().subscribe(
      (res: Array<{ name: string, status: number }>) =>
        res.forEach(value => this.nodeSelectorList.push(value.name))
    );
    this.k8sService.getLocate(this.service.service_project_name, this.service.service_name).subscribe(
      res => {
        if (res && res !== '') {
          this.dropdownDefaultText = res;
        }
      });
  }

  actionExecute() {
    this.k8sService.setLocate(this.nodeSelector, this.service.service_project_name, this.service.service_name).subscribe(
      () => this.messageEvent.emit('SERVICE.SERVICE_CONTROL_LOCATE_SUCCESSFUL'),
      (err) => this.errorEvent.emit(err));
  }

  setNodeSelector(selector: string) {
    this.nodeSelector = selector;
    this.actionIsEnabledEvent.emit(this.nodeSelector !== '' && this.nodeSelector !== this.dropdownDefaultText);
  }
}
