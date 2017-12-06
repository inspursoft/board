import { Component } from '@angular/core';

import { K8sService } from '../../service.k8s';

import { MessageService } from '../../../shared/message-service/message.service';

class NodeURL {
  url: string;
  identity: string;
  route: string;

  constructor(url: string, identity: string, route: string) {
    this.url = url;
    this.identity = identity;
    this.route = route;
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

  constructor(private k8sService: K8sService,
              private messageService: MessageService) {
  }

  openModal(serviceName: string, projectName: string, ownerName: string): void {
    this.getServiceDetail(serviceName, projectName, ownerName);
  }

  getServiceDetail(serviceName: string, projectName: string, ownerName: string): void {
    this.urlList = [];
    this.serviceDetail = "";
    this.serviceName = serviceName;
    this.k8sService.getServiceDetail(serviceName).then(res => {
      if (!res["details"]) {
        let arrNodePort = res["node_Port"] as Array<number>;
        this.k8sService.getNodesList().then(res => {
          let arrNode = res as Array<{node_name: string, node_ip: string, status: number}>;
          for (let i = 0; i < arrNode.length; i++) {
            let node = arrNode[i];
            if (node.status == 1) {
              let port = arrNodePort[Math.floor(Math.random() * arrNodePort.length)];
              let nodeInfo = {
                url: `http://${node.node_ip}:${port}`,
                identity: `${ownerName}_${projectName}_${serviceName}`,
                route: `http://${window.location.host}/deploy/${ownerName}/${projectName}/${serviceName}`
              };
              this.urlList.push(nodeInfo);
              this.k8sService.addServiceRoute(nodeInfo.url, nodeInfo.identity)
              break;
            }
          }
        });
      }
      this.serviceDetail = JSON.stringify(res);
      this.isOpenServiceDetail = true;
    }).catch(err => {
      this.isOpenServiceDetail = false;
      this.messageService.dispatchError(err);
    })
  }
}