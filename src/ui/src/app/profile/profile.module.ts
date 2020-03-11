import { NgModule } from '@angular/core';
import { ProfileComponent } from './profile.component';
import { CoreModule } from "../core/core.module";
import { RouterModule } from "@angular/router";

@NgModule({
  imports: [
    CoreModule,
    RouterModule.forChild([{path: '', component: ProfileComponent}])
  ],
  declarations: [
    ProfileComponent
  ]
})
export class ProfileModule {
}
