import { Component } from '@angular/core';
import { AppInitService } from './app.init.service';
import { TranslateService } from '@ngx-translate/core';
import { CookieService } from "ngx-cookie";

@Component({
  selector: 'board-app',
  templateUrl: './app.component.html'
})
export class AppComponent {

  cookieExpiry: Date = new Date(2018, 2, 1);

  constructor(private appInitService: AppInitService,
              private cookieService: CookieService,
              private translateService: TranslateService) {
    if (!cookieService.get('currentLang')) {
      console.log('No found cookie for current lang, will use the default browser language.');
      cookieService.put('currentLang', this.translateService.getBrowserCultureLang(), {expires: this.cookieExpiry});
    }
    this.appInitService.currentLang = cookieService.get('currentLang') || 'en-us';
    translateService.use(this.appInitService.currentLang);
    this.translateService.onLangChange.subscribe(() => {
      this.appInitService.currentLang = this.translateService.currentLang;
      cookieService.put('currentLang', this.appInitService.currentLang, {expires: this.cookieExpiry});
      console.log('Change lang to:' + this.appInitService.currentLang);
    });
  }
}
