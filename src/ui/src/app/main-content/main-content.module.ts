import { NgModule } from '@angular/core';
import { SharedModule } from '../shared/shared.module';
import { MainContentComponent } from './main-content.component';

@NgModule({
  imports: [
    SharedModule
  ],
  declarations: [
    MainContentComponent
  ],
  exports: [
    MainContentComponent
  ]
})
export class MainContentModule {}