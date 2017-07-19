import { Injectable } from '@angular/core';
import { Subject} from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';
import { Message } from './message';
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

}