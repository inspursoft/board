import { Component, OnInit } from '@angular/core';
import { ClrDatagridSortOrder, ClrDatagridStateInterface } from '@clr/angular';
import { ActivatedRoute } from '@angular/router';
import { TranslateService } from '@ngx-translate/core';
import { HttpErrorResponse } from '@angular/common/http';
import { UserService } from '../user-service/user-service';
import { editModel } from '../user-new-edit/user-new-edit.component';
import { AppInitService } from '../../../shared.service/app-init.service';
import { MessageService } from '../../../shared.service/message.service';
import { Message, RETURN_STATUS, User } from '../../../shared/shared.types';
import { UserPagination } from '../../admin.types';

@Component({
  selector: 'user-list',
  templateUrl: './user-list.component.html',
  styleUrls: ['./user-list.component.css']
})

export class UserListComponent implements OnInit {
  userListData: UserPagination;
  userListErrMsg = '';
  curUser: User;
  curEditModel = editModel.emNew;
  showNewUser = false;
  setUserSystemAdminWIP = false;
  isInLoading = false;
  totalRecordCount = 0;
  pageIndex = 1;
  pageSize = 15;
  authMode = '';
  currentUserID = 0;
  descSort = ClrDatagridSortOrder.DESC;
  oldStateInfo: ClrDatagridStateInterface;

  constructor(private route: ActivatedRoute,
              private userService: UserService,
              private appInitService: AppInitService,
              private translateService: TranslateService,
              private messageService: MessageService) {
    this.userListData = new UserPagination();
  }

  ngOnInit() {
    this.authMode = this.appInitService.systemInfo.authMode;
    this.currentUserID = this.appInitService.currentUser.userId;
  }

  refreshData(stateInfo: ClrDatagridStateInterface): void {
    if (stateInfo) {
      setTimeout(() => {
        this.isInLoading = true;
        this.oldStateInfo = stateInfo;
        this.userService.getUserList('', this.pageIndex, this.pageSize, stateInfo.sort.by as string, stateInfo.sort.reverse).subscribe(
          (res: UserPagination) => {
            this.userListData = res;
            this.totalRecordCount = res.pagination.TotalCount;
            this.isInLoading = false;
          },
          () => this.isInLoading = false
        );
      });
    }
  }

  addUser() {
    this.curUser = new User();
    this.curEditModel = editModel.emNew;
    this.showNewUser = true;
  }

  editUser(user: User) {
    if (user.userDeleted !== 1 && user.userId !== 1 && user.userId !== this.currentUserID) {
      this.curEditModel = editModel.emEdit;
      this.userService.getUser(user.userId).subscribe(res => {
        this.curUser = res;
        this.showNewUser = true;
      });
    }
  }

  deleteUser(user: User) {
    if (user.userDeleted !== 1 && user.userId !== 1 && user.userId !== this.currentUserID) {
      this.translateService.get('USER_CENTER.CONFIRM_DELETE_USER', [user.userName]).subscribe((res: string) => {
        this.messageService.showDeleteDialog(res, 'USER_CENTER.DELETE_USER').subscribe((message: Message) => {
          if (message.returnStatus === RETURN_STATUS.rsConfirm) {
            this.userService.deleteUser(user).subscribe(() => {
              this.refreshData(this.oldStateInfo);
              this.messageService.showAlert('USER_CENTER.DELETE_USER_SUCCESS');
            }, (error: HttpErrorResponse) => {
              if (error.status === 422) {
                this.messageService.showAlert('USER_CENTER.DELETE_USER_ERROR', {alertType: 'danger'});
              }
            });
          }
        });
      });
    }
  }

  setUserSystemAdmin(user: User, $event: MouseEvent) {
    this.setUserSystemAdminWIP = true;
    const oldUserSystemAdmin = user.userSystemAdmin;
    this.userService.setUserSystemAdmin(user.userId, oldUserSystemAdmin === 1 ? 0 : 1).subscribe(
      () => {
        this.setUserSystemAdminWIP = false;
        user.userSystemAdmin = oldUserSystemAdmin === 1 ? 0 : 1;
        if (user.userSystemAdmin === 1) {
          this.messageService.showAlert('USER_CENTER.SUCCESSFUL_SET_SYS_ADMIN');
        } else {
          this.messageService.showAlert('USER_CENTER.SUCCESSFUL_SET_NOT_SYS_ADMIN');
        }
      },
      () => {
        this.setUserSystemAdminWIP = false;
        ($event.srcElement as HTMLInputElement).checked = oldUserSystemAdmin === 1;
      });
  }
}
