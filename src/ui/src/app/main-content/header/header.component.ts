import { Component } from '@angular/core';
import { TranslateService, LangChangeEvent } from '@ngx-translate/core';


@Component({
  selector: 'header-content',
  templateUrl: 'header.component.html',
  styleUrls: [ 'header.component.css' ]
})
export class HeaderComponent {

  currentLang: string;

  get boardLogoIcon(): string {
    return '../../images/board.png';
  }

  constructor(private translateService: TranslateService) {
    this.currentLang = 'HEAD_NAV.LANG_EN_US';
  }

  changLanguage(lang: string) {
    this.translateService.use(lang);
    switch(lang) {
    case 'en-us':
      this.currentLang = 'HEAD_NAV.LANG_EN_US';
      break;
    case 'zh-cn': 
      this.currentLang = 'HEAD_NAV.LANG_ZH_CN';
      break;
    }
  }
}