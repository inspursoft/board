import { NgModule } from '@angular/core';
import { AppTokenService } from './app-token.service';
import { CsDialogComponent } from './cs-dialog/cs-dialog.component';
import { CsAlertComponent } from './cs-alert/cs-alert.component';
import { CsGlobalAlertComponent } from './cs-global-alert/cs-global-alert.component';
import { MessageService } from './message.service';
import { WebsocketService } from './websocket.service';
import { SharedService } from './shared.service';
import { SharedActionService } from './shared-action.service';
import { CoreModule } from '../core/core.module';
import { AppInitService } from './app-init.service';
import { AppGuardService } from './app-guard.service';
import { HttpInterceptorService } from './http-client-interceptor';

@NgModule({
  imports: [
    CoreModule
  ],
  declarations: [
    CsDialogComponent,
    CsAlertComponent,
    CsGlobalAlertComponent
  ],
  providers: [
    AppInitService,
    AppGuardService,
    AppTokenService,
    MessageService,
    WebsocketService,
    SharedService,
    SharedActionService,
    HttpInterceptorService,
  ],
  entryComponents: [
    CsDialogComponent,
    CsAlertComponent,
    CsGlobalAlertComponent
  ]
})
export class SharedServiceModule {

}
