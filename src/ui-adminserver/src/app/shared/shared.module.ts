import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { HeaderComponent } from './header/header.component';
import { ClarityModule } from '@clr/angular';
import { AppRoutingModule } from '../app-routing.module';
import { Error404Component } from './error-pages/error404/error404.component';
import { TranslateModule } from '@ngx-translate/core';
import { BackTopComponent } from './back-top/back-top.component';
import { ScrollTools } from './scroll.tools';
import { AlertComponent } from './message/alert/alert.component';
import { DialogComponent } from './message/dialog/dialog.component';
import { GlobalAlertComponent } from './message/global-alert/global-alert.component';
import { HttpInterceptorService } from './http-client-interceptor';
import { MessageService } from './message/message.service';

@NgModule({
  declarations: [
    HeaderComponent,
    Error404Component,
    BackTopComponent,
    AlertComponent,
    DialogComponent,
    GlobalAlertComponent
  ],
  entryComponents: [
    AlertComponent,
    DialogComponent,
    GlobalAlertComponent
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
    MessageService,
    HttpInterceptorService
  ],
})
export class SharedModule {
}
