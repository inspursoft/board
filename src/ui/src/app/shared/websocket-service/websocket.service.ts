import { Injectable } from '@angular/core';
import { Subject } from 'rxjs/Subject';
import { Observable } from 'rxjs/Observable';
import { Observer } from 'rxjs/Observer';

@Injectable()
export class WebsocketService {
  
  socket: Subject<MessageEvent>;

  connect(url: string): Subject<MessageEvent> {
    if (!this.socket) {
      this.socket = this.create(url);
    }
    return this.socket;
  }

  create(url: string): Subject<MessageEvent> {
  
    let ws = new WebSocket(url);
    let observable = Observable.create(
      (obs: Observer<MessageEvent>)=> {
        ws.onmessage = obs.next.bind(obs);
        ws.onerror = obs.error.bind(obs);
        ws.onclose = obs.complete.bind(obs);
        return ws.close.bind(ws);
      }
    );
    let observer = {
      next: (data: any) => {
        if (ws.readyState === WebSocket.OPEN) {
          ws.send(data);
        }
      }
    }
    return Subject.create(observer, observable);
  }
}