import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { SharedModule } from '../shared/shared.module';
import { ServiceComponent } from './service.component';
import { StepService } from './service-step.service';
import { K8sService } from './service.k8s';
import { ListServiceComponent } from './step0-list-service/list-service.component';
import { ChooseProjectComponent } from './step1-choose-project/choose-project.component';
import { ServiceWizardComponent } from './service-wizard/service-wizard.component';
import { VolumeMountsComponent } from './step2-config-container/volume-mounts/volume-mounts.component';
import { ServiceDetailComponent } from './step0-list-service/service-detail/service-detail.component';
import { ServiceControlComponent } from './step0-list-service/service-control/service-control.component';
import { ServiceCreateYamlComponent } from './step0-list-service/service-create-yaml/service-create-yaml.component';
import { ScaleComponent } from './step0-list-service/service-control/scale/scale.component';
import { UpdateComponent } from './step0-list-service/service-control/update/update.component';
import { LocateComponent } from './step0-list-service/service-control/locate/locate.component';
import { StatusComponent } from './step0-list-service/service-control/status/status.component';
import { ConfigSettingComponent } from './step3-config-setting/config-setting.component';
import { TestingComponent } from './step4-testing/testing.component';
import { DeployComponent } from './step5-deploy/deploy.component';
import { ConfigContainerComponent } from './step2-config-container/config-container.component';
import { SetAffinityComponent } from './step3-config-setting/set-affinity/set-affinity.component';
import { AffinityCardComponent } from './step3-config-setting/affinity-card/affinity-card.component';
import { AffinityCardListComponent } from './step3-config-setting/affinity-card-list/affinity-card-list.component';
import { LoadBalanceComponent } from './step0-list-service/service-control/loadBalance/loadBalance.component';
import { CoreModule } from '../core/core.module';
import { ServiceGuard } from './service-guard.service';
import { ConfigParamsComponent } from './step2-config-container/config-params/config-params.component';
import { ConsoleComponent } from './step0-list-service/service-control/console/console.component';

@NgModule({
  imports: [
    CoreModule,
    SharedModule,
    RouterModule.forChild([{path: '', component: ServiceComponent, canDeactivate: [ServiceGuard]}])
  ],
  declarations: [
    ServiceComponent,
    ListServiceComponent,
    ChooseProjectComponent,
    ConfigContainerComponent,
    ConfigSettingComponent,
    TestingComponent,
    DeployComponent,
    ServiceWizardComponent,
    VolumeMountsComponent,
    ServiceDetailComponent,
    ServiceControlComponent,
    ServiceCreateYamlComponent,
    ScaleComponent,
    UpdateComponent,
    LocateComponent,
    StatusComponent,
    LoadBalanceComponent,
    AffinityCardComponent,
    AffinityCardListComponent,
    SetAffinityComponent,
    ConfigParamsComponent,
    ConsoleComponent,
  ],
  entryComponents: [
    ListServiceComponent,
    ChooseProjectComponent,
    ConfigContainerComponent,
    ConfigSettingComponent,
    TestingComponent,
    ServiceControlComponent,
    DeployComponent,
    ServiceDetailComponent,
    SetAffinityComponent,
    VolumeMountsComponent,
    ConfigParamsComponent
  ],
  providers: [
    K8sService,
    StepService,
    ServiceGuard
  ]
})
export class ServiceModule {
}
