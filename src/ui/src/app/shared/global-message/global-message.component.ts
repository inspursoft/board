import { Component, OnDestroy } from '@angular/core';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs/Subscription';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/switchMap';

import { AppInitService } from '../../app.init.service';
import { MessageService } from '../message-service/message.service';
import { Message } from '../message-service/message';
import { MESSAGE_TYPE, DISMISS_GLOBAL_ALERT_INTERVAL } from '../shared.const';

@Component({
  selector: 'global-message',
  templateUrl: './global-message.component.html'
})
export class GlobalMessageComponent implements OnDestroy {

  globalMessageClosed: boolean;
  globalAnnoucedMessage: string;
  showAction: boolean;
  authMode: string = '';
  redirectionURL: string = '';
  
  _subscription: Subscription;

  constructor(
    private appInitService: AppInitService,
    private messageService: MessageService,
    private router: Router
  ) {
    this.globalMessageClosed = true;
    this.showAction = false;
    this._subscription = this.messageService
      .globalAnnounced$
      .switchMap(m=>Observable.of(m))
      .subscribe(m=>{
        this.globalMessageClosed = false;
        let globalMessage = <Message>m;
        this.globalAnnoucedMessage = globalMessage.message;
        if(globalMessage) {
          if(globalMessage.type === MESSAGE_TYPE.INVALID_USER) {
            this.showAction = true;
          }
        }
      });
    this.authMode = this.appInitService.systemInfo['auth_mode'];
    this.redirectionURL = this.appInitService.systemInfo['redirection_url'];
  }

  redirectToSignIn(): void {
    if(this.authMode === 'indata_auth') {
      window.location.href = this.redirectionURL;
      return;
    }
    this.router.navigateByUrl('/sign-in');
  }

  ngOnDestroy(): void {
    if(this._subscription) {
      this._subscription.unsubscribe();
    }
  }
}