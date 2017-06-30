import { NgModule } from '@angular/core';

import { SharedModule } from '../shared/shared.module';

import { DashboardComponent } from './dashboard.component';
import {AngularEchartsModule} from "angular2-echarts";
import {CommonModule} from "@angular/common";
import {CoreModule} from "../core/core.module";

@NgModule({
  imports: [
    CommonModule,
    CoreModule,
    AngularEchartsModule,
    SharedModule
  ],
  declarations:[ 
    DashboardComponent
  ]
})
export class DashboardModule {}