import {Component, OnInit} from '@angular/core';
import {AccountService} from "../account.service";
import {MessageService} from "../../shared/message-service/message.service";
import {ActivatedRoute, Router} from "@angular/router";
import {SignUp} from "../sign-up/sign-up";
import {BUTTON_STYLE, MESSAGE_TARGET} from "../../shared/shared.const";
import {Message} from "../../shared/message-service/message";
import {Subscription} from "rxjs/Subscription";
import {HttpErrorResponse} from "@angular/common/http";
import {AppInitService} from "../../app.init.service";
import {ParamMap} from "@angular/router/src/shared";

@Component({
  selector: 'reset-password',
  templateUrl: './reset-password.component.html',
  styleUrls: ['./reset-password.component.css']
})
export class ResetPasswordComponent implements OnInit {
  resetUuid: string;
  signUpModel: SignUp = new SignUp();
  private confirmSubscription: Subscription;

  constructor(
    private accountService: AccountService,
    private messageService: MessageService,
    private router: Router,
    private appInitService: AppInitService,
    private activatedRoute: ActivatedRoute) {
    this.appInitService.systemInfo = this.activatedRoute.snapshot.data['systeminfo'];
    this.confirmSubscription = this.messageService.messageConfirmed$.subscribe((msg: Message) => {
      if (msg.target == MESSAGE_TARGET.RESET_PASSWORD) {
        this.router.navigate(['/sign-in']);
      }
    });
  }

  ngOnInit() {
    this.activatedRoute.queryParamMap.subscribe((params: ParamMap) => this.resetUuid = params.get("reset_uuid"));
  }

  goBack(){
    this.router.navigate(['/sign-in']);
  }

  sendResetPassRequest() {
    this.accountService.resetPassword(this.signUpModel.password, this.resetUuid)
      .then(() => {
        let msg: Message = new Message();
        msg.title = "ACCOUNT.RESET_PASS_SUCCESS";
        msg.message = "ACCOUNT.RESET_PASS_SUCCESS_MSG";
        msg.buttons = BUTTON_STYLE.ONLY_CONFIRM;
        msg.target = MESSAGE_TARGET.RESET_PASSWORD;
        this.messageService.announceMessage(msg);
      })
      .catch((err: HttpErrorResponse) => {
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
