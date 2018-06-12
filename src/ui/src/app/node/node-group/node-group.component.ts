import { AfterViewInit, Component, OnDestroy, OnInit, ViewChild } from '@angular/core';
import { INodeGroup, NodeService } from "../node.service";
import { NodeCreateGroupComponent } from "../node-create-group/node-create-group.component";
import { MessageService } from "../../shared/message-service/message.service";
import { Message } from "../../shared/message-service/message";
import { Subscription } from "rxjs/Subscription";
import { BUTTON_STYLE, MESSAGE_TARGET, MESSAGE_TYPE } from "../../shared/shared.const";
import { HttpErrorResponse } from "@angular/common/http";

@Component({
  selector: 'node-group',
  templateUrl: './node-group.component.html',
  styleUrls: ['./node-group.component.css']
})
export class NodeGroupComponent implements OnInit, AfterViewInit, OnDestroy {
  nodeGroupList: Array<INodeGroup>;
  isInLoadWip: boolean = false;
  _deleteSubscription: Subscription;
  @ViewChild(NodeCreateGroupComponent) newGroupComponent: NodeCreateGroupComponent;

  constructor(private messageService: MessageService,
              private nodeService: NodeService) {
    this.nodeGroupList = Array<INodeGroup>();
    this._deleteSubscription = this.messageService.messageConfirmed$.subscribe((msg: Message) => {
      if (msg.target == MESSAGE_TARGET.DELETE_NODE_GROUP) {
        let inlineMessage: Message = new Message();
        this.nodeService.deleteNodeGroup(msg.data,msg.params[0]).subscribe(() => {
          inlineMessage.message = "NODE.NODE_GROUP_DELETE_SUCCESS";
          this.messageService.inlineAlertMessage(inlineMessage);
          this.refreshList();
        }, () => {
          inlineMessage.message = "NODE.NODE_GROUP_DELETE_FAILED";
          inlineMessage.type = MESSAGE_TYPE.COMMON_ERROR;
          this.messageService.inlineAlertMessage(inlineMessage);
        });
      }
    })
  }

  ngOnInit() {
    this.refreshList();
  }

  ngOnDestroy() {
    this._deleteSubscription.unsubscribe();
  }

  ngAfterViewInit() {

  }

  refreshList() {
    this.isInLoadWip = true;
    this.nodeService.getNodeGroups()
      .subscribe((res: Array<INodeGroup>) => {
        this.nodeGroupList = res;
        this.isInLoadWip = false
      }, (err: HttpErrorResponse) => {
        this.isInLoadWip = false;
        this.messageService.dispatchError(err);
      });
  }

  confirmToDeleteNodeGroup(groupName: string, groupId: number) {
    let announceMessage = new Message();
    announceMessage.title = 'NODE.NODE_GROUP_DELETE';
    announceMessage.message = 'NODE.CONFIRM_TO_DELETE_NODE_GROUP';
    announceMessage.params = [groupName];
    announceMessage.target = MESSAGE_TARGET.DELETE_NODE_GROUP;
    announceMessage.buttons = BUTTON_STYLE.DELETION;
    announceMessage.data = groupId;
    this.messageService.announceMessage(announceMessage);
  }

}
