import { NgModule } from '@angular/core';
import { FormsModule, ReactiveFormsModule } from '@angular/forms';
import { HttpClientModule} from '@angular/common/http'
import { ClarityModule } from '@clr/angular';
import { CookieModule } from 'ngx-cookie';
import { TranslateModule, TranslateLoader } from '@ngx-translate/core';
import { CustomTranslateLoader } from '../i18n/custom-translate-loader';
import { BrowserAnimationsModule } from "@angular/platform-browser/animations";
import { BrowserModule } from '@angular/platform-browser';

@NgModule({
  imports:[
    ClarityModule.forRoot(),
    CookieModule.forRoot(),
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
    HttpClientModule,
    FormsModule,
    ReactiveFormsModule,
    ClarityModule,
    CookieModule,
    TranslateModule
  ]
})
export class CoreModule {}