import { Component, EventEmitter, Input, Output, OnInit } from "@angular/core"
import { UserService } from "../../user-center/user-service/user-service";
import { User } from "../../user-center/user";
import { MessageService } from "../message-service/message.service";
import { Message } from "../message-service/message";

@Component({
  selector: "user-setting",
  templateUrl: "./account-setting.component.html",
  styleUrls: ["./account-setting.component.css"],
  providers: [UserService]
})
export class AccountSettingComponent implements OnInit {
  _isOpen: boolean = false;
  isWorkWip: boolean = false;
  curUser: User = new User();

  constructor(private userService: UserService,
              private messageService: MessageService) {
  }

  ngOnInit() {
    this.userService.getCurrentUser()
      .then(res => this.curUser = res)
      .catch(err =>this.messageService.dispatchError(err));
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

  submitAccountSetting() {
    this.isWorkWip = true;
    this.userService.usesChangeAccount(this.curUser)
      .then(() => {
        let m: Message = new Message();
        m.message = "ACCOUNT.ACCOUNT_SETTING_SUCCESS";
        this.messageService.inlineAlertMessage(m);
        this.isWorkWip = false;
        this.isOpen = false;
      })
      .catch(err => {
        this.isWorkWip = false;
        this.messageService.dispatchError(err);
      });
  }
}