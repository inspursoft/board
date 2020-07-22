import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { ProfileComponent } from './profile.component';
import { CoreModule } from '../core/core.module';

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
