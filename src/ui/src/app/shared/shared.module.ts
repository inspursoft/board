import { NgModule } from '@angular/core';

import { AngularEchartsModule } from 'angular2-echarts';

import { ChartComponent } from './chart/chart.component';
import { TranslateModule } from '@ngx-translate/core';

@NgModule({
  imports: [
    AngularEchartsModule,
    TranslateModule
  ],
  declarations: [
    ChartComponent  
  ],
  exports: [
    ChartComponent,
    TranslateModule
  ]
})
export class SharedModule {}