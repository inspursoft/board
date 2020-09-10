import { Component, ComponentFactoryResolver, OnDestroy, OnInit, ViewContainerRef } from '@angular/core';
import { interval, Subscription } from 'rxjs';
import { NodeLog, NodeLogs } from '../../resource.types';
import { ResourceService } from '../../services/resource.service';
import { MessageService } from '../../../shared/message/message.service';
import { HttpErrorResponse } from '@angular/common/http';
import { Message, ReturnStatus } from '../../../shared/message/message.types';
import { NodeLogComponent } from '../node-log/node-log.component';

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
    this.subscriptionUpdate = interval(5000).subscribe(() => this.retrieve());
  }

  ngOnDestroy(): void {
    this.subscriptionUpdate.unsubscribe();
    delete this.subscriptionUpdate;
  }

  retrieve() {
    this.resourceService.getNodeLogs(this.curPageIndex, this.curPageSize)
      .subscribe((res: NodeLogs) => this.nodeLogs = res);
  }

  stopExecuting(log: NodeLog) {
    this.messageService.showYesNoDialog('Node.Node_Logs_Stop_Ask', 'Node.Node_Logs_Stop')
      .subscribe((msg: Message) => {
        if (msg.returnStatus === ReturnStatus.rsConfirm) {
          this.resourceService.stopNodeAction(log).subscribe(
            () => {
              this.messageService.showAlert('Node.Node_Logs_Stop_Success');
              this.retrieve();
            },
          );
        }
      });
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
    const factory = this.resolver.resolveComponentFactory(NodeLogComponent);
    const detailRef = this.view.createComponent(factory);
    detailRef.instance.logInfo = log;
    detailRef.instance.openModal().subscribe(
      () => this.view.remove(this.view.indexOf(detailRef.hostView))
    );
  }
}
