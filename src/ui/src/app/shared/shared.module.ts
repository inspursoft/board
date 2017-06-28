import { NgModule } from '@angular/core';

import { CoreModule } from '../core/core.module';

import { AngularEchartsModule } from 'angular2-echarts';
import { HeaderComponent } from './header/header.component';
import { ChartComponent } from './chart/chart.component';
import { ConfirmationDialogComponent } from './confirmation-dialog/confirmation-dialog.component';
import { TranslateModule } from '@ngx-translate/core';
import { ClarityModule } from 'clarity-angular';

import { MessageService } from './service/message.service';

@NgModule({
  imports: [
    CoreModule,
    AngularEchartsModule,
    TranslateModule
  ],
  declarations: [
    HeaderComponent,
    ChartComponent,
    ConfirmationDialogComponent  
  ],
  exports: [
    CoreModule,
    HeaderComponent,
    ChartComponent,
    ConfirmationDialogComponent,
    TranslateModule
  ],
  providers: [
    MessageService
  ]
})
export class SharedModule {}