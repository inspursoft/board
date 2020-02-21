import { Component, ComponentFactoryResolver, OnInit, ViewContainerRef } from '@angular/core';
import { NodeActionsType, NodeList, NodeListType } from '../../resource.types';
import { ResourceService } from '../../services/resource.service';
import { MessageService } from "../../../shared/message/message.service";
import { Message, ReturnStatus } from "../../../shared/message/message.types";
import { NodeDetailComponent } from "../node-detail/node-detail.component";

@Component({
  selector: 'app-node-list',
  templateUrl: './node-list.component.html',
  styleUrls: ['./node-list.component.css']
})
export class NodeListComponent implements OnInit {
  nodeLists: NodeList;

  constructor(private resourceService: ResourceService,
              private messageService: MessageService,
              private view: ViewContainerRef,
              private resolver: ComponentFactoryResolver) {
    this.nodeLists = new NodeList({});
  }

  ngOnInit() {
    this.getNodeList();
  }

  getNodeList() {
    this.resourceService.getNodeList().subscribe((res: NodeList) => this.nodeLists = res);
  }

  addNode() {
    this.createNodeDetail('', NodeActionsType.Add);
  }

  deleteNode(nodeIp: string) {
    this.messageService.showDeleteDialog('确定要删除该节点么？', '删除节点').subscribe(
      (res: Message) => {
        if (res.returnStatus === ReturnStatus.rsConfirm) {
          this.createNodeDetail(nodeIp, NodeActionsType.Remove);
        }
      }
    );
  }

  showLog(node: NodeListType) {
    this.resourceService.getNodeLog(`${node.Ip}@${node.CreationTime}.txt`).subscribe();
  }

  createNodeDetail(nodeIp: string, action: NodeActionsType) {
    const factory = this.resolver.resolveComponentFactory(NodeDetailComponent);
    const detailRef = this.view.createComponent(factory);
    detailRef.instance.nodeIp = nodeIp;
    detailRef.instance.actionType = action;
    detailRef.instance.openModal().subscribe(
      () => this.view.remove(this.view.indexOf(detailRef.hostView))
    );
  }
}
