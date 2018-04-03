import { NgModule, APP_INITIALIZER, NO_ERRORS_SCHEMA } from '@angular/core';
import { AccountModule } from './account/account.module';
import { MainContentModule } from './main-content/main-content.module';
import { FeatureModule } from './common/feature.module';
import { AppComponent } from './app.component';
import { AppInitService, AppTokenService } from './app.init.service';
import { ROUTING, SystemInfoResolve } from './app.routing';
import { HTTP_INTERCEPTORS, HttpClientModule } from "@angular/common/http";
import { HttpClientInterceptor } from "./shared/http-interceptor/http-client-interceptor";

export function appInitServiceFactory(appInitService: AppInitService) {
  return () => (appInitService);
}

@NgModule({
  imports: [
    AccountModule,
    MainContentModule,
    FeatureModule,
    HttpClientModule,
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
      deps: [AppTokenService],
      multi: true
    },
    SystemInfoResolve
  ],

  bootstrap: [AppComponent]
})
export class AppModule {
}
