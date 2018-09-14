import { Component, EventEmitter, Input, Output, OnInit } from "@angular/core"
import { UserService } from "../../user-center/user-service/user-service";
import { User } from "../../user-center/user";
import { MessageService } from "../message-service/message.service";

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
    this.userService.getCurrentUser().then(res => this.curUser = res)
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
        this.isOpen = false;
        this.messageService.showAlert('ACCOUNT.ACCOUNT_SETTING_SUCCESS');
      })
      .catch(() => this.isOpen = false);
  }
}