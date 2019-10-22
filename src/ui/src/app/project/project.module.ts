import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { SharedModule } from '../shared/shared.module';
import { ProjectComponent } from './project.component';
import { ProjectService } from './project.service';
import { CoreModule } from '../core/core.module';

@NgModule({
  imports: [
    CoreModule,
    SharedModule,
    RouterModule.forChild([{path: '', component: ProjectComponent}])
  ],
  declarations: [
    ProjectComponent
  ],
  providers: [
    ProjectService
  ]
})
export class ProjectModule {
}
