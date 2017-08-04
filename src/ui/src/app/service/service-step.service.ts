import { Injectable } from '@angular/core';

import { ServiceStep } from './service-step';

import { ListServiceComponent } from './step0-list-service/list-service.component';
import { ChooseProjectComponent } from './step1-choose-project/choose-project.component';
import { SelectImageComponent } from './step2-select-image/select-image.component';
import { EditContainerComponent } from './step3-edit-container/edit-container.component';
import { ConfigSettingComponent } from './step4-config-setting/config-setting.component';
import { DeployTestingComponent } from './step5-deploy-testing/deploy-testing.component';

@Injectable()
export class StepService {
  getServiceSteps(): ServiceStep[] {
    return [
      new ServiceStep(ListServiceComponent,   {}),
      new ServiceStep(ChooseProjectComponent, { 'title': 'Choose Project' }),
      new ServiceStep(SelectImageComponent,   { 'title': 'Select Image'   }),
      new ServiceStep(EditContainerComponent, { 'title': 'Edit Container' }),
      new ServiceStep(ConfigSettingComponent, { 'title': 'Config Setting' }),
      new ServiceStep(DeployTestingComponent, { 'title': 'Deploy Testing' })
    ];
  }
}