import { Component, OnInit, OnDestroy, ViewChild } from '@angular/core';

import { Subscription } from 'rxjs/Subscription';

import { NodeDetailComponent } from './node-detail.component';
import { NodeService } from './node.service';

import { MESSAGE_TARGET, BUTTON_STYLE } from '../shared/shared.const';
import { Message } from '../shared/message-service/message';
import { MessageService } from '../shared/message-service/message.service';


class NodeStatus {
  name: string;
  status: boolean;

  constructor(name: string, status: boolean) {
    this.name = name;
    this.status = status;
  }
}

@Component({
  selector: 'node',
  templateUrl: 'node.component.html'
})
export class NodeComponent implements OnInit, OnDestroy {

  nodes: {[key: string]:any} = [];
  
  @ViewChild(NodeDetailComponent) nodeDetailModal;

  _subscription: Subscription;

  constructor(
    private nodeService: NodeService,
    private messageService: MessageService
  ){
    this._subscription = this.messageService.messageConfirmed$.subscribe(m=>{
      let confirmationMessage = <Message>m;
      if(confirmationMessage) {
        let nodeStatus = <NodeStatus>confirmationMessage.data;
        let m: Message = new Message();
        this.nodeService
          .toggleNodeStatus(nodeStatus.name, nodeStatus.status)
          .then(res=>{
            m.message = 'NODE.SUCCESSFUL_TOGGLE';
            this.messageService.inlineAlertMessage(m);
            this.retrieve();
          })
          .catch(err=>{
            m.message = 'NODE.FAILED_TO_TOGGLE';
            this.messageService.inlineAlertMessage(m);
          });
      }
    });
  }

  ngOnInit(): void {
    this.retrieve();
  }

  ngOnDestroy(): void {
    if(this._subscription) {
      this._subscription.unsubscribe();
    }
  }

  retrieve(): void {
    this.nodeService.getNodes()
    .then(nodes=>this.nodes = nodes)
    .catch(err=>this.messageService.dispatchError(err));
  }

  getStatus(status: number): string {
    switch(status) {
    case 1:
      return 'NODE.STATUS_RUNNING';
    case 2:
      return 'NODE.STATUS_UNSCHEDULABLE';
    case 3:
      return 'NODE.STATUS_UNKNOWN';
    }
  }

  openNodeDetail(nodeName: string): void {
    this.nodeDetailModal.openNodeDetailModal(nodeName);
  }

  confirmToToggleNodeStatus(nodeName: string, status: number): void {
    let announceMessage = new Message();
    announceMessage.title = 'NODE.TOGGLE_NODE';
    announceMessage.message = 'NODE.CONFIRM_TO_TOGGLE_NODE';
    announceMessage.params = [ nodeName ];
    announceMessage.target = MESSAGE_TARGET.DELETE_PROJECT;
    announceMessage.buttons = BUTTON_STYLE.CONFIRMATION;
    announceMessage.data = new NodeStatus(nodeName, !(status === 1));
    this.messageService.announceMessage(announceMessage);
  }
}