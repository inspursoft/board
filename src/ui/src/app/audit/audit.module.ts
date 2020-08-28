import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { SharedModule } from '../shared/shared.module';
import { ListAuditComponent } from './operation-audit-list/list-audit.component';
import { OperationAuditService } from './audit-service';
import { CoreModule } from '../core/core.module';

@NgModule({
  imports: [
    CoreModule,
    SharedModule,
    RouterModule.forChild([{path: '', component: ListAuditComponent}])
  ],
  providers: [
    OperationAuditService
  ],
  declarations: [
    ListAuditComponent
  ]
})
export class AuditModule {
}
