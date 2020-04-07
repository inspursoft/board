import { Component, ComponentFactoryResolver, OnDestroy, OnInit, ViewContainerRef } from '@angular/core';
import { interval, Subscription } from 'rxjs';
import { TranslateService } from '@ngx-translate/core';
import { NodeActionsType, NodeList, NodeListType, NodeLog } from '../../resource.types';
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

  deleteNode(node: NodeListType) {
    if (!node.isMaster) {
      this.translateService.get(['Node.Node_List_Remove_Ask', 'Node.Node_Logs_Stop_Ask'])
        .subscribe(translate => {
          const ask = Reflect.get(translate, 'Node.Node_List_Remove_Ask');
          const title = Reflect.get(translate, 'Node.Node_List_Remove_Node');
          this.messageService.showDeleteDialog(ask, title).subscribe(
            (res: Message) => {
              if (res.returnStatus === ReturnStatus.rsConfirm) {
                const logInfo = new NodeLog({});
                logInfo.ip = node.ip;
                this.createNodeDetail(logInfo, NodeActionsType.Remove);
              }
            }
          );
        });
    }
  }

  showLog(node: NodeListType) {
    if (node.origin === 1) {
      const logInfo = new NodeLog({});
      logInfo.ip = node.ip;
      logInfo.creationTime = node.creationTime;
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
