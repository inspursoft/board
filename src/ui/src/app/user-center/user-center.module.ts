import { NgModule, NO_ERRORS_SCHEMA } from '@angular/core';

import { UserList } from "./user-list/user-list.component";
import { SharedModule } from "../shared/shared.module";
import { UserService } from "./user-service/user-service";
import { UserCenterComponent } from './user-center.component';
import { NewEditUserComponent } from "./user-new-edit/user-new-edit.component";
import { MessageService } from "../shared/message-service/message.service";

@NgModule({
  imports: [ 
    SharedModule 
  ],
  declarations: [
    UserCenterComponent,
    NewEditUserComponent,
    UserList
  ],

  providers: [
    UserService
  ],
})
export class UserCenterModule { }
