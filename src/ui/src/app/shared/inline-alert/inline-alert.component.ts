import { Component, OnDestroy } from '@angular/core';
import { Subscription } from 'rxjs/Subscription';
import { Observable } from 'rxjs/Observable';
import 'rxjs/add/operator/switchMap';

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
    let hnd: any;
    this._subscription = this.messageService
      .inlineAlertAnnounced$
      .switchMap(m=>Observable.of(m))
      .subscribe(m=>{
        let inlineMessage = <Message>m;
        if(inlineMessage) {
          hnd = setTimeout(()=>{this.inlineAlertClosed = true; clearTimeout(hnd);}, DISMISS_INLINE_ALERT_INTERVAL); 
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