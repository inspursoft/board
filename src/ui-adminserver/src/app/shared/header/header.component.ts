import { Component, OnInit } from '@angular/core';
import { TranslateService } from '@ngx-translate/core';
import { AppInitService } from 'src/app/shared.service/app-init.service';
import { Router} from '@angular/router';

import 'src/assets/js/icon-translate/iconfont.js';

@Component({
  selector: 'app-header',
  templateUrl: './header.component.html',
  styleUrls: ['./header.component.css']
})
export class HeaderComponent implements OnInit {
  currentLang: string;
  display = true;

  constructor(private translateService: TranslateService,
              private appInitService: AppInitService,
              private router: Router) {
    this._assertLanguage(this.appInitService.currentLang);
  }

  ngOnInit() {
    this.display = location.pathname !== '/account/login' ? true : false;
  }

  _assertLanguage(lang: string) {
    lang = lang.toLowerCase();
    switch (lang) {
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

  logout() {
    window.sessionStorage.removeItem('token');
    this.router.navigateByUrl('account/login');
  }

}
