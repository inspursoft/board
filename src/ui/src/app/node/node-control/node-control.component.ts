import { ChangeDetectorRef, Component, OnInit } from '@angular/core';
import { INode, INodeGroup, NodeService } from "../node.service";
import "rxjs/add/operator/zip"
import "rxjs/add/operator/do"
import "rxjs/add/operator/catch"

@Component({
  selector: 'node-control',
  templateUrl: './node-control.component.html',
  styleUrls: ['./node-control.component.css']
})
export class NodeControlComponent implements OnInit {
  nodeControlOpened: boolean = false;
  nodeGroupList: Array<INodeGroup>;
  nodeGroupListSelect: Array<string>;
  nodeCurrent: INode;
  selectedAddNodeGroup: string;
  selectedDelNodeGroup: string;
  isActionWip: boolean = false;

  constructor(private nodeService: NodeService,
              private changeDetectorRef: ChangeDetectorRef) {
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
      .do((res: Array<INodeGroup>) => this.nodeGroupList = res);
    let obs2 = this.nodeService.getNodeGroupsOfOneNode(this.nodeCurrent.node_name)
      .do((res: Array<string>) => this.nodeGroupListSelect = res);
    obs1.zip(obs2).subscribe(
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
        () => this.refreshData(),
        () => this.nodeControlOpened = false);
    }
  }

  deleteNodeToNodeGroup(): void {
    if (!this.isActionWip &&
      this.selectedDelNodeGroup != '' &&
      this.nodeGroupListSelect.indexOf(this.selectedDelNodeGroup) >= 0) {
      this.isActionWip = true;
      this.nodeService.deleteNodeToNodeGroup(this.nodeCurrent.node_name, this.selectedDelNodeGroup)
        .subscribe(() => this.refreshData())
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
