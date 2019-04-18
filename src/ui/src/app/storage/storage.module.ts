import { NgModule } from "@angular/core";
import { StorageComponent } from "./storage.component";
import { SharedModule } from "../shared/shared.module";
import { CreatePvComponent } from "./pv/create-pv.component/create-pv.component";
import { StorageService } from "./storage.service";
import { MonitorsComponent } from "./pv/create-pv.component/monitors/monitors.component";
import { PvDetailComponent } from "./pv/pv-detail.compoent/pv-detail.component";
import { PvListComponent } from "./pv/pv-list.compoent/pv-list.component";
import { PvcListComponent } from "./pvc/pvc-list.component/pvc-list.component";
import { PvcDetailComponent } from "./pvc/pvc-detail.component/pvc-detail.component";
import { CoreModule } from "../core/core.module";
import { HttpInterceptorService } from "../shared.service/http-client-interceptor";
import { RouterModule } from "@angular/router";

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
    StorageService,
    HttpInterceptorService
  ],
  declarations: [
    StorageComponent,
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
