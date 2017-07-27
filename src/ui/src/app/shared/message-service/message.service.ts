import { Injectable } from '@angular/core';
import { Subject} from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';
import { Message } from './message';
import { MESSAGE_TYPE } from '../shared.const';
import { Response } from '@angular/http';

@Injectable()
export class MessageService {
  
  messageAnnouncedSource: Subject<Message> = new Subject<Message>();
  messageAnnounced$: Observable<Message> = this.messageAnnouncedSource.asObservable();

  messageConfirmedSource: Subject<Message> = new Subject<Message>();
  messageConfirmed$: Observable<Message> = this.messageConfirmedSource.asObservable();

  inlineAlertAnnouncedSource: Subject<Message> = new Subject<Message>();
  inlineAlertAnnounced$: Observable<Message> = this.inlineAlertAnnouncedSource.asObservable();

  globalAnnouncedSource: Subject<Message> = new Subject<Message>();
  globalAnnounced$: Observable<Message> = this.globalAnnouncedSource.asObservable();

  announceMessage(message: Message) {
    this.messageAnnouncedSource.next(message);
  }

  confirmMessage(message: Message) {
    this.messageConfirmedSource.next(message);
  }

  inlineAlertMessage(message: Message) {
    this.inlineAlertAnnouncedSource.next(message);
  }

  globalMessage(message: Message) {
    this.globalAnnouncedSource.next(message);
  }

  dispatchError(response: Response | Error, error: string) {
    let errMessage = new Message();
    if(response instanceof Response) {
      errMessage.message = error;
      switch(response.status){
      case 400:
      case 403:
      case 404:
      case 409:
      case 412:
        errMessage.type = MESSAGE_TYPE.COMMON_ERROR;
        break;
      case 401:
        errMessage.type = MESSAGE_TYPE.INVALID_USER;
        errMessage.message = 'ERROR.INVALID_USER';
        break;
      case 500:
        errMessage.type = MESSAGE_TYPE.INTERNAL_ERROR;
        errMessage.message = 'ERROR.INTERNAL_ERROR';
        break;
      }
      if(errMessage.type === MESSAGE_TYPE.COMMON_ERROR) {
        this.inlineAlertMessage(errMessage);
      } else {
        this.globalMessage(errMessage);
      } 
    } else if (response instanceof Error) {
      errMessage.message = response.message;
      this.globalMessage(errMessage);
    }
  }

}