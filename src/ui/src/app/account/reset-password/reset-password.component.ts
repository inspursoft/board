import { Component, OnInit } from '@angular/core';
import { AccountService } from "../account.service";
import { MessageService } from "../../shared.service/message.service";
import { ActivatedRoute, Router } from "@angular/router";
import { RouteSignIn } from "../../shared/shared.const";
import { HttpErrorResponse } from "@angular/common/http";
import { ParamMap } from "@angular/router/src/shared";
import { AppInitService } from "../../shared.service/app-init.service";
import { CsComponentBase } from "../../shared/cs-components-library/cs-component-base";
import { SignUp } from '../account.types';

@Component({
  selector: 'reset-password',
  templateUrl: './reset-password.component.html',
  styleUrls: ['./reset-password.component.css']
})
export class ResetPasswordComponent extends CsComponentBase implements OnInit {
  resetUuid: string;
  signUpModel: SignUp = new SignUp();
  sendRequestWIP = false;

  constructor(private accountService: AccountService,
              private messageService: MessageService,
              private router: Router,
              private appInitService: AppInitService,
              private activatedRoute: ActivatedRoute) {
    super();
  }

  ngOnInit() {
    if (this.appInitService.systemInfo.auth_mode != 'db_auth') {
      this.router.navigate([RouteSignIn]).then();
    } else {
      this.activatedRoute.queryParamMap.subscribe((params: ParamMap) => this.resetUuid = params.get("reset_uuid"));
    }
  }

  goBack() {
    this.router.navigate([RouteSignIn]).then();
  }

  sendResetPassRequest() {
    if (this.verifyInputExValid()) {
      this.sendRequestWIP = true;
      this.accountService.resetPassword(this.signUpModel.password, this.resetUuid).subscribe(
        () => this.messageService.showOnlyOkDialogObservable('ACCOUNT.RESET_PASS_SUCCESS_MSG', 'ACCOUNT.RESET_PASS_SUCCESS').subscribe(
          () => this.router.navigate([RouteSignIn]).then())
        , (err: HttpErrorResponse) => {
          this.sendRequestWIP = false;
          let rtnErrorMessage = (err: HttpErrorResponse): string => {
            if (/Invalid reset UUID/gm.test(err.error)) {
              return "ACCOUNT.INVALID_RESET_UUID"
            } else {
              return "ACCOUNT.RESET_PASS_ERR_MSG"
            }
          };
          let message = rtnErrorMessage(err);
          this.messageService.showOnlyOkDialog(message, 'ACCOUNT.RESET_PASS_ERR');
        });
    }
  }
}
