import { Component, EventEmitter, Output } from '@angular/core';
import { NodeService } from "../node.service";
import { MessageService } from "../../shared.service/message.service";
import { ValidationErrors } from "@angular/forms";
import { HttpErrorResponse } from "@angular/common/http";
import { CsModalChildBase } from "../../shared/cs-modal-base/cs-modal-child-base";
import { INodeGroup } from "../../shared/shared.types";
import { Observable, of } from "rxjs";
import { catchError, map } from "rxjs/operators";

class NodeGroup implements INodeGroup {
  nodegroup_id: number = 0;
  nodegroup_project: string = "";
  nodegroup_name: string = "";
  nodegroup_comment: string = "";
}

@Component({
  selector: 'node-create-group',
  templateUrl: './node-create-group.component.html',
  styleUrls: ['./node-create-group.component.css']
})
export class NodeCreateGroupComponent extends CsModalChildBase {
  newNodeGroupData: NodeGroup;
  patternNodeGroupName: RegExp = /^[a-zA-Z0-9][a-zA-Z0-9_.-]*[a-zA-Z0-9]$/;
  @Output("onAfterCommit") onAfterCommit: EventEmitter<INodeGroup>;

  constructor(private nodeService: NodeService,
              private messageService: MessageService) {
    super();
    this.onAfterCommit = new EventEmitter<INodeGroup>();
    this.newNodeGroupData = new NodeGroup();
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
          if (err.status == 409) {
            return of({nodeGroupExist: "NODE.NODE_GROUP_NAME_EXIST"})
          } else {
            return of(null)
          }
        }))
  }

  commitNodeGroup() {
    if (this.verifyInputExValid()) {
      this.nodeService.addNodeGroup(this.newNodeGroupData).subscribe(
        () => {
          this.onAfterCommit.emit(this.newNodeGroupData);
          this.messageService.showAlert('NODE.NODE_GROUP_CREATE_SUCCESS');
          this.modalOpened = false;
        },
        () => this.messageService.showAlert('NODE.NODE_GROUP_CREATE_FAILED', {alertType: 'danger', view: this.alertView}))
    }
  }
}
