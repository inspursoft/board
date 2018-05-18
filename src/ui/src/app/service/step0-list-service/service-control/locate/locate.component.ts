import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { Message } from "../../../../shared/message-service/message";
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
  @Output("onMessage") onMessage: EventEmitter<Message>;
  @Output("onError") onError: EventEmitter<any>;
  @Output("onActionIsEnabled") onActionIsEnabled: EventEmitter<boolean>;
  dropdownDefaultText: string = "";
  nodeSelectorList: Array<string>;
  nodeSelector: string = "";

  constructor(private k8sService: K8sService) {
    this.onMessage = new EventEmitter<Message>();
    this.onError = new EventEmitter<any>();
    this.nodeSelectorList = Array<string>();
    this.onActionIsEnabled = new EventEmitter<boolean>();
  }

  ngOnInit() {
    this.onActionIsEnabled.emit(false);
    this.k8sService.getNodeSelectors().subscribe(res => this.nodeSelectorList = res);
    this.k8sService.getLocate(this.service.service_project_name, this.service.service_name)
      .subscribe(res=>this.dropdownDefaultText = res);
  }

  actionExecute() {
     this.k8sService.setLocate(this.nodeSelector, this.service.service_project_name, this.service.service_name)
       .subscribe(() => {
         let msg: Message = new Message();
         msg.message = "SERVICE.SERVICE_CONTROL_LOCATE_SUCCESSFUL";
         msg.params = [this.service.service_name];
         this.onMessage.emit(msg);
       }, (err) => this.onError.emit(err));
  }

  setNodeSelector(selector: string){
    this.nodeSelector = selector;
    this.onActionIsEnabled.emit(this.nodeSelector != "" && this.nodeSelector != this.dropdownDefaultText);
  }
}
