import { ChangeDetectorRef, Component } from '@angular/core';
import { K8sService } from '../../service.k8s';
import { MessageService } from '../../../shared/message-service/message.service';
import { AppInitService } from '../../../app.init.service';
import { Service } from "../../service";
import { Subject } from "rxjs/Subject";
import { Observable } from "rxjs/Observable";
import { HttpErrorResponse } from "@angular/common/http";

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
  serviceDetail: Object = {};
  urlList: Array<NodeURL>;
  curService: Service;
  deploymentYamlFile: string = "";
  serviceYamlFile: string = "";
  closeNotification:Subject<any>;

  constructor(private appInitService: AppInitService,
              private k8sService: K8sService,
              private change:ChangeDetectorRef,
              private messageService: MessageService) {
    this.boardHost = this.appInitService.systemInfo['board_host'];
    this.closeNotification = new Subject<any>();
    this.change.detach();
  }

  get isOpenServiceDetail(): boolean {
    return this._isOpenServiceDetail;
  }

  set isOpenServiceDetail(value: boolean) {
    this._isOpenServiceDetail = value;
    if (!value){
      this.closeNotification.next();
    }
  }

  openModal(s: Service): Observable<any> {
    this.curService = s;
    this.getServiceDetail(s.service_id, s.service_project_name, s.service_owner_name);
    return this.closeNotification.asObservable();
  }

  getServiceDetail(serviceId: number, projectName: string, ownerName: string): void {
    this.urlList = [];
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
              this.k8sService.addServiceRoute(nodeInfo.url, nodeInfo.identity).then(() => {
              });
              break;
            }
          }
        });
      }
      this.serviceDetail = res;
      this.change.reattach();
      this.isOpenServiceDetail = true;
    }).catch(err => {
      this.isOpenServiceDetail = false;
      this.messageService.dispatchError(err);
    })
  }

  getDeploymentYamlFile() {
    if (this.deploymentYamlFile.length == 0) {
      this.change.detach();
      this.k8sService.getServiceYamlFile(this.curService.service_project_name, this.curService.service_name, "deployment")
        .then((res: string) => {
          this.deploymentYamlFile = res;
          this.change.reattach();
        })
        .catch((err: HttpErrorResponse) => {
          this.isOpenServiceDetail = false;
          this.messageService.dispatchError(err);
        })
    }
  }

  getServiceYamlFile() {
    if (this.serviceYamlFile.length == 0) {
      this.change.detach();
      this.k8sService.getServiceYamlFile(this.curService.service_project_name, this.curService.service_name, "service")
        .then((res: string) => {
          this.serviceYamlFile = res;
          this.change.reattach();
        })
        .catch((err: HttpErrorResponse) => {
          this.isOpenServiceDetail = false;
          this.messageService.dispatchError(err);
        })
    }
  }
}