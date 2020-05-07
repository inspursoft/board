import { Component, ComponentFactoryResolver, OnDestroy, OnInit, ViewContainerRef } from '@angular/core';
import { interval, Subscription } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import { NodeActionsType, NodeControlStatus, NodeList, NodeListType, NodeLog } from '../../resource.types';
import { ResourceService } from '../../services/resource.service';
import { MessageService } from '../../../shared/message/message.service';
import { Message, ReturnStatus } from '../../../shared/message/message.types';
import { NodeDetailComponent } from '../node-detail/node-detail.component';

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
    this.subscriptionUpdate = interval(3000).subscribe(() => this.getNodeList());
    this.getNodeList();
  }

  ngOnDestroy(): void {
    this.subscriptionUpdate.unsubscribe();
    delete this.subscriptionUpdate;
  }

  getNodeList() {
    this.resourceService.getNodeList().subscribe((res: NodeList) => this.nodeLists = res);
  }

  addNode() {
    const logInfo = new NodeLog({});
    this.createNodeDetail(logInfo, NodeActionsType.Add);
  }

  testDeleteNode(node: NodeListType) {
    this.resourceService.deleteNode(node.nodeName).subscribe(
      () => this.messageService.showAlert('删除成功～！')
    );
  }

  deleteNode(node: NodeListType) {
    this.resourceService.getNodeControlStatus(node.nodeName).subscribe(
      (res: NodeControlStatus) => {
        if (res.nodeDeletable) {
          const logInfo = new NodeLog({});
          logInfo.ip = node.ip;
          this.createNodeDetail(logInfo, NodeActionsType.Remove);
        } else {
          this.translateService.get(['Node.Node_Detail_Remove', 'Node.Node_Logs_Can_Not_Remove']).subscribe(
            translate => {
              const title = Reflect.get(translate, 'Node.Node_Detail_Remove');
              const msg = Reflect.get(translate, 'Node.Node_Logs_Can_Not_Remove');
              this.messageService.showDialog(msg, {title});
            }
          );
        }
      });
  }

  showLog(node: NodeListType) {
    if (node.origin === 1) {
      const logInfo = new NodeLog({});
      logInfo.ip = node.ip;
      logInfo.creationTime = node.logTime;
      this.createNodeDetail(logInfo, NodeActionsType.Log);
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

  createNodeDetail(logInfo: NodeLog, action: NodeActionsType) {
    const factory = this.resolver.resolveComponentFactory(NodeDetailComponent);
    const detailRef = this.view.createComponent(factory);
    detailRef.instance.logInfo = logInfo;
    detailRef.instance.actionType = action;
    detailRef.instance.openModal().subscribe(
      () => this.view.remove(this.view.indexOf(detailRef.hostView))
    );
  }
}
