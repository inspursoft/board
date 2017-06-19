import { NgModule } from '@angular/core';
import { CoreModule } from './core/core.module';

import { CommonModule } from './common/common.module';
import { AppComponent } from './app.component';

import { ROUTING } from './app.routing';


@NgModule({
  declarations: [
    AppComponent
  ],
  imports: [
    CoreModule,
    CommonModule,
    ROUTING
  ],
  providers: [],
  bootstrap: [ AppComponent ]
})
export class AppModule {}
