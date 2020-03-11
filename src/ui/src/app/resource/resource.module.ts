import { NgModule } from "@angular/core";
import { RouterModule } from "@angular/router";
import { SharedModule } from "../shared/shared.module";
import { ConfigMapListComponent } from "./config-map/config-map-list/config-map-list.component";
import { ResourceService } from "./resource.service";
import { CreateConfigMapComponent } from "./config-map/create-config-map/create-config-map.component";
import { ConfigMapDetailComponent } from "./config-map/config-map-detail/config-map-detail.component";
import { ConfigMapUpdateComponent } from "./config-map/config-map-update/config-map-update.component";
import { CoreModule } from "../core/core.module";

@NgModule({
  imports: [
    CoreModule,
    SharedModule,
    RouterModule.forChild([{path: 'config-map', component: ConfigMapListComponent}])
  ],
  declarations: [
    ConfigMapListComponent,
    CreateConfigMapComponent,
    ConfigMapDetailComponent,
    ConfigMapUpdateComponent
  ],
  entryComponents: [
    CreateConfigMapComponent,
    ConfigMapDetailComponent,
    ConfigMapUpdateComponent
  ],
  providers: [
    ResourceService
  ]
})
export class ResourceModule {

}
