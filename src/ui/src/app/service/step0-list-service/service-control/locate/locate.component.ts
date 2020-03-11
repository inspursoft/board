import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Service } from "../../../service";
import { K8sService } from "../../../service.k8s";

@Component({
  selector: 'locate',
  templateUrl: './locate.component.html',
  styleUrls: ['./locate.component.css']
})
export class LocateComponent implements OnInit {
  @Input('isActionInWIP') isActionInWIP: boolean;
  @Input('service') service: Service;
  @Output("onMessage") onMessage: EventEmitter<string>;
  @Output("onError") onError: EventEmitter<any>;
  @Output("onActionIsEnabled") onActionIsEnabled: EventEmitter<boolean>;
  dropdownDefaultText: string = "SERVICE.SERVICE_CONTROL_LOCATE_SELECT";
  nodeSelectorList: Array<string>;
  nodeSelector: string = "";

  constructor(private k8sService: K8sService) {
    this.onMessage = new EventEmitter<string>();
    this.onError = new EventEmitter<any>();
    this.nodeSelectorList = Array<string>();
    this.onActionIsEnabled = new EventEmitter<boolean>();
  }

  ngOnInit() {
    this.onActionIsEnabled.emit(false);
    this.k8sService.getNodeSelectors().subscribe((res: Array<{name: string, status: number}>) =>
      res.forEach(value => this.nodeSelectorList.push(value.name))
    );
    this.k8sService.getLocate(this.service.service_project_name, this.service.service_name)
      .subscribe(res => {
        if (res && res != ""){
          this.dropdownDefaultText = res
        }
      });
  }

  actionExecute() {
    this.k8sService.setLocate(this.nodeSelector, this.service.service_project_name, this.service.service_name).subscribe(
      () => this.onMessage.emit('SERVICE.SERVICE_CONTROL_LOCATE_SUCCESSFUL'),
      (err) => this.onError.emit(err));
  }

  setNodeSelector(selector: string){
    this.nodeSelector = selector;
    this.onActionIsEnabled.emit(this.nodeSelector != "" && this.nodeSelector != this.dropdownDefaultText);
  }
}
