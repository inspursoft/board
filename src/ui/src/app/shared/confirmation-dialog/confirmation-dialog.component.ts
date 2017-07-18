import { Component, OnDestroy } from '@angular/core';

import { Subscription } from 'rxjs/Subscription';

import { TranslateService } from '@ngx-translate/core';

import { MessageService } from '../message-service/message.service';
import { Message } from '../message-service/message';

@Component({
  selector: 'confirmation-dialog',
  templateUrl: 'confirmation-dialog.component.html'
})
export class ConfirmationDialogComponent implements OnDestroy {

  opened: boolean;
  confirmationMessage: Message = new Message();
  
  _subscription: Subscription

  constructor(
    private messageService: MessageService,
    private translateService: TranslateService) {
    this._subscription = this.messageService.messageAnnounced$.subscribe((message: any)=>{
      this.confirmationMessage = <Message>message;
      this.translateService.get(this.confirmationMessage.title)
        .subscribe(res=>this.confirmationMessage.title = res);
      this.translateService.get(this.confirmationMessage.message)
        .subscribe(res=>this.confirmationMessage.message = res);
      this.opened = true;
    });
  }

  confirm(): void {
    this.messageService.confirmMessage(this.confirmationMessage);
    this.opened = false;
  }

  ngOnDestroy(): void {
    if(this._subscription) {
      this._subscription.unsubscribe();
    }
  }
}