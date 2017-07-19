import { NgModule } from '@angular/core';

import { ProfileComponent } from './profile.component';
import { UserList } from "app/profile/user-center/user-list/user-list.component";
import { SharedModule } from "app/shared/shared.module";
import { Pagination } from "app/profile/user-center/pagination/pagination.component";
import { UserService } from "app/profile/user-center/user-service/user-service";
import { NewUser } from "app/profile/user-center/user-new-edit/user-new-edit.component";
import { MessageService } from "app/shared/message-service/message.service";

@NgModule({
  imports: [SharedModule],
  providers: [
    UserService,
    MessageService],
  declarations: [
    ProfileComponent,
    NewUser,
    UserList,
    Pagination]
})
export class ProfileModule { }
