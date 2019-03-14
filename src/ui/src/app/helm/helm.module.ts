import { NgModule } from "@angular/core";
import { SharedModule } from "../shared/shared.module";
import { RepoListComponent } from "./repo-list/repo-list.component";
import { HelmService } from "./helm.service";
import { HelmHostComponent } from "./helm-host/helm-host.component";
import { ChartListComponent } from "./chart-list/chart-list.component";
import { UploadChartComponent } from "./upload-chart/upload-chart.component";
import { ChartReleaseComponent } from "./chart-release/chart-release.component";
import { ChartReleaseListComponent } from "./chart-release-list/chart-release-list.component";
import { ChartReleaseDetailComponent } from "./chart-release-detail/chart-release-detail.component";

@NgModule({
  imports: [SharedModule],
  providers: [HelmService],
  declarations: [
    RepoListComponent,
    HelmHostComponent,
    ChartListComponent,
    UploadChartComponent,
    ChartReleaseComponent,
    ChartReleaseListComponent,
    ChartReleaseDetailComponent
  ],
  entryComponents: [
    RepoListComponent,
    ChartListComponent,
    UploadChartComponent,
    ChartReleaseComponent,
    ChartReleaseDetailComponent
  ]
})
export class HelmModule {

}