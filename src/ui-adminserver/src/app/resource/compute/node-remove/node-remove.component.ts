import { Component, ElementRef, OnDestroy, OnInit, TemplateRef, ViewChild, ViewContainerRef } from '@angular/core';
import {
  ActionStatus,
  NodeActionsType,
  NodeDetails,
  NodeLog,
  NodeLogStatus,
  NodePostData,
  NodePreparationData
} from '../../resource.types';
import { interval, Observable, Subscription } from 'rxjs';
import { MessageService } from '../../../shared/message/message.service';
import { TranslateService } from '@ngx-translate/core';
import { ResourceService } from '../../services/resource.service';
import { HttpErrorResponse } from '@angular/common/http';
import { ModalChildBase } from '../../../shared/cs-components-library/modal-child-base';

@Component({
  selector: 'app-node-remove',
  templateUrl: './node-remove.component.html',
  styleUrls: ['./node-remove.component.css']
})
export class NodeRemoveComponent extends ModalChildBase implements OnInit, OnDestroy {
  @ViewChild('consoleLogs', {read: ViewContainerRef}) consoleLogContainer: ViewContainerRef;
  @ViewChild('logTemplate') logTmp: TemplateRef<any>;
  @ViewChild('divElement') divElement: ElementRef;
  @ViewChild('msgViewContainer', {read: ViewContainerRef}) view: ViewContainerRef;
  preparationData: NodePreparationData;
  postData: NodePostData;
  actionStatus = ActionStatus.Ready;
  logInfo: NodeLog;
  refreshingLog = false;
  autoRefreshLogSubscription: Subscription;
  curNodeLogStatus: NodeLogStatus;

  constructor(private messageService: MessageService,
              private translateService: TranslateService,
              private resourceService: ResourceService) {
    super();
    this.preparationData = new NodePreparationData({});
    this.postData = new NodePostData();
  }

  ngOnInit() {
    this.autoRefreshLogSubscription = interval(3000).subscribe(() => {
      if (this.actionStatus === ActionStatus.Executing) {
        this.refreshLog();
      }
    });
    this.getPreparationData();
  }

  ngOnDestroy() {
    this.autoRefreshLogSubscription.unsubscribe();
    delete this.autoRefreshLogSubscription;
    super.ngOnDestroy();
  }

  get executeBtnCaption(): string {
    if (this.actionStatus === ActionStatus.Ready) {
      return 'Node.Node_Detail_Remove';
    } else {
      return 'BUTTON.OK';
    }
  }

  get btnClassName(): string {
    return this.actionStatus === ActionStatus.Ready ? 'btn-danger' : 'btn-default';
  }

  get executing(): boolean {
    return this.actionStatus === ActionStatus.Executing ||
      this.actionStatus === ActionStatus.Preparing;
  }

  get masterTitle(): Observable<string> {
    return this.translateService.get('Node.Node_Detail_Master_Title', [this.preparationData.masterIp]);
  }

  get hostUsernameTitle(): Observable<string> {
    return this.translateService.get('Node.Node_Detail_Host_Username', [this.preparationData.hostIp]);
  }

  get hostPasswordTitle(): Observable<string> {
    return this.translateService.get('Node.Node_Detail_Host_Password', [this.preparationData.hostIp]);
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

  getPreparationData() {
    this.resourceService.getNodePreparation().subscribe(
      (res: NodePreparationData) => this.preparationData = res,
      () => {
        this.messageService.cleanNotification();
        this.messageService.showGlobalMessage('Node.Node_Detail_Error_Failed_Request', {view: this.view});
      }
    );
  }

  removeNode() {
    if (this.verifyInputExValid()) {
      this.actionStatus = ActionStatus.Preparing;
      this.resourceService.removeNode(this.postData).subscribe(
        (res: NodeLog) => this.logInfo = res,
        (err: HttpErrorResponse) => {
          if (err.status === 406) {
            this.messageService.cleanNotification();
            this.messageService.showAlert('Node.Node_Detail_Error_Node_Locked',
              {view: this.view, alertType: 'warning'}
            );
          } else {
            this.messageService.cleanNotification();
            this.messageService.showGlobalMessage('Node.Node_Detail_Error_Invalid_Password', {view: this.view});
          }
          this.actionStatus = ActionStatus.Ready;
        },
        () => this.actionStatus = ActionStatus.Executing
      );
    }
  }

  execute() {
    if (this.actionStatus === ActionStatus.Ready) {
      this.removeNode();
    } else {
      this.modalOpened = false;
    }
  }

  refreshLog() {
    if (this.actionStatus === ActionStatus.Executing) {
      const el = this.divElement.nativeElement as HTMLDivElement;
      this.refreshingLog = true;
      this.resourceService.getNodeLogDetail(this.logInfo.ip, this.logInfo.creationTime).subscribe(
        (res: NodeDetails) => {
          this.refreshingLog = false;
          this.consoleLogContainer.clear();
          for (const nodeDetail of res.data) {
            this.consoleLogContainer.createEmbeddedView(this.logTmp,
              {message: nodeDetail.message, status: nodeDetail.status}
            );
            if (nodeDetail.status === NodeLogStatus.Failed ||
              nodeDetail.status === NodeLogStatus.Success) {
              this.actionStatus = ActionStatus.Finished;
              this.curNodeLogStatus = nodeDetail.status;
            }
          }
          el.scrollTop = el.scrollHeight;
        },
        () => {
          this.refreshingLog = false;
          this.messageService.cleanNotification();
          this.messageService.showGlobalMessage('Node.Node_Detail_Error_Failed_Request', {view: this.view});
        }
      );
    }
  }

}
