import { NgModule } from '@angular/core';
import { RouterModule, Routes } from '@angular/router';
import { CoreModule } from '../core/core.module';
import { SharedModule } from '../shared/shared.module';
import { NewEditUserComponent } from './user-center/user-new-edit/user-new-edit.component';
import { UserService } from './user-center/user-service/user-service';
import { UserListComponent } from './user-center/user-list/user-list.component';
import { SystemSettingComponent } from './system-setting/system-setting.component';
import { RouteSystemSetting, RouteUserCenters } from '../shared/shared.const';
import { AdminService } from './admin.service';

const routes: Routes = [
  {path: RouteUserCenters, component: UserListComponent},
  {path: RouteSystemSetting, component: SystemSettingComponent}
];

@NgModule({
  imports: [
    CoreModule,
    SharedModule,
    RouterModule.forChild(routes)
  ],
  providers: [
    AdminService,
    UserService
  ],
  declarations: [
    UserListComponent,
    NewEditUserComponent,
    SystemSettingComponent,
  ]
})
export class AdminModule {
}
