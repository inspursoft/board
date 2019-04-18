import { NgModule } from '@angular/core';
import { SharedModule } from "../shared/shared.module";
import { ListAuditComponent } from "./operation-audit-list/list-audit.component";
import { OperationAuditService } from "./audit-service";
import { CoreModule } from "../core/core.module";
import { HttpInterceptorService } from "../shared.service/http-client-interceptor";
import { RouterModule, Routes } from "@angular/router";

@NgModule({
  imports: [
    CoreModule,
    SharedModule,
    RouterModule.forChild([{path:'', component: ListAuditComponent}])
  ],
  providers: [
    OperationAuditService,
    HttpInterceptorService,
  ],
  declarations: [
    ListAuditComponent
  ]
})
export class AuditModule {
}
