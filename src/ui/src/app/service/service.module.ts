import { NgModule } from '@angular/core';
import { SharedModule } from '../shared/shared.module';

import { ServiceComponent } from './service.component';
import { ServiceGroupComponent } from './service-group.component';
import { ServiceHostDirective } from './service-host.directive';
import { StepService } from './service-step.service';

import { K8sService } from './service.k8s';

import { ListServiceComponent } from './step0-list-service/list-service.component';
import { ChooseProjectComponent } from './step1-choose-project/choose-project.component';
import { SelectImageComponent } from './step2-select-image/select-image.component';
import { EditContainerComponent } from './step3-edit-container/edit-container.component';
import { ConfigSettingComponent } from './step4-config-setting/config-setting.component';
import { DeployTestingComponent } from './step5-deploy-testing/deploy-testing.component';

@NgModule({
  imports: [
    SharedModule
  ],
  declarations: [ 
    ServiceComponent,
    ServiceGroupComponent,
    ServiceHostDirective,
    ListServiceComponent,
    ChooseProjectComponent,
    SelectImageComponent,
    EditContainerComponent,
    ConfigSettingComponent,
    DeployTestingComponent
  ],
  entryComponents: [
    ListServiceComponent,
    ChooseProjectComponent,
    SelectImageComponent,
    EditContainerComponent,
    ConfigSettingComponent,
    DeployTestingComponent
  ],
  providers: [
    K8sService,
    StepService
  ]
})
export class ServiceModule {}