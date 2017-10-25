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

  messageCanceledSource: Subject<boolean> = new Subject<boolean>();
  messageCanceled$: Observable<boolean> = this.messageCanceledSource.asObservable();

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

  cancelMessage() {
    this.messageCanceledSource.next(true);
  }

  inlineAlertMessage(message: Message) {
    this.inlineAlertAnnouncedSource.next(message);
  }

  globalMessage(message: Message) {
    this.globalAnnouncedSource.next(message);
  }

  _setErrorMessage(m: Message, defaultMessage: string, customMessage: string) {
    if(customMessage && customMessage.trim().length > 0) {
      m.message = customMessage;
    } else {
      m.message = defaultMessage;
    }
  }

  dispatchError(response: Response | Error, customMessage?: string) {
    let errMessage = new Message();
    if(response instanceof Response) {
      switch(response.status){
      case 401:
        errMessage.type = MESSAGE_TYPE.INVALID_USER;
        this._setErrorMessage(errMessage, 'ERROR.INVALID_USER', customMessage);
        break;
      case 404:
        errMessage.type = MESSAGE_TYPE.COMMON_ERROR;
        this._setErrorMessage(errMessage, 'ERROR.NOT_FOUND', customMessage);
        break;
      case 400:
        errMessage.type = MESSAGE_TYPE.COMMON_ERROR;
        this._setErrorMessage(errMessage, 'ERROR.BAD_REQUEST', customMessage);
        break;
      case 403:
        errMessage.type = MESSAGE_TYPE.COMMON_ERROR;
        this._setErrorMessage(errMessage, 'ERROR.INSUFFIENT_PRIVILEGES', customMessage);
        break;
      case 409:
        errMessage.type = MESSAGE_TYPE.COMMON_ERROR;
        this._setErrorMessage(errMessage, 'ERROR.CONFLICT_INPUT', customMessage);
        break;
      case 412:
        errMessage.type = MESSAGE_TYPE.COMMON_ERROR;
        this._setErrorMessage(errMessage, 'ERROR.PECONDITION_FAILED', customMessage);
        break;
      case 500:
      case 502:
      case 504:
        errMessage.type = MESSAGE_TYPE.INTERNAL_ERROR;
        this._setErrorMessage(errMessage, 'ERROR.INTERNAL_ERROR', customMessage);
        break;
      default:
        errMessage.type = MESSAGE_TYPE.INTERNAL_ERROR;
        this._setErrorMessage(errMessage, 'ERROR.UNKNOWN_ERROR', customMessage);
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
