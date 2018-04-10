import { NgModule } from '@angular/core';
import { GlobalSearchModule } from '../global-search/global-search.module';
import { DashboardModule } from '../dashboard/dashboard.module';
import { NodeModule } from '../node/node.module';
import { ServiceModule } from '../service/service.module';
import { ProjectModule } from '../project/project.module';
import { ImageModule } from '../image/image.module';
import { ProfileModule } from '../profile/profile.module';
import { UserCenterModule } from '../user-center/user-center.module';

@NgModule({
  exports: [
    GlobalSearchModule,
    DashboardModule,
    NodeModule,
    ServiceModule,
    ProjectModule,
    ImageModule,
    ProfileModule,
    UserCenterModule
  ]
})
export class FeatureModule {}