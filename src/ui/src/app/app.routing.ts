/*
 * Copyright (c) 2016 VMware, Inc. All Rights Reserved.
 * This software is released under MIT license.
 * The full license information can be found in LICENSE in the root directory of this project.
 */
import { ActivatedRouteSnapshot, Resolve, RouterStateSnapshot, Routes } from '@angular/router';
import { Injectable } from '@angular/core';
import { GlobalSearchComponent } from './global-search/global-search.component';
import { SignInComponent } from './account/sign-in/sign-in.component';
import { SignUpComponent } from './account/sign-up/sign-up.component';
import { MainContentComponent } from './main-content/main-content.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { NodeComponent } from './node/node.component';
import { ProjectComponent } from './project/project.component';
import { ImageListComponent } from './image/image-list/image-list.component';
import { ServiceComponent } from './service/service.component';
import { UserCenterComponent } from './user-center/user-center.component';
import { AuthGuard, ServiceGuard } from './shared/auth-guard.service';
import { ProfileComponent } from "./profile/profile.component";
import { AppInitService } from "./app.init.service";
import { MemberComponent } from "./shared/create-project/member/member.component";
import { ListAuditComponent } from "./audit/operation-audit-list/list-audit.component";
import { ResetPasswordComponent } from "./account/reset-password/reset-password.component";
import { ForgotPasswordComponent } from "./account/forgot-password/forgot-password.component";
import { TimeoutComponent } from "./shared/error-pages/timeout.component/timeout.component";
import { BadGatewayComponent } from "./shared/error-pages/bad-gateway.component/bad-gateway.component";
import { BoardLoadingComponent } from "./shared/error-pages/board-loading.component/board-loading.component";
import { KibanaComponent } from "./kibana/kibana/kibana.component";
import { GrafanaComponent } from "./grafana/grafana/grafana.component";
import {
  RouteGrafana,
  RouteHelm,
  RouteKibana,
  RoutePV,
  RoutePvc,
  RouteReleaseList,
  RouteRepoList,
  RouteStorage
} from "./shared/shared.const";
import { PvListComponent } from "./storage/pv/pv-list.compoent/pv-list.component";
import { Observable } from "rxjs/Observable";
import { PvcListComponent } from "./storage/pvc/pvc-list.component/pvc-list.component";
import { HelmHostComponent } from "./helm/helm-host/helm-host.component";
import { ChartReleaseListComponent } from "./helm/chart-release-list/chart-release-list.component";

@Injectable()
export class SystemInfoResolve implements Resolve<any> {
  constructor(private appInitService: AppInitService) {
  }

  resolve(route: ActivatedRouteSnapshot,
          state: RouterStateSnapshot): Observable<any> | Promise<any> | any {
    return this.appInitService.getSystemInfo();
  }
}

export const ROUTES: Routes = [
  {
    path: 'sign-in',
    component: SignInComponent,
    resolve: {
      systeminfo: SystemInfoResolve
    },
  },
  {
    path: 'sign-up',
    component: SignUpComponent,
    resolve: {
      systeminfo: SystemInfoResolve
    }
  },
  {
    path: 'reset-password',
    component: ResetPasswordComponent,
    resolve: {
      systeminfo: SystemInfoResolve
    }
  },
  {
    path: 'forgot-password',
    component: ForgotPasswordComponent,
    resolve: {
      systeminfo: SystemInfoResolve
    }
  },
  {path: 'timeout-page', component: TimeoutComponent},
  {path: 'bad-gateway-page', component: BadGatewayComponent},
  {path: 'board-loading-page', component: BoardLoadingComponent},
  {
    path: '', component: MainContentComponent,
    resolve: {
      systeminfo: SystemInfoResolve
    },
    canActivate: [AuthGuard],
    children: [
      {path: 'search', component: GlobalSearchComponent},
      {path: 'dashboard', component: DashboardComponent},
      {path: 'nodes', component: NodeComponent},
      {
        path: 'projects',
        children: [
          {path: '', component: ProjectComponent},
          {path: 'members', component: MemberComponent}
        ]
      },
      {path: 'images', component: ImageListComponent},
      {path: 'services', component: ServiceComponent, canDeactivate: [ServiceGuard]},
      {path: `${RouteHelm}/${RouteRepoList}`, component: HelmHostComponent},
      {path: `${RouteHelm}/${RouteReleaseList}`, component: ChartReleaseListComponent},
      {path: 'user-center', component: UserCenterComponent},
      {path: 'profile', component: ProfileComponent},
      {path: RouteKibana, component: KibanaComponent},
      {path: RouteGrafana, component: GrafanaComponent},
      {path: `${RouteStorage}/${RoutePV}`, component: PvListComponent},
      {path: `${RouteStorage}/${RoutePvc}`, component: PvcListComponent},
      {path: 'audit', component: ListAuditComponent}
    ]
  },
  {path: '', redirectTo: '/sign-in', pathMatch: 'full'},
  {path: '**', component: SignInComponent}
];
