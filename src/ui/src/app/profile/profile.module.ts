import { NgModule, NO_ERRORS_SCHEMA } from '@angular/core';
import { ProfileComponent } from './profile.component';
import { SharedModule } from "../shared/shared.module";

@NgModule({
  imports: [ SharedModule ],

  declarations: [
    ProfileComponent
  ]
})
export class ProfileModule { }
