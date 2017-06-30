import { Component, OnDestroy } from '@angular/core';

import { Subscription } from 'rxjs/Subscription';
import { MessageService } from '../service/message.service';
import { ConfirmationMessage } from '../service/confirmation-message';

@Component({
  selector: 'confirmation-dialog',
  templateUrl: 'confirmation-dialog.component.html',
  styleUrls: [ 'confirmation-dialog.component.css']
})
export class ConfirmationDialogComponent implements OnDestroy {

  opened: boolean;
  confirmationMessage: ConfirmationMessage = new ConfirmationMessage();
  
  _subscription: Subscription

  constructor(private messageService: MessageService) {
    this._subscription = this.messageService.messageAnnounced$.subscribe((message: any)=>{
      this.confirmationMessage = <ConfirmationMessage>message;
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