import { NgModule } from "@angular/core";
import { RouterModule } from "@angular/router";
import { CoreModule } from "../core/core.module";
import { SharedModule } from "../shared/shared.module";
import { JobListComponent } from "./job-list/job-list.component";
import { JobService } from "./job.service";
import { JobCreateComponent } from "./job-create/job-create.component";
import { JobContainerCreateComponent } from './job-container-create/job-container-create.component';
import { JobContainerConfigComponent } from './job-container-config/job-container-config.component';
import { JobVolumeMountsComponent } from "./job-volume-mounts/job-volume-mounts.component";
import { JobAffinityComponent } from "./job-affinity/job-affinity.component";
import { JobAffinityCardComponent } from "./job-affinity-card/job-affinity-card.component";
import { JobAffinityCardListComponent } from "./job-affinity-card-list/job-affinity-card-list.component";
import { JobDetailComponent } from './job-detail/job-detail.component';
import { JobLogsComponent } from './job-logs/job-logs.component';

@NgModule({
  imports: [
    CoreModule,
    SharedModule,
    RouterModule.forChild([{path: '', component: JobListComponent}])
  ],
  providers: [JobService],
  entryComponents: [
    JobContainerCreateComponent,
    JobContainerConfigComponent,
    JobVolumeMountsComponent,
    JobAffinityComponent,
    JobAffinityCardComponent,
    JobAffinityCardListComponent,
    JobDetailComponent,
    JobLogsComponent
  ],
  declarations: [
    JobListComponent,
    JobCreateComponent,
    JobContainerCreateComponent,
    JobContainerConfigComponent,
    JobVolumeMountsComponent,
    JobAffinityComponent,
    JobAffinityCardComponent,
    JobAffinityCardListComponent,
    JobDetailComponent,
    JobLogsComponent
  ]
})
export class JobModule {

}
