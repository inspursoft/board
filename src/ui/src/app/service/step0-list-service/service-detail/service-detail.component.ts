import { Component, OnInit, Input } from '@angular/core';

import { K8sService } from '../../service.k8s';

import { MessageService } from '../../../shared/message-service/message.service';
import { MESSAGE_TARGET, BUTTON_STYLE, MESSAGE_TYPE } from '../../../shared/shared.const';

class NodeURL {
  url: string;
  description: string;
  constructor(url: string, description: string) {
    this.url = url;
    this.description = description;
  }
}

@Component({
  selector: 'service-detail',
  templateUrl: 'service-detail.component.html'
})
export class ServiceDetailComponent {
  
  isOpenServiceDetail = false;
  serviceDetail: string = "";

  urlList: Array<NodeURL>;
  
  serviceName: string;

  constructor(
    private k8sService: K8sService,
    private messageService: MessageService
  ) {}

  openModal(serviceName: string, projectName: string, ownerName: string): void {
    this.isOpenServiceDetail = true;
    this.getServiceDetail(serviceName, projectName, ownerName);
  }

  getServiceDetail(serviceName: string, projectName: string, ownerName: string): void {
    this.urlList = [];
    this.serviceName = serviceName;
    this.k8sService.getServiceDetail(serviceName).then(res => {
      if (!res["details"]) {
        let arrNodePort = res["node_Port"] as Array<number>;
        this.k8sService.getNodesList().then(res => {
          let arrNode = res as Array<{node_name: string, node_ip: string, status: number}>;
          for(let i = 0; i < arrNode.length; i++){
            let node = arrNode[i];
            if (node.status == 1) {
              let port = arrNodePort[Math.floor(Math.random() * arrNodePort.length)];
              this.urlList.push(new NodeURL(`http://${node.node_ip}:${port}`, `http://${window.location.host}/${ownerName}/${projectName}/${serviceName}`));
              break;
            }
          }
        });
      }
      this.serviceDetail = JSON.stringify(res);
    }).catch(err => this.messageService.dispatchError(err))
  }
}