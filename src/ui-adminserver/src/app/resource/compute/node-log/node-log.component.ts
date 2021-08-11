import { Component, ElementRef, OnInit, TemplateRef, ViewChild, ViewContainerRef } from '@angular/core';
import {
  ActionStatus,
  NodeDetails,
  NodeLog,
  NodeLogStatus,
  NodePreparationData
} from '../../resource.types';
import { MessageService } from '../../../shared/message/message.service';
import { TranslateService } from '@ngx-translate/core';
import { ResourceService } from '../../services/resource.service';
import { ModalChildBase } from '../../../shared/cs-components-library/modal-child-base';

@Component({
  selector: 'app-node-log',
  templateUrl: './node-log.component.html',
  styleUrls: ['./node-log.component.css']
})
export class NodeLogComponent extends ModalChildBase implements OnInit {
  @ViewChild('consoleLogs', {read: ViewContainerRef}) consoleLogContainer: ViewContainerRef;
  @ViewChild('logTemplate') logTmp: TemplateRef<any>;
  @ViewChild('divElement') divElement: ElementRef;
  @ViewChild('msgViewContainer', {read: ViewContainerRef}) view: ViewContainerRef;
  preparationData: NodePreparationData;
  actionStatus = ActionStatus.Ready;
  logInfo: NodeLog;
  refreshingLog = false;
  curNodeLogStatus: NodeLogStatus;

  constructor(private messageService: MessageService,
              private translateService: TranslateService,
              private resourceService: ResourceService) {
    super();
    this.preparationData = new NodePreparationData({});
  }

  ngOnInit() {
    this.refreshLog();
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

  refreshLog() {
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
