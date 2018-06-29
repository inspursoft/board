import { Component, OnDestroy, OnInit } from '@angular/core';
import { MessageService } from "../../shared/message-service/message.service";
import { Router } from "@angular/router";
import { AccountService } from "../account.service";
import { Message } from "../../shared/message-service/message";
import { BUTTON_STYLE, MESSAGE_TARGET } from "../../shared/shared.const";
import { Subscription } from "rxjs/Subscription";
import { HttpErrorResponse } from "@angular/common/http";

@Component({
  selector: 'app-retrieve-pass',
  templateUrl: './retrieve-pass.component.html',
  styleUrls: ['./retrieve-pass.component.css']
})
export class RetrievePassComponent implements OnInit ,OnDestroy{
  private credential: string;
  protected confirmSubscription: Subscription;
  constructor(
    private accountService: AccountService,
    private messageService: MessageService,
    private router: Router
  ) {
    if (this.confirmSubscription) {
      this.confirmSubscription.unsubscribe();
    }
    this.confirmSubscription = this.messageService.messageConfirmed$.subscribe((msg: Message) => {
      if (msg.target == MESSAGE_TARGET.RETRIEVE_PASS) {
        this.router.navigate(['/sign-in']);
      }
    });
  }

  ngOnInit() {
  }

  goBack(): void {
    this.router.navigate(['/sign-in']);
  }

  sendRequest(): void {
    this.accountService.retrieveEmail(this.credential)
      .then(res=>{
        let msg:Message = new Message();
        msg.title = "ACCOUNT.SEND_REQUEST_SUCCESS";
        msg.message = "ACCOUNT.SEND_REQUEST_SUCCESS_MSG";
        msg.buttons = BUTTON_STYLE.ONLY_CONFIRM;
        msg.target = MESSAGE_TARGET.RETRIEVE_PASS;
        this.messageService.announceMessage(msg);
      })
      .catch((err:HttpErrorResponse)=>{
        let rtnMessage = function (err:HttpErrorResponse) {
          if(err.status == 404){
            return "ACCOUNT.USER_NOT_EXISTS"
          }else{
            return "ACCOUNT.SEND_REQUEST_ERR_MSG"
          }
        };
        let msg:Message = new Message();
        msg.title = "ACCOUNT.SEND_REQUEST_ERR";
        msg.message = rtnMessage(err);
        msg.buttons = BUTTON_STYLE.ONLY_CONFIRM;
        this.messageService.announceMessage(msg);
      });
  }

  ngOnDestroy(): void {
    if (this.confirmSubscription) {
      this.confirmSubscription.unsubscribe();
    }
  }
}
