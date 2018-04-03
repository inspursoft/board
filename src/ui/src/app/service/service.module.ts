import { NgModule } from '@angular/core';
import { SharedModule } from '../shared/shared.module';
import { ServiceComponent } from './service.component';
import { ServiceHostDirective } from './service-host.directive';
import { StepService } from './service-step.service';
import { K8sService } from './service.k8s';
import { ListServiceComponent } from './step0-list-service/list-service.component';
import { ChooseProjectComponent } from './step1-choose-project/choose-project.component';
import { SelectImageComponent } from './step2-select-image/select-image.component';
import { EditContainerComponent } from './step3-edit-container/edit-container.component';
import { ConfigSettingComponent } from './step4-config-setting/config-setting.component';
import { TestingComponent } from './step5-testing/testing.component';
import { ServiceWizardComponent } from "./service-wizard/service-wizard.component";
import { VolumeMountsComponent } from "./step3-edit-container/volume-mounts/volume-mounts.component";
import { DynamicFormModule } from '../shared/dynamic-form/dynamic-form.module';
import { ServiceDetailComponent } from './step0-list-service/service-detail/service-detail.component';
import { ServiceControlComponent } from "./step0-list-service/service-control/service-control.component";
import { ServiceCreateYamlComponent } from './step0-list-service/service-create-yaml/service-create-yaml.component';
import { DeployComponent } from "./step6-deploy/deploy.component";

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
    SelectImageComponent,
    EditContainerComponent,
    ConfigSettingComponent,
    TestingComponent,
    DeployComponent,
    ServiceWizardComponent,
    VolumeMountsComponent,
    ServiceDetailComponent,
    ServiceControlComponent,
    ServiceCreateYamlComponent
  ],
  entryComponents: [
    ListServiceComponent,
    ChooseProjectComponent,
    SelectImageComponent,
    EditContainerComponent,
    ConfigSettingComponent,
    TestingComponent,
    DeployComponent
  ],

  providers: [
    K8sService,
    StepService
  ]
})
export class ServiceModule {
}