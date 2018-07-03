import { Component, OnDestroy, OnInit } from '@angular/core';
import { MessageService } from "../../shared/message-service/message.service";
import { ActivatedRoute, Router } from "@angular/router";
import { AccountService } from "../account.service";
import { Message } from "../../shared/message-service/message";
import { BUTTON_STYLE, MESSAGE_TARGET } from "../../shared/shared.const";
import { Subscription } from "rxjs/Subscription";
import { HttpErrorResponse } from "@angular/common/http";
import { AppInitService } from "../../app.init.service";

@Component({
  selector: 'app-retrieve-pass',
  templateUrl: './retrieve-pass.component.html',
  styleUrls: ['./retrieve-pass.component.css']
})
export class RetrievePassComponent implements OnInit, OnDestroy {
  private credential: string = "";
  private sendRequestWIP: boolean = false;
  protected confirmSubscription: Subscription;

  constructor(
    private accountService: AccountService,
    private messageService: MessageService,
    private appInitService: AppInitService,
    private activatedRoute: ActivatedRoute,
    private router: Router) {
    this.appInitService.systemInfo = this.activatedRoute.snapshot.data['systeminfo'];
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
    this.sendRequestWIP = true;
    this.accountService.retrieveEmail(this.credential)
      .then(() => {
        this.sendRequestWIP = false;
        let msg: Message = new Message();
        msg.title = "ACCOUNT.SEND_REQUEST_SUCCESS";
        msg.message = "ACCOUNT.SEND_REQUEST_SUCCESS_MSG";
        msg.buttons = BUTTON_STYLE.ONLY_CONFIRM;
        msg.target = MESSAGE_TARGET.RETRIEVE_PASS;
        this.messageService.announceMessage(msg);
      })
      .catch((err: HttpErrorResponse) => {
        this.sendRequestWIP = false;
        let msg: Message = new Message();
        msg.title = "ACCOUNT.SEND_REQUEST_ERR";
        msg.message = err.status == 404 ? "ACCOUNT.USER_NOT_EXISTS" : "ACCOUNT.SEND_REQUEST_ERR_MSG";
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
