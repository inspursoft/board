import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, CanActivateChild, CanDeactivate, Router, RouterStateSnapshot } from '@angular/router';
import { Observable, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { AppInitService } from './app-init.service';
import { MessageService } from './message.service';

@Injectable()
export class AppGuardService implements CanActivate, CanActivateChild {

  constructor(private appInitService: AppInitService,
              private messageService: MessageService,
              private router: Router) {
  }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
    return this.appInitService.getCurrentUser(route.queryParamMap.get('token'))
      .pipe(map(() => {
        if (state.url === '/') {
          this.router.navigate(['/dashboard']).then();
        }
        return true;
      }), catchError(() => {
        if (state.url.indexOf('/search') === 0) {
          return of(true);
        } else {
          this.router.navigate(['/account/sign-in']).then();
          return of(true);
        }
      }));
  }

  canActivateChild(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
    return this.canActivate(route, state);
  }
}
