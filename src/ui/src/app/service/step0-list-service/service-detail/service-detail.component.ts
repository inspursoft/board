import { Component } from '@angular/core';
import { K8sService } from '../../service.k8s';
import { AppInitService } from '../../../app.init.service';
import { Service } from "../../service";
import { Subject } from "rxjs/Subject";
import { Observable } from "rxjs/Observable";
import "rxjs/add/operator/do"

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

const K8S_HOSTNAME_KEY = 'kubernetes.io/hostname';
const YAML_TYPE_DEPLOYMENT = 'deployment';
const YAML_TYPE_SERVICE = 'service';

@Component({
  selector: 'service-detail',
  styleUrls: ["./service-detail.component.css"],
  templateUrl: './service-detail.component.html'
})
export class ServiceDetailComponent {
  _isOpenServiceDetail: boolean = false;
  boardHost: string;
  serviceDetail: Object;
  urlList: Array<NodeURL>;
  curService: Service;
  deploymentYamlFile: string = "";
  serviceYamlFile: string = "";
  closeNotification:Subject<any>;
  k8sHostName: string = "";
  dns = '';

  constructor(private appInitService: AppInitService,
              private k8sService: K8sService) {
    this.boardHost = this.appInitService.systemInfo.board_host;
    this.closeNotification = new Subject<any>();
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

  openModal(service: Service): Observable<any> {
    this.curService = service;
    this.dns =`${this.curService.service_name}.${this.curService.service_project_name}.svc${this.appInitService.systemInfo.dns_suffix}`;
    this.getDeploymentYamlFile()
      .subscribe(() => this.getServiceDetail(service.service_id, service.service_project_name, service.service_owner_name));
    return this.closeNotification.asObservable();
  }

  getServiceDetail(serviceId: number, projectName: string, ownerName: string): void {
    this.urlList = [];
    this.k8sService.getServiceDetail(serviceId).subscribe(res => {
      if (!res["details"]) {
        let arrNodePort = res["node_Port"] as Array<number>;
        this.k8sService.getNodesList({"ping": true}).subscribe(res => {
          let arrNode = res as Array<{node_name: string, node_ip: string, status: number}>;
          for(let n = 0; n < arrNodePort.length; n++) {
            for (let i = 0; i < arrNode.length; i++) {
              let node = arrNode[i];
              let host = this.k8sHostName && this.k8sHostName.length > 0 ? this.k8sHostName : node.node_ip;
              let port = arrNodePort[n];
              let nodeInfo = {
                url: `http://${host}:${port}`,
                identity: `${ownerName}_${projectName}_${this.curService.service_name}_${port}`,
                route: `http://${host}:${port}/deploy/${ownerName}/${projectName}/${this.curService.service_name}`
              };
              this.urlList.push(nodeInfo);
              this.k8sService.addServiceRoute(nodeInfo.url, nodeInfo.identity).subscribe();
            }
          }
        });
      }
      this.serviceDetail = res;
      this.isOpenServiceDetail = true;
    }, () => this.isOpenServiceDetail = false);
  }

  getDeploymentYamlFile(): Observable<string> {
    return this.k8sService.getServiceYamlFile(this.curService.service_project_name, this.curService.service_name, YAML_TYPE_DEPLOYMENT)
      .do((res: string) => {
        this.deploymentYamlFile = res;
        let arr: Array<string> = res.split(/[\n]/g);
        let k8sHost = arr.find(value => value.startsWith(K8S_HOSTNAME_KEY));
        if (k8sHost && k8sHost.length > 0) {
          this.k8sHostName = k8sHost.split(':')[1].trim();
        }
      }, () => this.isOpenServiceDetail = false);
  }

  getServiceYamlFile() {
    if (this.serviceYamlFile.length == 0) {
      this.k8sService.getServiceYamlFile(this.curService.service_project_name, this.curService.service_name, YAML_TYPE_SERVICE).subscribe(
        (res: string) => this.serviceYamlFile = res,
        () => this.isOpenServiceDetail = false)
    }
  }
}