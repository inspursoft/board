import { Component } from '@angular/core';
import { Router } from '@angular/router';

import { ConfirmationMessage } from '../../shared/service/confirmation-message';
import { MessageService } from '../../shared/service/message.service';

import { Subscription } from 'rxjs/Subscription';

import { SignUp } from './sign-up';
import { Account } from '../account';
import { AccountService } from '../account.service';

@Component({
   templateUrl: './sign-up.component.html',
   styleUrls: [ './sign-up.component.css' ]
})
export class SignUpComponent {
  
  signUpModel: SignUp = new SignUp();
  _subscription: Subscription;

  constructor(
    private accountService: AccountService, 
    private messageService: MessageService,
    private router: Router) {}
  
  signUp(): void {
    let account: Account = {
      username: this.signUpModel.username,
      email: this.signUpModel.email,
      password: this.signUpModel.password,
      comment: this.signUpModel.comment
    };
    this.accountService
      .signUp(account)
      .then(res=>this.router.navigate(['/sign-in']))
      .catch(err=>{
        let m: ConfirmationMessage = new ConfirmationMessage();
        m.title = "Error";
        if(err && err.status === 409) {
          m.message = 'Username or password already exists.';
        } else {
          m.message = "Failed to sign up.";
        }
        this.messageService.announceMessage(m);
      });
  }

  goBack(): void {
    this.router.navigate(['/sign-in']);
  }
}