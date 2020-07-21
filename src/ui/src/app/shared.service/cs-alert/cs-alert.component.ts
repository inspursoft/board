import { Component, OnInit } from '@angular/core';
import { AlertMessage } from '../../shared/shared.types';
import { interval, Observable, Subject, Subscription } from 'rxjs';
import { animate, state, style, transition, trigger } from '@angular/animations';
import { DISMISS_ALERT_INTERVAL } from '../../shared/shared.const';

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
export class CsAlertComponent implements OnInit {
  isOpenValue = false;
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
    this.intervalSubscription = interval(1000).subscribe(() => {
      if (this.isRunningAnimation) {
        if (this.timeRemaining === 0) {
          this.animation = 'hidden';
          setTimeout(() => this.isOpen = false, 500);
        } else {
          this.timeRemaining--;
        }
      }
    });
  }

  get isOpen(): boolean {
    return this.isOpenValue;
  }

  set isOpen(value: boolean) {
    this.isOpenValue = value;
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
