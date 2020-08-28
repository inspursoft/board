import {
  Component,
  ComponentFactoryResolver,
  ElementRef,
  EventEmitter,
  OnDestroy,
  OnInit,
  Output,
  ViewChild,
  ViewContainerRef
} from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { NodeService } from '../node.service';
import { MessageService } from '../../shared.service/message.service';
import { NodeDetailComponent } from '../node-detail/node-detail.component';
import { Message, RETURN_STATUS } from '../../shared/shared.types';
import { CsModalParentBase } from '../../shared/cs-modal-base/cs-modal-parent-base';
import { NodeControlComponent } from '../node-control/node-control.component';
import { NodeStatus, NodeStatusType } from '../node.types';
import { interval, Subscription } from 'rxjs';

@Component({
  selector: 'app-node-list',
  templateUrl: './node-list.component.html',
  styleUrls: ['./node-list.component.css'],
})
export class NodeListComponent extends CsModalParentBase implements OnInit, OnDestroy {
  @ViewChild(NodeDetailComponent) nodeDetailModal;
  @Output() creatingNode: EventEmitter<Array<NodeStatus>>;
  autoRefreshSubscription: Subscription;
  nodeList: Array<NodeStatus>;
  isInLoadWip = false;

  constructor(private nodeService: NodeService,
              private translateService: TranslateService,
              private messageService: MessageService,
              public factoryResolver: ComponentFactoryResolver,
              public selfView: ViewContainerRef) {
    super(factoryResolver, selfView);
    this.nodeList = Array<NodeStatus>();
    this.creatingNode = new EventEmitter<Array<NodeStatus>>();
  }

  ngOnInit(): void {
    this.retrieve();
    this.autoRefreshSubscription = interval(60000).subscribe(() => this.retrieve());
  }

  ngOnDestroy(): void {
    this.autoRefreshSubscription.unsubscribe();
  }

  retrieve(): void {
    this.isInLoadWip = true;
    this.nodeService.getNodes().subscribe(
      (res: Array<NodeStatus>) => {
        this.nodeList = res.filter(value => value.nodeName !== '');
        this.isInLoadWip = false;
      },
      () => this.isInLoadWip = false
    );
  }

  getStatus(status: NodeStatusType): string {
    switch (status) {
      case NodeStatusType.Schedulable:
        return 'NODE.STATUS_SCHEDULABLE';
      case NodeStatusType.Unschedulable:
        return 'NODE.STATUS_UNSCHEDULABLE';
      case NodeStatusType.Unknown:
        return 'NODE.STATUS_UNKNOWN';
      case NodeStatusType.AutonomousOffline:
        return 'NODE.STATUS_OFFLINE';
    }
  }

  openNodeDetail(nodeName: string): void {
    this.nodeDetailModal.openNodeDetailModal(nodeName);
  }

  openNodeControl(node: NodeStatus): void {
    const instance = this.createNewModal(NodeControlComponent);
    instance.curNode = node;
  }

  confirmToToggleNodeStatus(node: NodeStatus): void {
    this.translateService.get('NODE.CONFIRM_TO_TOGGLE_NODE', [node.nodeName]).subscribe(res => {
      this.messageService.showYesNoDialog(res, 'NODE.TOGGLE_NODE').subscribe((message: Message) => {
        if (message.returnStatus === RETURN_STATUS.rsConfirm) {
          this.nodeService.toggleNodeStatus(node.nodeName, node.status !== NodeStatusType.Schedulable).subscribe(
            () => this.messageService.showAlert('NODE.SUCCESSFUL_TOGGLE'),
            () => this.messageService.showAlert('NODE.FAILED_TO_TOGGLE', {alertType: 'danger'}),
            () => this.retrieve()
          );
        }
      });
    });
  }

  showCreateNew() {
    this.creatingNode.emit(this.nodeList);
  }

  removeNode(node: NodeStatus) {
    this.translateService.get('NODE.CONFIRM_TO_REMOVE_NODE', [node.nodeIp]).subscribe(res => {
      this.messageService.showYesNoDialog(res).subscribe((message: Message) => {
        if (message.returnStatus === RETURN_STATUS.rsConfirm) {
          this.nodeService.removeEdgeNode(node.nodeName).subscribe(
            () => this.messageService.showAlert('NODE.REMOVE_NODE_SUCCESSFULLY'),
            () => this.messageService.showAlert('NODE.REMOVE_NODE_FAILED', {alertType: 'danger'}),
            () => this.retrieve()
          );
        }
      });
    });
  }
}
