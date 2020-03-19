import { Component, ComponentFactoryResolver, OnDestroy, OnInit, ViewContainerRef } from '@angular/core';
import { interval, Subscription } from 'rxjs';
import { NodeActionsType, NodeLog, NodeLogs } from '../../resource.types';
import { ResourceService } from '../../services/resource.service';
import { NodeDetailComponent } from '../node-detail/node-detail.component';

@Component({
  selector: 'app-node-logs',
  templateUrl: './node-logs.component.html',
  styleUrls: ['./node-logs.component.css']
})
export class NodeLogsComponent implements OnInit, OnDestroy {
  nodeLogs: NodeLogs;
  curPageIndex = 1;
  curPageSize = 10;
  subscriptionUpdate: Subscription;

  constructor(private resourceService: ResourceService,
              private resolver: ComponentFactoryResolver,
              private view: ViewContainerRef) {
    this.nodeLogs = new NodeLogs({});
  }

  ngOnInit() {
    this.subscriptionUpdate = interval(3000).subscribe(() => this.retrieve());
  }

  ngOnDestroy(): void {
    this.subscriptionUpdate.unsubscribe();
    delete this.subscriptionUpdate;
  }

  retrieve() {
    this.resourceService.getNodeLogs(this.curPageIndex, this.curPageSize)
      .subscribe((res: NodeLogs) => this.nodeLogs = res);
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
