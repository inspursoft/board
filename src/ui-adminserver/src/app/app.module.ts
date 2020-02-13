import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { ClarityModule } from '@clr/angular';
import { TranslateLoader, TranslateModule } from '@ngx-translate/core';

import { AppRoutingModule } from './app-routing.module';
import { AppComponent } from './app.component';
import { MainContentComponent } from './main-content/main-content.component';
import { SharedModule } from './shared/shared.module';
import { DashboardModule } from './dashboard/dashboard.module';
import { SharedServiceModule } from './shared.service/shared-service.module';
import { CustomTranslateLoader } from './i18n/custom-translate-loader';
import { ConfigurationModule } from './configuration/configuration.module';
import { HttpClientModule } from '@angular/common/http';
import { AccountModule } from './account/account.module';
import { ResourceModule } from './resource/resource.module';

@NgModule({
  declarations: [
    AppComponent,
    MainContentComponent
  ],
  imports: [
    BrowserModule,
    SharedModule,
    AppRoutingModule,
    AccountModule,
    ClarityModule,
    DashboardModule,
    ConfigurationModule,
    SharedServiceModule,
    TranslateModule.forRoot({
      loader: {
        provide: TranslateLoader,
        useClass: CustomTranslateLoader
      }
    }),
    HttpClientModule,
    ResourceModule,
  ],
  providers: [],
  bootstrap: [AppComponent]
})
export class AppModule { }
