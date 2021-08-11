import { Component, ElementRef, OnDestroy, OnInit, TemplateRef, ViewChild, ViewContainerRef } from '@angular/core';
import {
  ActionStatus,
  NodeDetails,
  NodeLog,
  NodeLogStatus,
  NodePostData,
  NodePreparationData
} from '../../resource.types';
import { interval, Observable, of, Subscription } from 'rxjs';
import { MessageService } from '../../../shared/message/message.service';
import { TranslateService } from '@ngx-translate/core';
import { ResourceService } from '../../services/resource.service';
import { HttpErrorResponse } from '@angular/common/http';
import { ModalChildBase } from '../../../shared/cs-components-library/modal-child-base';
import { AbstractControl, ValidationErrors } from '@angular/forms';
import { map } from 'rxjs/operators';

@Component({
  selector: 'app-node-create',
  templateUrl: './node-create.component.html',
  styleUrls: ['./node-create.component.css']
})
export class NodeCreateComponent extends ModalChildBase implements OnInit, OnDestroy {
  @ViewChild('consoleLogs', {read: ViewContainerRef}) consoleLogContainer: ViewContainerRef;
  @ViewChild('logTemplate') logTmp: TemplateRef<any>;
  @ViewChild('divElement') divElement: ElementRef;
  @ViewChild('msgViewContainer', {read: ViewContainerRef}) view: ViewContainerRef;
  preparationData: NodePreparationData;
  postData: NodePostData;
  title = 'Node.Node_Detail_Title_Add';
  actionStatus = ActionStatus.Ready;
  logInfo: NodeLog;
  refreshingLog = false;
  autoRefreshLogSubscription: Subscription;
  curNodeLogStatus: NodeLogStatus;
  newNodeList: Array<{ nodeIp: string, nodePassword: string, checked: boolean }>;

  constructor(private messageService: MessageService,
              private translateService: TranslateService,
              private resourceService: ResourceService) {
    super();
    this.preparationData = new NodePreparationData({});
    this.postData = new NodePostData();
    this.newNodeList = new Array<{ nodeIp: string, nodePassword: string, checked: boolean }>();
  }

  ngOnInit() {
    this.newNodeList.push({nodePassword: '', nodeIp: '', checked: false});
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

  get executeBtnCaption(): string {
    if (this.actionStatus === ActionStatus.Ready) {
      return 'Node.Node_Detail_Add';
    } else {
      return 'BUTTON.OK';
    }
  }

  get checkIpExist() {
    return this.checkIpExistFun.bind(this);
  }

  get checkPassword() {
    return this.checkPasswordFun.bind(this);
  }

  checkIpExistFun(control: AbstractControl): Observable<ValidationErrors | null> {
    const ip = control.value;
    if (this.newNodeList.find(value => value.nodeIp === ip && value.checked === false)) {
      return this.translateService.get('Node.Node_Detail_Error_Node_Repeat').pipe(
        map(msg => {
          return {ipExists: msg};
        })
      );
    } else {
      return of(null);
    }
  }

  checkPasswordFun(control: AbstractControl): Observable<ValidationErrors | null> {
    const ps = control.value as string;
    if (ps.indexOf('_') > -1) {
      return this.translateService.get('Node.Node_Detail_Error_Node_Reserve').pipe(
        map(msg => {
          return {charReserved: msg};
        })
      );
    } else {
      return of(null);
    }
  }

  removeNodeInfo(index: number): void {
    this.newNodeList.splice(index, 1);
  }

  addNodeInfo(): void {
    this.newNodeList.push({nodeIp: '', nodePassword: '', checked: false});
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

  generatePostData(): void {
    this.postData.nodeIp = '';
    this.postData.nodePassword = '';
    this.newNodeList.forEach((nodeInfo) => {
      nodeInfo.checked = true;
      this.postData.nodeIp += `${nodeInfo.nodeIp}_`;
      this.postData.nodePassword += `${nodeInfo.nodePassword}_`;
    });
    this.postData.nodeIp = this.postData.nodeIp.substr(0, this.postData.nodeIp.length - 1);
    this.postData.nodePassword = this.postData.nodePassword.substr(0, this.postData.nodePassword.length - 1);
  }

  execute() {
    if (this.actionStatus === ActionStatus.Ready) {
      this.addNode();
    } else {
      this.modalOpened = false;
    }
  }

  addNode() {
    if (this.newNodeList.length === 0) {
      this.messageService.showAlert('Node.Node_Detail_Error_Node_Limit',
        {view: this.view, alertType: 'warning'}
      );
    } else if (this.verifyInputExValid()) {
      this.generatePostData();
      this.actionStatus = ActionStatus.Preparing;
      this.resourceService.addNode(this.postData).subscribe(
        (res) => this.logInfo = res,
        (err: HttpErrorResponse) => {
          if (err.status === 406) {
            this.messageService.cleanNotification();
            this.messageService.showAlert('Node.Node_Detail_Error_Node_Locked',
              {view: this.view, alertType: 'warning'}
            );
          } else {
            this.messageService.cleanNotification();
            this.messageService.showGlobalMessage('Node.Node_Detail_Error_Bad_Input', {view: this.view});
          }
          this.actionStatus = ActionStatus.Ready;
        },
        () => this.actionStatus = ActionStatus.Executing
      );
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
