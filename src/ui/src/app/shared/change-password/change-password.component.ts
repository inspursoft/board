/**
 * Created by liyanq on 8/3/17.
 */

import { Component, Input, Output, EventEmitter, OnInit } from "@angular/core"
import { AppInitService } from "../../app.init.service";
import { MessageService } from "../message-service/message.service";
import { UserService } from "../../profile/user-center/user-service/user-service";
import { Message } from "../message-service/message";

@Component({
  selector: "change-password",
  styleUrls: ["./change-password.component.css"],
  templateUrl: "./change-password.component.html",
  providers: [UserService]
})
export class ChangePasswordComponent implements OnInit {
  _isOpen: boolean = false;
  _curUser: {[key: string]: any};
  isAlertClose: boolean = true;
  errMessage: string;
  curPassword: string = "";
  newPassword: string = "";
  newPasswordConfirm: string = "";

  constructor(private appInitService: AppInitService,
              private userService: UserService,
              private messageService: MessageService) {

  }

  ngOnInit() {
    this._curUser = this.appInitService.currentUser;//user id
  }

  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();

  @Input()
  get isOpen() {
    return this._isOpen;
  }

  set isOpen(open: boolean) {
    this._isOpen = open;
    this.isOpenChange.emit(this._isOpen);
  }

  submitChangePassword(): void {
    if (this._curUser && this._curUser["id"]) {
      this.userService.changeUserPassword(this._curUser["id"], this.curPassword, this.newPassword)
        .then(() => {
          let m: Message = new Message();
          m.message = "HEAD_NAV.CHANGE_PASSWORD_SUCCESS";
          this.messageService.inlineAlertMessage(m);
          this.isOpen = false;
        })
        .catch(err => {
          if (err && err["status"] && err["status"] == 403) {
            this.errMessage = "HEAD_NAV.OLD_PASSWORD_WRONG";
            this.isAlertClose = false;
          } else {
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