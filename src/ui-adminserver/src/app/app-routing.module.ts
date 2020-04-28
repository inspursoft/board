import { NgModule } from '@angular/core';
import { Routes, RouterModule } from '@angular/router';
import { MainContentComponent } from './main-content/main-content.component';
import { Error404Component } from './shared/error-pages/error404/error404.component';
import { PreviewerComponent } from './dashboard/previewer/previewer.component';
import { CfgCardsComponent } from './configuration/cfg-cards.component';
import { SignInComponent } from './account/sign-in/sign-in.component';
import { InstallationComponent } from './account/installation/installation.component';
import { SysStatusGuard } from './shared/sys-status.guard';

const childrenPath: Routes = [
  { path: 'dashboard', component: PreviewerComponent },
  { path: 'configuration', component: CfgCardsComponent },
  { path: 'resource', loadChildren: './resource/resource.module#ResourceModule' },
];

const accountPath: Routes = [
  { path: '', redirectTo: '/account/login', pathMatch: 'full' },
  { path: 'login', component: SignInComponent },
  // { path: 'sign-up', component: SignUpComponent },
];

const routes: Routes = [
  { path: 'account', canActivateChild: [SysStatusGuard], children: accountPath, pathMatch: 'prefix' },
  { path: '', redirectTo: '/dashboard', pathMatch: 'full' },
  { path: '', component: MainContentComponent, canActivate: [SysStatusGuard], children: childrenPath },
  { path: 'installation', component: InstallationComponent, pathMatch: 'full' },
  { path: '**', component: Error404Component },
];

@NgModule({
  imports: [RouterModule.forRoot(routes)],
  exports: [RouterModule]
})
export class AppRoutingModule { }
