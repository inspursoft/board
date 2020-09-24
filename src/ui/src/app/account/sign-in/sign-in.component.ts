import { Component, HostBinding, OnInit } from '@angular/core';
import { Router } from '@angular/router';
import { HttpErrorResponse } from '@angular/common/http';
import { TranslateService } from '@ngx-translate/core';
import { ResSignIn, ResSignInType, ReqSignIn } from '../account.types';
import { MessageService } from '../../shared.service/message.service';
import { AppInitService } from '../../shared.service/app-init.service';
import { AccountService } from '../account.service';
import { RouteDashboard, RouteForgotPassword, RouteSignUp } from '../../shared/shared.const';
import { AppTokenService } from '../../shared.service/app-token.service';
import { CsComponentBase } from '../../shared/cs-components-library/cs-component-base';

@Component({
  templateUrl: './sign-in.component.html',
  styleUrls: ['./sign-in.component.css']
})
export class SignInComponent extends CsComponentBase implements OnInit {
  @HostBinding('style.overflow-y') overflowY = 'hidden';
  isSignWIP = false;
  signInUser: ReqSignIn;
  authMode = '';
  redirectionURL = '';
  curErrorSignIn: ResSignIn;
  curShowVerifyCode = false;

  constructor(private appInitService: AppInitService,
              private messageService: MessageService,
              private accountService: AccountService,
              private appTokenService: AppTokenService,
              private translateService: TranslateService,
              private router: Router) {
    super();
    this.authMode = this.appInitService.systemInfo.authMode;
    this.redirectionURL = this.appInitService.systemInfo.redirectionUrl;
    this.curErrorSignIn = new ResSignIn();
    this.signInUser = new ReqSignIn();
  }

  ngOnInit(): void {
    if (this.authMode === 'indata_auth') {
      window.location.href = this.redirectionURL;
    }
  }

  get verifyPictureUrl() {
    return `/captcha/${this.signInUser.captchaId}.png`;
  }

  refreshVerifyPicture() {
    this.accountService.getCaptcha().subscribe((res: { captcha_id: string }) => this.signInUser.captchaId = res.captcha_id);
  }

  signIn(): void {
    if (this.verifyInputExValid()) {
      this.isSignWIP = true;
      this.accountService.signIn(this.signInUser).subscribe(res => {
        this.isSignWIP = false;
        this.messageService.showAlert('ACCOUNT.SUCCESS_TO_SIGN_IN');
        this.appTokenService.token = res.token;
        this.router.navigate([RouteDashboard], {queryParams: {token: this.appInitService.token}}).then();
      }, (err: HttpErrorResponse) => {
        this.isSignWIP = false;
        if (err.status === 400) {
          this.curErrorSignIn.assignFromRes(err.error);
          if (!this.curShowVerifyCode) {
            this.curShowVerifyCode = this.curErrorSignIn.retries >= 2;
          }
          switch (this.curErrorSignIn.type) {
            case ResSignInType.normal: {
              this.messageService.showOnlyOkDialog('ACCOUNT.INCORRECT_USERNAME_OR_PASSWORD', 'ACCOUNT.ERROR');
              if (this.curShowVerifyCode) {
                this.refreshVerifyPicture();
              }
              break;
            }
            case ResSignInType.overThreeTimes: {
              this.refreshVerifyPicture();
              this.messageService.showOnlyOkDialog('ACCOUNT.INCORRECT_VERIFY_CODE', 'ACCOUNT.ERROR');
              break;
            }
            case ResSignInType.temporarilyBlocked: {
              this.refreshVerifyPicture();
              this.translateService.get('ACCOUNT.USERNAME_TEMPORARY_BLOCKED',
                [this.signInUser.username, this.curErrorSignIn.value]).subscribe(msg =>
                this.messageService.showOnlyOkDialog(msg, 'ACCOUNT.ERROR')
              );
              break;
            }
          }
        } else if (err.status === 409) {
          this.messageService.showOnlyOkDialog('ACCOUNT.ALREADY_SIGNED_IN', 'ACCOUNT.ERROR');
        } else {
          this.messageService.showOnlyOkDialog('ACCOUNT.FAILED_TO_SIGN_IN', 'ACCOUNT.ERROR');
        }
      });
    }
  }

  signUp(): void {
    this.router.navigate([RouteSignUp]).then();
  }

  navigateForgotPassword(): void {
    this.router.navigate([RouteForgotPassword]).then();
  }
}
