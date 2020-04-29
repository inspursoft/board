import { Component, ComponentFactoryResolver, OnInit, ViewContainerRef } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { NodeService } from '../node.service';
import { NodeCreateGroupComponent } from '../node-create-group/node-create-group.component';
import { MessageService } from '../../shared.service/message.service';
import { CsModalParentBase } from '../../shared/cs-modal-base/cs-modal-parent-base';
import { Message, RETURN_STATUS } from '../../shared/shared.types';
import { NodeGroupStatus } from '../node.types';

@Component({
  selector: 'app-node-group',
  templateUrl: './node-group.component.html',
  styleUrls: ['./node-group.component.css']
})
export class NodeGroupComponent extends CsModalParentBase implements OnInit {
  nodeGroupList: Array<NodeGroupStatus>;
  isInLoadWip = false;

  constructor(private messageService: MessageService,
              private nodeService: NodeService,
              private resolver: ComponentFactoryResolver,
              private translateService: TranslateService,
              private view: ViewContainerRef) {
    super(resolver, view);
    this.nodeGroupList = Array<NodeGroupStatus>();
  }

  ngOnInit() {
    this.refreshList();
  }

  refreshList() {
    this.isInLoadWip = true;
    this.nodeService.getNodeGroups().subscribe((res: Array<NodeGroupStatus>) => {
      this.nodeGroupList = res;
      this.isInLoadWip = false;
    }, () => this.isInLoadWip = false);
  }

  confirmToDeleteNodeGroup(groupName: string, groupId: number) {
    this.translateService.get('NODE.CONFIRM_TO_DELETE_NODE_GROUP', [groupName]).subscribe(
      (msg: string) => {
        this.messageService.showDeleteDialog(msg, 'NODE.NODE_GROUP_DELETE').subscribe(
          (message: Message) => {
            if (message.returnStatus === RETURN_STATUS.rsConfirm) {
              this.nodeService.deleteNodeGroup(groupId, groupName).subscribe(
                () => {
                  this.messageService.showAlert('NODE.NODE_GROUP_DELETE_SUCCESS');
                  this.refreshList();
                }, () => this.messageService.showAlert('NODE_GROUP_DELETE_FAILED', {alertType: 'warning'}));
            }
          });
      });
  }

  showCreateNewGroup() {
    this.createNewModal(NodeCreateGroupComponent).afterCommit.subscribe(() => this.refreshList());
  }
}
