import { Component, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { SignIn } from './sign-in';
import { MessageService } from '../../shared.service/message.service';
import { AppInitService } from '../../shared.service/app-init.service';
import { AccountService } from '../account.service';
import { HttpErrorResponse } from '@angular/common/http';
import { RouteDashboard, RouteForgotPassword, RouteSignUp } from '../../shared/shared.const';
import { AppTokenService } from "../../shared.service/app-token.service";

@Component({
  templateUrl: './sign-in.component.html',
  styleUrls: [ './sign-in.component.css' ]
})
export class SignInComponent implements OnInit {
  isSignWIP = false;
  signInUser: SignIn = new SignIn();
  authMode = '';
  redirectionURL = '';

  constructor(private appInitService: AppInitService,
              private messageService: MessageService,
              private accountService: AccountService,
              private appTokenService: AppTokenService,
              private router: Router) {
    this.authMode = this.appInitService.systemInfo.auth_mode;
    this.redirectionURL = this.appInitService.systemInfo.redirection_url;
  }

  ngOnInit(): void {
    if (this.authMode === 'indata_auth') {
      window.location.href = this.redirectionURL;
    }
  }

  signIn(): void {
    this.isSignWIP = true;
    this.accountService.signIn(this.signInUser.username, this.signInUser.password).subscribe(res => {
      this.isSignWIP = false;
      this.messageService.showAlert('ACCOUNT.SUCCESS_TO_SIGN_IN');
      this.appTokenService.token = res.token;
      this.router.navigate([RouteDashboard], { queryParams: { token: this.appInitService.token }}).then();
    }, (err: HttpErrorResponse) => {
      this.isSignWIP = false;
      if (err.status === 400) {
        this.messageService.showOnlyOkDialog('ACCOUNT.INCORRECT_USERNAME_OR_PASSWORD', 'ACCOUNT.ERROR');
      } else if (err.status === 409) {
        this.messageService.showOnlyOkDialog('ACCOUNT.ALREADY_SIGNED_IN', 'ACCOUNT.ERROR');
      } else {
        this.messageService.showOnlyOkDialog('ACCOUNT.FAILED_TO_SIGN_IN', 'ACCOUNT.ERROR');
      }
    });
  }

  signUp(): void {
    this.router.navigate([RouteSignUp]).then();
  }

  navigateForgotPassword(): void {
    this.router.navigate([RouteForgotPassword]).then();
  }
}
