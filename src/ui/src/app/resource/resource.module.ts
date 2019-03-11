import { NgModule } from "@angular/core";
import { SharedModule } from "../shared/shared.module";
import { ConfigMapListComponent } from "./config-map/config-map-list/config-map-list.component";
import { ResourceService } from "./resource.service";
import { CreateConfigMapComponent } from "./config-map/create-config-map/create-config-map.component";
import { ConfigMapDetailComponent } from "./config-map/config-map-detail/config-map-detail.component";

@NgModule({
  imports: [SharedModule],
  declarations: [
    ConfigMapListComponent,
    CreateConfigMapComponent,
    ConfigMapDetailComponent
  ],
  entryComponents: [
    CreateConfigMapComponent,
    ConfigMapDetailComponent
  ],
  providers: [
    ResourceService
  ]
})
export class ResourceModule {

}