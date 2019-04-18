import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { SharedModule } from '../shared/shared.module';
import { ProjectComponent } from './project.component';
import { ProjectService } from './project.service';
import { CoreModule } from '../core/core.module';
import { HttpInterceptorService } from "../shared.service/http-client-interceptor";

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
    ProjectService,
    HttpInterceptorService
  ]
})
export class ProjectModule {
}
