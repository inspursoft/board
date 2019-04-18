import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanDeactivate, RouterStateSnapshot } from '@angular/router';
import { ServiceComponent } from './service.component';
import { Observable, Subject } from 'rxjs';
import { K8sService } from './service.k8s';
import { Message, RETURN_STATUS } from '../shared/shared.types';
import { MessageService } from '../shared.service/message.service';

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
        if (message.returnStatus === RETURN_STATUS.rsConfirm) {
          this.k8sService.cancelBuildService();
        }
        this.serviceSubject.next(message.returnStatus === RETURN_STATUS.rsConfirm);
      });
      const result = this.serviceSubject.asObservable();
      result.subscribe(isCanDeactivate => {
        return isCanDeactivate;
      });
      return result;
    }
    return true;
  }

}
