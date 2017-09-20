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


  toPercentage(num: number) {
    return Math.round(num * 100) / 100 + '%';
  }

  toGigaBytes(num: string, baseUnit?: string) {
    let denominator = 1024 * 1024 * 1024;
    if(baseUnit === 'KiB') {
      denominator = 1024 * 1024;
    }
    return Math.round(Number.parseInt(num) / denominator) + 'GB';
  }
}