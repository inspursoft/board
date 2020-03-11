import { NgModule } from '@angular/core';
import { SharedModule } from '../shared/shared.module';
import { SignInComponent } from './sign-in/sign-in.component';
import { SignUpComponent } from './sign-up/sign-up.component';
import { AccountService } from './account.service';
import { ResetPasswordComponent } from './reset-password/reset-password.component';
import { ForgotPasswordComponent } from './forgot-password/forgot-password.component';
import { CoreModule } from "../core/core.module";
import { RouterModule, Routes } from "@angular/router";
import { SystemInfoResolve } from "../app-routing.module";

const routes: Routes = [
  {
    path: 'sign-in',
    component: SignInComponent,
    resolve: {
      systeminfo: SystemInfoResolve
    }
  },
  {
    path: 'sign-up',
    component: SignUpComponent,
    resolve: {
      systeminfo: SystemInfoResolve
    }
  },
  {
    path: 'reset-password',
    component: ResetPasswordComponent,
    resolve: {
      systeminfo: SystemInfoResolve
    }
  },
  {
    path: 'forgot-password',
    component: ForgotPasswordComponent,
    resolve: {
      systeminfo: SystemInfoResolve
    }
  }];

@NgModule({
  imports: [
    CoreModule,
    SharedModule,
    RouterModule.forChild(routes)
  ],
  providers: [
    AccountService
  ],
  declarations: [
    SignInComponent,
    SignUpComponent,
    ResetPasswordComponent,
    ForgotPasswordComponent
  ]
})
export class AccountModule {
}
