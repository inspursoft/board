import { Injectable, Type } from '@angular/core';
import { ServiceStepBase } from './service-step';
import { ListServiceComponent } from './step0-list-service/list-service.component';
import { ChooseProjectComponent } from './step1-choose-project/choose-project.component';
import { ConfigContainerComponent } from "./step2-config-container/config-container.component";
import { ConfigSettingComponent } from "./step3-config-setting/config-setting.component";
import { TestingComponent } from "./step4-testing/testing.component";
import { DeployComponent } from "./step5-deploy/deploy.component";

@Injectable()
export class StepService {
  static getServiceSteps(): Array<Type<ServiceStepBase>> {
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