import { Component } from '@angular/core';
import { AppInitService } from './shared.service/app-init.service';
import { TranslateService, LangChangeEvent } from '@ngx-translate/core';
import { registerLocaleData } from '@angular/common';
import localeZhHans from '@angular/common/locales/zh-Hans';

@Component({
  selector: 'app-root',
  templateUrl: './app.component.html',
  styleUrls: ['./app.component.css']
})
export class AppComponent {
  constructor(private appInitService: AppInitService,
              private translateService: TranslateService) {

    if (!window.localStorage.getItem('currentLang')) {
      console.log('No found cookie for current lang, will use the default browser language.');
      window.localStorage.setItem('currentLang', this.translateService.getBrowserCultureLang());
    }

    this.appInitService.currentLang = window.localStorage.getItem('currentLang') || 'en-us';

    translateService.use(this.appInitService.currentLang);

    this.translateService.onLangChange.subscribe((res: LangChangeEvent) => {
      const oldLang = this.appInitService.currentLang;
      this.appInitService.currentLang = this.translateService.currentLang;
      window.localStorage.setItem('currentLang', this.appInitService.currentLang);
      if (res.lang.toLocaleLowerCase() !== oldLang.toLocaleLowerCase()) {
        window.location.reload(true);
      }
    });

    if (appInitService.currentLang.toLowerCase() === 'zh-cn') {
      registerLocaleData(localeZhHans, 'zh-Hans');
    }
  }
}
