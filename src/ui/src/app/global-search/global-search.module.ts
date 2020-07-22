import { NgModule } from '@angular/core';
import { GlobalSearchComponent } from './global-search.component';
import { GlobalSearchService } from './global-search.service';
import { SharedModule } from '../shared/shared.module';
import { CoreModule } from '../core/core.module';

@NgModule({
  imports: [
    SharedModule,
    CoreModule
  ],
  declarations: [GlobalSearchComponent],
  providers: [GlobalSearchService]
})
export class GlobalSearchModule {

}

