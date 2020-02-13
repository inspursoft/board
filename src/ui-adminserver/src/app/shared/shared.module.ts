import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HeaderComponent } from './header/header.component';
import { ClarityModule } from '@clr/angular';
import { AppRoutingModule } from '../app-routing.module';
import { Error404Component } from './error-pages/error404/error404.component';
import { TranslateModule } from '@ngx-translate/core';
import { BackTopComponent } from './back-top/back-top.component';
import { ScrollTools } from './scroll.tools';

@NgModule({
  declarations: [
    HeaderComponent,
    Error404Component,
    BackTopComponent
  ],
  imports: [
    CommonModule,
    ClarityModule,
    TranslateModule,
    AppRoutingModule,
  ],
  exports: [
    HeaderComponent,
    BackTopComponent,
  ],
  providers: [
    ScrollTools,
  ],
})
export class SharedModule { }
