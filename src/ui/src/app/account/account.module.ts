import { NgModule } from '@angular/core';
import { SharedModule } from '../shared/shared.module';
import { SignInComponent } from './sign-in/sign-in.component';
import { HeaderComponent } from '../shared/header/header.component';
@NgModule({
  imports: [
    SharedModule
  ],
  declarations: [ 
    SignInComponent
  ]
})
export class AccountModule {}