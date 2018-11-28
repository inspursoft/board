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

@NgModule({
  imports: [
    SharedModule
  ],
  entryComponents: [
    CreatePvComponent,
    PvDetailComponent
  ],
  providers: [
    StorageService,
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