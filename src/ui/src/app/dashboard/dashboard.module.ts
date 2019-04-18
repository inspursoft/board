import { NgModule } from '@angular/core';
import { SharedModule } from '../shared/shared.module';
import { DashboardComponent } from './dashboard.component';
import { TimeRangeScale } from './time-range-scale.component/time-range-scale.component';
import { DashboardService } from './dashboard.service';
import { CoreModule } from '../core/core.module';
import { RouterModule } from '@angular/router';
import { HttpInterceptorService } from "../shared.service/http-client-interceptor";
import { GrafanaComponent } from "./grafana/grafana.component";

@NgModule({
  imports: [
    CoreModule,
    SharedModule,
    RouterModule.forChild([{path: '', component: DashboardComponent}])
  ],
  declarations: [
    DashboardComponent,
    GrafanaComponent,
    TimeRangeScale
  ],
  providers: [
    DashboardService,
    HttpInterceptorService
  ]
})
export class DashboardModule {
}
