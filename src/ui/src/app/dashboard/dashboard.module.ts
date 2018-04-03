import { NgModule } from '@angular/core';
import { RouterModule } from "@angular/router";
import { SharedModule } from '../shared/shared.module';
import { DashboardComponent } from './dashboard.component';
import { TimeRangeScale } from "./time-range-scale.component/time-range-scale.component";
import { DashboardService } from "./dashboard.service";

@NgModule({
  imports: [
    SharedModule,
    RouterModule
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