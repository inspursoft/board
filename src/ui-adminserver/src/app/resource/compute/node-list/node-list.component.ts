import { Component, ComponentFactoryResolver, OnDestroy, OnInit, ViewContainerRef } from '@angular/core';
import { interval, Subscription } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import { NodeControlStatus, NodeList, NodeListType, NodeLog } from '../../resource.types';
import { ResourceService } from '../../services/resource.service';
import { MessageService } from '../../../shared/message/message.service';
import { NodeCreateComponent } from '../node-create/node-create.component';
import { NodeRemoveComponent } from '../node-remove/node-remove.component';
import { NodeLogComponent } from '../node-log/node-log.component';

@Component({
  selector: 'app-node-list',
  templateUrl: './node-list.component.html',
  styleUrls: ['./node-list.component.css']
})
export class NodeListComponent implements OnInit, OnDestroy {
  nodeLists: NodeList;
  subscriptionUpdate: Subscription;

  constructor(private resourceService: ResourceService,
              private messageService: MessageService,
              private view: ViewContainerRef,
              private translateService: TranslateService,
              private resolver: ComponentFactoryResolver) {
    this.nodeLists = new NodeList({});
  }

  ngOnInit() {
    this.subscriptionUpdate = interval(5000).subscribe(() => this.getNodeList());
    this.getNodeList();
  }

  ngOnDestroy(): void {
    this.subscriptionUpdate.unsubscribe();
    delete this.subscriptionUpdate;
  }

  getNodeList() {
    this.resourceService.getNodeList().subscribe((res: NodeList) => this.nodeLists = res);
  }

  addNodeAction() {
    const factory = this.resolver.resolveComponentFactory(NodeCreateComponent);
    const ref = this.view.createComponent(factory);
    ref.instance.openModal().subscribe(
      () => this.view.remove(this.view.indexOf(ref.hostView))
    );
  }

  removeNodeAction(nodeIp: string) {
    const factory = this.resolver.resolveComponentFactory(NodeRemoveComponent);
    const ref = this.view.createComponent(factory);
    ref.instance.postData.nodeIp = nodeIp;
    ref.instance.openModal().subscribe(
      () => this.view.remove(this.view.indexOf(ref.hostView))
    );
  }

  showNodeLogAction(logInfo: NodeLog): void {
    const factory = this.resolver.resolveComponentFactory(NodeLogComponent);
    const ref = this.view.createComponent(factory);
    ref.instance.logInfo = logInfo;
    ref.instance.openModal().subscribe(
      () => this.view.remove(this.view.indexOf(ref.hostView))
    );
  }

  removeNode(node: NodeListType) {
    this.resourceService.getNodeControlStatus(node.nodeName).subscribe(
      (res: NodeControlStatus) => {
        if (res.nodeUnschedulable) {
          if (res.nodeDeletable) {
            this.removeNodeAction(node.ip);
          } else {
            this.translateService.get(['Node.Node_Detail_Remove', 'Node.Node_Logs_Can_Not_Remove']).subscribe(
              translate => {
                const title = Reflect.get(translate, 'Node.Node_Detail_Remove');
                const msg = Reflect.get(translate, 'Node.Node_Logs_Can_Not_Remove');
                this.messageService.showDialog(msg, {title});
              }
            );
          }
        } else {
          this.translateService.get(['Node.Node_Detail_Remove', 'Node.Node_Logs_Can_Not_Remove_1']).subscribe(
            translate => {
              const title = Reflect.get(translate, 'Node.Node_Detail_Remove');
              const msg = Reflect.get(translate, 'Node.Node_Logs_Can_Not_Remove_1');
              this.messageService.showDialog(msg, {title});
            }
          );
        }
      }
    );
  }

  showLog(node: NodeListType) {
    if (node.origin === 1) {
      const logInfo = new NodeLog({});
      logInfo.ip = node.ip;
      logInfo.creationTime = node.logTime;
      this.showNodeLogAction(logInfo);
    }
  }

  getStatus(status: number): string {
    switch (status) {
      case 1:
        return 'Node.Node_List_Status_Schedulable';
      case 2:
        return 'Node.Node_List_Status_Unschedulable';
      case 3:
        return 'Node.Node_List_Status_Unknown';
    }
  }
}
