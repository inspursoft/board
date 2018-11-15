import { NgModule } from '@angular/core';
import { SharedModule } from '../shared/shared.module';
import { DashboardComponent } from './dashboard.component';
import { TimeRangeScale } from "./time-range-scale.component/time-range-scale.component";
import { DashboardService } from "./dashboard.service";

@NgModule({
  imports: [
    SharedModule
  ],
  declarations: [
    DashboardComponent,
    TimeRangeScale
  ],
  providers: [
    DashboardService
  ]
})
export class DashboardModule {
}