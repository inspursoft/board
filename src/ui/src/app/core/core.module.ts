import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpClientModule } from '@angular/common/http'
import { ClarityModule } from '@clr/angular';
import { RouterModule } from "@angular/router";
import { CommonModule } from "@angular/common";
import { TranslateModule } from "@ngx-translate/core";

@NgModule({
  exports:[
    CommonModule,
    HttpClientModule,
    FormsModule,
    ReactiveFormsModule,
    RouterModule,
    ClarityModule,
    TranslateModule
  ]
})
export class CoreModule {}
