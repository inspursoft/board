import { Injectable } from '@angular/core';
import { Subject} from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';

@Injectable()
export class MessageService {
  
  messageAnnouncedSource: Subject<any> = new Subject<any>();
  messageAnnounced$: Observable<any> = this.messageAnnouncedSource.asObservable();

  messageConfirmedSource: Subject<any> = new Subject<any>();
  messageConfirmed$: Observable<any> = this.messageConfirmedSource.asObservable();

  announceMessage(message: any) {
    this.messageAnnouncedSource.next(message);
  }

  confirmMessage(message: any) {
    this.messageConfirmedSource.next(message);
  }

}