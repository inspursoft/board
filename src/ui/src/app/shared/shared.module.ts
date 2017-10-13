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
import { CsDropdownComponent } from "../service/cs-dropdown/cs-dropdown.component";

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
    AccountSettingComponent
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
    ValidateOnBlurDirective
  ],
  providers: [
    AuthGuard,
    ServiceGuard,
    MessageService
  ]
})
export class SharedModule {}
