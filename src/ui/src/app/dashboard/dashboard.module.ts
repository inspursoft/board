import { NgModule } from '@angular/core';
import { SharedModule } from '../shared/shared.module';
import { DashboardComponent } from './dashboard.component';
import { DashboardService } from "app/dashboard/dashboard.service";
import { TimeRangeScale } from "app/dashboard/time-range-scale.component/time-range-scale.component";
import { RouterModule } from "@angular/router";

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