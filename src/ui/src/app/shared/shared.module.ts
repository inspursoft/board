import { NgModule } from '@angular/core';
import { CoreModule } from '../core/core.module';
import { AngularEchartsModule } from 'angular2-echarts';
import { ChartComponent } from './chart/chart.component';
import { HeaderComponent } from './header/header.component';
import { ConfirmationDialogComponent } from './confirmation-dialog/confirmation-dialog.component';

import { CheckItemExistingDirective } from './directives/check-item-existing.directive';
import { CheckItemIdenticalDirective } from './directives/check-item-identical.directive';
import { CheckItemPatternDirective } from "app/shared/directives/check-item-pattern.directive";
import { MessageService } from './service/message.service';
import { FixAlert } from "app/shared/FixAlert/fix-alert.component";

@NgModule({
  imports: [
    CoreModule,
    AngularEchartsModule
  ],
  declarations: [
    ChartComponent,
    HeaderComponent, 
    ConfirmationDialogComponent,
    CheckItemExistingDirective,
    CheckItemIdenticalDirective,
    CheckItemPatternDirective,
    FixAlert
  ],
  exports: [
    CoreModule,
    AngularEchartsModule,
    ChartComponent,
    HeaderComponent,
    ConfirmationDialogComponent,
    CheckItemExistingDirective,
    CheckItemIdenticalDirective,
    CheckItemPatternDirective,
    FixAlert
  ],
  providers: [
    MessageService
  ]
})
export class SharedModule {}