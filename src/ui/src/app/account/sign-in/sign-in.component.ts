import { Component, OnDestroy } from '@angular/core';
import { Router } from '@angular/router';
import { SignIn } from './sign-in';
import { Message } from '../../shared/message-service/message';
import { MessageService } from '../../shared/message-service/message.service';

import { Subscription } from 'rxjs/Subscription';
import { AppInitService } from '../../app.init.service';
import { AccountService } from '../account.service';

@Component({
  templateUrl: './sign-in.component.html',
  styleUrls: [ './sign-in.component.css' ]
})
export class SignInComponent implements OnDestroy {

  signInUser: SignIn = new SignIn();
  
  _subscription: Subscription;

  constructor(
    private appInitService: AppInitService,
    private messageService: MessageService, 
    private accountService: AccountService,
    private router: Router) {
    this._subscription = this.messageService.messageConfirmed$.subscribe((message: any)=>{
      let confirmationMessage = <Message>message;
      console.error('Received:' + JSON.stringify(confirmationMessage));
    })
  }

  signIn(): void {
    this.accountService
      .signIn(this.signInUser.username, this.signInUser.password)
      .then(res=>{
          this.appInitService.token = res.token;
          this.router.navigate(['/dashboard']);
      })
      .catch(err=>{
        let announceMessage = new Message();
          announceMessage.title = 'ACCOUNT.ERROR';
        if(err && err.status === 400) {
          announceMessage.message = 'ACCOUNT.INCORRECT_USERNAME_OR_PASSWORD';
        } else {
          announceMessage.message = 'ACCOUNT.FAILED_TO_SIGN_IN' + (err && err.status);
        }
        this.messageService.announceMessage(announceMessage);
      });
  }

  signUp(): void {
    this.router.navigate(['/sign-up']);
  }

  ngOnDestroy(): void {
    if(this._subscription) {
      this._subscription.unsubscribe();
    }
  }
}