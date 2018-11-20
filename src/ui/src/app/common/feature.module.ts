import { NgModule } from '@angular/core';
import { GlobalSearchModule } from '../global-search/global-search.module';
import { DashboardModule } from '../dashboard/dashboard.module';
import { NodeModule } from '../node/node.module';
import { ServiceModule } from '../service/service.module';
import { ProjectModule } from '../project/project.module';
import { ImageModule } from '../image/image.module';
import { ProfileModule } from '../profile/profile.module';
import { UserCenterModule } from '../user-center/user-center.module';
import { AuditModule } from "../audit/audit.module";
import { AccountModule } from "../account/account.module";
import { MainContentModule } from "../main-content/main-content.module";
import { KibanaModule } from "../kibana/kibana.module";
import { GrafanaModule } from "../grafana/grafana.module";
import { StorageModule } from "../storage/storage.module";

@NgModule({
  exports: [
    MainContentModule,
    GlobalSearchModule,
    DashboardModule,
    NodeModule,
    ServiceModule,
    ProjectModule,
    ImageModule,
    ProfileModule,
    UserCenterModule,
    AuditModule,
    AccountModule,
    KibanaModule,
    GrafanaModule,
    StorageModule
  ]
})
export class FeatureModule {}