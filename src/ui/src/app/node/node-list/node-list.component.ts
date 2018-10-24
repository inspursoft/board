import { Component, OnInit, ViewChild } from '@angular/core';
import { NodeService } from "../node.service";
import { MessageService } from "../../shared/message-service/message.service";
import { NodeDetailComponent } from "../node-detail/node-detail.component";
import { NodeControlComponent } from "../node-control/node-control.component";
import { INode, Message, RETURN_STATUS } from "../../shared/shared.types";
import { TranslateService } from "@ngx-translate/core";

@Component({
  selector: 'node-list',
  templateUrl: './node-list.component.html',
  styleUrls: ['./node-list.component.css']
})
export class NodeListComponent implements OnInit {
  @ViewChild(NodeDetailComponent) nodeDetailModal;
  @ViewChild(NodeControlComponent) nodeControl;
  nodeList: Array<INode> = [];
  isInLoadWip: boolean = false;

  constructor(private nodeService: NodeService,
              private translateService: TranslateService,
              private messageService: MessageService) {
  }

  ngOnInit(): void {
    this.retrieve();
  }

  retrieve(): void {
    this.isInLoadWip = true;
    this.nodeService.getNodes().subscribe((res: Array<INode>) => {
        this.nodeList = res;
        this.isInLoadWip = false;
      },
      () => this.isInLoadWip = false);
  }

  getStatus(status: number): string {
    switch (status) {
      case 1:
        return 'NODE.STATUS_SCHEDULABLE';
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
    this.translateService.get('NODE.CONFIRM_TO_TOGGLE_NODE', [node.node_name]).subscribe(res => {
      this.messageService.showYesNoDialog(res, 'NODE.TOGGLE_NODE').subscribe((message: Message) => {
        if (message.returnStatus == RETURN_STATUS.rsConfirm) {
          this.nodeService.toggleNodeStatus(node.node_name, node.status != 1).subscribe(
            () => this.messageService.showAlert('NODE.SUCCESSFUL_TOGGLE'),
            () => this.messageService.showAlert('NODE.FAILED_TO_TOGGLE', {alertType: 'alert-danger'}),
            () => this.retrieve())
        }
      })
    });
  }
}
