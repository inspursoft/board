import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule } from '@angular/http';

import { ClarityModule } from 'clarity-angular';
import { TranslateModule, TranslateLoader } from '@ngx-translate/core';

import { CustomTranslateLoader } from '../i18n/custom-translate-loader';

import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { BrowserModule } from '@angular/platform-browser';

@NgModule({
  imports:[
   ClarityModule.forRoot(),
    TranslateModule.forRoot({
      loader: {
        provide: TranslateLoader,
        useClass: CustomTranslateLoader
      }
    }),
  ],
  exports:[
    BrowserAnimationsModule,
    BrowserModule,
    HttpModule,
    FormsModule,
    ClarityModule,
    TranslateModule
  ]
})
export class CoreModule {}