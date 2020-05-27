import { ChangeDetectorRef, Component, Input, OnInit, ViewContainerRef } from '@angular/core';
import { NodeService } from '../node.service';
import { MessageService } from '../../shared.service/message.service';
import { tap } from 'rxjs/operators';
import { zip } from 'rxjs';
import { CsModalChildBase } from '../../shared/cs-modal-base/cs-modal-child-base';
import { NodeGroupStatus, NodeStatus } from '../node.types';

@Component({
  selector: 'app-node-group-control',
  templateUrl: './node-group-control.component.html',
  styleUrls: ['./node-group-control.component.css']
})
export class NodeGroupControlComponent extends CsModalChildBase implements OnInit {
  @Input() nodeCurrent: NodeStatus;
  nodeGroupList: Array<NodeGroupStatus>;
  nodeGroupListSelect: Array<string>;
  selectedAddNodeGroup = '';
  selectedDelNodeGroup = '';
  isActionWip = false;

  constructor(private nodeService: NodeService,
              private messageService: MessageService,
              private view: ViewContainerRef,
              private changeDetectorRef: ChangeDetectorRef) {
    super();
    this.nodeGroupList = Array<NodeGroupStatus>();
    this.nodeGroupListSelect = Array<string>();
    this.changeDetectorRef.detach();
  }

  ngOnInit() {
    this.refreshData();
  }

  removeAlreadySelected() {
    this.nodeGroupListSelect.forEach(value => {
      let indexGroup = 0;
      const nodeGroup = this.nodeGroupList.find(
        (value2, index) => {
          indexGroup = index;
          return value2.name === value;
        });
      if (nodeGroup) {
        this.nodeGroupList.splice(indexGroup, 1);
      }
    });
  }

  refreshData() {
    const obs1 = this.nodeService.getNodeGroups()
      .pipe(tap((res: Array<NodeGroupStatus>) => this.nodeGroupList = res));
    const obs2 = this.nodeService.getNodeGroupsOfOneNode(this.nodeCurrent.nodeName)
      .pipe(tap((res: Array<string>) => this.nodeGroupListSelect = res));
    zip(obs1, obs2).subscribe(
      () => {
        this.removeAlreadySelected();
        this.isActionWip = false;
        this.selectedDelNodeGroup = '';
        this.selectedAddNodeGroup = '';
        this.changeDetectorRef.reattach();
      },
      () => this.modalOpened = false
    );
  }

  addNodeToNodeGroup(): void {
    if (!this.isActionWip &&
      this.selectedAddNodeGroup !== '' &&
      this.nodeGroupListSelect.indexOf(this.selectedAddNodeGroup) < 0) {
      this.isActionWip = true;
      this.nodeService.addNodeToNodeGroup(this.nodeCurrent.nodeName, this.selectedAddNodeGroup).subscribe(
        () => this.messageService.showAlert('NODE.NODE_GROUP_ADD_SUCCESS', {view: this.view}),
        () => this.modalOpened = false,
        () => this.refreshData());
    }
  }

  deleteNodeToNodeGroup(): void {
    if (!this.isActionWip &&
      this.selectedDelNodeGroup !== '' &&
      this.nodeGroupListSelect.indexOf(this.selectedDelNodeGroup) >= 0) {
      this.isActionWip = true;
      this.nodeService.deleteNodeToNodeGroup(this.nodeCurrent.nodeName, this.selectedDelNodeGroup).subscribe(
        () => this.messageService.showAlert('NODE.NODE_GROUP_REMOVE_SUCCESS', {view: this.view}),
        () => this.modalOpened = false,
        () => this.refreshData()
      );
    }
  }

  setNodeGroupSelectedToAdd(nodeGroup: NodeGroupStatus): void {
    if (this.nodeGroupListSelect.indexOf(nodeGroup.name) < 0) {
      if (this.selectedAddNodeGroup === nodeGroup.name) {
        this.selectedAddNodeGroup = '';
      } else {
        this.selectedAddNodeGroup = nodeGroup.name;
      }
    }
  }

  setNodeGroupSelectedToDel(nodeGroupName: string): void {
    if (this.selectedDelNodeGroup === nodeGroupName) {
      this.selectedDelNodeGroup = '';
    } else {
      this.selectedDelNodeGroup = nodeGroupName;
    }
  }
}
