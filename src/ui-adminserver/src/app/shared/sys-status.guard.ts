import { Injectable } from '@angular/core';
import { CanActivate, CanActivateChild, ActivatedRouteSnapshot, RouterStateSnapshot, UrlTree, Router } from '@angular/router';
import { Observable } from 'rxjs';
import { AppInitService } from '../shared.service/app-init.service';
import { InitStatus, InitStatusCode } from '../shared.service/app-init.type';
import { HttpErrorResponse } from '@angular/common/http';

@Injectable({
  providedIn: 'root'
})
export class SysStatusGuard implements CanActivate, CanActivateChild {
  constructor(private appInitService: AppInitService,
              private router: Router) { }

  canActivate(
    next: ActivatedRouteSnapshot,
    state: RouterStateSnapshot): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
    return this.checkSysStatus().then((res) => {
      if (!res) {
        this.router.navigateByUrl('installation');
      }
      return res;
    });
  }
  canActivateChild(
    next: ActivatedRouteSnapshot,
    state: RouterStateSnapshot): Observable<boolean | UrlTree> | Promise<boolean | UrlTree> | boolean | UrlTree {
    return this.checkSysStatus().then((res) => {
      if (!res) {
        this.router.navigateByUrl('installation');
      }
      return res;
    });
  }

  async checkSysStatus() {
    let result = false;
    this.appInitService.getSystemStatus().subscribe(
      (res: InitStatus) => {
        result = (res.status === InitStatusCode.InitStatusThird);
      },
      (err: HttpErrorResponse) => {
        console.log('error:' + err.message);
      }
    );
    return result;
  }
}
