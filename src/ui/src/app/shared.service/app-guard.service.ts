import { Injectable } from '@angular/core';
import { ActivatedRouteSnapshot, CanActivate, CanActivateChild, CanDeactivate, Router, RouterStateSnapshot } from '@angular/router';
import { Observable, of } from 'rxjs';
import { catchError, map } from 'rxjs/operators';
import { AppInitService } from './app-init.service';
import { MessageService } from './message.service';
import { AppTokenService } from './app-token.service';
import { RouteInitialize, RouteSignIn } from '../shared/shared.const';

@Injectable()
export class AppGuardService implements CanActivate, CanActivateChild {

  constructor(private appInitService: AppInitService,
              private messageService: MessageService,
              private appTokenService: AppTokenService,
              private router: Router) {
  }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
    return this.appInitService.getCurrentUser(route.queryParamMap.get('token'))
      .pipe(map(() => {
        if (state.url === '/') {
          this.router.navigate(['/dashboard'], {queryParams: {token: this.appTokenService.token}}).then();
        }
        return true;
      }), catchError(() => {
        if (state.url.indexOf('/search') === 0) {
          this.messageService.cleanNotification();
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

@Injectable()
export class AppInitializeGuard implements CanActivate, CanActivateChild {

  constructor(private appInitService: AppInitService,
              private messageService: MessageService,
              private appTokenService: AppTokenService,
              private router: Router) {
  }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
    return this.appInitService.getSystemInfo()
      .pipe(
        map(() => true),
        catchError(() => {
          this.router.navigate([RouteInitialize]).then();
          return of(true);
        })
      );
  }

  canActivateChild(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
    return this.canActivate(route, state);
  }
}

@Injectable()
export class AppInitializePageGuard implements CanActivate, CanActivateChild {

  constructor(private appInitService: AppInitService,
              private router: Router) {
  }

  canActivate(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
    return this.appInitService.getSystemInfo()
      .pipe(
        map(() => {
          this.router.navigate([RouteSignIn]).then();
          return false;
        }),
        catchError(() => of(true))
      );
  }

  canActivateChild(route: ActivatedRouteSnapshot, state: RouterStateSnapshot): Observable<boolean> | Promise<boolean> | boolean {
    return this.canActivate(route, state);
  }
}

