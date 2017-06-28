import { NgModule } from '@angular/core';

import { CommonModule } from './common/common.module';
import { AppComponent } from './app.component';

import { ROUTING } from './app.routing';

@NgModule({
  imports: [
    CommonModule,
    ROUTING
  ],
  declarations: [
    AppComponent
  ],
  providers: [],
  bootstrap: [ AppComponent ]
})
export class AppModule {}
