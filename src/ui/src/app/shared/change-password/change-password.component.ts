/**
 * Created by liyanq on 8/3/17.
 */

import { Component, Input, Output, EventEmitter } from "@angular/core"
import { AppInitService } from "../../app.init.service";
import { MessageService } from "../message-service/message.service";
import { UserService } from "../../user-center/user-service/user-service";
import { Message } from "../message-service/message";
import { CsComponentBase } from "../cs-components-library/cs-component-base";

@Component({
  selector: "change-password",
  styleUrls: ["./change-password.component.css"],
  templateUrl: "./change-password.component.html",
  providers: [UserService]
})
export class ChangePasswordComponent extends CsComponentBase{
  _isOpen: boolean = false;
  isAlertClose: boolean = true;
  errMessage: string;
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
            let m: Message = new Message();
            m.message = "HEAD_NAV.CHANGE_PASSWORD_SUCCESS";
            this.messageService.inlineAlertMessage(m);
            this.isWorkWip = false;
            this.isOpen = false;
          })
          .catch(err => {
            this.isWorkWip = false;
            if (err && err["status"] && err["status"] == 403) {
              this.errMessage = "HEAD_NAV.OLD_PASSWORD_WRONG";
              this.isAlertClose = false;
            } else if (err && err["status"] && err["status"] == 401) {
              this.errMessage = "ERROR.INVALID_USER";
              this.isAlertClose = false;
            }
            else {
              this.isOpen = false;
              this.messageService.dispatchError(err);
            }
          })
      } else {
        this.errMessage = "ERROR.INVALID_USER";
        this.isAlertClose = false;
      }
    }
  }
}