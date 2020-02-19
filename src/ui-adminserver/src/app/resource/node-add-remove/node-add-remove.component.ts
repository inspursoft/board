import { Component, ElementRef, OnDestroy, OnInit, TemplateRef, ViewChild, ViewContainerRef } from '@angular/core';
import { ModalChildBase } from '../../shared/cs-components-library/modal-child-base';
import { MessageService } from '../../shared/message/message.service';
import { NodeActionsType, NodeReadyStatus, NodeLogResponse, WsNodeResponseStatus } from '../resource.types';
import { Subject } from 'rxjs';
import { ResourceService } from '../services/resource.service';

@Component({
  templateUrl: './node-add-remove.component.html',
  styleUrls: ['./node-add-remove.component.css']
})
export class NodeAddRemoveComponent extends ModalChildBase implements OnInit, OnDestroy {
  @ViewChild('consoleLogs', {read: ViewContainerRef}) consoleLogContainer: ViewContainerRef;
  @ViewChild('logTemplate') logTmp: TemplateRef<any>;
  @ViewChild('divElement') divElement: ElementRef;
  @ViewChild('msgViewContainer', {read: ViewContainerRef}) view: ViewContainerRef;
  actionType: NodeActionsType = NodeActionsType.Add;
  title = 'Node.Node_Form_Title_Add';
  nodeIp = '';
  successNotification: Subject<any>;
  readyStatus = NodeReadyStatus.Ready;

  constructor(private messageService: MessageService,
              private resourceService: ResourceService) {
    super();
    this.successNotification = new Subject();
  }

  ngOnInit() {
    if (this.actionType === NodeActionsType.Remove) {
      this.title = 'Node.Node_Form_Title_Remove';
    }
  }

  ngOnDestroy() {
    this.successNotification.unsubscribe();
    delete this.successNotification;
    super.ngOnDestroy();
  }

  get executeBtnCaption(): string {
    if (this.readyStatus === NodeReadyStatus.Ready) {
      return this.actionType === NodeActionsType.Add ? 'Node.Node_Form_Add' : 'Node.Node_Form_Remove';
    } else {
      return 'BUTTON.OK';
    }
  }

  get btnClassName(): string {
    if (this.readyStatus === NodeReadyStatus.Ready) {
      return this.actionType === NodeActionsType.Add ? 'btn-primary' : 'btn-danger';
    } else {
      return 'btn-default';
    }
  }

  getLogStyle(status: WsNodeResponseStatus): { [key: string]: string } {
    switch (status) {
      case WsNodeResponseStatus.Normal:
        return {color: 'white', fontSize: '14px'};
      case WsNodeResponseStatus.Error:
        return {color: 'red', fontSize: '16px'};
      case WsNodeResponseStatus.Failed:
        return {color: '#ff551b', fontSize: '16px'};
      case WsNodeResponseStatus.Start:
        return {color: '#08ff22', fontSize: '14px'};
      case WsNodeResponseStatus.Success:
        return {color: 'lightgreen', fontSize: '18px'};
      case WsNodeResponseStatus.Warning:
        return {color: 'yellow', fontSize: '16px'};
      default:
        return {color: 'white', fontSize: '14px'};
    }
  }

  cancel() {
    this.modalOpened = false;
  }

  execute() {
    if (this.readyStatus === NodeReadyStatus.Ready) {
      const el = this.divElement.nativeElement as HTMLDivElement;
      this.resourceService.addRemoveNode(this.actionType, this.nodeIp).subscribe(
        (res: NodeLogResponse) => {
          this.consoleLogContainer.createEmbeddedView(this.logTmp,
            {message: res.message, status: res.status}
          );
          this.readyStatus = NodeReadyStatus.Chatting;
          el.scrollTop = el.scrollHeight;
        },
        (err: any) => {
          this.readyStatus = NodeReadyStatus.Closed;
          if (err instanceof CloseEvent) {
            this.consoleLogContainer.createEmbeddedView(this.logTmp,
              {message: `Websocket connection closed.`, status: WsNodeResponseStatus.Warning}
            );
            el.scrollTop = el.scrollHeight;
          } else {
            const msg = 'Websocket connection failed.';
            this.messageService.showAlert(msg, {alertType: 'danger', view: this.view});
          }
        },
        () => {
          this.readyStatus = NodeReadyStatus.Closed;
          this.consoleLogContainer.createEmbeddedView(this.logTmp,
            {message: `Websocket connection closed.`, status: WsNodeResponseStatus.Warning}
          );
          el.scrollTop = el.scrollHeight;
        }
      );
    } else {
      if (this.readyStatus === NodeReadyStatus.Closed) {
        this.successNotification.next();
      }
      this.modalOpened = false;
    }
  }
}
