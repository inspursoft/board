import { NgModule } from "@angular/core";
import { RouterModule } from "@angular/router";
import { KibanaComponent } from "./kibana/kibana.component";
import { KibanaService } from "./kibana.service";
import { CoreModule } from "../core/core.module";
import { SharedModule } from "../shared/shared.module";

@NgModule({
  imports: [
    CoreModule,
    SharedModule,
    RouterModule.forChild([{path: '', component: KibanaComponent}])
  ],
  declarations: [
    KibanaComponent
  ],
  exports: [
    KibanaComponent
  ],
  providers: [
    KibanaService
  ]
})
export class KibanaModule {

}
