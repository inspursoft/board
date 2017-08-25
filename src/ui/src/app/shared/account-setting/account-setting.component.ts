import { Component, EventEmitter, Input, Output,OnInit } from "@angular/core"
import { AppInitService } from "../../app.init.service";
import { UserService } from "../../profile/user-center/user-service/user-service";
import { User } from "../../profile/user-center/user";
import { MessageService } from "../message-service/message.service";

@Component({
  selector: "user-setting",
  templateUrl: "./account-setting.component.html",
  styleUrls: ["./account-setting.component.css"],
  providers: [UserService]
})
export class AccountSettingComponent implements OnInit{
  _isOpen: boolean = false;
  curUser: User;
  isAlertClose: boolean = true;
  errMessage: string;

  constructor(private appInitService: AppInitService,
              private userService: UserService,
              private messageService: MessageService) {
  }

  ngOnInit(){
    let curUserID = this.appInitService.currentUser["user_id"];
    this.userService.getUser(curUserID).then(res => {
      console.log(res);
      this.curUser = res;
    }).catch(err => {
      if (err && err["status"] && err["status"] == 404) {
        this.errMessage = "ERROR.INVALID_USER";
        this.isAlertClose = false;
      } else if (err && err["status"] && err["status"] == 401) {
        this.errMessage = "ERROR.INVALID_USER";
        this.isAlertClose = false;
      }
      else {
        this.isOpen = false;
        this.messageService.dispatchError(err);
      }
    });
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
}