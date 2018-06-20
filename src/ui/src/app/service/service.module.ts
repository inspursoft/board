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
import { ScaleComponent } from './step0-list-service/service-control/scale/scale.component';
import { UpdateComponent } from './step0-list-service/service-control/update/update.component';
import { LocateComponent } from './step0-list-service/service-control/locate/locate.component';
import { StatusComponent } from './step0-list-service/service-control/status/status.component';
import { CsSyntaxHighlighterComponent } from "../shared/cs-components-library/cs-syntax-highlighter/cs-syntax-highlighter.component";
import { ApPrismModule } from '@angular-package/prism';

@NgModule({
  imports: [
    SharedModule,
    DynamicFormModule,
    ApPrismModule
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
    ServiceCreateYamlComponent,
    ScaleComponent,
    UpdateComponent,
    LocateComponent,
    StatusComponent
  ],
  entryComponents: [
    ListServiceComponent,
    ChooseProjectComponent,
    SelectImageComponent,
    EditContainerComponent,
    ConfigSettingComponent,
    TestingComponent,
    ServiceControlComponent,
    DeployComponent,
    ServiceDetailComponent,
    CsSyntaxHighlighterComponent
  ],
  providers: [
    K8sService,
    StepService
  ]
})
export class ServiceModule {
}