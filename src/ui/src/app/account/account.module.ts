import { NgModule } from '@angular/core';
import { SharedModule } from '../shared/shared.module';
import { SignInComponent } from './sign-in/sign-in.component';
import { SignUpComponent } from './sign-up/sign-up.component';
import { AccountService } from './account.service';

@NgModule({
  imports: [
    SharedModule
  ],

  declarations: [ 
    SignInComponent,
    SignUpComponent
  ],
  providers: [
    AccountService
  ]
})
export class AccountModule {}