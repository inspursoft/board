import { Component, EventEmitter, Input, Output, OnInit } from "@angular/core"
import { AppInitService } from "../../app.init.service";
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
  curUser: User = new User();
  isAlertClose: boolean = true;
  errMessage: string;

  constructor(private appInitService: AppInitService,
              private userService: UserService,
              private messageService: MessageService) {
  }

  ngOnInit() {
    this.userService.getCurrentUser()
      .then(res => {
        this.curUser = res;
      })
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
    this.userService.usesChangeAccount(this.curUser)
      .then(() => {
        let m: Message = new Message();
        m.message = "ACCOUNT.ACCOUNT_SETTING_SUCCESS";
        this.messageService.inlineAlertMessage(m);
        this.isOpen = false;
      })
      .catch(err => {
        if (err){
          if(err.status === 409) {
            this.isAlertClose = false;
            this.errMessage = "ACCOUNT.EMAIL_ALREADY_EXISTS";
          } else {
            this.messageService.dispatchError(err);
          }
        }
      });
  }
}