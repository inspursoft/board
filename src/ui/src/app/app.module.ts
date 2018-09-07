import { NgModule, APP_INITIALIZER } from '@angular/core';
import { FeatureModule } from './common/feature.module';
import { AppComponent } from './app.component';
import { AppInitService, AppTokenService } from './app.init.service';
import { ROUTING, SystemInfoResolve } from './app.routing';
import { HTTP_INTERCEPTORS } from "@angular/common/http";
import { HttpClientInterceptor } from "./shared/http-interceptor/http-client-interceptor";
import { SharedModule } from "./shared/shared.module";
import { MessageService } from "./shared/message-service/message.service";

export function appInitServiceFactory(appInitService: AppInitService) {
  return () => (appInitService);
}

@NgModule({
  imports: [
    FeatureModule,
    SharedModule,
    ROUTING
  ],
  declarations: [
    AppComponent
  ],
  providers: [
    AppTokenService,
    AppInitService,
    {
      provide: APP_INITIALIZER,
      useFactory: appInitServiceFactory,
      deps: [ AppInitService ],
      multi: true
    },
    {
      provide: HTTP_INTERCEPTORS,
      useClass: HttpClientInterceptor,
      deps: [AppTokenService,MessageService],
      multi: true
    },
    SystemInfoResolve
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
