import { NgModule, APP_INITIALIZER } from '@angular/core';

import { AccountModule } from './account/account.module';
import { MainContentModule } from './main-content/main-content.module';
import { FeatureModule } from './common/feature.module';
import { AppComponent } from './app.component';

import { AppInitService } from './app.init.service';

import { ROUTING, SystemInfoResolve } from './app.routing';

export function appInitServiceFactory(appInitService: AppInitService) {
  return () => (appInitService);
}

@NgModule({
  imports: [
    AccountModule,
    MainContentModule,
    FeatureModule,
    ROUTING
  ],
  declarations: [
    AppComponent
  ],
  providers: [
    AppInitService,
    {
      provide: APP_INITIALIZER,
      useFactory: appInitServiceFactory,
      deps: [ AppInitService ],
      multi: true
    },
    SystemInfoResolve
  ],
  bootstrap: [ AppComponent ]
})
export class AppModule {}
