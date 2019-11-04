import { Injectable, NgModule } from '@angular/core';
import { Routes, RouterModule, Resolve, ActivatedRouteSnapshot, RouterStateSnapshot } from '@angular/router';
import { MainContentComponent } from './main-content/main-content.component';
import { GlobalSearchComponent } from './global-search/global-search.component';
import { AppInitService } from './shared.service/app-init.service';
import { Observable } from 'rxjs';
import { AppGuardService } from "./shared.service/app-guard.service";

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
  {path: 'board/account', loadChildren: './account/account.module#AccountModule', pathMatch: 'prefix'},
  {path: '', redirectTo: '/board/account/sign-in', pathMatch: 'full'},
  {
    path: '', component: MainContentComponent, resolve: {systeminfo: SystemInfoResolve},
    canActivate: [AppGuardService],
    children: [
      {path: 'board/search', component: GlobalSearchComponent},
      {path: 'board/dashboard', loadChildren: './dashboard/dashboard.module#DashboardModule'},
      {path: 'board/nodes', loadChildren: './node/node.module#NodeModule'},
      {path: 'board/services', loadChildren: './service/service.module#ServiceModule'},
      {path: 'board/audit', loadChildren: './audit/audit.module#AuditModule'},
      {path: 'board/user-center', loadChildren: './user-center/user-center.module#UserCenterModule'},
      {path: 'board/projects', loadChildren: './project/project.module#ProjectModule'},
      {path: 'board/training-job', loadChildren: './job/job.module#JobModule'},
      {path: 'board/resource', loadChildren: './resource/resource.module#ResourceModule'},
      {path: 'board/helm', loadChildren: './helm/helm.module#HelmModule'},
      {path: 'board/profile', loadChildren: './profile/profile.module#ProfileModule'},
      {path: 'board/storage', loadChildren: './storage/storage.module#StorageModule'},
      {path: 'board/kibana-url', loadChildren:'./kibana/kibana.module#KibanaModule'},
      {path: 'board/images', loadChildren: './image/image.module#ImageModule'},
    ]
  },
  {path: '**', redirectTo: '/board/account/sign-in'},
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule {
}
