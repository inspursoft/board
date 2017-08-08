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
import { AuthGuard } from './auth-guard.service';
import { CheckItemPatternDirective } from "./directives/check-item-pattern.directive";
import { ChangePasswordComponent } from "./change-password/change-password.component";

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
    ChangePasswordComponent
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
    CheckItemPatternDirective
  ],
  providers: [
    AuthGuard,
    MessageService
  ]
})
export class SharedModule {}
