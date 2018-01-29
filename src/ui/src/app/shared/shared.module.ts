import { NgModule } from '@angular/core';
import { CoreModule } from '../core/core.module';
import { AngularEchartsModule } from 'angular2-echarts';
import { ChartComponent } from './chart/chart.component';
import { HeaderComponent } from './header/header.component';
import { ConfirmationDialogComponent } from './confirmation-dialog/confirmation-dialog.component';
import { InlineAlertComponent } from './inline-alert/inline-alert.component';
import { GlobalMessageComponent } from './global-message/global-message.component';
import { CheckItemExistingDirective } from './directives/check-item-existing.directive';
import { CheckItemIdenticalDirective } from './directives/check-item-identical.directive';

import { MessageService } from './message-service/message.service';
import { AuthGuard, ServiceGuard } from './auth-guard.service';
import { CheckItemPatternDirective } from "./directives/check-item-pattern.directive";
import { ChangePasswordComponent } from "./change-password/change-password.component";
import { AccountSettingComponent } from "./account-setting/account-setting.component";

import { ValidateOnBlurDirective } from './directives/validate-onblur.directive';
import { CsDropdownComponent } from "./cs-components-library/cs-dropdown/cs-dropdown.component";
import { WebsocketService } from './websocket-service/websocket.service';
import { CsSearchInput } from "./cs-components-library/cs-search-input/cs-search-input.component";
import { CheckboxRevert } from "./directives/checkbox-revert.directive";
import { CsInputComponent } from "./cs-components-library/cs-input/cs-input.component";
import { CsInputArrayComponent } from "./cs-components-library/cs-input-array/cs-input-array.component";
import { CreateImageComponent } from "../image/image-create/image-create.component";
import { EnvironmentValueComponent } from "./environment-value/environment-value.component";
import { SizePipe } from "./pipes/size-pipe";
import { CsGuideComponent } from "./cs-components-library/cs-guide/cs-guide.component";

@NgModule({
  imports: [
    CoreModule,
    AngularEchartsModule
  ],
  declarations: [
    ChartComponent,
    HeaderComponent, 
    ConfirmationDialogComponent,
    InlineAlertComponent,
    GlobalMessageComponent,
    CheckItemExistingDirective,
    CheckItemIdenticalDirective,
    CheckItemPatternDirective,
    ValidateOnBlurDirective,
    ChangePasswordComponent,
    CsDropdownComponent,
    AccountSettingComponent,
    CreateImageComponent,
    EnvironmentValueComponent,
    CsSearchInput,
    CsInputComponent,
    CsInputArrayComponent,
    CheckboxRevert,
    SizePipe,
    CsGuideComponent

  ],
  exports: [
    CoreModule,
    AngularEchartsModule,
    ChartComponent,
    HeaderComponent,
    ConfirmationDialogComponent,
    InlineAlertComponent,
    GlobalMessageComponent,
    CheckItemExistingDirective,
    CheckItemIdenticalDirective,
    CsDropdownComponent,
    CheckItemPatternDirective,
    ValidateOnBlurDirective,
    CreateImageComponent,
    EnvironmentValueComponent,
    CsSearchInput,
    CsInputComponent,
    CsInputArrayComponent,
    CsGuideComponent,
    CheckboxRevert,
    SizePipe
  ],
  providers: [
    AuthGuard,
    ServiceGuard,
    MessageService,
    WebsocketService
  ]
})
export class SharedModule {}
