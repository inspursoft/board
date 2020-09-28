import { Component } from '@angular/core';
import { Observable, Subject } from 'rxjs';
import { tap } from 'rxjs/operators';
import { K8sService } from '../../service.k8s';
import { AppInitService } from '../../../shared.service/app-init.service';
import { Service, ServiceDetailInfo, ServiceType } from '../../service.types';

class NodeURL {
  constructor(public url: string, public identity: string, public route: string) {
  }
}

const K8S_HOSTNAME_KEY = 'kubernetes.io/hostname';
const YAML_TYPE_DEPLOYMENT = 'deployment';
const YAML_TYPE_STATEFUL_SET = 'statefulset';
const YAML_TYPE_SERVICE = 'service';

@Component({
  styleUrls: ['./service-detail.component.css'],
  templateUrl: './service-detail.component.html'
})
export class ServiceDetailComponent {
  isOpenServiceDetailValue = false;
  serviceDetail: ServiceDetailInfo;
  urlList: Array<NodeURL>;
  curService: Service;
  deploymentYamlFile = '';
  serviceYamlFile = '';
  closeNotification: Subject<any>;
  k8sHostName = '';
  dns = '';

  constructor(private appInitService: AppInitService,
              private k8sService: K8sService) {
    this.closeNotification = new Subject<any>();
    this.urlList = Array<NodeURL>();
  }

  get isOpenServiceDetail(): boolean {
    return this.isOpenServiceDetailValue;
  }

  set isOpenServiceDetail(value: boolean) {
    this.isOpenServiceDetailValue = value;
    if (!value) {
      this.closeNotification.next();
    }
  }

  openModal(service: Service): Observable<any> {
    this.curService = service;
    this.dns = `${this.curService.serviceName}.${this.curService.serviceProjectName}.svc${this.appInitService.systemInfo.dnsSuffix}`;
    this.getDeploymentYamlFile()
      .subscribe(() => this.getServiceDetail(service.serviceId, service.serviceProjectName, service.serviceOwnerName));
    this.isOpenServiceDetail = true;
    return this.closeNotification.asObservable();
  }

  getServiceDetail(serviceId: number, projectName: string, ownerName: string): void {
    if (this.curService.serviceType !== ServiceType.ServiceTypeStatefulSet &&
      this.curService.serviceType !== ServiceType.ServiceTypeClusterIP) {
      this.k8sService.getServiceDetail(serviceId).subscribe((serviceDetail: ServiceDetailInfo) => {
        if (serviceDetail.serviceContainers.length > 0) {
          serviceDetail.serviceContainers.forEach(container => {
            if (container.nodeIp !== '') {
              serviceDetail.nodePorts.forEach(port => {
                const nodeInfo = {
                  url: `http://${container.nodeIp}:${port}`,
                  identity: `${ownerName}_${projectName}_${this.curService.serviceName}_${port}`,
                  route: `http://${container.nodeIp}:${port}/deploy/${ownerName}/${projectName}/${this.curService.serviceName}`
                };
                this.urlList.push(nodeInfo);
                this.k8sService.addServiceRoute(nodeInfo.url, nodeInfo.identity).subscribe();
              });
            }
          });
          // this.k8sService.getNodesList({ping: true}).subscribe((nodeList: Array<ServiceNode>) => {
          //   for (const nodePort of serviceDetail.nodePorts) {
          //     for (const node of nodeList) {
          //       const host = this.k8sHostName && this.k8sHostName.length > 0 ? this.k8sHostName : node.nodeIp;
          //       const nodeInfo = {
          //         url: `http://${host}:${nodePort}`,
          //         identity: `${ownerName}_${projectName}_${this.curService.serviceName}_${nodePort}`,
          //         route: `http://${host}:${nodePort}/deploy/${ownerName}/${projectName}/${this.curService.serviceName}`
          //       };
          //       this.urlList.push(nodeInfo);
          //       this.k8sService.addServiceRoute(nodeInfo.url, nodeInfo.identity).subscribe();
          //     }
          //   }
          // });
        }
        this.serviceDetail = serviceDetail;
      }, () => this.isOpenServiceDetail = true);
    }
  }

  getDeploymentYamlFile(): Observable<string> {
    const yamlType = this.curService.serviceType === ServiceType.ServiceTypeStatefulSet ? YAML_TYPE_STATEFUL_SET : YAML_TYPE_DEPLOYMENT;
    return this.k8sService.getServiceYamlFile(this.curService.serviceProjectName, this.curService.serviceName, yamlType)
      .pipe(tap((yamlStr: string) => {
        this.deploymentYamlFile = yamlStr;
        const yamlArr: Array<string> = yamlStr.split(/[\n]/g);
        const k8sHost = yamlArr.find(value => value.startsWith(K8S_HOSTNAME_KEY));
        if (k8sHost && k8sHost.length > 0) {
          this.k8sHostName = k8sHost.split(':')[1].trim();
        }
      }, () => this.isOpenServiceDetail = false));
  }

  getServiceYamlFile() {
    if (this.serviceYamlFile.length === 0) {
      this.k8sService.getServiceYamlFile(this.curService.serviceProjectName, this.curService.serviceName, YAML_TYPE_SERVICE).subscribe(
        (res: string) => this.serviceYamlFile = res,
        () => this.isOpenServiceDetail = false);
    }
  }
}
