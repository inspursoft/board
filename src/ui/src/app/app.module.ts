import { APP_INITIALIZER, LOCALE_ID, NgModule } from '@angular/core';
import { FeatureModule } from './common/feature.module';
import { AppComponent } from './app.component';
import { AppInitService, AppTokenService } from './app.init.service';
import { ROUTES, SystemInfoResolve } from './app.routing';
import { HTTP_INTERCEPTORS } from "@angular/common/http";
import { HttpClientInterceptor } from "./shared/http-interceptor/http-client-interceptor";
import { SharedModule } from "./shared/shared.module";
import { MessageService } from "./shared/message-service/message.service";
import { RouterModule } from "@angular/router";
import { CookieModule } from "ngx-cookie";
import { CustomTranslateLoader } from "./i18n/custom-translate-loader";
import { TranslateLoader, TranslateModule } from "@ngx-translate/core";
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { BrowserModule } from "@angular/platform-browser";

export function appInitServiceFactory(appInitService: AppInitService) {
  return () => (appInitService);
}

export function localIdServiceFactory(appInitService: AppInitService) {
  if (appInitService.currentLang == 'zh-cn') {
    return 'zh-Hans';
  } else {
    return 'en-US';
  }
}

@NgModule({
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    SharedModule,
    FeatureModule,
    CookieModule.forRoot(),
    TranslateModule.forRoot({
      loader: {
        provide: TranslateLoader,
        useClass: CustomTranslateLoader
      }
    }),
    RouterModule.forRoot(ROUTES)
  ],
  declarations: [
    AppComponent
  ],
  providers: [
    AppTokenService,
    AppInitService,
    MessageService,
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
    {
      provide: LOCALE_ID,
      useFactory: localIdServiceFactory,
      deps: [AppInitService]
    },
    SystemInfoResolve
  ],
  bootstrap: [AppComponent]
})
export class AppModule {
}
