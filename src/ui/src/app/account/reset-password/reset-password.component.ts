import { Component, OnInit, QueryList, ViewChildren } from '@angular/core';
import { AccountService } from "../account.service";
import { MessageService } from "../../shared/message-service/message.service";
import { ActivatedRoute, Router } from "@angular/router";
import { SignUp } from "../sign-up/sign-up";
import { BUTTON_STYLE, MESSAGE_TARGET } from "../../shared/shared.const";
import { Message } from "../../shared/message-service/message";
import { Subscription } from "rxjs/Subscription";
import { HttpErrorResponse } from "@angular/common/http";
import { ParamMap } from "@angular/router/src/shared";
import { AppInitService } from "../../app.init.service";
import { CsComponentBase } from "../../shared/cs-components-library/cs-component-base";

@Component({
  selector: 'reset-password',
  templateUrl: './reset-password.component.html',
  styleUrls: ['./reset-password.component.css']
})
export class ResetPasswordComponent extends CsComponentBase implements OnInit {
  resetUuid: string;
  signUpModel: SignUp = new SignUp();
  sendRequestWIP: boolean = false;
  private confirmSubscription: Subscription;

  constructor(private accountService: AccountService,
              private messageService: MessageService,
              private router: Router,
              private appInitService: AppInitService,
              private activatedRoute: ActivatedRoute) {
    super();
    this.confirmSubscription = this.messageService.messageConfirmed$.subscribe((msg: Message) => {
      if (msg.target == MESSAGE_TARGET.RESET_PASSWORD) {
        this.router.navigate(['/sign-in']);
      }
    });
  }

  ngOnInit() {
    if(this.appInitService.systemInfo["auth_mode"] != 'db_auth') {
      this.router.navigate(['/sign-in']);
    } else {
      this.activatedRoute.queryParamMap.subscribe((params: ParamMap) => this.resetUuid = params.get("reset_uuid"));
    }
  }

  goBack(){
    this.router.navigate(['/sign-in']);
  }


  sendResetPassRequest() {
    if (this.verifyInputValid()) {
      this.sendRequestWIP = true;
      this.accountService.resetPassword(this.signUpModel.password, this.resetUuid)
        .then(() => {
          this.sendRequestWIP = false;
          let msg: Message = new Message();
          msg.title = "ACCOUNT.RESET_PASS_SUCCESS";
          msg.message = "ACCOUNT.RESET_PASS_SUCCESS_MSG";
          msg.buttons = BUTTON_STYLE.ONLY_CONFIRM;
          msg.target = MESSAGE_TARGET.RESET_PASSWORD;
          this.messageService.announceMessage(msg);
        })
        .catch((err: HttpErrorResponse) => {
          this.sendRequestWIP = false;
          let msg: Message = new Message();
          let rtnErrorMessage = (err: HttpErrorResponse): string => {
            if (/Invalid reset UUID/gm.test(err.error)) {
              return "ACCOUNT.INVALID_RESET_UUID"
            } else {
              return "ACCOUNT.RESET_PASS_ERR_MSG"
            }
          };
          msg.title = "ACCOUNT.RESET_PASS_ERR";
          msg.message = rtnErrorMessage(err);
          msg.buttons = BUTTON_STYLE.ONLY_CONFIRM;
          this.messageService.announceMessage(msg);
        });
    }
  }
}
