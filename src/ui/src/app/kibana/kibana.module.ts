import { NgModule } from "@angular/core";
import { SharedModule } from "../shared/shared.module";
import { KibanaComponent } from "./kibana/kibana.component";

@NgModule({
  imports: [SharedModule],
  declarations: [KibanaComponent],
  exports: [KibanaComponent],
})
export class KibanaModule {

}