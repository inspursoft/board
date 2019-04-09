import { Injectable, NgModule } from '@angular/core';
import { Routes, RouterModule, Resolve, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { SignInComponent } from './account/sign-in/sign-in.component';
import { SignUpComponent } from './account/sign-up/sign-up.component';
import { ResetPasswordComponent } from './account/reset-password/reset-password.component';
import { ForgotPasswordComponent } from './account/forgot-password/forgot-password.component';
import { TimeoutComponent } from './shared/error-pages/timeout.component/timeout.component';
import { BadGatewayComponent } from './shared/error-pages/bad-gateway.component/bad-gateway.component';
import { BoardLoadingComponent } from './shared/error-pages/board-loading.component/board-loading.component';
import { MainContentComponent } from './main-content/main-content.component';
import { AuthGuard, ServiceGuard } from './shared/auth-guard.service';
import { GlobalSearchComponent } from './global-search/global-search.component';
import { DashboardComponent } from './dashboard/dashboard.component';
import { NodeComponent } from './node/node.component';
import { ProjectComponent } from './project/project.component';
import { ImageListComponent } from './image/image-list/image-list.component';
import { ConfigMapListComponent } from './resource/config-map/config-map-list/config-map-list.component';
import { ServiceComponent } from './service/service.component';
import { HelmHostComponent } from './helm/helm-host/helm-host.component';
import { ChartReleaseListComponent } from './helm/chart-release-list/chart-release-list.component';
import { UserCenterComponent } from './user-center/user-center.component';
import { ProfileComponent } from './profile/profile.component';
import { KibanaComponent } from './kibana/kibana/kibana.component';
import { GrafanaComponent } from './grafana/grafana/grafana.component';
import { PvListComponent } from './storage/pv/pv-list.compoent/pv-list.component';
import { PvcListComponent } from './storage/pvc/pvc-list.component/pvc-list.component';
import { ListAuditComponent } from './audit/operation-audit-list/list-audit.component';
import { AppInitService } from './app.init.service';
import { Observable } from 'rxjs';

@Injectable()
export class SystemInfoResolve implements Resolve<any> {
  constructor(private appInitService: AppInitService) {
  }

  resolve(route: ActivatedRouteSnapshot,
          state: RouterStateSnapshot): Observable<any> | Promise<any> | any {
    return this.appInitService.getSystemInfo();
  }
}

const routes: Routes = [{
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
      {path: 'projects', component: ProjectComponent},
      {path: 'images', component: ImageListComponent},
      {
        path: 'resource', children: [
          {path: 'config-map', component: ConfigMapListComponent}
        ]
      },
      {path: 'services', component: ServiceComponent, canDeactivate: [ServiceGuard]},
      {
        path: 'helm', children: [
          {path: 'repo-list', component: HelmHostComponent},
          {path: 'release-list', component: ChartReleaseListComponent}
        ]
      },
      {path: 'user-center', component: UserCenterComponent},
      {path: 'profile', component: ProfileComponent},
      {path: 'kibana-url', component: KibanaComponent},
      {path: 'grafana', component: GrafanaComponent},
      {
        path: 'storage', children: [
          {path: 'pv', component: PvListComponent},
          {path: 'pvc', component: PvcListComponent}
        ]
      },
      {path: 'audit', component: ListAuditComponent}
    ]
  },
  {path: '', redirectTo: '/sign-in', pathMatch: 'full'},
  {path: '**', component: SignInComponent}
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
