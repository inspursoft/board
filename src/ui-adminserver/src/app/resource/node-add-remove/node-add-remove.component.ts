import { Component, OnDestroy, OnInit, TemplateRef, ViewChild, ViewContainerRef } from '@angular/core';
import { ModalChildMessage } from '../../shared/cs-components-library/modal-child-base';
import { MessageService } from '../../shared/message/message.service';
import { NodeActionsType, NodeLogResponse, WsNodeResponseStatus } from '../resource.types';
import { Subject } from 'rxjs';
import { ResourceService } from '../services/resource.service';

@Component({
  templateUrl: './node-add-remove.component.html',
  styleUrls: ['./node-add-remove.component.css']
})
export class NodeAddRemoveComponent extends ModalChildMessage implements OnInit, OnDestroy {
  @ViewChild('consoleLogs', {read: ViewContainerRef}) consoleLogContainer: ViewContainerRef;
  @ViewChild('logTemplate') logTmp: TemplateRef<any>;
  actionType: NodeActionsType = NodeActionsType.Add;
  title = 'Node.Node_Form_Title_Add';
  nodeIp = '';
  successNotification: Subject<any>;
  isExecuting = false;

  constructor(protected messageService: MessageService,
              protected view: ViewContainerRef,
              private resourceService: ResourceService) {
    super(messageService);
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

  get alertView(): ViewContainerRef {
    return this.view;
  }

  get executeActionName(): string {
    return this.actionType === NodeActionsType.Add ? 'Node.Node_Form_Add' : 'Node.Node_Form_Remove';
  }

  get btnClassName(): string {
    return this.actionType === NodeActionsType.Add ? 'btn-primary' : 'btn-danger';
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
    this.resourceService.addRemoveNode(this.actionType, this.nodeIp).subscribe(
      (res: NodeLogResponse) => {
        this.consoleLogContainer.createEmbeddedView(this.logTmp,
          {message: res.message, status: res.status});
        this.isExecuting = true;
      },
      () => this.isExecuting = false,
      () => this.isExecuting = false
    );
  }
}
