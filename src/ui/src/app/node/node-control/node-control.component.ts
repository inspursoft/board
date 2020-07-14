import { Component, OnInit, ViewChild } from '@angular/core';
import { CsModalChildBase } from '../../shared/cs-modal-base/cs-modal-child-base';
import { NodeServiceControlComponent } from '../node-service-control/node-service-control.component';
import { NodeService } from '../node.service';
import { NodeStatus, NodeStatusType } from '../node.types';
import { AppInitService } from '../../shared.service/app-init.service';

@Component({
  selector: 'app-node-control',
  templateUrl: './node-control.component.html',
  styleUrls: ['./node-control.component.css']
})
export class NodeControlComponent extends CsModalChildBase implements OnInit {
  @ViewChild(NodeServiceControlComponent) serviceControl: NodeServiceControlComponent;
  curNode: NodeStatus;
  serviceInstanceCount = 0;
  nodeDeletable = false;
  tabServiceControlActive = false;
  isActionWip = false;

  constructor(private nodeService: NodeService,
              private appInitService: AppInitService) {
    super();
  }

  ngOnInit(): void {

  }

  get showDeleteNodeTip(): boolean {
    return !this.nodeDeletable &&
      !this.isActionWip &&
      !this.curNode.isEdge &&
      !this.curNode.isMaster &&
      this.tabServiceControlActive;
  }

  get btnDrainDisable(): boolean {
    return this.curNode.status !== NodeStatusType.Unschedulable ||
      this.isActionWip ||
      this.tabServiceControlActive === false;
  }

  get adminServerDeleteNodeUrl(): string {
    if (this.showDeleteNodeTip) {
      return `javascript:void(0)`;
    } else {
      return `http://${this.appInitService.systemInfo.board_host}:8082/resource/node-list`;
    }
  }

  drainService() {
    this.isActionWip = true;
    this.nodeService.drainNodeService(this.curNode.nodeName, this.serviceInstanceCount).subscribe(
      () => this.serviceControl.refreshData(),
      () => this.isActionWip = false,
      () => this.isActionWip = false
    );
  }
}
