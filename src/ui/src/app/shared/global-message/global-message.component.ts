import { Component, OnDestroy } from '@angular/core';
import { Router } from '@angular/router';
import { Subscription } from 'rxjs/Subscription';
import { MessageService } from '../message-service/message.service';
import { Message } from '../message-service/message';
import { MESSAGE_TYPE } from '../shared.const';

@Component({
  selector: 'global-message',
  templateUrl: './global-message.component.html'
})
export class GlobalMessageComponent implements OnDestroy {

  globalAnnoucedMessage: string;
  showAction: boolean;

  _subscription: Subscription;

  constructor(
    private messageService: MessageService,
    private router: Router
  ) {
    this.showAction = false;
    this._subscription = this.messageService
      .globalAnnounced$
      .subscribe(m=>{
        let globalMessage = <Message>m;
        this.globalAnnoucedMessage = globalMessage.message;
        if(globalMessage) {
          if(globalMessage.type === MESSAGE_TYPE.INVALID_USER) {
            this.showAction = true;
          }
        }
      });
  }

  redirectToSignIn(): void {
    this.router.navigateByUrl('/sign-in');
  }

  ngOnDestroy(): void {
    if(this._subscription) {
      this._subscription.unsubscribe();
    }
  }
}