import { Component, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { MessageService } from '../message-service/message.service';
import { Message } from '../message-service/message';
import { DISMISS_INLINE_ALERT_INTERVAL } from '../shared.const';

@Component({
  selector: 'inline-alert',
  templateUrl: './inline-alert.component.html'
})
export class InlineAlertComponent implements OnDestroy {
  
  inlineAlertClosed: boolean;
  inlineAnnouncedMessage: string;

  _subscription: Subscription;
 
  constructor(private messageService: MessageService){
    this.inlineAlertClosed = true;
    this._subscription = this.messageService
      .inlineAlertAnnounced$
      .subscribe(m=>{
        setTimeout(()=>this.inlineAlertClosed = true, DISMISS_INLINE_ALERT_INTERVAL);
        let inlineMessage = <Message>m;
        this.inlineAnnouncedMessage = inlineMessage.message;
        this.inlineAlertClosed = false;
      });   
  }  

  ngOnDestroy(): void {
    if(this._subscription) {
      this._subscription.unsubscribe();
    }
  }
}