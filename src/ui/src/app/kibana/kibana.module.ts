import { NgModule } from "@angular/core";
import { SharedModule } from "../shared/shared.module";
import { KibanaComponent } from "./kibana/kibana.component";
import { KibanaService } from "./kibana.service";

@NgModule({
  imports: [SharedModule],
  declarations: [KibanaComponent],
  exports: [KibanaComponent],
  providers:[KibanaService]
})
export class KibanaModule {

}