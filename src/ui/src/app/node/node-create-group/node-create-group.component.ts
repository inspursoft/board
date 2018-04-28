import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { INodeGroup, NodeService } from "../node.service";
import { Message } from "../../shared/message-service/message";
import { MessageService } from "../../shared/message-service/message.service";
import { MESSAGE_TYPE } from "../../shared/shared.const";
import { ValidationErrors } from "@angular/forms";
import { HttpErrorResponse } from "@angular/common/http";

class NodeGroup implements INodeGroup {
  nodegroup_id: number = 0;
  nodegroup_project: string = "";
  nodegroup_name: string = "";
  nodegroup_comment: string = "";

  constructor() {
  }
}

@Component({
  selector: 'node-create-group',
  templateUrl: './node-create-group.component.html',
  styleUrls: ['./node-create-group.component.css']
})
export class NodeCreateGroupComponent implements OnInit {
  isOpen: boolean = false;
  newNodeGroupData: NodeGroup;
  patternNodeGroupName: RegExp = /^[a-zA-Z0-9][a-zA-Z0-9_.-]*[a-zA-Z0-9]$/;
  @Output("onAfterCommit") onAfterCommit: EventEmitter<INodeGroup>;

  constructor(private nodeService: NodeService,
              private messageService: MessageService) {
    this.onAfterCommit = new EventEmitter<INodeGroup>();
    this.newNodeGroupData = new NodeGroup();
  }

  ngOnInit() {
  }

  showModal() {
    this.isOpen = true;
  }

  get checkNodeGroupNameFun() {
    return this.checkNodeGroupName.bind(this);
  }

  checkNodeGroupName(control: HTMLInputElement): Promise<ValidationErrors | null> {
    return this.nodeService.checkNodeGroupExist(control.value)
      .toPromise()
      .then(() => null)
      .catch(err => {
        if (err && err instanceof HttpErrorResponse && (err as HttpErrorResponse).status == 409) {
          return {nodeGroupExist: "NODE.NODE_GROUP_NAME_EXIST"}
        }
        this.messageService.dispatchError(err);
      });
  }

  commitNodeGroup() {
    let msg: Message = new Message();
    this.nodeService.addNodeGroup(this.newNodeGroupData).subscribe(() => {
      this.onAfterCommit.emit(this.newNodeGroupData);
      this.isOpen = false;
      msg.message = "NODE.NODE_GROUP_CREATE_SUCCESS";
      this.messageService.inlineAlertMessage(msg)
    }, () => {
      this.isOpen = false;
      msg.message = "NODE.NODE_GROUP_CREATE_FAILED";
      msg.type = MESSAGE_TYPE.COMMON_ERROR;
      this.messageService.inlineAlertMessage(msg);
    })
  }
}
