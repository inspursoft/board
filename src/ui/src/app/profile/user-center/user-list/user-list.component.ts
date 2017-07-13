import { Component, OnInit, OnDestroy } from "@angular/core";
import { UserService } from "app/profile/user-center/user-service/user-service";
import { user, User } from "app/profile/user-center/user";
import { editModel } from "../user-new-edit/user-new-edit.component"
import { ConfirmationMessage } from "app/shared/service/confirmation-message";
import { MessageService } from "app/shared/service/message.service";
import { Subscription } from "rxjs/Subscription";

@Component({
  selector: "user-list",
  templateUrl: "./user-list.component.html",
  styleUrls: ["./user-list.component.css"]
})

export class UserList implements OnInit, OnDestroy {
  userListData: Array<user> = new Array<user>();
  userListErrMsg: string = "";
  curUser: user;
  curEditModel: editModel = editModel.emNew;
  showNewUser: boolean = false;
  _deleteSubscription: Subscription;
  constructor(
    private userService: UserService,
    private messageService: MessageService) { }

  refreshData(
    username?: string,
    user_list_page: number = 0,
    user_list_page_size: number = 0): void {
    this.userService.getUserList().then(
      res => (this.userListData = res),
      (reason: Response) => {
        this.userListErrMsg = `status:${reason.status}`;
      }
    )
  }

  addUser() {
    this.curUser = new User();
    this.curEditModel = editModel.emNew;
    this.showNewUser = true;
  }

  editUser(user: user) {
    this.curEditModel = editModel.emEdit;
    this.curUser = user;
    this.showNewUser = true;
  }

  deleteUser(user: user) {
    let m: ConfirmationMessage = new ConfirmationMessage();
    m.title = "删除用户";
    m.data = user;
    m.message = `确定删除[${user.user_name}]么?`;
    this.messageService.announceMessage(m);
  }

  ngOnInit() {
    this._deleteSubscription = this.messageService.messageConfirmed$.subscribe(next => {
      this.userService.delete(next.data).then(res => this.refreshData());
    });
    this.refreshData();
  }

  ngOnDestroy(): void {
    if (this._deleteSubscription) {
      this._deleteSubscription.unsubscribe();
    }
  }
}
