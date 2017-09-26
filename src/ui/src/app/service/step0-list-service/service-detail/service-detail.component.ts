import { Component, OnInit } from '@angular/core';

import { K8sService } from '../../service.k8s';

import { MessageService } from '../../../shared/message-service/message.service';
import { MESSAGE_TARGET, BUTTON_STYLE, MESSAGE_TYPE } from '../../../shared/shared.const';


@Component({
  selector: 'service-detail',
  templateUrl: 'service-detail.component.html'
})
export class ServiceDetailComponent {
  
  isOpenServiceDetail = false;
  serviceDetail: string = "";

  urlList: Array<string>;
  
  serviceName: string;

  constructor(
    private k8sService: K8sService,
    private messageService: MessageService
  ) {}

  openModal(serviceName: string): void {
    this.isOpenServiceDetail = true;
    this.getServiceDetail(serviceName);
  }

  getServiceDetail(serviceName: string): void {
    this.urlList = [];
    this.serviceName = serviceName;
    this.k8sService.getServiceDetail(serviceName).then(res => {
      if (!res["details"]) {
        let arrNodePort = res["node_Port"] as Array<number>;
        this.k8sService.getNodesList().then(res => {
          let arrNode = res as Array<{node_name: string, node_ip: string, status: number}>;
          arrNode.forEach(node => {
            if (node.status == 1) {
              arrNodePort.forEach(port => {
                this.urlList.push(`http://${node.node_ip}:${port}`);
              });
            }
          });
        });
      }
      this.serviceDetail = JSON.stringify(res);
    }).catch(err => this.messageService.dispatchError(err))
  }
}