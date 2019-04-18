import { APP_INITIALIZER, LOCALE_ID, NgModule } from '@angular/core';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
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
import { SharedModule } from "./shared/shared.module";
import { GlobalSearchComponent } from "./global-search/global-search.component";

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

@NgModule({
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    CoreModule,
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
  ],
  declarations: [
    AppComponent,
    GlobalSearchComponent,
    MainContentComponent
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
    }
  ],
  bootstrap: [AppComponent]
})
export class AppModule {

}
