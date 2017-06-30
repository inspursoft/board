import { Component, OnDestroy } from '@angular/core';
import { Router } from '@angular/router';
import { SignIn } from './sign-in';
import { ConfirmationMessage } from '../../shared/service/confirmation-message';
import { MessageService } from '../../shared/service/message.service';

import { Subscription } from 'rxjs/Subscription';

import { AccountService } from '../account.service';

@Component({
  templateUrl: './sign-in.component.html',
  styleUrls: [ './sign-in.component.css' ]
})
export class SignInComponent implements OnDestroy {

  signInUser: SignIn = new SignIn();
  
  _subscription: Subscription;

  constructor(
    private messageService: MessageService, 
    private accountService: AccountService,
    private router: Router) {
    this._subscription = this.messageService.messageConfirmed$.subscribe((message: any)=>{
      let confirm: ConfirmationMessage = <ConfirmationMessage>message;
      console.error('Received:' + JSON.stringify(confirm));
    })
  }

  signIn(): void {
    this.accountService
      .signIn(this.signInUser.username, this.signInUser.password)
      .then(res=>{
          this.router.navigate(['/dashboard']);
      })
      .catch(err=>{
        let m: ConfirmationMessage = new ConfirmationMessage();
        m.title = 'Error';
        if(err && err.status === 400) {
          m.message = 'Incorrect username or password';
        } else {
          m.message = 'Sign in failed:' + (err && err.status);
        }
        this.messageService.announceMessage(m);
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