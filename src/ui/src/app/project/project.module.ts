import { NgModule } from '@angular/core';
import { SharedModule } from '../shared/shared.module';
import { ProjectComponent } from './project.component';
import { ProjectService } from './project.service';

@NgModule({
  imports: [
    SharedModule
  ],
  declarations: [ 
    ProjectComponent
  ],
  providers: [
    ProjectService
  ]
})
export class ProjectModule {}