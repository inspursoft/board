import { NgModule } from '@angular/core';

import { AccountModule } from './account/account.module';
import { MainContentModule } from './main-content/main-content.module';
import { FeatureModule } from './common/feature.module';
import { AppComponent } from './app.component';

import { ROUTING } from './app.routing';

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
  providers: [],
  bootstrap: [ AppComponent ]
})
export class AppModule {}
