import { Injectable, OnDestroy, OnInit } from '@angular/core';
import {
  CanActivate, CanActivateChild, Router,
  ActivatedRouteSnapshot,
  RouterStateSnapshot, CanDeactivate
}from '@angular/router';
import { AppInitService } from '../app.init.service';
import { Message } from './message-service/message';
import { MessageService } from './message-service/message.service';
import { ServiceComponent } from "../service/service.component";
import { Observable } from "rxjs/Observable";
import { Subscription } from "rxjs/Subscription";
import { BUTTON_STYLE } from "./shared.const";
import { Subject } from "rxjs/Subject";

@Injectable()
export class AuthGuard implements CanActivate, CanActivateChild {

  constructor(private appInitService: AppInitService,
              private messageService: MessageService,
              private router: Router) {
  }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Promise<boolean> | boolean {
    return new Promise<boolean>((resolve, reject) => {
      this.appInitService
        .getCurrentUser(route.queryParamMap.get("token"))
        .then(res => {
          if (state.url === '/') {
            this.router.navigate(['/dashboard']);
          }
          resolve(true);
        })
        .catch(err => {
          if (state.url.indexOf('/search') === 0) {
            return resolve(true);
          }
          this.router.navigate(['/sign-in']);
          resolve(false);
        });
    });
  }

  canActivateChild(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Promise<boolean> | boolean {
    return this.canActivate(route, state);
  }
}

@Injectable()
export class ServiceGuard implements OnDestroy, CanDeactivate<ServiceComponent> {
  serviceSubject: Subject<boolean> = new Subject<boolean>();
  _confirmSubscription: Subscription;
  _cancelSubscription: Subscription;

  constructor(private messageService: MessageService) {
    this._confirmSubscription = this.messageService.messageConfirmed$.subscribe(next => {
      this.serviceSubject.next(true);
    });
    this._cancelSubscription = this.messageService.messageCanceled$.subscribe(next => {
      this.serviceSubject.next(false);
    });
  }

  ngOnDestroy() {
    this._confirmSubscription.unsubscribe();
    this._cancelSubscription.unsubscribe();
  }

  canDeactivate(component: ServiceComponent,
                currentRoute: ActivatedRouteSnapshot,
                currentState: RouterStateSnapshot,
                nextState?: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
    if (component.currentStepIndex > 1) {
      let m: Message = new Message();
      m.title = "SERVICE.ASK_TITLE";
      m.buttons = BUTTON_STYLE.YES_NO;
      m.message = "SERVICE.ASK_TEXT";
      this.messageService.announceMessage(m);
      let result = this.serviceSubject.asObservable();
      result.subscribe(isCanDeactivate => {
        return isCanDeactivate;
      });
      return result;
    }
    return true;
  }

}