import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, CanActivateChild, CanDeactivate, Router, RouterStateSnapshot } from '@angular/router';
import { AppInitService } from '../app.init.service';
import { MessageService } from './message-service/message.service';
import { ServiceComponent } from "../service/service.component";
import { Observable } from "rxjs/Observable";
import { Subject } from "rxjs/Subject";
import { K8sService } from "../service/service.k8s";
import { Message, RETURN_STATUS } from "./shared.types";

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