import { Component } from '@angular/core';
import { AlertMessage } from '../../shared.types';
import { Observable, Subject } from 'rxjs';
import { animate, state, style, transition, trigger } from '@angular/animations';
import { DISMISS_ALERT_INTERVAL } from "../../shared.const";

@Component({
  templateUrl: './cs-alert.component.html',
  styleUrls: ['./cs-alert.component.css'],
  animations: [
    trigger('open', [
      state('hidden', style({height: '0'})),
      state('show', style({height: '50px'})),
      transition('hidden <=> show', animate(500))
    ])
  ]
})
export class CsAlertComponent {
  _isOpen: boolean = false;
  curMessage: AlertMessage;
  onCloseEvent: Subject<any>;
  animation: string;

  constructor() {
    this.onCloseEvent = new Subject<any>();
  }

  get isOpen(): boolean {
    return this._isOpen;
  }

  set isOpen(value: boolean) {
    this._isOpen = value;
    if (!value) {
      this.onCloseEvent.next();
    }
  }

  public openAlert(message: AlertMessage): Observable<any> {
    this.curMessage = message;
    this.isOpen = true;
    this.animation = 'hidden';
    setTimeout(() => this.animation = 'show');
    setTimeout(() => this.animation = 'hidden', DISMISS_ALERT_INTERVAL);
    setTimeout(() => this.isOpen = false, DISMISS_ALERT_INTERVAL + 500);
    return this.onCloseEvent.asObservable();
  }
}
