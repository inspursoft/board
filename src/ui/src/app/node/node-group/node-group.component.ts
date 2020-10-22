import { Component, ComponentFactoryResolver, OnInit, ViewContainerRef } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { NodeService } from '../node.service';
import { NodeCreateGroupComponent } from '../node-create-group/node-create-group.component';
import { MessageService } from '../../shared.service/message.service';
import { CsModalParentBase } from '../../shared/cs-modal-base/cs-modal-parent-base';
import { Message, RETURN_STATUS } from '../../shared/shared.types';
import { NodeGroupStatus } from '../node.types';
import { Observable, of } from 'rxjs';
import { NodeGroupEditComponent } from '../node-group-edit/node-group-edit.component';

@Component({
  selector: 'app-node-group',
  templateUrl: './node-group.component.html',
  styleUrls: ['./node-group.component.css']
})
export class NodeGroupComponent extends CsModalParentBase implements OnInit {
  nodeGroupList: Array<NodeGroupStatus>;
  isInLoadWip = false;
  nodeListMap: Map<number, string>;

  constructor(private messageService: MessageService,
              private nodeService: NodeService,
              private resolver: ComponentFactoryResolver,
              private translateService: TranslateService,
              private view: ViewContainerRef) {
    super(resolver, view);
    this.nodeGroupList = Array<NodeGroupStatus>();
    this.nodeListMap = new Map<number, string>();
  }

  ngOnInit() {
    this.refreshList();
  }

  refreshList() {
    this.isInLoadWip = true;
    this.nodeListMap.clear();
    this.nodeService.getNodeGroups().subscribe((res: Array<NodeGroupStatus>) => {
      this.nodeGroupList = res;
      res.forEach(value => {
        this.nodeService.getGroupMembers(value.id).subscribe(
          (nodeList: Array<string>) => {
            let nodes = '';
            nodeList.forEach(node => nodes += `${node}<br>`);
            this.nodeListMap.set(value.id, nodes);
          }
        );
      });
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

  editGroup(group: NodeGroupStatus): void {
    const ref = this.createNewModal(NodeGroupEditComponent);
    ref.nodeGroup = new NodeGroupStatus(group.res);
    ref.nodeGroup.initFromRes();
    ref.afterUpdate.subscribe(() => this.refreshList());
  }
}
