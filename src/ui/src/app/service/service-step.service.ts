import { Injectable, Type } from '@angular/core';
import { ServiceStepBase } from './service-step';
import { ListServiceComponent } from './step0-list-service/list-service.component';
import { ChooseProjectComponent } from './step1-choose-project/choose-project.component';
import { SelectImageComponent } from './step2-select-image/select-image.component';
import { EditContainerComponent } from './step3-edit-container/edit-container.component';
import { ConfigSettingComponent } from './step4-config-setting/config-setting.component';
import { TestingComponent } from './step5-testing/testing.component';


@Injectable()
export class StepService {
  static getServiceSteps(): Array<Type<ServiceStepBase>> {
    return [
      ListServiceComponent,
      ChooseProjectComponent,
      SelectImageComponent,
      EditContainerComponent,
      ConfigSettingComponent,
      TestingComponent
    ];
  }
}