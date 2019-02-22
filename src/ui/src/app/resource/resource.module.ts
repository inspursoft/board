import { NgModule } from "@angular/core";
import { SharedModule } from "../shared/shared.module";
import { ConfigMapListComponent } from "./config-map/config-map-list/config-map-list.component";
import { ResourceService } from "./resource.service";

@NgModule({
  imports: [SharedModule],
  declarations: [
    ConfigMapListComponent
  ],
  entryComponents:[

  ],
  providers:[
   ResourceService
  ]
})
export class ResourceModule {

}