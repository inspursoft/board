import { Injectable, OnDestroy } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, CanActivateChild, CanDeactivate, Router, RouterStateSnapshot } from '@angular/router';
import { AppInitService } from '../app.init.service';
import { Message } from './message-service/message';
import { MessageService } from './message-service/message.service';
import { ServiceComponent } from "../service/service.component";
import { Observable } from "rxjs/Observable";
import { Subscription } from "rxjs/Subscription";
import { BUTTON_STYLE, MESSAGE_TARGET } from "./shared.const";
import { Subject } from "rxjs/Subject";
import { K8sService } from "../service/service.k8s";

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
            resolve(true);
          }
          resolve(true);
        })
        .catch(err => {
          if (state.url.indexOf('/search') === 0) {
            resolve(true);
          } else {
            this.router.navigate(['/sign-in']);
            resolve(true);
          }
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
  confirmSubscription: Subscription;
  cancelSubscription: Subscription;

  constructor(private messageService: MessageService,
              private k8sService: K8sService) {
    this.confirmSubscription = this.messageService.messageConfirmed$.subscribe((msg: Message) => {
      if (msg.target === MESSAGE_TARGET.CANCEL_BUILD_SERVICE_GUARD) {
        this.k8sService.cancelBuildService();
        this.serviceSubject.next(true);
      }
    });
    this.cancelSubscription = this.messageService.messageCanceled$.subscribe((msg: Message) => {
      if (msg.target === MESSAGE_TARGET.CANCEL_BUILD_SERVICE_GUARD) {
        this.serviceSubject.next(false);
      }
    });
  }

  ngOnDestroy() {
    this.cancelSubscription.unsubscribe();
    this.confirmSubscription.unsubscribe();
  }

  canDeactivate(component: ServiceComponent,
                currentRoute: ActivatedRouteSnapshot,
                currentState: RouterStateSnapshot,
                nextState?: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
    if (component.currentStepIndex > 0) {
      let msg: Message = new Message();
      msg.title = "SERVICE.ASK_TITLE";
      msg.buttons = BUTTON_STYLE.YES_NO;
      msg.message = "SERVICE.ASK_TEXT";
      msg.target = MESSAGE_TARGET.CANCEL_BUILD_SERVICE_GUARD;
      this.messageService.announceMessage(msg);
      let result = this.serviceSubject.asObservable();
      result.subscribe(isCanDeactivate => {
        return isCanDeactivate;
      });
      return result;
    }
    return true;
  }

}