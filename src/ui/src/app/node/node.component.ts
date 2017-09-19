import { Component, OnInit, ViewChild } from '@angular/core';

import { NodeDetailComponent } from './node-detail.component';
import { NodeService } from './node.service';
import { MessageService } from '../shared/message-service/message.service';


@Component({
  selector: 'node',
  templateUrl: 'node.component.html'
})
export class NodeComponent implements OnInit {

  nodes: {[key: string]:any} = [];

  @ViewChild(NodeDetailComponent) nodeDetailModal;

  constructor(
    private nodeService: NodeService,
    private messageService: MessageService
  ){}

  ngOnInit(): void {
    this.nodeService.getNodes()
      .then(nodes=>this.nodes = nodes)
      .catch(err=>this.messageService.dispatchError(err));
  }

  openNodeDetail(nodeName: string): void {
    this.nodeDetailModal.openNodeDetailModal();
  }

}