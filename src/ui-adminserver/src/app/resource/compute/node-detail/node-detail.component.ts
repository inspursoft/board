import { Component, ElementRef, OnInit, TemplateRef, ViewChild, ViewContainerRef } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { ActionStatus, NodeActionsType, NodeDetails, NodeLog, NodeLogStatus } from '../../resource.types';
import { MessageService } from '../../../shared/message/message.service';
import { ResourceService } from '../../services/resource.service';
import { ModalChildBase } from '../../../shared/cs-components-library/modal-child-base';

@Component({
  selector: 'app-node-log-detail',
  templateUrl: './node-detail.component.html',
  styleUrls: ['./node-detail.component.css']
})
export class NodeDetailComponent extends ModalChildBase implements OnInit {
  @ViewChild('consoleLogs', {read: ViewContainerRef}) consoleLogContainer: ViewContainerRef;
  @ViewChild('logTemplate') logTmp: TemplateRef<any>;
  @ViewChild('divElement') divElement: ElementRef;
  @ViewChild('msgViewContainer', {read: ViewContainerRef}) view: ViewContainerRef;
  actionType: NodeActionsType = NodeActionsType.Add;
  title = 'Node.Node_Detail_Title_Add';
  actionStatus = ActionStatus.Ready;
  logInfo: NodeLog;
  refreshingLog = false;

  constructor(private messageService: MessageService,
              private resourceService: ResourceService) {
    super();
  }

  ngOnInit() {
    if (this.actionType === NodeActionsType.Remove) {
      this.title = 'Node.Node_Detail_Title_Remove';
      this.removeNode();
    } else if (this.actionType === NodeActionsType.Log) {
      this.title = 'Node.Node_Detail_Title_Log';
      this.refreshLog();
    }
  }

  get executeBtnCaption(): string {
    if (this.actionType === NodeActionsType.Add && this.actionStatus === ActionStatus.Ready) {
      return 'Node.Node_Detail_Add';
    } else {
      return 'BUTTON.OK';
    }
  }

  get cancelBtnCaption(): string {
    return this.actionStatus === ActionStatus.Ready && this.actionType === NodeActionsType.Add ?
      'BUTTON.CANCEL' : 'Node.Node_Detail_Refresh';
  }

  get btnClassName(): string {
    if (this.actionType === NodeActionsType.Add) {
      return this.actionStatus === ActionStatus.Ready ? 'btn-primary' : 'btn-default';
    } else {
      return 'btn-default';
    }
  }

  get executing(): boolean {
    if (this.actionType === NodeActionsType.Log) {
      return false;
    }
    return this.actionStatus === ActionStatus.Executing;

  }

  getLogStyle(status: NodeLogStatus): { [key: string]: string } {
    switch (status) {
      case NodeLogStatus.Normal:
        return {color: 'white', fontSize: '14px'};
      case NodeLogStatus.Error:
        return {color: 'red', fontSize: '16px'};
      case NodeLogStatus.Failed:
        return {color: '#ff551b', fontSize: '16px'};
      case NodeLogStatus.Start:
        return {color: '#08ff22', fontSize: '14px'};
      case NodeLogStatus.Success:
        return {color: 'lightgreen', fontSize: '18px'};
      case NodeLogStatus.Warning:
        return {color: 'yellow', fontSize: '16px'};
      default:
        return {color: 'white', fontSize: '14px'};
    }
  }

  removeNode() {
    this.resourceService.removeNode(this.logInfo.ip).subscribe(
      (res: NodeLog) => this.logInfo = res,
      (err: HttpErrorResponse) => {
        this.messageService.cleanNotification();
        this.messageService.showGlobalMessage(err.message, {view: this.view});
      },
      () => {
        this.actionStatus = ActionStatus.Executing;
        this.refreshLog();
      }
    );
  }

  cancel() {
    this.actionStatus === ActionStatus.Ready && this.actionType === NodeActionsType.Add ?
      this.modalOpened = false : this.refreshLog();
  }

  execute() {
    if (this.actionStatus === ActionStatus.Ready && this.actionType === NodeActionsType.Add) {
      this.resourceService.addNode(this.logInfo.ip).subscribe(
        (res: NodeLog) => this.logInfo = res,
        (err: HttpErrorResponse) => {
          this.messageService.cleanNotification();
          this.messageService.showAlert(err.message, {alertType: 'danger', view: this.view});
        },
        () => {
          this.actionStatus = ActionStatus.Executing;
          this.refreshLog();
        }
      );
    } else {
      this.modalOpened = false;
    }
  }

  refreshLog() {
    const logFileName = `${this.logInfo.ip}@${this.logInfo.creationTime}.txt`;
    const el = this.divElement.nativeElement as HTMLDivElement;
    this.refreshingLog = true;
    this.resourceService.getNodeLog(logFileName).subscribe(
      (res: NodeDetails) => {
        this.refreshingLog = false;
        this.consoleLogContainer.clear();
        for (const nodeDetail of res.originData) {
          this.consoleLogContainer.createEmbeddedView(this.logTmp,
            {message: nodeDetail.message, status: nodeDetail.status}
          );
          if (nodeDetail.status === NodeLogStatus.Failed ||
            nodeDetail.status === NodeLogStatus.Success) {
            this.actionStatus = ActionStatus.Finished;
          }
        }
        el.scrollTop = el.scrollHeight;
      },
      (err: HttpErrorResponse) => {
        this.refreshingLog = false;
        this.messageService.cleanNotification();
        this.messageService.showAlert(err.message, {alertType: 'danger', view: this.view});
      });
  }

}
