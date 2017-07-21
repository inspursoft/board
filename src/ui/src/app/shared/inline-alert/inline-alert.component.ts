import { Component, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { TranslateService } from '@ngx-translate/core';

import { MessageService } from '../message-service/message.service';
import { Message } from '../message-service/message';
import { DISMISS_INLINE_ALERT_INTERVAL, MESSAGE_TYPE } from '../shared.const';

@Component({
  selector: 'inline-alert',
  templateUrl: './inline-alert.component.html'
})
export class InlineAlertComponent implements OnDestroy {
  
  inlineAlertClosed: boolean;
  inlineAnnouncedMessage: string;
  inlineAlertType: string = 'alert-success';

  _subscription: Subscription;
 
  constructor(
    private messageService: MessageService,
    private translateService: TranslateService
  ){
    this.inlineAlertClosed = true;
    this._subscription = this.messageService
      .inlineAlertAnnounced$
      .subscribe(m=>{
        setTimeout(()=>this.inlineAlertClosed = true, DISMISS_INLINE_ALERT_INTERVAL);
        let inlineMessage = <Message>m;
        if(inlineMessage) {
          this.translateService.get(inlineMessage.message, inlineMessage.params || [])
            .subscribe(res=>{
              this.inlineAnnouncedMessage = res;
              if(inlineMessage.type) {
                this.inlineAlertType = 'alert-danger';
              }
              this.inlineAlertClosed = false;
            });
        }
      });   
  }  

  ngOnDestroy(): void {
    if(this._subscription) {
      this._subscription.unsubscribe();
    }
  }
}