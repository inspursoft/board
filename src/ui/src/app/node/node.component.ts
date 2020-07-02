import { Component, ViewChild } from '@angular/core';
import { NodeListComponent } from './node-list/node-list.component';
import { NodeStatus } from './node.types';

@Component({
  templateUrl: './node.component.html'
})
export class NodeComponent {
  @ViewChild(NodeListComponent) nodeListComponent;
  isCreatingNode = false;
  nodeList: Array<NodeStatus>;

  constructor() {
    this.nodeList = new Array<NodeStatus>();
  }

  setCreatingNode(nodeList: Array<NodeStatus>) {
    this.nodeList = nodeList;
    this.isCreatingNode = true;
  }
}
