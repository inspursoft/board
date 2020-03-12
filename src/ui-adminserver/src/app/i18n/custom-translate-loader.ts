import { TranslateLoader } from '@ngx-translate/core';

import { Observable, of } from 'rxjs';
import { LANG_EN_US } from './en-us';
import { LANG_ZH_CN } from './zh-cn';

export class CustomTranslateLoader implements TranslateLoader {
  
  supportedLangs: {[key: string]: any} = {};

  constructor() {
    this.supportedLangs['en-us'] = LANG_EN_US;
    this.supportedLangs['zh-cn'] = LANG_ZH_CN;
  }
  getTranslation(lang: string): Observable<any> {
    lang = lang.toLowerCase();
    console.log('Current lang is:' + lang);
    return of(this.supportedLangs[lang]);
  }
}
