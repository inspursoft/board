import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { BrowserModule } from '@angular/platform-browser';
import { NgModule } from '@angular/core';
import { FormsModule } from '@angular/forms';
import { HttpModule, Http } from '@angular/http';

import { ClarityModule } from 'clarity-angular';
import { TranslateModule, TranslateLoader } from '@ngx-translate/core';

import { CustomTranslateLoader } from '../i18n/custom-translate-loader';

@NgModule({
  imports:[
    ClarityModule.forRoot(),
    TranslateModule.forRoot({
      loader: {
        provide: TranslateLoader,
        useClass: CustomTranslateLoader
      }
    })
  ],
  exports:[
    BrowserAnimationsModule,
    BrowserModule,
    FormsModule,
    HttpModule,
    ClarityModule,
    TranslateModule
  ]
})
export class CoreModule {}