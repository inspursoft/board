import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, CanActivateChild, CanDeactivate, Router, RouterStateSnapshot } from '@angular/router';
import { AppInitService } from '../app.init.service';
import { MessageService } from './message-service/message.service';
import { ServiceComponent } from "../service/service.component";
import { Observable } from "rxjs/Observable";
import { Subject } from "rxjs/Subject";
import { K8sService } from "../service/service.k8s";
import { Message, RETURN_STATUS } from "./shared.types";
import "rxjs/add/operator/map"
import "rxjs/add/operator/catch"
import "rxjs/add/observable/of"

@Injectable()
export class AuthGuard implements CanActivate, CanActivateChild {

  constructor(private appInitService: AppInitService,
              private messageService: MessageService,
              private router: Router) {
  }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
    return this.appInitService.getCurrentUser(route.queryParamMap.get('token'))
      .map(() => {
        if (state.url === '/') {
          this.router.navigate(['/dashboard']).then();
        }
        return true;
      })
      .catch(() => {
        if (state.url.indexOf('/search') === 0) {
          return Observable.of(true);
        } else {
          this.router.navigate(['/sign-in']).then();
          return Observable.of(true);
        }
      })
  }

  canActivateChild(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
    return this.canActivate(route, state);
  }
}

@Injectable()
export class ServiceGuard implements CanDeactivate<ServiceComponent> {
  serviceSubject: Subject<boolean> = new Subject<boolean>();

  constructor(private messageService: MessageService,
              private k8sService: K8sService) {
  }

  canDeactivate(component: ServiceComponent,
                currentRoute: ActivatedRouteSnapshot,
                currentState: RouterStateSnapshot,
                nextState?: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
    if (component.currentStepIndex > 0) {
      this.messageService.showYesNoDialog('SERVICE.ASK_TEXT', 'SERVICE.ASK_TITLE').subscribe((message: Message) => {
        if (message.returnStatus == RETURN_STATUS.rsConfirm) {
          this.k8sService.cancelBuildService();
        }
        this.serviceSubject.next(message.returnStatus == RETURN_STATUS.rsConfirm);
      });
      let result = this.serviceSubject.asObservable();
      result.subscribe(isCanDeactivate => {
        return isCanDeactivate;
      });
      return result;
    }
    return true;
  }

}