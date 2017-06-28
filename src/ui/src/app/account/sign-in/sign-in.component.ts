import { Component, ViewChild, OnDestroy } from '@angular/core';

import { SignIn } from './sign-in';
import { ConfirmationMessage } from '../../shared/service/confirmation-message';
import { MessageService } from '../../shared/service/message.service';

import { Subscription } from 'rxjs/Subscription';

@Component({
  templateUrl: './sign-in.component.html',
  styleUrls: [ './sign-in.component.css' ]
})
export class SignInComponent implements OnDestroy {

  @ViewChild('signInForm') currentForm;

  signInUser: SignIn = new SignIn();
  
  _subscription: Subscription;

  constructor(private messageService: MessageService) {
    this._subscription = this.messageService.messageConfirmed$.subscribe((message: any)=>{
      let confirm: ConfirmationMessage = <ConfirmationMessage>message;
      alert('Received:' + JSON.stringify(confirm));
    })
  }

  signIn(): void {
    let m: ConfirmationMessage = new ConfirmationMessage();
    m.title = 'Sign In';
    m.message = 'Sign in success.';
    this.messageService.announceMessage(m);
    this.currentForm.reset();
  }

  ngOnDestroy(): void {
    if(this._subscription) {
      this._subscription.unsubscribe();
    }
  }
}