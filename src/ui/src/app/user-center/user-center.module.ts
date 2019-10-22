import { NgModule } from '@angular/core';
import { RouterModule } from "@angular/router";
import { UserList } from "./user-list/user-list.component";
import { SharedModule } from "../shared/shared.module";
import { UserService } from "./user-service/user-service";
import { UserCenterComponent } from './user-center.component';
import { NewEditUserComponent } from "./user-new-edit/user-new-edit.component";
import { CoreModule } from "../core/core.module";

@NgModule({
  imports: [
    CoreModule,
    SharedModule,
    RouterModule.forChild([{path: '', component: UserList}])
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
export class UserCenterModule {
}
