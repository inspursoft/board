import { Component, ComponentFactoryResolver, OnDestroy, OnInit, ViewContainerRef } from '@angular/core';
import { NodeActionsType, NodeList, NodeListType, NodeLog } from '../../resource.types';
import { ResourceService } from '../../services/resource.service';
import { MessageService } from '../../../shared/message/message.service';
import { Message, ReturnStatus } from '../../../shared/message/message.types';
import { NodeDetailComponent } from '../node-detail/node-detail.component';
import { interval, Subscription } from 'rxjs';

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

  deleteNode(nodeIp: string) {
    this.messageService.showDeleteDialog('确定要删除该节点么？', '删除节点').subscribe(
      (res: Message) => {
        if (res.returnStatus === ReturnStatus.rsConfirm) {
          const logInfo = new NodeLog({});
          logInfo.ip = nodeIp;
          this.createNodeDetail(logInfo, NodeActionsType.Remove);
        }
      }
    );
  }

  showLog(node: NodeListType) {
    const logInfo = new NodeLog({});
    logInfo.ip = node.ip;
    logInfo.creationTime = node.creationTime;
    this.createNodeDetail(logInfo, NodeActionsType.Log);
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
