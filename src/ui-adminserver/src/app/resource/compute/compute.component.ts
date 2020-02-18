import { Component, ComponentFactoryResolver, OnInit, ViewContainerRef } from '@angular/core';
import { ResponseArrayNode } from '../resource.types';
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

  constructor(private resolver: ComponentFactoryResolver,
              private messageService: MessageService,
              private resourceService: ResourceService) {
    this.nodes = new ResponseArrayNode({});
  }

  ngOnInit() {
    this.resourceService.getNodeList().subscribe(
      res => this.nodes = res,
      (err: HttpErrorResponse) => this.nodeLoadingInfo = err.message
    );
  }

  fetchNodes() {

  }

  addNode() {
    const nodeFactory = this.resolver.resolveComponentFactory(NodeAddRemoveComponent);
    const nodeComponentRef = this.messageService.dialogView.createComponent(nodeFactory);
    nodeComponentRef.instance.openModal().subscribe(() =>
      this.messageService.dialogView.remove(this.messageService.dialogView.indexOf(nodeComponentRef.hostView))
    );
    nodeComponentRef.instance.successNotification.subscribe(() => this.fetchNodes());
  }

}
