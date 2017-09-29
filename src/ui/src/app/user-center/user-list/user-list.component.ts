import { Component, OnInit, OnDestroy } from "@angular/core";
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
  userCountPerPage: number = 10;
  curUser: User;
  curEditModel: editModel = editModel.emNew;
  showNewUser: boolean = false;
  setUserSystemAdminIng: boolean = false;
  setUserProjectAdminIng: boolean = false;

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
    this.refreshData();
  }

  ngOnDestroy(): void {
    if (this._deleteSubscription) {
      this._deleteSubscription.unsubscribe();
    }
  }

  get currentUserID(): number {
    return this.appInitService.currentUser["user_id"];
  }

  refreshData(username?: string,
              user_list_page: number = 0,
              user_list_page_size: number = 0): void {
    this.userService.getUserList(username, user_list_page, user_list_page_size)
      .then(res => {
        this.userListData = res.filter(value => {
          return value.user_name != "admin";
        });
      })
      .catch(err => this.messageService.dispatchError(err, ''));
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
      .catch(err => this.messageService.dispatchError(err));
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

  setUserSystemAdmin(user: User) {
    this.setUserSystemAdminIng = true;
    let userSystemAdmin = user.user_system_admin == 1 ? 0 : 1;
    this.userService.setUserSystemAdmin(user.user_id, userSystemAdmin)
      .then(() => {
        this.setUserSystemAdminIng = false;
        user.user_system_admin = userSystemAdmin;
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
        this.messageService.dispatchError(err);
      })
  }

  setUserProjectAdmin(user: User) {
    this.setUserProjectAdminIng = true;
    let userProjectAdmin = user.user_project_admin == 1 ? 0 : 1;
    this.userService.setUserProjectAdmin(user.user_id, userProjectAdmin)
      .then(() => {
        this.setUserProjectAdminIng = false;
        user.user_project_admin = userProjectAdmin;
        let m: Message = new Message();
        if (user.user_project_admin === 1) {
          m.message = "USER_CENTER.SUCCESSFUL_SET_PROJECT_ADMIN";
        } else {
          m.message = "USER_CENTER.SUCCESSFUL_SET_NOT_PROJECT_ADMIN";
        }
        this.messageService.inlineAlertMessage(m);
      })
      .catch(err => {
        this.setUserProjectAdminIng = false;
        this.messageService.dispatchError(err)
      })
  }
}
