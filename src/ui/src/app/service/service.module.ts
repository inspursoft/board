import { NgModule } from '@angular/core';
import { SharedModule } from '../shared/shared.module';
import { ServiceComponent } from './service.component';
import { ServiceHostDirective } from './service-host.directive';
import { StepService } from './service-step.service';
import { K8sService } from './service.k8s';
import { ListServiceComponent } from './step0-list-service/list-service.component';
import { ChooseProjectComponent } from './step1-choose-project/choose-project.component';
import { ServiceWizardComponent } from "./service-wizard/service-wizard.component";
import { VolumeMountsComponent } from "./step2-config-container/volume-mounts/volume-mounts.component";
import { DynamicFormModule } from '../shared/dynamic-form/dynamic-form.module';
import { ServiceDetailComponent } from './step0-list-service/service-detail/service-detail.component';
import { ServiceControlComponent } from "./step0-list-service/service-control/service-control.component";
import { ServiceCreateYamlComponent } from './step0-list-service/service-create-yaml/service-create-yaml.component';
import { ScaleComponent } from './step0-list-service/service-control/scale/scale.component';
import { UpdateComponent } from './step0-list-service/service-control/update/update.component';
import { LocateComponent } from './step0-list-service/service-control/locate/locate.component';
import { StatusComponent } from './step0-list-service/service-control/status/status.component';
import { ConfigSettingComponent } from "./step3-config-setting/config-setting.component";
import { TestingComponent } from "./step4-testing/testing.component";
import { DeployComponent } from "./step5-deploy/deploy.component";
import { ConfigContainerComponent } from "./step2-config-container/config-container.component";
import { SetExternalComponent } from "./step3-config-setting/set-external-port/set-external.component";
import { ConfigCardComponent } from "./step3-config-setting/config-card/config-card.component";
import { ConfigCardListComponent } from "./step3-config-setting/config-card-list/config-card-list.component";


@NgModule({
  imports: [
    SharedModule,
    DynamicFormModule
  ],
  declarations: [
    ServiceComponent,
    ServiceHostDirective,
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
    ConfigCardComponent,
    ConfigCardListComponent,
    SetExternalComponent,
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
    SetExternalComponent
  ],
  providers: [
    K8sService,
    StepService
  ]
})
export class ServiceModule {
}