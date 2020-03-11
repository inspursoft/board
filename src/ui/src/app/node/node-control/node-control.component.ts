import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { NodeService } from "../node.service";
import { CsModalChildBase } from "../../shared/cs-modal-base/cs-modal-child-base";
import { MessageService } from "../../shared.service/message.service";
import { INode, INodeGroup } from "../../shared/shared.types";
import { tap } from "rxjs/operators";
import { zip } from "rxjs";

@Component({
  selector: 'node-control',
  templateUrl: './node-control.component.html',
  styleUrls: ['./node-control.component.css']
})
export class NodeControlComponent extends CsModalChildBase implements OnInit {
  nodeControlOpened: boolean = false;
  nodeGroupList: Array<INodeGroup>;
  nodeGroupListSelect: Array<string>;
  nodeCurrent: INode;
  selectedAddNodeGroup: string;
  selectedDelNodeGroup: string;
  isActionWip: boolean = false;

  constructor(private nodeService: NodeService,
              private messageService: MessageService,
              private changeDetectorRef: ChangeDetectorRef) {
    super();
    this.nodeGroupList = Array<INodeGroup>();
    this.nodeGroupListSelect = Array<string>();
    this.changeDetectorRef.detach();
  }

  ngOnInit() {
  }

  removeAlreadySelected() {
    this.nodeGroupListSelect.forEach(value => {
      let indexGroup = 0;
      let nodeGroup = this.nodeGroupList.find(
        (value2, index) => {
          indexGroup = index;
          return value2.nodegroup_name == value
        });
      if (nodeGroup) {
        this.nodeGroupList.splice(indexGroup, 1);
      }
    })
  }

  refreshData() {
    let obs1 = this.nodeService.getNodeGroups()
      .pipe(tap((res: Array<INodeGroup>) => this.nodeGroupList = res));
    let obs2 = this.nodeService.getNodeGroupsOfOneNode(this.nodeCurrent.node_name)
      .pipe(tap((res: Array<string>) => this.nodeGroupListSelect = res));
    zip(obs1,obs2).subscribe(
      () => {
        this.removeAlreadySelected();
        this.isActionWip = false;
        this.selectedDelNodeGroup = "";
        this.selectedAddNodeGroup = "";
        this.changeDetectorRef.reattach();
      },
      () => this.nodeControlOpened = false);
  }

  openNodeControlModal(node: INode): void {
    this.changeDetectorRef.detach();
    this.nodeCurrent = node;
    this.nodeControlOpened = true;
    this.refreshData();
  }

  addNodeToNodeGroup(): void {
    if (!this.isActionWip &&
      this.selectedAddNodeGroup != '' &&
      this.nodeGroupListSelect.indexOf(this.selectedAddNodeGroup) < 0) {
      this.isActionWip = true;
      this.nodeService.addNodeToNodeGroup(this.nodeCurrent.node_name, this.selectedAddNodeGroup).subscribe(
        () => this.messageService.showAlert('NODE.NODE_GROUP_ADD_SUCCESS',{view: this.alertView}),
        () => this.nodeControlOpened = false,
        () => this.refreshData());
    }
  }

  deleteNodeToNodeGroup(): void {
    if (!this.isActionWip &&
      this.selectedDelNodeGroup != '' &&
      this.nodeGroupListSelect.indexOf(this.selectedDelNodeGroup) >= 0) {
      this.isActionWip = true;
      this.nodeService.deleteNodeToNodeGroup(this.nodeCurrent.node_name, this.selectedDelNodeGroup).subscribe(
        () => this.messageService.showAlert('NODE.NODE_GROUP_REMOVE_SUCCESS',{view: this.alertView}),
        () => this.nodeControlOpened = false,
        () => this.refreshData())
    }
  }

  setNodeGroupSelectedToAdd(nodeGroup: INodeGroup): void {
    if (this.nodeGroupListSelect.indexOf(nodeGroup.nodegroup_name) < 0) {
      if (this.selectedAddNodeGroup == nodeGroup.nodegroup_name) {
        this.selectedAddNodeGroup = "";
      } else {
        this.selectedAddNodeGroup = nodeGroup.nodegroup_name;
      }
    }
  }

  setNodeGroupSelectedToDel(nodeGroupName: string): void {
    if (this.selectedDelNodeGroup == nodeGroupName) {
      this.selectedDelNodeGroup = "";
    } else {
      this.selectedDelNodeGroup = nodeGroupName;
    }
  }
}
