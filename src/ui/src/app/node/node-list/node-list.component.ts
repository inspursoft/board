import { Component, OnInit, ViewChild } from '@angular/core';
import { HttpErrorResponse } from "@angular/common/http";
import { Subscription } from "rxjs/Subscription";
import { Message } from "../../shared/message-service/message";
import { BUTTON_STYLE, MESSAGE_TARGET, MESSAGE_TYPE } from "../../shared/shared.const";
import { INode, NodeService } from "../node.service";
import { MessageService } from "../../shared/message-service/message.service";
import { NodeDetailComponent } from "../node-detail/node-detail.component";
import { NodeControlComponent } from "../node-control/node-control.component";

@Component({
  selector: 'node-list',
  templateUrl: './node-list.component.html',
  styleUrls: ['./node-list.component.css']
})
export class NodeListComponent implements OnInit {
  @ViewChild(NodeDetailComponent) nodeDetailModal;
  @ViewChild(NodeControlComponent) nodeControl;
  private _subscription: Subscription;
  nodeList: Array<INode> = [];
  isInLoadWip: boolean = false;

  constructor(private nodeService: NodeService,
              private messageService: MessageService) {
    this._subscription = this.messageService.messageConfirmed$
      .subscribe((message: Message) => {
        if (message.target == MESSAGE_TARGET.TOGGLE_NODE) {
          let node: INode = message.data;
          let m: Message = new Message();
          this.nodeService
            .toggleNodeStatus(node.node_name, node.status != 1)
            .subscribe(() => {
              m.message = 'NODE.SUCCESSFUL_TOGGLE';
              this.messageService.inlineAlertMessage(m);
              this.retrieve();
            }, () => {
              m.message = 'NODE.FAILED_TO_TOGGLE';
              m.type = MESSAGE_TYPE.COMMON_ERROR;
              this.messageService.inlineAlertMessage(m);
            });
        }
      });
  }

  ngOnInit(): void {
    this.retrieve();
  }

  ngOnDestroy(): void {
    if (this._subscription) {
      this._subscription.unsubscribe();
    }
  }

  retrieve(): void {
    this.isInLoadWip = true;
    this.nodeService.getNodes().subscribe((res: Array<INode>) => {
        this.nodeList = res;
        this.isInLoadWip = false;
      },
      (error: HttpErrorResponse) => {
        this.messageService.dispatchError(error);
        this.isInLoadWip = false;
      });
  }

  getStatus(status: number): string {
    switch (status) {
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

  openNodeControl(node:INode):void{
    this.nodeControl.openNodeControlModal(node);
  }

  confirmToToggleNodeStatus(node: INode): void {
    let announceMessage = new Message();
    announceMessage.title = 'NODE.TOGGLE_NODE';
    announceMessage.message = 'NODE.CONFIRM_TO_TOGGLE_NODE';
    announceMessage.params = Array.from([node.node_name]);
    announceMessage.target = MESSAGE_TARGET.TOGGLE_NODE;
    announceMessage.buttons = BUTTON_STYLE.CONFIRMATION;
    announceMessage.data = node;
    this.messageService.announceMessage(announceMessage);
  }
}
