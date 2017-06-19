import { NgModule } from '@angular/core';
import { RouterModule } from '@angular/router';
import { SharedModule } from '../shared/shared.module';
import { HeaderComponent } from './header/header.component';
import { MainContentComponent } from './main-content.component';
import { ClarityModule } from 'clarity-angular';

@NgModule({
  imports: [
    SharedModule,
    ClarityModule,
    RouterModule
  ],
  declarations: [
    HeaderComponent,
    MainContentComponent
  ],
  exports: [
    MainContentComponent
  ]
})
export class MainContentModule {}