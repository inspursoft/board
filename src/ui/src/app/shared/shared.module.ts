import { NgModule } from '@angular/core';

import { AngularEchartsModule } from 'angular2-echarts';
import { HeaderComponent } from './header/header.component';
import { ChartComponent } from './chart/chart.component';
import { TranslateModule } from '@ngx-translate/core';
import { ClarityModule } from 'clarity-angular';


@NgModule({
  imports: [
    ClarityModule,
    AngularEchartsModule,
    TranslateModule
  ],
  declarations: [
    HeaderComponent,
    ChartComponent  
  ],
  exports: [
    HeaderComponent,
    ChartComponent,
    TranslateModule,
    ClarityModule
  ]
})
export class SharedModule {}