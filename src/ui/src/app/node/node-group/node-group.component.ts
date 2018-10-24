import { Component, ComponentFactoryResolver, OnInit, ViewContainerRef } from '@angular/core';
import { NodeService } from "../node.service";
import { NodeCreateGroupComponent } from "../node-create-group/node-create-group.component";
import { MessageService } from "../../shared/message-service/message.service";
import { CsModalParentBase } from "../../shared/cs-modal-base/cs-modal-parent-base";
import { INodeGroup, Message, RETURN_STATUS } from "../../shared/shared.types";
import { TranslateService } from "@ngx-translate/core";

@Component({
  selector: 'node-group',
  templateUrl: './node-group.component.html',
  styleUrls: ['./node-group.component.css']
})
export class NodeGroupComponent extends CsModalParentBase implements OnInit {
  nodeGroupList: Array<INodeGroup>;
  isInLoadWip: boolean = false;

  constructor(private messageService: MessageService,
              private nodeService: NodeService,
              private resolver: ComponentFactoryResolver,
              private translateService: TranslateService,
              private view: ViewContainerRef) {
    super(resolver, view);
    this.nodeGroupList = Array<INodeGroup>();
  }

  ngOnInit() {
    this.refreshList();
  }

  refreshList() {
    this.isInLoadWip = true;
    this.nodeService.getNodeGroups().subscribe((res: Array<INodeGroup>) => {
        this.nodeGroupList = res;
        this.isInLoadWip = false
    }, () => this.isInLoadWip = false);
  }

  confirmToDeleteNodeGroup(groupName: string, groupId: number) {
    this.translateService.get('NODE.CONFIRM_TO_DELETE_NODE_GROUP', [groupName]).subscribe((msg: string) => {
      this.messageService.showDeleteDialog(msg, 'NODE.NODE_GROUP_DELETE').subscribe((message: Message) => {
        if (message.returnStatus == RETURN_STATUS.rsConfirm) {
          this.nodeService.deleteNodeGroup(groupId, groupName).subscribe(() => {
            this.messageService.showAlert('NODE.NODE_GROUP_DELETE_SUCCESS');
            this.refreshList();
          }, () => this.messageService.showAlert('NODE_GROUP_DELETE_FAILED', {alertType: 'alert-warning'}))
        }
      });
    })
  }

  showCreateNewGroup() {
    this.createNewModal(NodeCreateGroupComponent).onAfterCommit.subscribe(() => this.refreshList());
  }
}
