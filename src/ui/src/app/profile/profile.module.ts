import { NgModule } from '@angular/core';

import { ProfileComponent } from './profile.component';
import { SharedModule } from "app/shared/shared.module";

@NgModule({
  imports: [ SharedModule ],
  declarations: [
    ProfileComponent
  ]
})
export class ProfileModule { }
