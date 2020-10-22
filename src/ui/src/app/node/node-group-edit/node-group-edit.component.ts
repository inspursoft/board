import { Component, EventEmitter, OnInit, Output } from '@angular/core';
import { CsModalChildBase } from '../../shared/cs-modal-base/cs-modal-child-base';
import { NodeGroupStatus } from '../node.types';
import { NodeService } from '../node.service';
import { MessageService } from '../../shared.service/message.service';

@Component({
  selector: 'app-node-group-edit',
  templateUrl: './node-group-edit.component.html',
  styleUrls: ['./node-group-edit.component.css']
})
export class NodeGroupEditComponent extends CsModalChildBase implements OnInit {
  afterUpdate: EventEmitter<any>;
  patternNodeGroupName: RegExp = /^[a-zA-Z0-9][a-zA-Z0-9_.-]*[a-zA-Z0-9]$/;
  nodeGroup: NodeGroupStatus;
  isActionWip = false;

  constructor(private nodeService: NodeService,
              private messageService: MessageService) {
    super();
    this.afterUpdate = new EventEmitter<any>();
  }

  ngOnInit() {
  }

  updateNodeGroup(): void {
    this.isActionWip = true;
    this.nodeService.updateGroup(this.nodeGroup).subscribe(
      () => {
        this.messageService.showAlert('NODE.NODE_GROUP_UPDATE_SUCCESS');
        this.afterUpdate.emit();
      },
      () => {
        this.modalOpened = false;
        this.messageService.showAlert('NODE.NODE_GROUP_UPDATE_FAILED', {alertType: 'danger'});
      },
      () => this.modalOpened = false
    );
  }
}
