import { NgModule } from '@angular/core';
import { SharedModule } from '../shared/shared.module';
import { SignInComponent } from './sign-in/sign-in.component';
import { SignUpComponent } from './sign-up/sign-up.component';
import { AccountService } from './account.service';
import { ResetPassComponent } from './reset-pass/reset-pass.component';
import { RetrievePassComponent } from './retrieve-pass/retrieve-pass.component';

@NgModule({
  imports: [SharedModule],
  declarations: [ 
    SignInComponent,
    SignUpComponent,
    ResetPassComponent,
    RetrievePassComponent
  ],
  providers: [
    AccountService
  ]
})
export class AccountModule {}