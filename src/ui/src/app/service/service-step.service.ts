import { Injectable, Type } from '@angular/core';
import { ListServiceComponent } from './step0-list-service/list-service.component';
import { ChooseProjectComponent } from './step1-choose-project/choose-project.component';
import { ConfigContainerComponent } from './step2-config-container/config-container.component';
import { ConfigSettingComponent } from './step3-config-setting/config-setting.component';
import { TestingComponent } from './step4-testing/testing.component';
import { DeployComponent } from './step5-deploy/deploy.component';
import { ServiceStepComponentBase } from './service-step';

@Injectable()
export class StepService {
  static getServiceSteps(): Array<Type<ServiceStepComponentBase>> {
    return [
      ListServiceComponent,
      ChooseProjectComponent,
      ConfigContainerComponent,
      ConfigSettingComponent,
      TestingComponent,
      DeployComponent
    ];
  }
}
