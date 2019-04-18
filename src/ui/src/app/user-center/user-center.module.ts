import { NgModule } from '@angular/core';
import { UserList } from "./user-list/user-list.component";
import { SharedModule } from "../shared/shared.module";
import { UserService } from "./user-service/user-service";
import { UserCenterComponent } from './user-center.component';
import { NewEditUserComponent } from "./user-new-edit/user-new-edit.component";
import { CoreModule } from "../core/core.module";
import { RouterModule } from "@angular/router";
import { HttpInterceptorService } from "../shared.service/http-client-interceptor";

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
    UserService,
    HttpInterceptorService
  ],
})
export class UserCenterModule {
}
