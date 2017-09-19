import { Component } from '@angular/core';

import { MessageService } from '../shared/message-service/message.service';

import { AppInitService } from '../app.init.service';
import { NodeService } from './node.service';

@Component({
  selector: 'node-detail',
  templateUrl: './node-detail.component.html'
})
export class NodeDetailComponent {
  
  nodeDetailOpened: boolean;
  node: {[key: string]: any} = {};

  constructor(
    private appInitService: AppInitService,
    private nodeService: NodeService,
    private messageService: MessageService
  ){}

  openNodeDetailModal(nodeName: string): void {
    this.nodeDetailOpened = true;
    this.nodeService.getNodeByName(nodeName)
      .then(node=>this.node = node)
      .catch(err=>this.messageService.dispatchError(err));
  }

}