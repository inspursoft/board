/**
 * Created by liyanq on 8/3/17.
 */

import { Component, Input, Output, EventEmitter } from "@angular/core"
import { AppInitService } from "../../app.init.service";
import { MessageService } from "../message-service/message.service";
import { UserService } from "../../user-center/user-service/user-service";
import { CsModalChildBase } from "../cs-modal-base/cs-modal-child-base";

@Component({
  selector: "change-password",
  styleUrls: ["./change-password.component.css"],
  templateUrl: "./change-password.component.html",
  providers: [UserService]
})
export class ChangePasswordComponent extends CsModalChildBase{
  _isOpen: boolean = false;
  curPassword: string = "";
  newPassword: string = "";
  newPasswordConfirm: string = "";
  isWorkWip: boolean = false;
  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();

  constructor(private appInitService: AppInitService,
              private userService: UserService,
              private messageService: MessageService) {
    super();
  }

  @Input()
  get isOpen() {
    return this._isOpen;
  }

  set isOpen(open: boolean) {
    this._isOpen = open;
    this.isOpenChange.emit(this._isOpen);
  }

  submitChangePassword(): void {
    if (this.verifyInputValid()) {
      let curUser = this.appInitService.currentUser;
      if (curUser && curUser["user_id"]) {
        this.isWorkWip = true;
        this.userService.changeUserPassword(curUser["user_id"], this.curPassword, this.newPassword)
          .then(() => {
            this.isOpen = false;
            this.messageService.showAlert('HEAD_NAV.CHANGE_PASSWORD_SUCCESS')
          })
          .catch(err => {
            this.isWorkWip = false;
            if (err && err["status"] && err["status"] == 403) {
              this.messageService.showAlert('HEAD_NAV.OLD_PASSWORD_WRONG',{alertType:'alert-warning',view: this.alertView});
            } else if (err && err["status"] && err["status"] == 401) {
              this.messageService.showAlert('ERROR.HTTP_401',{alertType:'alert-warning',view: this.alertView});
            }
            else {
              this.isOpen = false;
            }
          })
      } else {
        this.messageService.showAlert('ERROR.HTTP_401',{alertType:'alert-warning',view: this.alertView});
      }
    }
  }
}