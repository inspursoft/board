import { Component, EventEmitter, Output } from '@angular/core';
import { ValidationErrors } from '@angular/forms';
import { HttpErrorResponse } from '@angular/common/http';
import { Observable, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { NodeService } from '../node.service';
import { MessageService } from '../../shared.service/message.service';
import { CsModalChildBase } from '../../shared/cs-modal-base/cs-modal-child-base';
import { NodeGroupStatus } from '../node.types';

@Component({
  selector: 'app-node-create-group',
  templateUrl: './node-create-group.component.html',
  styleUrls: ['./node-create-group.component.css']
})
export class NodeCreateGroupComponent extends CsModalChildBase {
  newNodeGroupData: NodeGroupStatus;
  patternNodeGroupName: RegExp = /^[a-zA-Z0-9][a-zA-Z0-9_.-]*[a-zA-Z0-9]$/;
  @Output() afterCommit: EventEmitter<NodeGroupStatus>;

  constructor(private nodeService: NodeService,
              private messageService: MessageService) {
    super();
    this.afterCommit = new EventEmitter<NodeGroupStatus>();
    this.newNodeGroupData = new NodeGroupStatus({});
  }

  get checkNodeGroupNameFun() {
    return this.checkNodeGroupName.bind(this);
  }

  checkNodeGroupName(control: HTMLInputElement): Observable<ValidationErrors | null> {
    return this.nodeService.checkNodeGroupExist(control.value)
      .pipe(
        map(() => null),
        catchError((err: HttpErrorResponse) => {
          this.messageService.cleanNotification();
          if (err.status === 409) {
            return of({nodeGroupExist: 'NODE.NODE_GROUP_NAME_EXIST'});
          } else {
            return of(null);
          }
        })
      );
  }

  commitNodeGroup() {
    if (this.verifyInputExValid()) {
      this.nodeService.addNodeGroup(this.newNodeGroupData).subscribe(
        () => {
          this.afterCommit.emit(this.newNodeGroupData);
          this.messageService.showAlert('NODE.NODE_GROUP_CREATE_SUCCESS');
          this.modalOpened = false;
        },
        () => this.messageService.showAlert('NODE.NODE_GROUP_CREATE_FAILED', {alertType: 'danger', view: this.alertView}));
    }
  }
}
