import { Component } from '@angular/core';
import { Router } from '@angular/router';
import { AppInitService } from './app.init.service';
import { TranslateService } from '@ngx-translate/core';
import { NgXCookies } from 'ngx-cookies';

@Component({
    selector: 'board-app',
    templateUrl: './app.component.html'
})
export class AppComponent {

    cookieExpiry: number =  60*24*365;

    constructor(
      private appInitService: AppInitService,
      private translateService: TranslateService
    ) {
      if(!NgXCookies.exists('currentLang')) {
        console.log('No found cookie for current lang, will use the default browser language.');
        NgXCookies.setCookie('currentLang', this.translateService.getBrowserCultureLang(), this.cookieExpiry);
      } 
      this.appInitService.currentLang = NgXCookies.getCookie('currentLang') || 'en-us';
      translateService.use(this.appInitService.currentLang);
      this.translateService.onLangChange.subscribe(()=>{
        this.appInitService.currentLang = this.translateService.currentLang;
        NgXCookies.setCookie('currentLang', this.appInitService.currentLang, this.cookieExpiry);
        console.log('Change lang to:' + this.appInitService.currentLang);
      });
    }
}
