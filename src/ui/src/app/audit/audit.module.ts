import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { SharedModule } from "../shared/shared.module";
import { ListAuditComponent } from "./operation-audit-list/list-audit.component";
import { OperationAuditService } from "./audit-service";

@NgModule({
  imports: [
    CommonModule,
    SharedModule
  ],
  providers: [OperationAuditService],
  declarations: [
    ListAuditComponent
  ]
})
export class AuditModule {
}
