import { Component, Input } from '@angular/core';
import { Router } from '@angular/router';

import { TranslateService, LangChangeEvent } from '@ngx-translate/core';


@Component({
  selector: 'header-content',
  templateUrl: 'header.component.html',
  styleUrls: [ 'header.component.css' ]
})
export class HeaderComponent {

  currentLang: string;
  @Input() isSignIn: boolean;
  @Input() hasSignedIn: boolean;

  get brandLogoUrl(): string {
    return this.isSignIn ? '../../images/board-blue.jpg': '../../../images/board.png';
  }

  constructor(
    private router: Router,
    private translateService: TranslateService) {
    let lang: string = this.translateService.getBrowserCultureLang();
    this._assertLanguage(lang);
  }

  _assertLanguage(lang: string) {
    lang = lang.toLowerCase();
    switch(lang) {
    case 'en':
    case 'en-us':
      lang = 'en-us';
      this.currentLang = 'HEAD_NAV.LANG_EN_US';
      break;

    case 'zh':
    case 'zh-cn': 
      lang = 'zh-cn';
      this.currentLang = 'HEAD_NAV.LANG_ZH_CN';
      break;
    }
    this.translateService.use(lang);
  }

  changLanguage(lang: string) {
    this._assertLanguage(lang);
  }

  logOut() {
    this.router.navigate(['/sign-in']);
  }
}