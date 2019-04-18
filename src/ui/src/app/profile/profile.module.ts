import { NgModule } from '@angular/core';
import { ProfileComponent } from './profile.component';
import { CoreModule } from "../core/core.module";
import { RouterModule } from "@angular/router";
import { HttpInterceptorService } from "../shared.service/http-client-interceptor";

@NgModule({
  imports: [
    CoreModule,
    RouterModule.forChild([{path: '', component: ProfileComponent}])
  ],
  declarations: [
    ProfileComponent
  ],
  providers: [
    HttpInterceptorService
  ]
})
export class ProfileModule {
}
