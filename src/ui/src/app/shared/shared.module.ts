import { NgModule } from '@angular/core';
import { CoreModule } from '../core/core.module';
import { AngularEchartsModule } from 'angular2-echarts';
import { ChartComponent } from './chart/chart.component';
import { HeaderComponent } from './header/header.component';
import { ConfirmationDialogComponent } from './confirmation-dialog/confirmation-dialog.component';
import { InlineAlertComponent } from './inline-alert/inline-alert.component';

import { CheckItemExistingDirective } from './directives/check-item-existing.directive';
import { CheckItemIdenticalDirective } from './directives/check-item-identical.directive';

import { MessageService } from './message-service/message.service';

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
    CheckItemExistingDirective,
    CheckItemIdenticalDirective
  ],
  exports: [
    CoreModule,
    AngularEchartsModule,
    ChartComponent,
    HeaderComponent,
    ConfirmationDialogComponent,
    InlineAlertComponent,
    CheckItemExistingDirective,
    CheckItemIdenticalDirective
  ],
  providers: [
    MessageService
  ]
})
export class SharedModule {}