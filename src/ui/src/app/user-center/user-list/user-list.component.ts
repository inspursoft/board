import { Component, OnInit, OnDestroy } from "@angular/core";

import { State } from "clarity-angular";

import { UserService } from "../user-service/user-service";
import { User } from "../user";
import { editModel } from "../user-new-edit/user-new-edit.component"
import { Message } from "app/shared/message-service/message";
import { MessageService } from "app/shared/message-service/message.service";
import { Subscription } from "rxjs/Subscription";
import { BUTTON_STYLE } from "app/shared/shared.const"
import { AppInitService } from "../../app.init.service";

@Component({
  selector: "user-list",
  templateUrl: "./user-list.component.html",
  styleUrls: ["./user-list.component.css"]
})

export class UserList implements OnInit, OnDestroy {
  _deleteSubscription: Subscription;
  userListData: Array<User> = Array<User>();
  userListErrMsg: string = "";

  curUser: User;
  curEditModel: editModel = editModel.emNew;
  showNewUser: boolean = false;
  setUserSystemAdminIng: boolean = false;
  isInLoading: boolean = false;
  checkboxRevertInfo: {isNeeded: boolean; value: boolean;};

  totalRecordCount: number;
  pageIndex: number = 1;
  pageSize: number = 15;

  constructor(private userService: UserService,
              private appInitService: AppInitService,
              private messageService: MessageService) {
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
        .catch(err => this.messageService.dispatchError(err));
    });
  }

  ngOnDestroy(): void {
    if (this._deleteSubscription) {
      this._deleteSubscription.unsubscribe();
    }
  }

  get currentUserID(): number {
    return this.appInitService.currentUser["user_id"];
  }

  refreshData(state?: State): void {
    setTimeout(()=>{
      this.isInLoading = true;
      this.userService.getUserList('', this.pageIndex, this.pageSize)
        .then(res => {
          this.totalRecordCount = res.pagination.total_count;
          this.userListData = res.user_list;
          this.isInLoading = false;
        })
        .catch(err => {
          this.messageService.dispatchError(err, '');
          this.isInLoading = false;
        });
    });
  }

  addUser() {
    this.curUser = new User();
    this.curEditModel = editModel.emNew;
    this.showNewUser = true;
  }

  editUser(user: User) {
    if (user.user_deleted != 1 && user.user_id != 1 && user.user_id != this.currentUserID) {
      this.curEditModel = editModel.emEdit;
      this.userService.getUser(user.user_id)
        .then(user => {
          this.curUser = user;
          this.showNewUser = true;
        })
        .catch(err => this.messageService.dispatchError(err));
    }
  }

  deleteUser(user: User) {
    if (user.user_deleted != 1 && user.user_id != 1 && user.user_id != this.currentUserID) {
      let m: Message = new Message();
      m.title = "USER_CENTER.DELETE_USER";
      m.buttons = BUTTON_STYLE.DELETION;
      m.data = user;
      m.params = [user.user_name];
      m.message = "USER_CENTER.CONFIRM_DELETE_USER";
      this.messageService.announceMessage(m);
    }
  }

  setUserSystemAdmin(user: User) {
    this.setUserSystemAdminIng = true;
    let oldUserSystemAdmin = user.user_system_admin;
    this.userService.setUserSystemAdmin(user.user_id, oldUserSystemAdmin == 1 ? 0 : 1)
      .then(() => {
        this.setUserSystemAdminIng = false;
        user.user_system_admin = oldUserSystemAdmin == 1 ? 0 : 1;
        let m: Message = new Message();
        if (user.user_system_admin === 1) {
          m.message = "USER_CENTER.SUCCESSFUL_SET_SYS_ADMIN";
          user.user_project_admin = 1;
        } else {
          m.message = "USER_CENTER.SUCCESSFUL_SET_NOT_SYS_ADMIN";
        }
        this.messageService.inlineAlertMessage(m);
      })
      .catch(err => {
        this.setUserSystemAdminIng = false;
        this.checkboxRevertInfo = {isNeeded: true, value: oldUserSystemAdmin == 1};
        this.messageService.dispatchError(err);
      })
  }
}
