import { APP_INITIALIZER, LOCALE_ID, NgModule } from '@angular/core';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { HttpClientModule, HttpClientXsrfModule } from '@angular/common/http';
import { TranslateLoader, TranslateModule } from '@ngx-translate/core';
import { AppComponent } from './app.component';
import { AppInitService } from './shared.service/app-init.service';
import { CookieModule } from 'ngx-cookie';
import { CustomTranslateLoader } from './i18n/custom-translate-loader';
import { BrowserModule } from '@angular/platform-browser';
import { AppRoutingModule, SystemInfoResolve } from './app-routing.module';
import { MainContentComponent } from './main-content/main-content.component';
import { SharedServiceModule } from './shared.service/shared-service.module';
import { CoreModule } from './core/core.module';
import { SharedModule } from './shared/shared.module';
import { COMPONENTS_CUR_LANG } from 'board-components-library';
import { InitializePageComponent } from './initialize-page/initialize-page.component';
import { GlobalSearchModule } from './global-search/global-search.module';

export function appInitServiceFactory(appInitService: AppInitService) {
  return () => (appInitService);
}

export function localIdServiceFactory(appInitService: AppInitService) {
  if (appInitService.currentLang === 'zh-cn') {
    return 'zh-Hans';
  } else {
    return 'en-US';
  }
}

export function InitBoardLibraryLang(appInitService: AppInitService): string {
  return appInitService.currentLang === 'zh-cn' ? 'zh' : 'en';
}


@NgModule({
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    CoreModule,
    HttpClientModule,
    HttpClientXsrfModule.withOptions({cookieName: 'token'}),
    TranslateModule.forRoot({
      loader: {
        provide: TranslateLoader,
        useClass: CustomTranslateLoader
      }
    }),
    CookieModule.forRoot(),
    SharedModule,
    SharedServiceModule,
    AppRoutingModule,
    GlobalSearchModule,
  ],
  declarations: [
    AppComponent,
    MainContentComponent,
    InitializePageComponent
  ],
  providers: [
    SystemInfoResolve,
    {
      provide: APP_INITIALIZER,
      useFactory: appInitServiceFactory,
      deps: [AppInitService],
      multi: true
    },
    {
      provide: LOCALE_ID,
      useFactory: localIdServiceFactory,
      deps: [AppInitService]
    },
    {
      provide: COMPONENTS_CUR_LANG,
      useFactory: InitBoardLibraryLang,
      deps: [AppInitService]
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule {

}
