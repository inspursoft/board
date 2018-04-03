import { NgModule } from '@angular/core';

import { SharedModule } from '../shared/shared.module';
import { GlobalSearchComponent } from './global-search.component';
import { GlobalSearchService } from './global-search.service';

@NgModule({
  imports: [
    SharedModule
  ],

  declarations: [
    GlobalSearchComponent
  ],
  exports: [
    GlobalSearchComponent
  ],
  providers: [
    GlobalSearchService
  ]
})
export class GlobalSearchModule {

}