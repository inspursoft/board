import { Component, EventEmitter, Input, OnInit, Output } from '@angular/core';
import { UserService } from '../../admin/user-center/user-service/user-service';
import { User } from '../shared.types';
import { AppInitService } from '../../shared.service/app-init.service';
import { MessageService } from '../../shared.service/message.service';

@Component({
  selector: 'app-user-setting',
  templateUrl: './account-setting.component.html',
  styleUrls: ['./account-setting.component.css'],
  providers: [UserService]
})
export class AccountSettingComponent implements OnInit {

  constructor(private appInitService: AppInitService,
              private userService: UserService,
              private messageService: MessageService) {
  }

  @Input()
  get isOpen() {
    return this.isOpenValue;
  }

  set isOpen(open: boolean) {
    this.isOpenValue = open;
    this.isOpenChange.emit(this.isOpenValue);
  }
  isOpenValue = false;
  isWorkWip = false;
  curUser: User = new User();

  @Output() isOpenChange: EventEmitter<boolean> = new EventEmitter<boolean>();

  ngOnInit() {
    this.curUser = this.appInitService.currentUser;
  }

  submitAccountSetting() {
    this.isWorkWip = true;
    this.userService.usesChangeAccount(this.curUser).subscribe(
      () => this.isOpen = false,
      () => this.isOpen = false,
      () => this.messageService.showAlert('ACCOUNT.ACCOUNT_SETTING_SUCCESS'));
  }
}
