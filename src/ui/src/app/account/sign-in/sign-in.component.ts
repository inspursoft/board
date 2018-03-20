import { Component, OnInit, OnDestroy } from '@angular/core';
import { Router, ActivatedRoute } from '@angular/router';
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
    this._subscription = this.messageService.messageConfirmed$.subscribe((message: any)=>{
      let confirmationMessage = <Message>message;
      console.error('Received:' + JSON.stringify(confirmationMessage));
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
      .catch(err=>{
        this.isSignWIP = false;
        let announceMessage = new Message();
          announceMessage.title = 'ACCOUNT.ERROR';
        if(err) {
          switch(err.status){
          case 400:
            announceMessage.message = 'ACCOUNT.INCORRECT_USERNAME_OR_PASSWORD';
            break;
          case 409:
            announceMessage.message = 'ACCOUNT.ALREADY_SIGNED_IN';
            break;
          }
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