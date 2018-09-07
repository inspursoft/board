import { Component, OnDestroy, OnInit } from "@angular/core";
import { ClrDatagridSortOrder, ClrDatagridStateInterface } from "@clr/angular";
import { ActivatedRoute } from "@angular/router";
import { Subscription } from "rxjs/Subscription";
import { UserService } from "../user-service/user-service";
import { User } from "../user";
import { editModel } from "../user-new-edit/user-new-edit.component"
import { AppInitService } from "../../app.init.service";
import { MessageService } from "../../shared/message-service/message.service";
import { Message, RETURN_STATUS } from "../../shared/shared.types";
import { TranslateService } from "@ngx-translate/core";

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
  setUserSystemAdminWIP: boolean = false;
  isInLoading: boolean = false;
  totalRecordCount: number;
  pageIndex: number = 1;
  pageSize: number = 15;
  authMode: string = '';
  currentUserID: number;
  descSort = ClrDatagridSortOrder.DESC;
  oldStateInfo: ClrDatagridStateInterface;

  constructor(private route: ActivatedRoute,
              private userService: UserService,
              private appInitService: AppInitService,
              private translateService: TranslateService,
              private messageService: MessageService) {
      this.authMode = this.appInitService.systemInfo['auth_mode'];
  }

  ngOnInit() {
    this.currentUserID = this.appInitService.currentUser["user_id"];
  }

  ngOnDestroy(): void {
    if (this._deleteSubscription) {
      this._deleteSubscription.unsubscribe();
    }
  }

  refreshData(stateInfo: ClrDatagridStateInterface): void {
    if (stateInfo) {
      setTimeout(()=>{
        this.isInLoading = true;
        this.oldStateInfo = stateInfo;
        this.userService.getUserList('', this.pageIndex, this.pageSize, stateInfo.sort.by as string, stateInfo.sort.reverse)
          .then(res => {
            this.totalRecordCount = res["pagination"]["total_count"];
            this.userListData = res["user_list"];
            this.isInLoading = false;
          })
          .catch(() => this.isInLoading = false);
      });
    }
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
    }
  }

  deleteUser(user: User) {
    if (user.user_deleted != 1 && user.user_id != 1 && user.user_id != this.currentUserID) {
      this.translateService.get('USER_CENTER.CONFIRM_DELETE_USER', [user.user_name]).subscribe((res: string) => {
        this.messageService.showDeleteDialog(res, 'USER_CENTER.DELETE_USER').subscribe((message: Message) => {
          if (message.returnStatus == RETURN_STATUS.rsConfirm) {
            this.userService.deleteUser(user).then(() => {
              this.refreshData(this.oldStateInfo);
              this.messageService.showAlert('USER_CENTER.DELETE_USER_SUCCESS');
            })
          }
        })
      });
    }
  }

  setUserSystemAdmin(user: User, $event:MouseEvent) {
    this.setUserSystemAdminWIP = true;
    let oldUserSystemAdmin = user.user_system_admin;
    this.userService.setUserSystemAdmin(user.user_id, oldUserSystemAdmin == 1 ? 0 : 1)
      .then(() => {
        this.setUserSystemAdminWIP = false;
        user.user_system_admin = oldUserSystemAdmin == 1 ? 0 : 1;
        if (user.user_system_admin === 1) {
          this.messageService.showAlert('USER_CENTER.SUCCESSFUL_SET_SYS_ADMIN');
          user.user_project_admin = 1;
        } else {
          this.messageService.showAlert('USER_CENTER.SUCCESSFUL_SET_NOT_SYS_ADMIN');
        }
      })
      .catch(() => {
        this.setUserSystemAdminWIP = false;
        ($event.srcElement as HTMLInputElement).checked = oldUserSystemAdmin == 1;
      })
  }
}
