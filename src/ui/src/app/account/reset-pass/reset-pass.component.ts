import { Component, OnInit } from '@angular/core';
import { AccountService } from "../account.service";
import { MessageService } from "../../shared/message-service/message.service";
import { ActivatedRoute, Router } from "@angular/router";
import { SignUp } from "../sign-up/sign-up";
import { BUTTON_STYLE, MESSAGE_TARGET } from "../../shared/shared.const";
import { Message } from "../../shared/message-service/message";
import { Subscription } from "rxjs/Subscription";
import { HttpErrorResponse } from "@angular/common/http";

@Component({
  selector: 'app-reset-pass',
  templateUrl: './reset-pass.component.html',
  styleUrls: ['./reset-pass.component.css']
})
export class ResetPassComponent implements OnInit {
  private resetUuid: string;
  private signUpModel: SignUp = new SignUp();
  private confirmSubscription:Subscription;

  constructor(
    private accountService: AccountService,
    private messageService: MessageService,
    private router: Router,
    private route: ActivatedRoute,
  ) {
    this.confirmSubscription = this.messageService.messageConfirmed$.subscribe((msg: Message) => {
      if (msg.target == MESSAGE_TARGET.RESET_PASS) {
        this.router.navigate(['/sign-in']);
      }
    });
  }

  ngOnInit() {
    this.route.queryParamMap.subscribe(params => {
      this.resetUuid = params.get("reset_uuid")
    });
  }

  sendResetPassRequest() {
    this.accountService.resetPass(this.signUpModel.password, this.resetUuid)
      .then(() => {
        let msg: Message = new Message();
        msg.title = "ACCOUNT.RESET_PASS_SUCCESS";
        msg.message = "ACCOUNT.RESET_PASS_SUCCESS_MSG";
        msg.buttons = BUTTON_STYLE.ONLY_CONFIRM;
        msg.target = MESSAGE_TARGET.RESET_PASS;
        this.messageService.announceMessage(msg);
      })
      .catch((err:HttpErrorResponse) => {
        let msg: Message = new Message();
        let status = err.status;
        let rtnErrorMessage = (err:HttpErrorResponse):string =>{
          if(/Invalid reset UUID/gm.test(err.error)){
            return "ACCOUNT.INVALID_RESET_UUID"
          }else{
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
