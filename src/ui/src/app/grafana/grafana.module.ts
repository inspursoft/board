import { NgModule } from "@angular/core";
import { SharedModule } from "../shared/shared.module";
import { GrafanaComponent } from "./grafana/grafana.component";
import { GrafanaService } from "./grafana.service";


@NgModule({
  imports: [SharedModule],
  declarations: [GrafanaComponent],
  exports: [GrafanaComponent],
  providers:[GrafanaService]
})
export class GrafanaModule {

}