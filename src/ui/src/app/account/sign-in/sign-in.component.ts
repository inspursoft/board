import { Component, OnDestroy, OnInit } from '@angular/core';
import { ActivatedRoute, Router } from '@angular/router';
import { SignIn } from './sign-in';
import { Message } from '../../shared/message-service/message';
import { MessageService } from '../../shared/message-service/message.service';

import { Subscription } from 'rxjs/Subscription';
import { AppInitService } from '../../app.init.service';
import { AccountService } from '../account.service';
import { BUTTON_STYLE, MESSAGE_TARGET } from "../../shared/shared.const";
import { HttpErrorResponse } from "@angular/common/http";

@Component({
  templateUrl: './sign-in.component.html',
  styleUrls: [ './sign-in.component.css' ]
})
export class SignInComponent implements OnInit, OnDestroy {
  isSignWIP: boolean = false;
  signInUser: SignIn = new SignIn();
  authMode: string = '';
  redirectionURL: string = '';

  _subscription: Subscription;

  constructor(
    private appInitService: AppInitService,
    private messageService: MessageService,
    private accountService: AccountService,
    private router: Router,
    private route: ActivatedRoute
  ) {
    this._subscription = this.messageService.messageConfirmed$.subscribe((msg: Message) => {
      if (msg.target == MESSAGE_TARGET.SIGN_IN_ERROR) {
        console.error('Received:' + JSON.stringify(msg.message));
      }
    });
    this.appInitService.systemInfo = this.route.snapshot.data['systeminfo'];
    this.authMode = this.appInitService.systemInfo['auth_mode'];
    this.redirectionURL = this.appInitService.systemInfo['redirection_url'];
  }

  ngOnInit(): void {
    if(this.authMode === 'indata_auth') {
      window.location.href = this.redirectionURL;
    }
  }

  signIn(): void {
    this.isSignWIP = true;
    this.accountService
      .signIn(this.signInUser.username, this.signInUser.password)
      .then(res=>{
          this.isSignWIP = false;
          this.appInitService.token = res.token;
          this.router.navigate(['/dashboard'], { queryParams: { token: this.appInitService.token }});
      })
      .catch((err: HttpErrorResponse) => {
        this.isSignWIP = false;
        let announceMessage = new Message();
        announceMessage.title = 'ACCOUNT.ERROR';
        announceMessage.target = MESSAGE_TARGET.SIGN_IN_ERROR;
        announceMessage.buttons = BUTTON_STYLE.ONLY_CONFIRM;
        if(err) {
          switch(err.status){
            case 400: {
              announceMessage.message = 'ACCOUNT.INCORRECT_USERNAME_OR_PASSWORD';
              break;
            }
            case 409: {
              announceMessage.message = 'ACCOUNT.ALREADY_SIGNED_IN';
              break;
            }
            default: {
              announceMessage.message = 'ACCOUNT.FAILED_TO_SIGN_IN';
            }
          }
        }
        this.messageService.announceMessage(announceMessage);
      });
  }

  signUp(): void {
    this.router.navigate(['/sign-up']);
  }
  resetPass(): void {
    this.router.navigate(['/reset-password']);
  }
  retrievePass(): void {
    this.router.navigate(['/retrieve-pass']);
  }

  ngOnDestroy(): void {
    if(this._subscription) {
      this._subscription.unsubscribe();
    }
  }
}