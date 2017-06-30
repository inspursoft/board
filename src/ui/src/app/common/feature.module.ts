import { NgModule } from '@angular/core';

import { DashboardModule } from '../dashboard/dashboard.module';
import { NodeModule } from '../node/node.module';
import { ServiceModule } from '../service/service.module';
import { ProjectModule } from '../project/project.module';
import { ImageModule } from '../image/image.module';
import { AdminOptionModule } from '../admin-option/admin-option.module';
import { ProfileModule } from '../profile/profile.module';

@NgModule({
  exports: [

    DashboardModule,
    NodeModule,
    ServiceModule,
    ProjectModule,
    ImageModule,
    AdminOptionModule,
    ProfileModule
  ]
})
export class FeatureModule {}