import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { ListAuditComponent } from "./step0-list-audit/list-audit.component";
import { SharedModule } from "../shared/shared.module";

@NgModule({
  imports: [
    CommonModule,
    SharedModule
  ],
  declarations: [
    ListAuditComponent
  ]
})
export class AuditModule { }
