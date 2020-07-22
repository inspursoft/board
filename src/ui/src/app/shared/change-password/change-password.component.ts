/**
 * Created by liyanq on 8/3/17.
 */

import { Component, EventEmitter, Input, Output } from '@angular/core';
import { HttpErrorResponse } from '@angular/common/http';
import { AppInitService } from '../../shared.service/app-init.service';
import { UserService } from '../../admin/user-center/user-service/user-service';
import { CsModalChildBase } from '../cs-modal-base/cs-modal-child-base';
import { MessageService } from '../../shared.service/message.service';

@Component({
  selector: 'app-change-password',
  styleUrls: ['./change-password.component.css'],
  templateUrl: './change-password.component.html',
  providers: [UserService]
})
export class ChangePasswordComponent extends CsModalChildBase {
  isOpenValue = false;
  curPassword = '';
  newPassword = '';
  newPasswordConfirm = '';
  isWorkWip = false;
  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(private appInitService: AppInitService,
              private userService: UserService,
              private messageService: MessageService) {
    super();
  }

  @Input()
  get isOpen() {
    return this.isOpenValue;
  }

  set isOpen(open: boolean) {
    this.isOpenValue = open;
    this.isOpenChange.emit(this.isOpenValue);
  }

  submitChangePassword(): void {
    if (this.verifyInputExValid()) {
      const curUser = this.appInitService.currentUser;
      if (curUser.user_id > 0) {
        this.isWorkWip = true;
        this.userService.changeUserPassword(curUser.user_id, this.curPassword, this.newPassword).subscribe(() => {
            this.isOpen = false;
            this.messageService.showAlert('HEAD_NAV.CHANGE_PASSWORD_SUCCESS');
          },
          (err: HttpErrorResponse) => {
            this.isWorkWip = false;
            if (err && err.status && err.status === 403) {
              this.messageService.showAlert('HEAD_NAV.OLD_PASSWORD_WRONG', {alertType: 'warning', view: this.alertView});
            } else if (err && err.status && err.status === 401) {
              this.messageService.showAlert('ERROR.HTTP_401', {alertType: 'warning', view: this.alertView});
            } else {
              this.isOpen = false;
            }
          });
      } else {
        this.messageService.showAlert('ERROR.HTTP_401', {alertType: 'warning', view: this.alertView});
      }
    }
  }
}
