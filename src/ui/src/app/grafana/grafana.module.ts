import { NgModule } from "@angular/core";
import { SharedModule } from "../shared/shared.module";
import { GrafanaComponent } from "./grafana/grafana.component";


@NgModule({
  imports: [SharedModule],
  declarations: [GrafanaComponent],
  exports: [GrafanaComponent],
})
export class GrafanaModule {

}