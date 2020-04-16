import { Component, OnInit, ViewChild } from '@angular/core';
import { CsModalChildBase } from '../../shared/cs-modal-base/cs-modal-child-base';
import { NodeServiceControlComponent } from '../node-service-control/node-service-control.component';
import { NodeService } from '../node.service';
import { NodeStatus, NodeStatusType } from '../node.types';

@Component({
  selector: 'app-node-control',
  templateUrl: './node-control.component.html',
  styleUrls: ['./node-control.component.css']
})
export class NodeControlComponent extends CsModalChildBase implements OnInit {
  @ViewChild(NodeServiceControlComponent) serviceControl: NodeServiceControlComponent;
  curNode: NodeStatus;
  tabServiceControlActive = false;
  isActionWip = false;

  constructor(private nodeService: NodeService) {
    super();
  }

  ngOnInit(): void {

  }

  get btnDrainDisable(): boolean {
    return this.curNode.status !== NodeStatusType.Schedulable ||
      this.tabServiceControlActive === false;
  }

  drainService() {
    this.isActionWip = true;
    this.nodeService.drainNodeService(this.curNode.nodeName).subscribe(
      () => this.serviceControl.retrieve({page: {from: 0, to: 5}}),
      () => this.isActionWip = false,
      () => this.isActionWip = false
    );
  }
}
