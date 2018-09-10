import { Component, OnDestroy } from '@angular/core';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs/Subscription';
import 'rxjs/add/operator/switchMap';

import { AppInitService } from '../../app.init.service';
import { MessageService } from '../message-service/message.service';
import { Message } from '../message-service/message';
import { MESSAGE_TYPE } from '../shared.const';
import { HttpErrorResponse } from "@angular/common/http";

@Component({
  selector: 'global-message',
  templateUrl: './global-message.component.html'
})
export class GlobalMessageComponent implements OnDestroy {

  globalMessageClosed: boolean;
  globalAnnoucedMessage: string = "";
  showAction: boolean = false;
  showDetail: boolean = false;
  detailModalOpen: boolean = false;
  errorObject: HttpErrorResponse | Error;
  authMode: string = '';
  redirectionURL: string = '';
  
  _subscription: Subscription;

  constructor(private appInitService: AppInitService,
              private messageService: MessageService,
              private router: Router) {
    this.globalMessageClosed = true;
    this.showAction = false;
    this._subscription = this.messageService.globalAnnounced$
      .subscribe((msg: Message) => {
        this.globalMessageClosed = false;
        this.errorObject = msg.errorObject;
        this.globalAnnoucedMessage = msg.message;
        this.showDetail = msg.type === MESSAGE_TYPE.SHOW_DETAIL;
        this.showAction = msg.type === MESSAGE_TYPE.INVALID_USER;
      });
    this.authMode = this.appInitService.systemInfo['auth_mode'];
    this.redirectionURL = this.appInitService.systemInfo['redirection_url'];
  }

  get errorDetailMsg(): string {
    let result: string = "";
    if (this.errorObject && this.errorObject instanceof HttpErrorResponse) {
      result = (this.errorObject as HttpErrorResponse).message
    } else if (this.errorObject) {
      result = (this.errorObject as Error).message;
    }
    return result;
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