import { Component, ComponentFactoryResolver, Input, OnInit, ViewContainerRef } from '@angular/core';
import { NodeActionsType, ResponseArrayNode } from '../resource.types';
import { ResourceService } from '../services/resource.service';
import { HttpErrorResponse } from "@angular/common/http";
import { NodeAddRemoveComponent } from "../node-add-remove/node-add-remove.component";
import { MessageService } from "../../shared/message/message.service";

@Component({
  selector: 'app-compute',
  templateUrl: './compute.component.html',
  styleUrls: ['./compute.component.css']
})
export class ComputeComponent implements OnInit {
  nodes: ResponseArrayNode;
  nodeLoadingInfo = 'Loading...';
  @Input('viewContainer') viewContainer: ViewContainerRef;

  constructor(private resolver: ComponentFactoryResolver,
              private messageService: MessageService,
              private resourceService: ResourceService) {
    this.nodes = new ResponseArrayNode({});
  }

  ngOnInit() {

  }

  fetchNodes() {
    this.resourceService.getNodeList().subscribe(
      res => this.nodes = res,
      (err: HttpErrorResponse) => this.nodeLoadingInfo = err.message,
      () => this.nodeLoadingInfo = 'Node list is empty.'
    );
  }

  addRemoveNode(actionType: NodeActionsType, nodeIp: string) {
    const nodeFactory = this.resolver.resolveComponentFactory(NodeAddRemoveComponent);
    const nodeComponentRef = this.viewContainer.createComponent(nodeFactory);
    nodeComponentRef.instance.actionType = actionType;
    nodeComponentRef.instance.nodeIp = nodeIp;
    nodeComponentRef.instance.openModal().subscribe(() =>
      this.viewContainer.remove(this.viewContainer.indexOf(nodeComponentRef.hostView))
    );
    nodeComponentRef.instance.successNotification.subscribe(() => this.fetchNodes());
  }

}
