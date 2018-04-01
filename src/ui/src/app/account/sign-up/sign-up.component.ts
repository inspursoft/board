import { Component } from '@angular/core';
import { Router } from '@angular/router';

import { Message } from '../../shared/message-service/message';
import { MessageService } from '../../shared/message-service/message.service';

import { Subscription } from 'rxjs/Subscription';

import { SignUp } from './sign-up';
import { Account } from '../account';
import { AccountService } from '../account.service';

@Component({
   templateUrl: './sign-up.component.html',
   styleUrls: [ './sign-up.component.css' ]
})
export class SignUpComponent {
  isSignUpWIP:boolean = false;
  signUpModel: SignUp = new SignUp();
  _subscription: Subscription;

  constructor(
    private accountService: AccountService, 
    private messageService: MessageService,
    private router: Router) {}
  
  signUp(): void {
    this.isSignUpWIP = true;
    let account: Account = {
      username: this.signUpModel.username,
      email: this.signUpModel.email,
      password: this.signUpModel.password,
      realname: this.signUpModel.realname,
      comment: this.signUpModel.comment
    };
    this.accountService
      .signUp(account)
      .then(res=>{
        this.isSignUpWIP = false;
        this.router.navigate(['/sign-in']);
      })
      .catch(err=>{
        this.isSignUpWIP = false;
        let confirmationMessage = new Message();
        confirmationMessage.title = "ACCOUNT.ERROR";
        if(err && err.status === 409) {
          confirmationMessage.message = 'ACCOUNT.USERNAME_ALREADY_EXISTS';
        } else {
          confirmationMessage.message = "ACCOUNT.FAILED_TO_SIGN_UP";
        }
        this.messageService.announceMessage(confirmationMessage);
      });
  }

  goBack(): void {
    this.router.navigate(['/sign-in']);
  }
}