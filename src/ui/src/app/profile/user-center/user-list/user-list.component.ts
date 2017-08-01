import { Component, OnInit, OnDestroy } from "@angular/core";
import { UserService } from "app/profile/user-center/user-service/user-service";
import { User } from "app/profile/user-center/user";
import { editModel } from "../user-new-edit/user-new-edit.component"
import { Message } from "app/shared/message-service/message";
import { MessageService } from "app/shared/message-service/message.service";
import { Subscription } from "rxjs/Subscription";
import { BUTTON_STYLE } from "app/shared/shared.const"

@Component({
  selector: "user-list",
  templateUrl: "./user-list.component.html",
  styleUrls: ["./user-list.component.css"]
})

export class UserList implements OnInit, OnDestroy {
  _deleteSubscription: Subscription;
  userListData: Array<User> = Array<User>();
  userListErrMsg: string = "";
  userCountPerPage: number = 2;
  curUser: User;
  curEditModel: editModel = editModel.emNew;
  curPage: number = 1;
  showNewUser: boolean = false;

  constructor(private userService: UserService,
              private messageService: MessageService) {
  }

  refreshData(username?: string,
              user_list_page: number = 0,
              user_list_page_size: number = 0): void {
    this.userService.getUserList(username, user_list_page, user_list_page_size)
      .then(res => this.userListData = res)
      .catch(err => {
        if (err) {
          switch (err.status) {
            case 400:
              this.userListErrMsg = 'USER_CENTER.INVALID_USER_ID';
              break;
            default:
              this.messageService.dispatchError(err, '');
          }
        }
      });
  }

  addUser() {
    this.curUser = new User();
    this.curEditModel = editModel.emNew;
    this.showNewUser = true;
  }

  editUser(user: User) {
    this.curEditModel = editModel.emEdit;
    this.userService.getUser(user.user_id)
      .then(user => {
        this.curUser = user;
        this.showNewUser = true;
      })
      .catch((err) => {
        if (err) {
          switch (err.status) {
            case 400:
              this.userListErrMsg = 'USER_CENTER.INVALID_USER_ID';
              break;
            default:
              this.messageService.dispatchError(err, '');
          }
        }
      });
  }

  pageChange(page: number) {
    this.curPage = page;
  }

  deleteUser(user: User) {
    let m: Message = new Message();
    m.title = "USER_CENTER.DELETE_USER";
    m.buttons = BUTTON_STYLE.DELETION;
    m.data = user;
    m.params = [user.user_name];
    m.message = "USER_CENTER.CONFIRM_DELETE_USER";
    this.messageService.announceMessage(m);
  }

  ngOnInit() {
    this._deleteSubscription = this.messageService.messageConfirmed$.subscribe(next => {
      this.userService.deleteUser(next.data)
        .then(() => {
          this.refreshData();
          let m: Message = new Message();
          m.message = "USER_CENTER.DELETE_USER_SUCCESS";
          this.messageService.inlineAlertMessage(m);
        })
        .catch((err) => {
          if (err) {
            switch (err.status) {
              case 400:
                this.userListErrMsg = 'USER_CENTER.INVALID_USER_ID';
                break;
              default:
                this.messageService.dispatchError(err, '');
            }
          }
        });
    });
    this.refreshData();
  }

  ngOnDestroy(): void {
    if (this._deleteSubscription) {
      this._deleteSubscription.unsubscribe();
    }
  }
}
