import { Injectable, NgModule } from '@angular/core';
import { Routes, RouterModule, Resolve, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { MainContentComponent } from './main-content/main-content.component';
import { GlobalSearchComponent } from './global-search/global-search.component';
import { AppInitService } from './shared.service/app-init.service';
import { Observable } from 'rxjs';
import { AppGuardService, AppInitializeGuard, AppInitializePageGuard } from './shared.service/app-guard.service';
import { RouteInitialize } from './shared/shared.const';
import { InitializePageComponent } from './initialize-page/initialize-page.component';

@Injectable()
export class SystemInfoResolve implements Resolve<any> {
  constructor(private appInitService: AppInitService) {
  }

  resolve(route: ActivatedRouteSnapshot,
          state: RouterStateSnapshot): Observable<any> | Promise<any> | any {
    return this.appInitService.getSystemInfo();
  }
}

const routes: Routes = [
  {
    path: 'account',
    loadChildren: './account/account.module#AccountModule',
    pathMatch: 'prefix',
    canActivate: [AppInitializeGuard]
  },
  {path: '', redirectTo: '/account/sign-in', pathMatch: 'full'},
  {
    path: RouteInitialize,
    component: InitializePageComponent,
    canActivate: [AppInitializePageGuard]
  },
  {
    path: '', component: MainContentComponent, resolve: {systeminfo: SystemInfoResolve},
    canActivate: [AppGuardService],
    children: [
      {path: 'search', component: GlobalSearchComponent},
      {path: 'dashboard', loadChildren: './dashboard/dashboard.module#DashboardModule'},
      {path: 'nodes', loadChildren: './node/node.module#NodeModule'},
      {path: 'services', loadChildren: './service/service.module#ServiceModule'},
      {path: 'audit', loadChildren: './audit/audit.module#AuditModule'},
      {path: 'admin', loadChildren: './admin/admin.module#AdminModule'},
      {path: 'projects', loadChildren: './project/project.module#ProjectModule'},
      {path: 'training-job', loadChildren: './job/job.module#JobModule'},
      {path: 'resource', loadChildren: './resource/resource.module#ResourceModule'},
      {path: 'helm', loadChildren: './helm/helm.module#HelmModule'},
      {path: 'profile', loadChildren: './profile/profile.module#ProfileModule'},
      {path: 'storage', loadChildren: './storage/storage.module#StorageModule'},
      {path: 'kibana-url', loadChildren: './kibana/kibana.module#KibanaModule'},
      {path: 'images', loadChildren: './image/image.module#ImageModule'},
    ]
  },
  {path: '**', redirectTo: '/account/sign-in'},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
