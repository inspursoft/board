import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { ClarityModule } from '@clr/angular';
import { RouterModule } from "@angular/router";
import { CommonModule } from "@angular/common";
import { TranslateModule } from "@ngx-translate/core";
import { NgxEchartsModule } from "ngx-echarts";
import { BoardComponentsLibraryModule } from "board-components-library";

@NgModule({
  exports:[
    CommonModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule,
    ClarityModule,
    TranslateModule,
    NgxEchartsModule,
    BoardComponentsLibraryModule
  ]
})
export class CoreModule {}
