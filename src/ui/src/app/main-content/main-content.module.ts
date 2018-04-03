import { NgModule, NO_ERRORS_SCHEMA } from '@angular/core';
import { RouterModule } from '@angular/router';
import { SharedModule } from '../shared/shared.module';
import { MainContentComponent } from './main-content.component';

@NgModule({
  imports: [
    SharedModule,
    RouterModule
  ],
  declarations: [
    MainContentComponent
  ],

  exports: [
    MainContentComponent
  ]
})
export class MainContentModule {}