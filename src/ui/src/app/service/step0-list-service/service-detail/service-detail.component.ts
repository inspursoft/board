import { Component } from '@angular/core';
import { K8sService } from '../../service.k8s';
import { MessageService } from '../../../shared/message-service/message.service';
import { AppInitService } from '../../../app.init.service';
import { Service } from "../../service";

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
  styleUrls: ["./service-detail.component.css"],
  templateUrl: './service-detail.component.html'
})
export class ServiceDetailComponent {
  _isOpenServiceDetail: boolean = false;
  boardHost: string;
  serviceDetail: string = "";
  urlList: Array<NodeURL>;
  curService: Service;
  deploymentYamlFile: string = "";
  deploymentYamlWIP: boolean = false;
  isShowDeploymentYaml: boolean = false;
  serviceYamlFile: string = "";
  serviceYamlWIP: boolean = false;
  isShowServiceYaml: boolean = false;

  constructor(private appInitService: AppInitService,
              private k8sService: K8sService,
              private messageService: MessageService) {
    this.boardHost = this.appInitService.systemInfo['board_host'];
  }

  get isOpenServiceDetail(): boolean {
    return this._isOpenServiceDetail;
  }

  set isOpenServiceDetail(value: boolean) {
    this._isOpenServiceDetail = value;
    this.isShowServiceYaml = false;
    this.isShowDeploymentYaml = false;
  }

  openModal(s: Service): void {
    this.curService = s;
    this.getServiceDetail(s.service_id, s.service_project_name, s.service_owner);
  }

  getServiceDetail(serviceId: number, projectName: string, ownerName: string): void {
    this.urlList = [];
    this.serviceDetail = "";
    this.k8sService.getServiceDetail(serviceId).then(res => {
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
                identity: `${ownerName}_${projectName}_${this.curService.service_name}`,
                route: `http://${this.boardHost}/deploy/${ownerName}/${projectName}/${this.curService.service_name}`
              };
              this.urlList.push(nodeInfo);
              this.k8sService.addServiceRoute(nodeInfo.url, nodeInfo.identity);
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

  getDeploymentYamlFile() {
    this.isShowDeploymentYaml = !this.isShowDeploymentYaml;
    if (this.isShowDeploymentYaml) {
      this.deploymentYamlWIP = true;
      this.k8sService.getServiceYamlFile(this.curService.service_project_name, this.curService.service_name, "deployment")
        .then((res: string) => {
          this.deploymentYamlWIP = false;
          this.deploymentYamlFile = res;
        })
        .catch(err => {
          this.deploymentYamlWIP = false;
          this.isOpenServiceDetail = false;
          this.messageService.dispatchError(err);
        })
    }
  }

  getServiceYamlFile() {
    this.isShowServiceYaml = !this.isShowServiceYaml;
    if (this.isShowServiceYaml) {
      this.serviceYamlWIP = true;
      this.k8sService.getServiceYamlFile(this.curService.service_project_name, this.curService.service_name, "service")
        .then((res: string) => {
          this.serviceYamlWIP = false;
          this.serviceYamlFile = res;
        })
        .catch(err => {
          this.serviceYamlWIP = false;
          this.isOpenServiceDetail = false;
          this.messageService.dispatchError(err);
        })
    }
  }
}