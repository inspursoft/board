import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { SharedModule } from '../shared/shared.module';
import { CreatePvComponent } from './pv/create-pv/create-pv.component';
import { StorageService } from './storage.service';
import { MonitorsComponent } from './pv/create-pv/monitors/monitors.component';
import { PvDetailComponent } from './pv/pv-detail/pv-detail.component';
import { PvListComponent } from './pv/pv-list/pv-list.component';
import { PvcListComponent } from './pvc/pvc-list/pvc-list.component';
import { PvcDetailComponent } from './pvc/pvc-detail/pvc-detail.component';
import { CoreModule } from '../core/core.module';

@NgModule({
  imports: [
    CoreModule,
    SharedModule,
    RouterModule.forChild([
      {path: 'pv', component: PvListComponent},
      {path: 'pvc', component: PvcListComponent}
    ])
  ],
  entryComponents: [
    CreatePvComponent,
    PvDetailComponent,
    PvcDetailComponent
  ],
  providers: [
    StorageService
  ],
  declarations: [
    PvListComponent,
    PvcListComponent,
    CreatePvComponent,
    PvDetailComponent,
    PvcDetailComponent,
    MonitorsComponent
  ]
})
export class StorageModule {

}
