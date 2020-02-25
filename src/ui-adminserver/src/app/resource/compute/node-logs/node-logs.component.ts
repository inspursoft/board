import { Component, ComponentFactoryResolver, OnInit, ViewContainerRef } from '@angular/core';
import { NodeActionsType, NodeLog, NodeLogs } from '../../resource.types';
import { ResourceService } from '../../services/resource.service';
import { NodeDetailComponent } from '../node-detail/node-detail.component';

@Component({
  selector: 'app-node-logs',
  templateUrl: './node-logs.component.html',
  styleUrls: ['./node-logs.component.css']
})
export class NodeLogsComponent implements OnInit {
  loadingWIP = false;
  nodeLogs: NodeLogs;
  curPageIndex = 1;
  curPageSize = 15;

  constructor(private resourceService: ResourceService,
              private resolver: ComponentFactoryResolver,
              private view: ViewContainerRef) {
    this.nodeLogs = new NodeLogs({});
  }

  ngOnInit() {
    this.retrieve();
  }

  retrieve() {
    this.loadingWIP = true;
    this.resourceService.getNodeLogs().subscribe(
      (res: NodeLogs) => this.nodeLogs = res,
      () => this.loadingWIP = false,
      () => this.loadingWIP = false
    );
  }

  deleteLog(log: NodeLog) {

  }

  showLogDetail(log: NodeLog) {
    const factory = this.resolver.resolveComponentFactory(NodeDetailComponent);
    const detailRef = this.view.createComponent(factory);
    detailRef.instance.actionType = NodeActionsType.Log;
    detailRef.instance.logInfo = log;
    detailRef.instance.openModal().subscribe(
      () => this.view.remove(this.view.indexOf(detailRef.hostView))
    );
  }
}
