import { Injectable } from '@angular/core';

import { ServiceStep } from './service-step';

import { ListServiceComponent } from './step0-list-service/list-service.component';
import { ChooseProjectComponent } from './step1-choose-project/choose-project.component';
import { SelectImageComponent } from './step2-select-image/select-image.component';
import { EditContainerComponent } from './step3-edit-container/edit-container.component';
import { ConfigSettingComponent } from './step4-config-setting/config-setting.component';
import { TestingComponent } from './step5-testing/testing.component';
import { DeployComponent } from "./step6-deploy/deploy.component";

@Injectable()
export class StepService {
  getServiceSteps(): ServiceStep[] {
    return [
      new ServiceStep(ListServiceComponent,   {}),
      new ServiceStep(ChooseProjectComponent, { 'title': 'Choose Project' }),
      new ServiceStep(SelectImageComponent,   { 'title': 'Select Image'   }),
      new ServiceStep(EditContainerComponent, { 'title': 'Edit Container' }),
      new ServiceStep(ConfigSettingComponent, { 'title': 'Config Setting' }),
      new ServiceStep(TestingComponent, { 'title': 'Deploy Testing' }),
      new ServiceStep(DeployComponent, { 'title': 'Deploy Testing' })
    ];
  }
}