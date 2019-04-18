import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http'
import { ClarityModule } from '@clr/angular';
import { RouterModule } from "@angular/router";
import { CommonModule } from "@angular/common";
import { TranslateModule } from "@ngx-translate/core";
import { NgxEchartsModule } from "ngx-echarts";

@NgModule({
  exports:[
    CommonModule,
    HttpClientModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule,
    ClarityModule,
    TranslateModule,
    NgxEchartsModule
  ]
})
export class CoreModule {}
