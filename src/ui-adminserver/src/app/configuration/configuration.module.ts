import { NgModule } from '@angular/core';
import { CommonModule } from '@angular/common';
import { CfgCardsComponent } from './cfg-cards.component';
import { RouterModule } from '@angular/router';
import { ToTopDirective } from '../shared/to-top.directive';
import { ClarityModule } from '@clr/angular';
import { BoardComponentsLibraryModule } from 'board-components-library';
import { BrowserModule } from '@angular/platform-browser';
import { BrowserAnimationsModule } from '@angular/platform-browser/animations';
import { SharedModule } from '../shared/shared.module';
import { TranslateModule } from '@ngx-translate/core';
import { FormsModule } from '@angular/forms';

@NgModule({
  declarations: [
    CfgCardsComponent,
    ToTopDirective
  ],
  providers: [ ],
  imports: [
    BrowserModule,
    BrowserAnimationsModule,
    CommonModule,
    FormsModule,
    ClarityModule,
    BoardComponentsLibraryModule,
    SharedModule,
    TranslateModule,
    RouterModule.forChild([{ path: '', component: CfgCardsComponent }])
  ]
})
export class ConfigurationModule { }
