import { NgModule, NO_ERRORS_SCHEMA } from '@angular/core';
import { SharedModule } from '../shared/shared.module';
import { ProjectComponent } from './project.component';
import { CreateProjectComponent } from './create-project/create-project.component';

import { MemberComponent } from './member/member.component';
import { ProjectService } from './project.service';

@NgModule({
  imports: [
    SharedModule
  ],
  declarations: [ 
    ProjectComponent,
    CreateProjectComponent,
    MemberComponent
  ],

  providers: [
    ProjectService
  ]
})
export class ProjectModule {}