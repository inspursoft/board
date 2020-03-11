import { Injectable } from '@angular/core';
import { Observable, Observer, Subject } from 'rxjs';

@Injectable()
export class WebsocketService {

  socket: Subject<MessageEvent>;

  connect(url: string): Subject<MessageEvent> {
    return this.create(url);
  }

  create(url: string): Subject<MessageEvent> {
    const ws = new WebSocket(url);
    const observable = Observable.create(
      (obs: Observer<MessageEvent>) => {
        ws.onmessage = obs.next.bind(obs);
        ws.onerror = obs.error.bind(obs);
        ws.onclose = obs.complete.bind(obs);
        return ws.close.bind(ws);
      }
    );
    const observer = {
      next: (data: any) => {
        if (ws.readyState === WebSocket.OPEN) {
          ws.send(data);
        }
      }
    };
    return Subject.create(observer, observable);
  }
}
