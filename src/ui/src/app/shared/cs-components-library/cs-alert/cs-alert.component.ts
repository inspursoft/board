import { Component, OnInit } from '@angular/core';
import { AlertMessage } from '../../shared.types';
import { Observable, Subject } from 'rxjs';
import { animate, state, style, transition, trigger } from '@angular/animations';
import { DISMISS_ALERT_INTERVAL } from "../../shared.const";
import { Subscription } from "rxjs/Subscription";
import "rxjs/add/observable/interval"

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
export class CsAlertComponent implements OnInit{
  _isOpen = false;
  curMessage: AlertMessage;
  onCloseEvent: Subject<any>;
  animation: string;
  isRunningAnimation = true;
  timeRemaining: number;
  intervalSubscription: Subscription;

  constructor() {
    this.onCloseEvent = new Subject<any>();
  }

  ngOnInit(): void {
    this.timeRemaining = DISMISS_ALERT_INTERVAL;
    this.intervalSubscription = Observable.interval(1000).subscribe(()=>{
      if (this.isRunningAnimation){
        if (this.timeRemaining == 0){
          this.animation = 'hidden';
          setTimeout(() => this.isOpen = false, 500);
        } else {
          this.timeRemaining --;
        }
      }
    })
  }

  get isOpen(): boolean {
    return this._isOpen;
  }

  set isOpen(value: boolean) {
    this._isOpen = value;
    if (!value) {
      this.onCloseEvent.next();
      this.intervalSubscription.unsubscribe();
    }
  }

  public openAlert(message: AlertMessage): Observable<any> {
    this.curMessage = message;
    this.isOpen = true;
    this.animation = 'hidden';
    setTimeout(() => this.animation = 'show');
    return this.onCloseEvent.asObservable();
  }
}
