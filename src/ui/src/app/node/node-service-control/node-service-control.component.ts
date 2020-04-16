import { Component, Input, OnInit } from '@angular/core';
import { ClrDatagridStateInterface } from '@clr/angular';
import { NodeService } from '../node.service';
import { CsModalChildBase } from '../../shared/cs-modal-base/cs-modal-child-base';
import { NodeControlStatus, NodeStatus, ServiceInstance } from '../node.types';

@Component({
  selector: 'app-node-service-control',
  templateUrl: './node-service-control.component.html',
  styleUrls: ['./node-service-control.component.css']
})
export class NodeServiceControlComponent extends CsModalChildBase implements OnInit {
  @Input() nodeCurrent: NodeStatus;
  nodeControlStatus: NodeControlStatus;
  serviceInstanceList: Array<ServiceInstance>;
  curPageIndex = 1;
  curPageSize = 6;

  constructor(private nodeService: NodeService) {
    super();
    this.nodeControlStatus = new NodeControlStatus({});
    this.serviceInstanceList = Array<ServiceInstance>();
  }

  ngOnInit() {
    this.nodeService.getNodeControlStatus(this.nodeCurrent.nodeName).subscribe(
      (res: NodeControlStatus) => {
        this.nodeControlStatus = res;
        this.retrieve({page: {from: 0, to: 5}});
      }
    );
  }

  get phaseStyle(): { [p: string]: string } {
    switch (this.nodeControlStatus.nodePhase) {
      case 'Pending':
        return {color: 'darkorange'};
      case 'Running':
        return {color: 'green'};
      case 'Terminal':
        return {color: 'red'};
      default:
        return {color: 'black'};
    }
  }

  retrieve(page: ClrDatagridStateInterface) {
    if (Reflect.has(page, 'page')) {
      const from = page.page.from;
      const to = page.page.to;
      this.serviceInstanceList = this.nodeControlStatus.serviceInstances.slice(from, to);
    }
  }

  cancel() {
    this.modalOpened = false;
  }
}
