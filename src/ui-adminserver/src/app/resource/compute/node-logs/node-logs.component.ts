import { Component, ComponentFactoryResolver, OnDestroy, OnInit, ViewContainerRef } from '@angular/core';
import { interval, Subscription } from 'rxjs';
import { NodeActionsType, NodeLog, NodeLogs } from '../../resource.types';
import { ResourceService } from '../../services/resource.service';
import { NodeDetailComponent } from '../node-detail/node-detail.component';
import { MessageService } from '../../../shared/message/message.service';
import { HttpErrorResponse } from '@angular/common/http';
import { Message, ReturnStatus } from '../../../shared/message/message.types';

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
              private messageService: MessageService,
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

  deleteLogInfo(log: NodeLog) {
    this.messageService.showDeleteDialog('Node.Node_Logs_Delete_Ask').subscribe((msg: Message) => {
      if (msg.returnStatus === ReturnStatus.rsConfirm) {
        this.resourceService.deleteNodeLog(log.creationTime).subscribe(
          () => this.messageService.showAlert('Node.Node_Logs_Delete_Log_Success'),
          (res: HttpErrorResponse) => {
            if (res.status === 409) {
              this.messageService.cleanNotification();
              this.messageService.showAlert('Node.Node_Logs_Delete_Log_In_Use', {alertType: 'warning'});
            }
          },
          () => this.retrieve()
        );
      }
    });
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
